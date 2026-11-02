package fishing_spot

import (
	"time"

	"github.com/google/uuid"
)

// FishingSpot represents a fishing spot location.
type FishingSpot struct {
	ID               uuid.UUID          `json:"id"`
	UserID           uuid.UUID          `json:"user_id"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Address          string             `json:"address"`
	Latitude         float64            `json:"latitude"`
	Longitude        float64            `json:"longitude"`
	Photos           []Photo            `json:"photos,omitempty" gorm:"-"`
	FishTypes        []FishType         `json:"fish_types,omitempty" gorm:"-"`
	Techniques       []FishingTechnique `json:"techniques,omitempty" gorm:"-"`
	RecommendedBaits []RecommendedBait  `json:"recommended_baits,omitempty" gorm:"-"`
	CommentsCount    int64              `json:"comments_count,omitempty" gorm:"-"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// Photo represents an uploaded image of a spot stored in MinIO.
type Photo struct {
	ID        uuid.UUID `json:"id"`
	SpotID    uuid.UUID `json:"spot_id"`
	ObjectKey string    `json:"object_key"`
	URL       string    `json:"url,omitempty" gorm:"-"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

// FishType represents a catalogued fish species.
type FishType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// FishingTechnique represents a fishing method.
type FishingTechnique struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// RecommendedBait represents suggested baits for a specific spot.
type RecommendedBait struct {
	ID          uuid.UUID `json:"id"`
	SpotID      uuid.UUID `json:"spot_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
}

// Comment represents community reviews on a spot.
type Comment struct {
	ID        uuid.UUID `json:"id"`
	SpotID    uuid.UUID `json:"spot_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name,omitempty" gorm:"-"`
	UserAvatar *string  `json:"user_avatar,omitempty" gorm:"-"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
