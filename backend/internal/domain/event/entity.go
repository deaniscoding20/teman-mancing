package event

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Difficulty levels
const (
	DifficultyPemula   = "pemula"
	DifficultyMenengah = "menengah"
	DifficultyAhli     = "ahli"
)

// Visibility
const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

// Event statuses
const (
	StatusOpen      = "open"
	StatusClosed    = "closed"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

// Participant roles
const (
	RoleCreator     = "creator"
	RoleCoHost      = "co-host"
	RoleParticipant = "participant"
)

// Participant statuses
const (
	ParticipantStatusPending   = "pending"
	ParticipantStatusApproved  = "approved"
	ParticipantStatusRejected  = "rejected"
	ParticipantStatusCancelled = "cancelled"
)

var (
	ErrEventCancelled     = errors.New("event has been cancelled")
	ErrEventClosed        = errors.New("event is closed for registration")
	ErrDeadlinePassed     = errors.New("event registration deadline has passed")
	ErrEventCapacityFull  = errors.New("event participant capacity has been reached")
	ErrAlreadyParticipant = errors.New("user is already a participant in this event")
)

// Event represents a fishing gathering (Janjian Mancing).
type Event struct {
	ID                   uuid.UUID     `json:"id"`
	CreatorID            uuid.UUID     `json:"creator_id"`
	CreatorName          string        `json:"creator_name,omitempty" gorm:"-"`
	CreatorAvatar        *string       `json:"creator_avatar,omitempty" gorm:"-"`
	Name                 string        `json:"name"`
	Description          string        `json:"description"`
	CoverImageKey        *string       `json:"cover_image_key,omitempty"`
	CoverImageURL        *string       `json:"cover_image_url,omitempty" gorm:"-"`
	FishingSpotID        *uuid.UUID    `json:"fishing_spot_id,omitempty"`
	DestinationName      string        `json:"destination_name"`
	DestinationAddress   string        `json:"destination_address"`
	Latitude             float64       `json:"latitude"`
	Longitude            float64       `json:"longitude"`
	EventDate            string        `json:"event_date"` // YYYY-MM-DD
	StartTime            string        `json:"start_time"` // HH:MM:SS
	EndTime              *string       `json:"end_time,omitempty"`
	MeetingPoint         string        `json:"meeting_point"`
	Difficulty           string        `json:"difficulty"`
	HasParkingFee        bool          `json:"has_parking_fee"`
	ParkingFee           *float64      `json:"parking_fee,omitempty"`
	RequiresBoat         bool          `json:"requires_boat"`
	BoatRentalFee        *float64      `json:"boat_rental_fee,omitempty"`
	AdditionalCost       *string       `json:"additional_cost,omitempty"`
	Visibility           string        `json:"visibility"`
	MaxParticipants      *int          `json:"max_participants,omitempty"`
	ApprovedParticipants int           `json:"approved_participants" gorm:"-"`
	RegistrationDeadline *time.Time    `json:"registration_deadline,omitempty"`
	Status               string        `json:"status"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
}

// Participant represents a user joining an event.
type Participant struct {
	ID         uuid.UUID `json:"id"`
	EventID    uuid.UUID `json:"event_id"`
	UserID     uuid.UUID `json:"user_id"`
	UserName   string    `json:"user_name,omitempty" gorm:"-"`
	UserAvatar *string   `json:"user_avatar,omitempty" gorm:"-"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	JoinedAt   time.Time `json:"joined_at"`
}

// CanAcceptParticipant checks if the event can accept another participant.
func (e *Event) CanAcceptParticipant(currentApprovedCount int) error {
	if e.Status == StatusCancelled {
		return ErrEventCancelled
	}
	if e.Status == StatusClosed || e.Status == StatusCompleted {
		return ErrEventClosed
	}
	if e.RegistrationDeadline != nil && time.Now().UTC().After(*e.RegistrationDeadline) {
		return ErrDeadlinePassed
	}
	if e.MaxParticipants != nil && currentApprovedCount >= *e.MaxParticipants {
		return ErrEventCapacityFull
	}
	return nil
}
