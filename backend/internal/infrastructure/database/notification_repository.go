package database

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domainnotif "teman-mancing/backend/internal/domain/notification"
)

type GormNotificationRepository struct {
	db *gorm.DB
}

func NewGormNotificationRepository(db *gorm.DB) *GormNotificationRepository {
	return &GormNotificationRepository{db: db}
}

func (r *GormNotificationRepository) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domainnotif.Notification, int64, error) {
	var list []*domainnotif.Notification
	var total int64

	q := r.db.WithContext(ctx).Table("notifications").Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *GormNotificationRepository) Create(ctx context.Context, notif *domainnotif.Notification) error {
	return r.db.WithContext(ctx).Table("notifications").Create(notif).Error
}

func (r *GormNotificationRepository) MarkAsRead(ctx context.Context, notifID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Table("notifications").
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("is_read", true).Error
}

func (r *GormNotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Table("notifications").
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

func (r *GormNotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("notifications").
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

// ==========================================================
// InMemoryNotificationRepository for tests & dev fallback
// ==========================================================
type InMemoryNotificationRepository struct {
	mu     sync.RWMutex
	notifs map[uuid.UUID]*domainnotif.Notification
}

func NewInMemoryNotificationRepository() *InMemoryNotificationRepository {
	return &InMemoryNotificationRepository{
		notifs: make(map[uuid.UUID]*domainnotif.Notification),
	}
}

func (r *InMemoryNotificationRepository) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domainnotif.Notification, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domainnotif.Notification
	for _, n := range r.notifs {
		if n.UserID == userID {
			result = append(result, n)
		}
	}

	total := int64(len(result))
	if offset >= len(result) {
		return []*domainnotif.Notification{}, total, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], total, nil
}

func (r *InMemoryNotificationRepository) Create(ctx context.Context, notif *domainnotif.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifs[notif.ID] = notif
	return nil
}

func (r *InMemoryNotificationRepository) MarkAsRead(ctx context.Context, notifID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n, ok := r.notifs[notifID]; ok && n.UserID == userID {
		n.IsRead = true
	}
	return nil
}

func (r *InMemoryNotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, n := range r.notifs {
		if n.UserID == userID {
			n.IsRead = true
		}
	}
	return nil
}

func (r *InMemoryNotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, n := range r.notifs {
		if n.UserID == userID && !n.IsRead {
			count++
		}
	}
	return count, nil
}
