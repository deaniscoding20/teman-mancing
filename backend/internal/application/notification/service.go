package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainnotif "teman-mancing/backend/internal/domain/notification"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domainnotif.Notification, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.List(ctx, userID, limit, offset)
}

func (s *Service) CreateNotification(ctx context.Context, userID uuid.UUID, notifType, title, message string, refID *uuid.UUID) (*domainnotif.Notification, error) {
	notif := &domainnotif.Notification{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        notifType,
		Title:       title,
		Message:     message,
		ReferenceID: refID,
		IsRead:      false,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, notif); err != nil {
		return nil, err
	}
	return notif, nil
}

func (s *Service) MarkAsRead(ctx context.Context, notifID, userID uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, notifID, userID)
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *Service) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.repo.CountUnread(ctx, userID)
}
