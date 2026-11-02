package database

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	appevent "teman-mancing/backend/internal/application/event"
	domainevent "teman-mancing/backend/internal/domain/event"
)

type GormEventRepository struct {
	db *gorm.DB
}

func NewGormEventRepository(db *gorm.DB) *GormEventRepository {
	return &GormEventRepository{db: db}
}

func (r *GormEventRepository) List(ctx context.Context, filter appevent.ListEventsFilter) ([]*domainevent.Event, int64, error) {
	var events []*domainevent.Event
	var total int64

	q := r.db.WithContext(ctx).Table("events")

	if filter.Query != "" {
		term := "%" + strings.ToLower(filter.Query) + "%"
		q = q.Where("LOWER(name) LIKE ? OR LOWER(destination_name) LIKE ? OR LOWER(meeting_point) LIKE ?", term, term, term)
	}
	if filter.Difficulty != "" {
		q = q.Where("difficulty = ?", filter.Difficulty)
	}
	if filter.Visibility != "" {
		q = q.Where("visibility = ?", filter.Visibility)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := q.Order("event_date ASC, start_time ASC").Offset(offset).Limit(filter.Limit).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	// Fetch participant counts and creator details
	for _, ev := range events {
		var count int64
		_ = r.db.WithContext(ctx).Table("event_participants").Where("event_id = ? AND status = 'approved'", ev.ID).Count(&count)
		ev.ApprovedParticipants = int(count)

		var u struct {
			Name      string
			AvatarKey *string
		}
		if err := r.db.WithContext(ctx).Table("users").Select("name, avatar_key").Where("id = ?", ev.CreatorID).Scan(&u).Error; err == nil {
			ev.CreatorName = u.Name
			ev.CreatorAvatar = u.AvatarKey
		}
	}

	return events, total, nil
}

func (r *GormEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainevent.Event, error) {
	var ev domainevent.Event
	err := r.db.WithContext(ctx).Table("events").Where("id = ?", id).First(&ev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var count int64
	_ = r.db.WithContext(ctx).Table("event_participants").Where("event_id = ? AND status = 'approved'", id).Count(&count)
	ev.ApprovedParticipants = int(count)

	var u struct {
		Name      string
		AvatarKey *string
	}
	if err := r.db.WithContext(ctx).Table("users").Select("name, avatar_key").Where("id = ?", ev.CreatorID).Scan(&u).Error; err == nil {
		ev.CreatorName = u.Name
		ev.CreatorAvatar = u.AvatarKey
	}

	return &ev, nil
}

func (r *GormEventRepository) Create(ctx context.Context, event *domainevent.Event) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create Event record
		if err := tx.Table("events").Create(event).Error; err != nil {
			return err
		}

		// 2. Add creator as first approved participant
		creatorParticipant := domainevent.Participant{
			ID:       uuid.New(),
			EventID:  event.ID,
			UserID:   event.CreatorID,
			Role:     domainevent.RoleCreator,
			Status:   domainevent.ParticipantStatusApproved,
			JoinedAt: time.Now().UTC(),
		}
		if err := tx.Table("event_participants").Create(&creatorParticipant).Error; err != nil {
			return err
		}

		// 3. Create group chat conversation for this event (README Section 9.2, 11.12)
		convID := uuid.New()
		conv := map[string]any{
			"id":         convID,
			"type":       "group",
			"event_id":   event.ID,
			"created_at": time.Now().UTC(),
		}
		if err := tx.Table("conversations").Create(&conv).Error; err != nil {
			return err
		}

		// 4. Add creator to conversation participants
		convPart := map[string]any{
			"id":              uuid.New(),
			"conversation_id": convID,
			"user_id":         event.CreatorID,
			"joined_at":       time.Now().UTC(),
		}
		return tx.Table("conversation_participants").Create(&convPart).Error
	})
}

func (r *GormEventRepository) Update(ctx context.Context, event *domainevent.Event) error {
	return r.db.WithContext(ctx).Table("events").Save(event).Error
}

func (r *GormEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Table("events").Where("id = ?", id).Delete(&domainevent.Event{}).Error
}

// JoinEventTx prevents race conditions using database row-level locking (README Section 12.1, 12.2)
func (r *GormEventRepository) JoinEventTx(ctx context.Context, eventID, userID uuid.UUID, isPrivate bool) (*domainevent.Participant, error) {
	var participant domainevent.Participant

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Row-level lock on event record
		var ev domainevent.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("events").Where("id = ?", eventID).First(&ev).Error; err != nil {
			return err
		}

		// 2. Prevent duplicate join first (README Section 11.11 UNIQUE)
		var existing domainevent.Participant
		err := tx.Table("event_participants").Where("event_id = ? AND user_id = ?", eventID, userID).First(&existing).Error
		if err == nil {
			return domainevent.ErrAlreadyParticipant
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 3. Validate event capacity & deadline
		var approvedCount int64
		if err := tx.Table("event_participants").Where("event_id = ? AND status = 'approved'", eventID).Count(&approvedCount).Error; err != nil {
			return err
		}

		if err := ev.CanAcceptParticipant(int(approvedCount)); err != nil {
			return err
		}

		// 4. Determine status
		status := domainevent.ParticipantStatusApproved
		if isPrivate {
			status = domainevent.ParticipantStatusPending
		}

		participant = domainevent.Participant{
			ID:       uuid.New(),
			EventID:  eventID,
			UserID:   userID,
			Role:     domainevent.RoleParticipant,
			Status:   status,
			JoinedAt: time.Now().UTC(),
		}

		if err := tx.Table("event_participants").Create(&participant).Error; err != nil {
			return err
		}

		// 5. If approved, add to event group chat conversation
		if status == domainevent.ParticipantStatusApproved {
			var convID uuid.UUID
			if err := tx.Table("conversations").Select("id").Where("event_id = ?", eventID).Scan(&convID).Error; err == nil && convID != uuid.Nil {
				convPart := map[string]any{
					"id":              uuid.New(),
					"conversation_id": convID,
					"user_id":         userID,
					"joined_at":       time.Now().UTC(),
				}
				_ = tx.Table("conversation_participants").Create(&convPart).Error
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *GormEventRepository) LeaveEvent(ctx context.Context, eventID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Table("event_participants").Where("event_id = ? AND user_id = ?", eventID, userID).Delete(&domainevent.Participant{}).Error
}

func (r *GormEventRepository) ListParticipants(ctx context.Context, eventID uuid.UUID, status string) ([]*domainevent.Participant, error) {
	var list []*domainevent.Participant
	q := r.db.WithContext(ctx).Table("event_participants").Where("event_id = ?", eventID)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Order("joined_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}

	// Enrich with user name and avatar
	for _, p := range list {
		var u struct {
			Name      string
			AvatarKey *string
		}
		if err := r.db.WithContext(ctx).Table("users").Select("name, avatar_key").Where("id = ?", p.UserID).Scan(&u).Error; err == nil {
			p.UserName = u.Name
			p.UserAvatar = u.AvatarKey
		}
	}

	return list, nil
}

func (r *GormEventRepository) GetParticipant(ctx context.Context, eventID, userID uuid.UUID) (*domainevent.Participant, error) {
	var p domainevent.Participant
	err := r.db.WithContext(ctx).Table("event_participants").Where("event_id = ? AND user_id = ?", eventID, userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *GormEventRepository) UpdateParticipantStatus(ctx context.Context, eventID, userID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("event_participants").Where("event_id = ? AND user_id = ?", eventID, userID).Update("status", status).Error; err != nil {
			return err
		}

		// When approved, add to conversation
		if status == domainevent.ParticipantStatusApproved {
			var convID uuid.UUID
			if err := tx.Table("conversations").Select("id").Where("event_id = ?", eventID).Scan(&convID).Error; err == nil && convID != uuid.Nil {
				convPart := map[string]any{
					"id":              uuid.New(),
					"conversation_id": convID,
					"user_id":         userID,
					"joined_at":       time.Now().UTC(),
				}
				_ = tx.Table("conversation_participants").Create(&convPart).Error
			}
		}
		return nil
	})
}

func (r *GormEventRepository) CountApprovedParticipants(ctx context.Context, eventID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("event_participants").Where("event_id = ? AND status = 'approved'", eventID).Count(&count).Error
	return int(count), err
}

// ==========================================================
// InMemoryEventRepository for fast testing & offline fallback
// ==========================================================
type InMemoryEventRepository struct {
	mu           sync.Mutex
	events       map[uuid.UUID]*domainevent.Event
	participants map[string]*domainevent.Participant // key: eventID_userID
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		events:       make(map[uuid.UUID]*domainevent.Event),
		participants: make(map[string]*domainevent.Participant),
	}
}

func (r *InMemoryEventRepository) List(ctx context.Context, filter appevent.ListEventsFilter) ([]*domainevent.Event, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []*domainevent.Event
	for _, ev := range r.events {
		if filter.Query != "" {
			q := strings.ToLower(filter.Query)
			if !strings.Contains(strings.ToLower(ev.Name), q) && !strings.Contains(strings.ToLower(ev.DestinationName), q) {
				continue
			}
		}
		if filter.Difficulty != "" && ev.Difficulty != filter.Difficulty {
			continue
		}
		if filter.Visibility != "" && ev.Visibility != filter.Visibility {
			continue
		}
		result = append(result, ev)
	}

	total := int64(len(result))
	start := (filter.Page - 1) * filter.Limit
	if start >= len(result) {
		return []*domainevent.Event{}, total, nil
	}
	end := start + filter.Limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], total, nil
}

func (r *InMemoryEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainevent.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ev, ok := r.events[id]
	if !ok {
		return nil, nil
	}
	return ev, nil
}

func (r *InMemoryEventRepository) Create(ctx context.Context, event *domainevent.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[event.ID] = event

	// Creator as participant
	key := event.ID.String() + "_" + event.CreatorID.String()
	r.participants[key] = &domainevent.Participant{
		ID:       uuid.New(),
		EventID:  event.ID,
		UserID:   event.CreatorID,
		Role:     domainevent.RoleCreator,
		Status:   domainevent.ParticipantStatusApproved,
		JoinedAt: time.Now().UTC(),
	}
	return nil
}

func (r *InMemoryEventRepository) Update(ctx context.Context, event *domainevent.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[event.ID] = event
	return nil
}

func (r *InMemoryEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.events, id)
	return nil
}

func (r *InMemoryEventRepository) JoinEventTx(ctx context.Context, eventID, userID uuid.UUID, isPrivate bool) (*domainevent.Participant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ev, ok := r.events[eventID]
	if !ok {
		return nil, appevent.ErrEventNotFound
	}

	key := eventID.String() + "_" + userID.String()
	if _, exists := r.participants[key]; exists {
		return nil, domainevent.ErrAlreadyParticipant
	}

	// Count approved
	approvedCount := 0
	for _, p := range r.participants {
		if p.EventID == eventID && p.Status == domainevent.ParticipantStatusApproved {
			approvedCount++
		}
	}

	if err := ev.CanAcceptParticipant(approvedCount); err != nil {
		return nil, err
	}

	status := domainevent.ParticipantStatusApproved
	if isPrivate {
		status = domainevent.ParticipantStatusPending
	} else {
		ev.ApprovedParticipants++
	}

	p := &domainevent.Participant{
		ID:       uuid.New(),
		EventID:  eventID,
		UserID:   userID,
		Role:     domainevent.RoleParticipant,
		Status:   status,
		JoinedAt: time.Now().UTC(),
	}
	r.participants[key] = p
	return p, nil
}

func (r *InMemoryEventRepository) LeaveEvent(ctx context.Context, eventID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := eventID.String() + "_" + userID.String()
	if p, ok := r.participants[key]; ok {
		if p.Status == domainevent.ParticipantStatusApproved {
			if ev, ok := r.events[eventID]; ok && ev.ApprovedParticipants > 0 {
				ev.ApprovedParticipants--
			}
		}
		delete(r.participants, key)
	}
	return nil
}

func (r *InMemoryEventRepository) ListParticipants(ctx context.Context, eventID uuid.UUID, status string) ([]*domainevent.Participant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*domainevent.Participant
	for _, p := range r.participants {
		if p.EventID == eventID {
			if status == "" || p.Status == status {
				list = append(list, p)
			}
		}
	}
	return list, nil
}

func (r *InMemoryEventRepository) GetParticipant(ctx context.Context, eventID, userID uuid.UUID) (*domainevent.Participant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := eventID.String() + "_" + userID.String()
	p, ok := r.participants[key]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *InMemoryEventRepository) UpdateParticipantStatus(ctx context.Context, eventID, userID uuid.UUID, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := eventID.String() + "_" + userID.String()
	p, ok := r.participants[key]
	if !ok {
		return appevent.ErrParticipantNotFound
	}
	prevStatus := p.Status
	p.Status = status
	if prevStatus != domainevent.ParticipantStatusApproved && status == domainevent.ParticipantStatusApproved {
		if ev, ok := r.events[eventID]; ok {
			ev.ApprovedParticipants++
		}
	}
	return nil
}

func (r *InMemoryEventRepository) CountApprovedParticipants(ctx context.Context, eventID uuid.UUID) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for _, p := range r.participants {
		if p.EventID == eventID && p.Status == domainevent.ParticipantStatusApproved {
			count++
		}
	}
	return count, nil
}
