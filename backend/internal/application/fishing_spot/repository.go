package fishing_spot

import (
	"context"

	"github.com/google/uuid"
	domainspot "teman-mancing/backend/internal/domain/fishing_spot"
)

// ListSpotsFilter contains query filtering options for spot catalog.
type ListSpotsFilter struct {
	Page      int
	Limit     int
	Query     string
	UserID    *uuid.UUID
	FishType  string
	Technique string
	Sort      string
}

// Repository defines database operations for fishing spots and related entities.
type Repository interface {
	List(ctx context.Context, filter ListSpotsFilter) ([]*domainspot.FishingSpot, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainspot.FishingSpot, error)
	Create(ctx context.Context, spot *domainspot.FishingSpot) error
	Update(ctx context.Context, spot *domainspot.FishingSpot) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Photos
	AddPhoto(ctx context.Context, photo *domainspot.Photo) error
	DeletePhoto(ctx context.Context, photoID uuid.UUID) error
	GetPhotoByID(ctx context.Context, photoID uuid.UUID) (*domainspot.Photo, error)

	// Comments
	ListComments(ctx context.Context, spotID uuid.UUID, page, limit int) ([]*domainspot.Comment, int64, error)
	GetCommentByID(ctx context.Context, commentID uuid.UUID) (*domainspot.Comment, error)
	GetSpotCommenterIDs(ctx context.Context, spotID uuid.UUID) ([]uuid.UUID, error)
	CreateComment(ctx context.Context, comment *domainspot.Comment) error
	DeleteComment(ctx context.Context, commentID uuid.UUID) error

	// Master Catalogs
	ListFishTypes(ctx context.Context) ([]domainspot.FishType, error)
	ListTechniques(ctx context.Context) ([]domainspot.FishingTechnique, error)
}
