package chat

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	TypeDirect = "direct"
	TypeGroup  = "group"
)

var (
	ErrNotConversationMember = errors.New("user is not an authorized participant of this conversation")
	ErrEmptyMessage          = errors.New("message content cannot be empty")
)

// Conversation represents a chat room (direct message or event group chat).
type Conversation struct {
	ID           uuid.UUID                 `json:"id"`
	Type         string                    `json:"type"`
	EventID      *uuid.UUID                `json:"event_id,omitempty"`
	EventName    *string                   `json:"event_name,omitempty"`
	EventCover   *string                   `json:"event_cover,omitempty"`
	Participants []ConversationParticipant `json:"participants,omitempty"`
	LastMessage  *Message                  `json:"last_message,omitempty"`
	UnreadCount  int                       `json:"unread_count"`
	CreatedAt    time.Time                 `json:"created_at"`
}

// ConversationParticipant tracks room membership.
type ConversationParticipant struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	UserID         uuid.UUID `json:"user_id"`
	UserName       string    `json:"user_name,omitempty"`
	UserAvatar     *string   `json:"user_avatar,omitempty"`
	JoinedAt       time.Time `json:"joined_at"`
}

// Message represents a persisted chat message.
type Message struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	SenderName     string    `json:"sender_name,omitempty"`
	SenderAvatar   *string   `json:"sender_avatar,omitempty"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
