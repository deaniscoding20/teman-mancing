package event

import (
	"context"

	"github.com/google/uuid"
	domainevent "teman-mancing/backend/internal/domain/event"
)

type ListEventsFilter struct {
	Page          int
	Limit         int
	Query         string
	CreatorID     *uuid.UUID
	ParticipantID *uuid.UUID
	Date          string
	Difficulty    string
	Visibility    string
	Status        string
}

type Repository interface {
	List(ctx context.Context, filter ListEventsFilter) ([]*domainevent.Event, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainevent.Event, error)
	Create(ctx context.Context, event *domainevent.Event) error
	Update(ctx context.Context, event *domainevent.Event) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Participation & ACID Locking
	JoinEventTx(ctx context.Context, eventID, userID uuid.UUID, isPrivate bool) (*domainevent.Participant, error)
	LeaveEvent(ctx context.Context, eventID, userID uuid.UUID) error
	ListParticipants(ctx context.Context, eventID uuid.UUID, status string) ([]*domainevent.Participant, error)
	GetParticipant(ctx context.Context, eventID, userID uuid.UUID) (*domainevent.Participant, error)
	UpdateParticipantStatus(ctx context.Context, eventID, userID uuid.UUID, status string) error
	CountApprovedParticipants(ctx context.Context, eventID uuid.UUID) (int, error)
}
