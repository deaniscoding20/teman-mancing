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
	Photos           []Photo            `json:"photos,omitempty"`
	FishTypes        []FishType         `json:"fish_types,omitempty"`
	Techniques       []FishingTechnique `json:"techniques,omitempty"`
	RecommendedBaits []RecommendedBait  `json:"recommended_baits,omitempty"`
	CommentsCount    int64              `json:"comments_count,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// Photo represents an uploaded image of a spot stored in MinIO.
type Photo struct {
	ID        uuid.UUID `json:"id"`
	SpotID    uuid.UUID `json:"spot_id"`
	ObjectKey string    `json:"object_key"`
	URL       string    `json:"url,omitempty"`
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
	UserName  string    `json:"user_name,omitempty"`
	UserAvatar *string  `json:"user_avatar,omitempty"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
