package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domainchat "teman-mancing/backend/internal/domain/chat"
)

type GormChatRepository struct {
	db *gorm.DB
}

func NewGormChatRepository(db *gorm.DB) *GormChatRepository {
	return &GormChatRepository{db: db}
}

func (r *GormChatRepository) ListConversations(ctx context.Context, userID uuid.UUID) ([]*domainchat.Conversation, error) {
	var convs []*domainchat.Conversation

	// Find conversation IDs where user is participant
	var convIDs []uuid.UUID
	err := r.db.WithContext(ctx).Table("conversation_participants").
		Where("user_id = ?", userID).
		Pluck("conversation_id", &convIDs).Error
	if err != nil {
		return nil, err
	}

	if len(convIDs) == 0 {
		return []*domainchat.Conversation{}, nil
	}

	if err := r.db.WithContext(ctx).Table("conversations").
		Where("id IN ?", convIDs).
		Order("created_at DESC").
		Find(&convs).Error; err != nil {
		return nil, err
	}

	for _, c := range convs {
		// If group conversation with event_id, fetch event name and cover
		if c.Type == domainchat.TypeGroup && c.EventID != nil {
			var ev struct {
				Name          string
				CoverImageKey *string
			}
			if err := r.db.WithContext(ctx).Table("events").Select("name, cover_image_key").Where("id = ?", c.EventID).Scan(&ev).Error; err == nil {
				c.EventName = &ev.Name
				c.EventCover = ev.CoverImageKey
			}
		}

		// Fetch last message
		var lastMsg domainchat.Message
		if err := r.db.WithContext(ctx).Table("messages").
			Where("conversation_id = ?", c.ID).
			Order("created_at DESC").
			First(&lastMsg).Error; err == nil {
			c.LastMessage = &lastMsg
		}
	}

	return convs, nil
}

func (r *GormChatRepository) GetConversationByID(ctx context.Context, convID uuid.UUID) (*domainchat.Conversation, error) {
	var conv domainchat.Conversation
	err := r.db.WithContext(ctx).Table("conversations").Where("id = ?", convID).First(&conv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *GormChatRepository) GetConversationByEventID(ctx context.Context, eventID uuid.UUID) (*domainchat.Conversation, error) {
	var conv domainchat.Conversation
	err := r.db.WithContext(ctx).Table("conversations").Where("event_id = ?", eventID).First(&conv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *GormChatRepository) CreateGroupConversation(ctx context.Context, eventID *uuid.UUID, name *string, creatorID uuid.UUID) (*domainchat.Conversation, error) {
	convID := uuid.New()
	now := time.Now().UTC()

	conv := &domainchat.Conversation{
		ID:        convID,
		Type:      domainchat.TypeGroup,
		EventID:   eventID,
		EventName: name,
		CreatedAt: now,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("conversations").Create(conv).Error; err != nil {
			return err
		}
		part := map[string]any{
			"id":              uuid.New(),
			"conversation_id": convID,
			"user_id":         creatorID,
			"joined_at":       now,
		}
		return tx.Table("conversation_participants").Create(&part).Error
	})
	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *GormChatRepository) GetDirectConversation(ctx context.Context, user1, user2 uuid.UUID) (*domainchat.Conversation, error) {
	var convID uuid.UUID
	query := `
		SELECT cp1.conversation_id 
		FROM conversation_participants cp1
		JOIN conversation_participants cp2 ON cp1.conversation_id = cp2.conversation_id
		JOIN conversations c ON c.id = cp1.conversation_id
		WHERE c.type = 'direct' AND cp1.user_id = ? AND cp2.user_id = ?
		LIMIT 1
	`
	err := r.db.WithContext(ctx).Raw(query, user1, user2).Scan(&convID).Error
	if err != nil || convID == uuid.Nil {
		return nil, nil
	}

	return r.GetConversationByID(ctx, convID)
}

func (r *GormChatRepository) CreateDirectConversation(ctx context.Context, user1, user2 uuid.UUID) (*domainchat.Conversation, error) {
	convID := uuid.New()
	now := time.Now().UTC()

	conv := &domainchat.Conversation{
		ID:        convID,
		Type:      domainchat.TypeDirect,
		CreatedAt: now,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("conversations").Create(conv).Error; err != nil {
			return err
		}

		p1 := map[string]any{"id": uuid.New(), "conversation_id": convID, "user_id": user1, "joined_at": now}
		p2 := map[string]any{"id": uuid.New(), "conversation_id": convID, "user_id": user2, "joined_at": now}

		if err := tx.Table("conversation_participants").Create(&p1).Error; err != nil {
			return err
		}
		return tx.Table("conversation_participants").Create(&p2).Error
	})

	if err != nil {
		return nil, err
	}
	return conv, nil
}

func (r *GormChatRepository) IsParticipant(ctx context.Context, convID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("conversation_participants").
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *GormChatRepository) AddParticipant(ctx context.Context, convID, userID uuid.UUID) error {
	part := map[string]any{
		"id":              uuid.New(),
		"conversation_id": convID,
		"user_id":         userID,
		"joined_at":       time.Now().UTC(),
	}
	return r.db.WithContext(ctx).Table("conversation_participants").Create(&part).Error
}

func (r *GormChatRepository) RemoveParticipant(ctx context.Context, convID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Table("conversation_participants").
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Delete(nil).Error
}

func (r *GormChatRepository) SaveMessage(ctx context.Context, msg *domainchat.Message) error {
	return r.db.WithContext(ctx).Table("messages").Create(msg).Error
}

func (r *GormChatRepository) ListMessages(ctx context.Context, convID uuid.UUID, limit, offset int) ([]*domainchat.Message, error) {
	var msgs []*domainchat.Message
	err := r.db.WithContext(ctx).Table("messages").
		Where("conversation_id = ?", convID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}

	// Enrich with sender names
	for _, m := range msgs {
		var u struct {
			Name      string
			AvatarKey *string
		}
		if err := r.db.WithContext(ctx).Table("users").Select("name, avatar_key").Where("id = ?", m.SenderID).Scan(&u).Error; err == nil {
			m.SenderName = u.Name
			m.SenderAvatar = u.AvatarKey
		}
	}

	return msgs, nil
}

// ==========================================================
// InMemoryChatRepository for tests & dev fallback
// ==========================================================
type InMemoryChatRepository struct {
	mu           sync.RWMutex
	convs        map[uuid.UUID]*domainchat.Conversation
	participants map[uuid.UUID][]uuid.UUID // convID -> []userID
	messages     map[uuid.UUID][]*domainchat.Message // convID -> messages
}

func NewInMemoryChatRepository() *InMemoryChatRepository {
	return &InMemoryChatRepository{
		convs:        make(map[uuid.UUID]*domainchat.Conversation),
		participants: make(map[uuid.UUID][]uuid.UUID),
		messages:     make(map[uuid.UUID][]*domainchat.Message),
	}
}

func (r *InMemoryChatRepository) ListConversations(ctx context.Context, userID uuid.UUID) ([]*domainchat.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domainchat.Conversation
	for convID, users := range r.participants {
		for _, u := range users {
			if u == userID {
				if c, ok := r.convs[convID]; ok {
					msgs := r.messages[convID]
					if len(msgs) > 0 {
						c.LastMessage = msgs[len(msgs)-1]
					}
					result = append(result, c)
				}
				break
			}
		}
	}
	return result, nil
}

func (r *InMemoryChatRepository) GetConversationByID(ctx context.Context, convID uuid.UUID) (*domainchat.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.convs[convID]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemoryChatRepository) GetConversationByEventID(ctx context.Context, eventID uuid.UUID) (*domainchat.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.convs {
		if c.EventID != nil && *c.EventID == eventID {
			return c, nil
		}
	}
	return nil, nil
}

func (r *InMemoryChatRepository) CreateGroupConversation(ctx context.Context, eventID *uuid.UUID, name *string, creatorID uuid.UUID) (*domainchat.Conversation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	convID := uuid.New()
	conv := &domainchat.Conversation{
		ID:        convID,
		Type:      domainchat.TypeGroup,
		EventID:   eventID,
		EventName: name,
		CreatedAt: time.Now().UTC(),
	}
	r.convs[convID] = conv
	r.participants[convID] = []uuid.UUID{creatorID}
	return conv, nil
}

func (r *InMemoryChatRepository) GetDirectConversation(ctx context.Context, user1, user2 uuid.UUID) (*domainchat.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for convID, users := range r.participants {
		c := r.convs[convID]
		if c != nil && c.Type == domainchat.TypeDirect && len(users) == 2 {
			hasU1 := users[0] == user1 || users[1] == user1
			hasU2 := users[0] == user2 || users[1] == user2
			if hasU1 && hasU2 {
				return c, nil
			}
		}
	}
	return nil, nil
}

func (r *InMemoryChatRepository) CreateDirectConversation(ctx context.Context, user1, user2 uuid.UUID) (*domainchat.Conversation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	convID := uuid.New()
	conv := &domainchat.Conversation{
		ID:        convID,
		Type:      domainchat.TypeDirect,
		CreatedAt: time.Now().UTC(),
	}
	r.convs[convID] = conv
	r.participants[convID] = []uuid.UUID{user1, user2}
	return conv, nil
}

func (r *InMemoryChatRepository) IsParticipant(ctx context.Context, convID, userID uuid.UUID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users, ok := r.participants[convID]
	if !ok {
		return false, nil
	}
	for _, u := range users {
		if u == userID {
			return true, nil
		}
	}
	return false, nil
}

func (r *InMemoryChatRepository) AddParticipant(ctx context.Context, convID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.participants[convID] = append(r.participants[convID], userID)
	return nil
}

func (r *InMemoryChatRepository) RemoveParticipant(ctx context.Context, convID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	users := r.participants[convID]
	var updated []uuid.UUID
	for _, u := range users {
		if u != userID {
			updated = append(updated, u)
		}
	}
	r.participants[convID] = updated
	return nil
}

func (r *InMemoryChatRepository) SaveMessage(ctx context.Context, msg *domainchat.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages[msg.ConversationID] = append(r.messages[msg.ConversationID], msg)
	return nil
}

func (r *InMemoryChatRepository) ListMessages(ctx context.Context, convID uuid.UUID, limit, offset int) ([]*domainchat.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgs := r.messages[convID]
	if offset >= len(msgs) {
		return []*domainchat.Message{}, nil
	}
	end := offset + limit
	if end > len(msgs) {
		end = len(msgs)
	}
	return msgs[offset:end], nil
}
