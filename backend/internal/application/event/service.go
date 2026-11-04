package event

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	domainevent "teman-mancing/backend/internal/domain/event"
)

var (
	ErrEventNotFound       = errors.New("event not found")
	ErrUnauthorized        = errors.New("you are not authorized to manage this event")
	ErrParticipantNotFound = errors.New("participant not found")
)

type CreateEventInput struct {
	CreatorID            uuid.UUID  `json:"creator_id"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	CoverImageKey        *string    `json:"cover_image_key,omitempty"`
	FishingSpotID        *uuid.UUID `json:"fishing_spot_id,omitempty"`
	DestinationName      string     `json:"destination_name"`
	DestinationAddress   string     `json:"destination_address"`
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
	EventDate            string     `json:"event_date"` // YYYY-MM-DD
	StartTime            string     `json:"start_time"` // HH:MM:SS
	EndTime              *string    `json:"end_time,omitempty"`
	MeetingPoint         string     `json:"meeting_point"`
	MeetingAddress       *string    `json:"meeting_address,omitempty"`
	MeetingLatitude      *float64   `json:"meeting_latitude,omitempty"`
	MeetingLongitude     *float64   `json:"meeting_longitude,omitempty"`
	SameMeetingLocation  bool       `json:"same_meeting_location"`
	LocationNotes        *string    `json:"location_notes,omitempty"`
	Difficulty           string     `json:"difficulty"`
	HasParkingFee        bool       `json:"has_parking_fee"`
	ParkingFee           *float64   `json:"parking_fee,omitempty"`
	RequiresBoat         bool       `json:"requires_boat"`
	BoatRentalFee        *float64   `json:"boat_rental_fee,omitempty"`
	AdditionalCost       *string    `json:"additional_cost,omitempty"`
	Visibility           string     `json:"visibility"`
	MaxParticipants      *int       `json:"max_participants,omitempty"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty"`
}

type UpdateEventInput struct {
	Name                 *string    `json:"name"`
	Description          *string    `json:"description"`
	DestinationName      *string    `json:"destination_name"`
	DestinationAddress   *string    `json:"destination_address"`
	Latitude             *float64   `json:"latitude"`
	Longitude            *float64   `json:"longitude"`
	MeetingPoint         *string    `json:"meeting_point"`
	MeetingAddress       *string    `json:"meeting_address"`
	MeetingLatitude      *float64   `json:"meeting_latitude"`
	MeetingLongitude     *float64   `json:"meeting_longitude"`
	SameMeetingLocation  *bool      `json:"same_meeting_location"`
	LocationNotes        *string    `json:"location_notes"`
	Difficulty           *string    `json:"difficulty"`
	HasParkingFee        *bool      `json:"has_parking_fee"`
	ParkingFee           *float64   `json:"parking_fee"`
	RequiresBoat         *bool      `json:"requires_boat"`
	BoatRentalFee        *float64   `json:"boat_rental_fee"`
	AdditionalCost       *string    `json:"additional_cost"`
	MaxParticipants      *int       `json:"max_participants"`
	Status               *string    `json:"status"`
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListEvents(ctx context.Context, filter ListEventsFilter) ([]*domainevent.Event, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return s.repo.List(ctx, filter)
}

func (s *Service) GetEvent(ctx context.Context, id uuid.UUID) (*domainevent.Event, error) {
	ev, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ev == nil {
		return nil, ErrEventNotFound
	}
	return ev, nil
}

func (s *Service) CreateEvent(ctx context.Context, input CreateEventInput) (*domainevent.Event, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("event name is required")
	}
	if strings.TrimSpace(input.DestinationName) == "" {
		return nil, errors.New("destination name is required")
	}
	if strings.TrimSpace(input.MeetingPoint) == "" {
		return nil, errors.New("meeting point is required")
	}
	if input.Difficulty == "" {
		input.Difficulty = domainevent.DifficultyPemula
	}
	if input.Visibility == "" {
		input.Visibility = domainevent.VisibilityPublic
	}

	eventID := uuid.New()
	now := time.Now().UTC()

	ev := &domainevent.Event{
		ID:                   eventID,
		CreatorID:            input.CreatorID,
		Name:                 input.Name,
		Description:          input.Description,
		CoverImageKey:        input.CoverImageKey,
		FishingSpotID:        input.FishingSpotID,
		DestinationName:      input.DestinationName,
		DestinationAddress:   input.DestinationAddress,
		Latitude:             input.Latitude,
		Longitude:            input.Longitude,
		EventDate:            input.EventDate,
		StartTime:            input.StartTime,
		EndTime:              input.EndTime,
		MeetingPoint:         input.MeetingPoint,
		MeetingAddress:       input.MeetingAddress,
		MeetingLatitude:      input.MeetingLatitude,
		MeetingLongitude:     input.MeetingLongitude,
		SameMeetingLocation:  input.SameMeetingLocation,
		LocationNotes:        input.LocationNotes,
		Difficulty:           input.Difficulty,
		HasParkingFee:        input.HasParkingFee,
		ParkingFee:           input.ParkingFee,
		RequiresBoat:         input.RequiresBoat,
		BoatRentalFee:        input.BoatRentalFee,
		AdditionalCost:       input.AdditionalCost,
		Visibility:           input.Visibility,
		MaxParticipants:      input.MaxParticipants,
		ApprovedParticipants: 1, // Creator is the first approved participant
		RegistrationDeadline: input.RegistrationDeadline,
		Status:               domainevent.StatusOpen,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := s.repo.Create(ctx, ev); err != nil {
		return nil, err
	}

	return ev, nil
}

func (s *Service) UpdateEvent(ctx context.Context, eventID, actorID uuid.UUID, input UpdateEventInput) (*domainevent.Event, error) {
	ev, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if ev == nil {
		return nil, ErrEventNotFound
	}

	// Ownership check: only creator can manage event (README Section 12.4)
	if ev.CreatorID != actorID {
		return nil, ErrUnauthorized
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) != "" {
		ev.Name = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil {
		ev.Description = *input.Description
	}
	if input.DestinationName != nil && strings.TrimSpace(*input.DestinationName) != "" {
		ev.DestinationName = strings.TrimSpace(*input.DestinationName)
	}
	if input.DestinationAddress != nil {
		ev.DestinationAddress = *input.DestinationAddress
	}
	if input.MeetingPoint != nil {
		ev.MeetingPoint = *input.MeetingPoint
	}
	if input.MeetingAddress != nil {
		ev.MeetingAddress = input.MeetingAddress
	}
	if input.MeetingLatitude != nil {
		ev.MeetingLatitude = input.MeetingLatitude
	}
	if input.MeetingLongitude != nil {
		ev.MeetingLongitude = input.MeetingLongitude
	}
	if input.SameMeetingLocation != nil {
		ev.SameMeetingLocation = *input.SameMeetingLocation
	}
	if input.LocationNotes != nil {
		ev.LocationNotes = input.LocationNotes
	}
	if input.Difficulty != nil {
		ev.Difficulty = *input.Difficulty
	}
	if input.HasParkingFee != nil {
		ev.HasParkingFee = *input.HasParkingFee
	}
	if input.ParkingFee != nil {
		ev.ParkingFee = input.ParkingFee
	}
	if input.RequiresBoat != nil {
		ev.RequiresBoat = *input.RequiresBoat
	}
	if input.BoatRentalFee != nil {
		ev.BoatRentalFee = input.BoatRentalFee
	}
	if input.AdditionalCost != nil {
		ev.AdditionalCost = input.AdditionalCost
	}
	if input.MaxParticipants != nil {
		ev.MaxParticipants = input.MaxParticipants
	}
	if input.Status != nil {
		ev.Status = *input.Status
	}

	ev.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, ev); err != nil {
		return nil, err
	}

	return ev, nil
}

func (s *Service) CancelEvent(ctx context.Context, eventID, actorID uuid.UUID) error {
	ev, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if ev == nil {
		return ErrEventNotFound
	}

	if ev.CreatorID != actorID {
		return ErrUnauthorized
	}

	ev.Status = domainevent.StatusCancelled
	ev.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, ev)
}

// JoinEvent handles transactional membership with capacity and race condition locks (README Section 12.1, 12.2)
func (s *Service) JoinEvent(ctx context.Context, eventID, userID uuid.UUID) (*domainevent.Participant, error) {
	ev, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if ev == nil {
		return nil, ErrEventNotFound
	}

	isPrivate := ev.Visibility == domainevent.VisibilityPrivate
	return s.repo.JoinEventTx(ctx, eventID, userID, isPrivate)
}

func (s *Service) LeaveEvent(ctx context.Context, eventID, userID uuid.UUID) error {
	ev, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if ev == nil {
		return ErrEventNotFound
	}

	// Creator cannot leave own event without cancelling
	if ev.CreatorID == userID {
		return errors.New("event creator cannot leave; you may cancel the event instead")
	}

	return s.repo.LeaveEvent(ctx, eventID, userID)
}

func (s *Service) ListParticipants(ctx context.Context, eventID uuid.UUID, status string) ([]*domainevent.Participant, error) {
	return s.repo.ListParticipants(ctx, eventID, status)
}

func (s *Service) ApproveJoinRequest(ctx context.Context, eventID, targetUserID, actorID uuid.UUID) error {
	ev, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if ev == nil {
		return ErrEventNotFound
	}

	if ev.CreatorID != actorID {
		return ErrUnauthorized
	}

	// Check current approved capacity
	approvedCount, err := s.repo.CountApprovedParticipants(ctx, eventID)
	if err != nil {
		return err
	}

	if err := ev.CanAcceptParticipant(approvedCount); err != nil {
		return err
	}

	return s.repo.UpdateParticipantStatus(ctx, eventID, targetUserID, domainevent.ParticipantStatusApproved)
}

func (s *Service) RejectJoinRequest(ctx context.Context, eventID, targetUserID, actorID uuid.UUID) error {
	ev, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if ev == nil {
		return ErrEventNotFound
	}

	if ev.CreatorID != actorID {
		return ErrUnauthorized
	}

	return s.repo.UpdateParticipantStatus(ctx, eventID, targetUserID, domainevent.ParticipantStatusRejected)
}
