package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appchat "teman-mancing/backend/internal/application/chat"
	appevent "teman-mancing/backend/internal/application/event"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/middleware"
	domainevent "teman-mancing/backend/internal/domain/event"
	infraWS "teman-mancing/backend/internal/infrastructure/websocket"
)

type ChatHandler struct {
	chatService *appchat.Service
	eventRepo   appevent.Repository
	hub         *infraWS.Hub
}

func NewChatHandler(chatService *appchat.Service, eventRepo appevent.Repository, hub *infraWS.Hub) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		eventRepo:   eventRepo,
		hub:         hub,
	}
}

type CreateConversationRequest struct {
	Type          string     `json:"type" binding:"required"` // direct or group
	ParticipantID *uuid.UUID `json:"participant_id"`
	EventID       *uuid.UUID `json:"event_id"`
}

type SendMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

// GET /conversations
func (h *ChatHandler) ListConversations(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	convs, err := h.chatService.ListUserConversations(c.Request.Context(), userID)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to retrieve conversations")
		return
	}

	dto.Success(c, http.StatusOK, "Conversations retrieved", convs)
}

// POST /conversations
func (h *ChatHandler) CreateConversation(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid request body", dto.FieldError{Field: "validation", Message: err.Error()})
		return
	}

	if req.Type == "group" {
		if req.EventID == nil {
			dto.Error(c, http.StatusBadRequest, "event_id wajib disertakan untuk percakapan grup")
			return
		}
		if h.eventRepo != nil {
			ev, err := h.eventRepo.GetByID(c.Request.Context(), *req.EventID)
			if err != nil || ev == nil {
				dto.Error(c, http.StatusNotFound, "Event tidak ditemukan")
				return
			}
			p, err := h.eventRepo.GetParticipant(c.Request.Context(), *req.EventID, userID)
			if err != nil || p == nil || p.Status != domainevent.ParticipantStatusApproved {
				dto.Error(c, http.StatusForbidden, "Hanya peserta yang telah disetujui yang dapat mengakses grup chat event")
				return
			}
			conv, err := h.chatService.GetOrCreateEventGroupConversation(c.Request.Context(), *req.EventID, ev.Name, ev.CreatorID, userID)
			if err != nil {
				dto.Error(c, http.StatusInternalServerError, "Gagal membuka grup chat event")
				return
			}
			dto.Success(c, http.StatusOK, "Grup chat event siap", conv)
			return
		}
	}

	if req.ParticipantID == nil {
		dto.Error(c, http.StatusBadRequest, "participant_id wajib disertakan untuk percakapan direct")
		return
	}

	conv, err := h.chatService.GetOrCreateDirectConversation(c.Request.Context(), userID, *req.ParticipantID)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusCreated, "Conversation ready", conv)
}

// GET /conversations/:id/messages
func (h *ChatHandler) ListMessages(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	convID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.chatService.GetMessages(c.Request.Context(), convID, userID, limit, offset)
	if err != nil {
		if errors.Is(err, appchat.ErrNotParticipant) {
			dto.Error(c, http.StatusForbidden, "You are not an authorized member of this conversation")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to retrieve messages")
		return
	}

	dto.Success(c, http.StatusOK, "Messages retrieved", messages)
}

// POST /conversations/:id/messages
func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}
	userName, _ := middleware.GetAuthUsername(c)

	convID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Message body cannot be empty")
		return
	}

	msg, err := h.chatService.SendMessage(c.Request.Context(), convID, userID, userName, req.Message)
	if err != nil {
		if errors.Is(err, appchat.ErrNotParticipant) {
			dto.Error(c, http.StatusForbidden, "You are not an authorized member of this conversation")
			return
		}
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if h.hub != nil {
		h.hub.BroadcastMessage(msg)
	}

	dto.Success(c, http.StatusCreated, "Message sent successfully", msg)
}

// GET /ws
func (h *ChatHandler) WebSocket(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "Unauthorized WebSocket connection")
		return
	}
	userName, _ := middleware.GetAuthUsername(c)

	h.hub.HandleWS(c.Writer, c.Request, userID, userName)
}
