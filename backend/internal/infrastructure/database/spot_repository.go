package database

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"

	appspot "teman-mancing/backend/internal/application/fishing_spot"
	domainspot "teman-mancing/backend/internal/domain/fishing_spot"
)

// GormSpotRepository handles PostgreSQL operations for fishing spots.
type GormSpotRepository struct {
	db *gorm.DB
}

func NewGormSpotRepository(db *gorm.DB) *GormSpotRepository {
	return &GormSpotRepository{db: db}
}

func (r *GormSpotRepository) List(ctx context.Context, filter appspot.ListSpotsFilter) ([]*domainspot.FishingSpot, int64, error) {
	var spots []*domainspot.FishingSpot
	var total int64

	q := r.db.WithContext(ctx).Table("fishing_spots")

	if filter.Query != "" {
		term := "%" + strings.ToLower(filter.Query) + "%"
		q = q.Where("LOWER(name) LIKE ? OR LOWER(address) LIKE ? OR LOWER(description) LIKE ?", term, term, term)
	}

	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := q.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&spots).Error; err != nil {
		return nil, 0, err
	}

	// Enrich spots with primary photos and comments count
	for _, s := range spots {
		var photo domainspot.Photo
		if err := r.db.WithContext(ctx).Table("fishing_spot_photos").Where("spot_id = ?", s.ID).Order("is_primary DESC, created_at ASC").First(&photo).Error; err == nil {
			photo.URL = fmt.Sprintf("/api/v1/media/fishing-spots/%s", photo.ObjectKey)
			s.Photos = []domainspot.Photo{photo}
			s.PrimaryPhotoURL = &photo.URL
		}
		var count int64
		_ = r.db.WithContext(ctx).Table("spot_comments").Where("spot_id = ?", s.ID).Count(&count)
		s.CommentsCount = count
	}

	return spots, total, nil
}

func (r *GormSpotRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainspot.FishingSpot, error) {
	var spot domainspot.FishingSpot
	err := r.db.WithContext(ctx).Table("fishing_spots").Where("id = ?", id).First(&spot).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch photos
	var photos []domainspot.Photo
	_ = r.db.WithContext(ctx).Table("fishing_spot_photos").Where("spot_id = ?", id).Order("is_primary DESC, created_at ASC").Find(&photos)
	for i := range photos {
		photos[i].URL = fmt.Sprintf("/api/v1/media/fishing-spots/%s", photos[i].ObjectKey)
		if (photos[i].IsPrimary || spot.PrimaryPhotoURL == nil) && spot.PrimaryPhotoURL == nil {
			spot.PrimaryPhotoURL = &photos[i].URL
		}
	}
	spot.Photos = photos

	// Fetch baits
	var baits []domainspot.RecommendedBait
	_ = r.db.WithContext(ctx).Table("recommended_baits").Where("spot_id = ?", id).Find(&baits)
	spot.RecommendedBaits = baits

	var count int64
	_ = r.db.WithContext(ctx).Table("spot_comments").Where("spot_id = ?", id).Count(&count)
	spot.CommentsCount = count

	// Fetch author details (Issue #19)
	var author domainspot.AuthorInfo
	if err := r.db.WithContext(ctx).Table("users").
		Select("id, name, username, avatar_key").
		Where("id = ?", spot.UserID).
		Scan(&author).Error; err == nil && author.ID != uuid.Nil {
		if author.AvatarKey != nil && *author.AvatarKey != "" {
			avatarURL := fmt.Sprintf("/api/v1/media/avatars/%s", *author.AvatarKey)
			author.AvatarURL = &avatarURL
		}
		spot.Author = &author
	}

	return &spot, nil
}

func (r *GormSpotRepository) Create(ctx context.Context, spot *domainspot.FishingSpot) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("fishing_spots").Create(spot).Error; err != nil {
			return err
		}
		for _, b := range spot.RecommendedBaits {
			if err := tx.Table("recommended_baits").Create(&b).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormSpotRepository) Update(ctx context.Context, spot *domainspot.FishingSpot) error {
	return r.db.WithContext(ctx).Table("fishing_spots").Save(spot).Error
}

func (r *GormSpotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Table("fishing_spots").Where("id = ?", id).Delete(&domainspot.FishingSpot{}).Error
}

func (r *GormSpotRepository) AddPhoto(ctx context.Context, photo *domainspot.Photo) error {
	return r.db.WithContext(ctx).Table("fishing_spot_photos").Create(photo).Error
}

func (r *GormSpotRepository) DeletePhoto(ctx context.Context, photoID uuid.UUID) error {
	return r.db.WithContext(ctx).Table("fishing_spot_photos").Where("id = ?", photoID).Delete(&domainspot.Photo{}).Error
}

func (r *GormSpotRepository) GetPhotoByID(ctx context.Context, photoID uuid.UUID) (*domainspot.Photo, error) {
	var photo domainspot.Photo
	err := r.db.WithContext(ctx).Table("fishing_spot_photos").Where("id = ?", photoID).First(&photo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *GormSpotRepository) ListComments(ctx context.Context, spotID uuid.UUID, page, limit int) ([]*domainspot.Comment, int64, error) {
	var rootComments []*domainspot.Comment
	var totalRoots int64

	q := r.db.WithContext(ctx).Table("spot_comments").Where("spot_id = ? AND parent_id IS NULL", spotID)
	if err := q.Count(&totalRoots).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&rootComments).Error; err != nil {
		return nil, 0, err
	}

	if len(rootComments) == 0 {
		return []*domainspot.Comment{}, totalRoots, nil
	}

	// Collect root comment IDs
	rootIDs := make([]uuid.UUID, len(rootComments))
	rootMap := make(map[uuid.UUID]*domainspot.Comment, len(rootComments))
	for i, rc := range rootComments {
		rc.Replies = make([]*domainspot.Comment, 0)
		rootIDs[i] = rc.ID
		rootMap[rc.ID] = rc
	}

	// Fetch replies belonging to these root comments
	var replies []*domainspot.Comment
	if err := r.db.WithContext(ctx).Table("spot_comments").
		Where("parent_id IN (?)", rootIDs).
		Order("created_at ASC").
		Find(&replies).Error; err != nil {
		return nil, 0, err
	}

	for _, rep := range replies {
		rep.Replies = make([]*domainspot.Comment, 0)
		if rep.ParentID != nil {
			if parent, exists := rootMap[*rep.ParentID]; exists {
				parent.Replies = append(parent.Replies, rep)
			}
		}
	}

	// Enrich user info for all comments (roots + replies)
	allComments := append(rootComments, replies...)
	userIDs := make(map[uuid.UUID]bool)
	for _, c := range allComments {
		userIDs[c.UserID] = true
	}

	if len(userIDs) > 0 {
		var uList []struct {
			ID        uuid.UUID
			Name      string
			AvatarKey *string
		}
		var ids []uuid.UUID
		for uid := range userIDs {
			ids = append(ids, uid)
		}
		if err := r.db.WithContext(ctx).Table("users").Select("id, name, avatar_key").Where("id IN (?)", ids).Scan(&uList).Error; err == nil {
			uMap := make(map[uuid.UUID]struct {
				Name      string
				AvatarKey *string
			})
			for _, u := range uList {
				uMap[u.ID] = struct {
					Name      string
					AvatarKey *string
				}{Name: u.Name, AvatarKey: u.AvatarKey}
			}
			for _, c := range allComments {
				if u, ok := uMap[c.UserID]; ok {
					c.UserName = u.Name
					c.UserAvatar = u.AvatarKey
				}
			}
		}
	}

	return rootComments, totalRoots, nil
}

func (r *GormSpotRepository) GetSpotCommenterIDs(ctx context.Context, spotID uuid.UUID) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID
	err := r.db.WithContext(ctx).Table("spot_comments").
		Distinct("user_id").
		Where("spot_id = ?", spotID).
		Pluck("user_id", &userIDs).Error
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (r *GormSpotRepository) GetCommentByID(ctx context.Context, commentID uuid.UUID) (*domainspot.Comment, error) {
	var c domainspot.Comment
	err := r.db.WithContext(ctx).Table("spot_comments").Where("id = ?", commentID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GormSpotRepository) CreateComment(ctx context.Context, comment *domainspot.Comment) error {
	return r.db.WithContext(ctx).Table("spot_comments").Create(comment).Error
}

func (r *GormSpotRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	return r.db.WithContext(ctx).Table("spot_comments").Where("id = ?", commentID).Delete(&domainspot.Comment{}).Error
}

func (r *GormSpotRepository) ListFishTypes(ctx context.Context) ([]domainspot.FishType, error) {
	var list []domainspot.FishType
	err := r.db.WithContext(ctx).Table("fish_types").Order("name ASC").Find(&list).Error
	return list, err
}

func (r *GormSpotRepository) ListTechniques(ctx context.Context) ([]domainspot.FishingTechnique, error) {
	var list []domainspot.FishingTechnique
	err := r.db.WithContext(ctx).Table("fishing_techniques").Order("name ASC").Find(&list).Error
	return list, err
}

// ==========================================================
// InMemorySpotRepository for fast testing & offline fallback
// ==========================================================
type InMemorySpotRepository struct {
	mu         sync.RWMutex
	spots      map[uuid.UUID]*domainspot.FishingSpot
	photos     map[uuid.UUID]*domainspot.Photo
	comments   map[uuid.UUID]*domainspot.Comment
	fishTypes  []domainspot.FishType
	techniques []domainspot.FishingTechnique
}

func NewInMemorySpotRepository() *InMemorySpotRepository {
	return &InMemorySpotRepository{
		spots:    make(map[uuid.UUID]*domainspot.FishingSpot),
		photos:   make(map[uuid.UUID]*domainspot.Photo),
		comments: make(map[uuid.UUID]*domainspot.Comment),
		fishTypes: []domainspot.FishType{
			{ID: uuid.New(), Name: "Kakap Putih"},
			{ID: uuid.New(), Name: "Nila"},
			{ID: uuid.New(), Name: "Kerapu"},
			{ID: uuid.New(), Name: "Tenggiri"},
		},
		techniques: []domainspot.FishingTechnique{
			{ID: uuid.New(), Name: "Casting"},
			{ID: uuid.New(), Name: "Dasaran"},
			{ID: uuid.New(), Name: "Jigging"},
		},
	}
}

func (r *InMemorySpotRepository) List(ctx context.Context, filter appspot.ListSpotsFilter) ([]*domainspot.FishingSpot, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domainspot.FishingSpot
	for _, s := range r.spots {
		if filter.Query != "" {
			q := strings.ToLower(filter.Query)
			if !strings.Contains(strings.ToLower(s.Name), q) && !strings.Contains(strings.ToLower(s.Address), q) {
				continue
			}
		}
		if filter.UserID != nil && s.UserID != *filter.UserID {
			continue
		}
		result = append(result, s)
	}

	total := int64(len(result))
	start := (filter.Page - 1) * filter.Limit
	if start >= len(result) {
		return []*domainspot.FishingSpot{}, total, nil
	}
	end := start + filter.Limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], total, nil
}

func (r *InMemorySpotRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainspot.FishingSpot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.spots[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *InMemorySpotRepository) Create(ctx context.Context, spot *domainspot.FishingSpot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.spots[spot.ID] = spot
	return nil
}

func (r *InMemorySpotRepository) Update(ctx context.Context, spot *domainspot.FishingSpot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.spots[spot.ID] = spot
	return nil
}

func (r *InMemorySpotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.spots, id)
	return nil
}

func (r *InMemorySpotRepository) AddPhoto(ctx context.Context, photo *domainspot.Photo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if photo.URL == "" && photo.ObjectKey != "" {
		photo.URL = fmt.Sprintf("/api/v1/media/fishing-spots/%s", photo.ObjectKey)
	}
	r.photos[photo.ID] = photo
	if s, ok := r.spots[photo.SpotID]; ok {
		s.Photos = append(s.Photos, *photo)
		if photo.IsPrimary || s.PrimaryPhotoURL == nil {
			s.PrimaryPhotoURL = &photo.URL
		}
	}
	return nil
}

func (r *InMemorySpotRepository) DeletePhoto(ctx context.Context, photoID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.photos[photoID]
	if ok {
		if s, exists := r.spots[p.SpotID]; exists {
			var updated []domainspot.Photo
			for _, item := range s.Photos {
				if item.ID != photoID {
					updated = append(updated, item)
				}
			}
			s.Photos = updated
		}
		delete(r.photos, photoID)
	}
	return nil
}

func (r *InMemorySpotRepository) GetPhotoByID(ctx context.Context, photoID uuid.UUID) (*domainspot.Photo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.photos[photoID]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *InMemorySpotRepository) ListComments(ctx context.Context, spotID uuid.UUID, page, limit int) ([]*domainspot.Comment, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var roots []*domainspot.Comment
	var replies []*domainspot.Comment

	for _, c := range r.comments {
		if c.SpotID == spotID {
			// Clone comment to safely populate replies without mutating store directly
			cloned := *c
			cloned.Replies = make([]*domainspot.Comment, 0)
			if c.ParentID == nil {
				roots = append(roots, &cloned)
			} else {
				replies = append(replies, &cloned)
			}
		}
	}

	// Sort roots by created_at DESC
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].CreatedAt.After(roots[j].CreatedAt)
	})

	// Sort replies by created_at ASC
	sort.Slice(replies, func(i, j int) bool {
		return replies[i].CreatedAt.Before(replies[j].CreatedAt)
	})

	rootMap := make(map[uuid.UUID]*domainspot.Comment)
	for _, root := range roots {
		rootMap[root.ID] = root
	}

	for _, rep := range replies {
		if rep.ParentID != nil {
			if parent, exists := rootMap[*rep.ParentID]; exists {
				parent.Replies = append(parent.Replies, rep)
			}
		}
	}

	total := int64(len(roots))
	start := (page - 1) * limit
	if start >= len(roots) {
		return []*domainspot.Comment{}, total, nil
	}
	end := start + limit
	if end > len(roots) {
		end = len(roots)
	}

	return roots[start:end], total, nil
}

func (r *InMemorySpotRepository) GetSpotCommenterIDs(ctx context.Context, spotID uuid.UUID) ([]uuid.UUID, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idMap := make(map[uuid.UUID]bool)
	for _, c := range r.comments {
		if c.SpotID == spotID {
			idMap[c.UserID] = true
		}
	}

	var userIDs []uuid.UUID
	for id := range idMap {
		userIDs = append(userIDs, id)
	}
	return userIDs, nil
}

func (r *InMemorySpotRepository) GetCommentByID(ctx context.Context, commentID uuid.UUID) (*domainspot.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.comments[commentID]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemorySpotRepository) CreateComment(ctx context.Context, comment *domainspot.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.comments[comment.ID] = comment
	if s, ok := r.spots[comment.SpotID]; ok {
		s.CommentsCount++
	}
	return nil
}

func (r *InMemorySpotRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.comments[commentID]
	if ok {
		if s, exists := r.spots[c.SpotID]; exists {
			s.CommentsCount--
		}
		delete(r.comments, commentID)
		// Cascade delete child replies
		for id, rep := range r.comments {
			if rep.ParentID != nil && *rep.ParentID == commentID {
				delete(r.comments, id)
				if s, exists := r.spots[rep.SpotID]; exists {
					s.CommentsCount--
				}
			}
		}
	}
	return nil
}

func (r *InMemorySpotRepository) ListFishTypes(ctx context.Context) ([]domainspot.FishType, error) {
	return r.fishTypes, nil
}

func (r *InMemorySpotRepository) ListTechniques(ctx context.Context) ([]domainspot.FishingTechnique, error) {
	return r.techniques, nil
}
