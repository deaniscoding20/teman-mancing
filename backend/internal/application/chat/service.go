package chat

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	domainchat "teman-mancing/backend/internal/domain/chat"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrNotParticipant       = errors.New("you are not an authorized member of this conversation")
	ErrEmptyMessage         = errors.New("message content cannot be empty")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListUserConversations(ctx context.Context, userID uuid.UUID) ([]*domainchat.Conversation, error) {
	return s.repo.ListConversations(ctx, userID)
}

func (s *Service) GetOrCreateDirectConversation(ctx context.Context, currentUserID, targetUserID uuid.UUID) (*domainchat.Conversation, error) {
	if currentUserID == targetUserID {
		return nil, errors.New("cannot create conversation with yourself")
	}

	existing, err := s.repo.GetDirectConversation(ctx, currentUserID, targetUserID)
	if err == nil && existing != nil {
		return existing, nil
	}

	return s.repo.CreateDirectConversation(ctx, currentUserID, targetUserID)
}

func (s *Service) GetOrCreateEventGroupConversation(ctx context.Context, eventID uuid.UUID, eventName string, creatorID, actorID uuid.UUID) (*domainchat.Conversation, error) {
	conv, err := s.repo.GetConversationByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if conv == nil {
		conv, err = s.repo.CreateGroupConversation(ctx, &eventID, &eventName, creatorID)
		if err != nil {
			return nil, err
		}
	}

	// Ensure actor is a participant
	isMember, _ := s.repo.IsParticipant(ctx, conv.ID, actorID)
	if !isMember {
		_ = s.repo.AddParticipant(ctx, conv.ID, actorID)
	}

	return conv, nil
}

func (s *Service) GetMessages(ctx context.Context, convID, actorID uuid.UUID, limit, offset int) ([]*domainchat.Message, error) {
	// 1. Strict membership validation (README Section 12.3)
	isMember, err := s.repo.IsParticipant(ctx, convID, actorID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotParticipant
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.ListMessages(ctx, convID, limit, offset)
}

func (s *Service) SendMessage(ctx context.Context, convID, senderID uuid.UUID, senderName string, text string) (*domainchat.Message, error) {
	cleanText := strings.TrimSpace(text)
	if cleanText == "" {
		return nil, ErrEmptyMessage
	}

	// 1. Strict membership validation (README Section 12.3)
	isMember, err := s.repo.IsParticipant(ctx, convID, senderID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotParticipant
	}

	now := time.Now().UTC()
	msg := &domainchat.Message{
		ID:             uuid.New(),
		ConversationID: convID,
		SenderID:       senderID,
		SenderName:     senderName,
		Message:        cleanText,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}
