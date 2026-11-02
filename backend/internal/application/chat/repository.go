package chat

import (
	"context"

	"github.com/google/uuid"
	domainchat "teman-mancing/backend/internal/domain/chat"
)

type Repository interface {
	ListConversations(ctx context.Context, userID uuid.UUID) ([]*domainchat.Conversation, error)
	GetConversationByID(ctx context.Context, convID uuid.UUID) (*domainchat.Conversation, error)
	GetDirectConversation(ctx context.Context, user1, user2 uuid.UUID) (*domainchat.Conversation, error)
	CreateDirectConversation(ctx context.Context, user1, user2 uuid.UUID) (*domainchat.Conversation, error)

	// Membership check (README Section 12.3: Chat Authorization)
	IsParticipant(ctx context.Context, convID, userID uuid.UUID) (bool, error)
	AddParticipant(ctx context.Context, convID, userID uuid.UUID) error
	RemoveParticipant(ctx context.Context, convID, userID uuid.UUID) error

	// Messages
	SaveMessage(ctx context.Context, msg *domainchat.Message) error
	ListMessages(ctx context.Context, convID uuid.UUID, limit, offset int) ([]*domainchat.Message, error)
}
