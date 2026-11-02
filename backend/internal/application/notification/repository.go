package notification

import (
	"context"

	"github.com/google/uuid"
	domainnotif "teman-mancing/backend/internal/domain/notification"
)

type Repository interface {
	List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domainnotif.Notification, int64, error)
	Create(ctx context.Context, notif *domainnotif.Notification) error
	MarkAsRead(ctx context.Context, notifID, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
}
