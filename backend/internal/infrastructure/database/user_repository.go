package database

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	appuser "teman-mancing/backend/internal/application/user"
	domainuser "teman-mancing/backend/internal/domain/user"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, u *domainuser.User) error {
	return r.db.WithContext(ctx).Table("users").Create(u).Error
}

func (r *GormUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	var u domainuser.User
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	var u domainuser.User
	email = strings.ToLower(strings.TrimSpace(email))
	err := r.db.WithContext(ctx).Table("users").Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*domainuser.User, error) {
	var u domainuser.User
	username = strings.ToLower(strings.TrimSpace(username))
	err := r.db.WithContext(ctx).Table("users").Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepository) Update(ctx context.Context, u *domainuser.User) error {
	return r.db.WithContext(ctx).Table("users").Save(u).Error
}

func (r *GormUserRepository) GetProfileView(ctx context.Context, id uuid.UUID) (*appuser.UserProfileView, error) {
	u, err := r.GetByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}

	var spotsCount int64
	var eventsCreatedCount int64
	var eventsJoinedCount int64

	_ = r.db.WithContext(ctx).Table("fishing_spots").Where("user_id = ?", id).Count(&spotsCount)
	_ = r.db.WithContext(ctx).Table("events").Where("creator_id = ?", id).Count(&eventsCreatedCount)
	_ = r.db.WithContext(ctx).Table("event_participants").Where("user_id = ? AND status = 'approved'", id).Count(&eventsJoinedCount)

	return &appuser.UserProfileView{
		User: u,
		Stats: appuser.UserStats{
			SpotsCount:         spotsCount,
			EventsCreatedCount: eventsCreatedCount,
			EventsJoinedCount:  eventsJoinedCount,
		},
	}, nil
}
