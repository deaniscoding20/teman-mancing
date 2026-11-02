package user

import (
	"context"

	"github.com/google/uuid"
	"teman-mancing/backend/internal/domain/user"
)

// UserStats bundles counter metrics for user profile.
type UserStats struct {
	SpotsCount         int64 `json:"spots_count"`
	EventsCreatedCount int64 `json:"events_created_count"`
	EventsJoinedCount  int64 `json:"events_joined_count"`
}

// UserProfileView extends user entity with stats.
type UserProfileView struct {
	User  *user.User
	Stats UserStats
}

// Repository defines storage operations for users.
type Repository interface {
	Create(ctx context.Context, u *user.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
	GetProfileView(ctx context.Context, id uuid.UUID) (*UserProfileView, error)
}
