package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appchat "teman-mancing/backend/internal/application/chat"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/middleware"
	infraWS "teman-mancing/backend/internal/infrastructure/websocket"
)

type ChatHandler struct {
	chatService *appchat.Service
	hub         *infraWS.Hub
}

func NewChatHandler(chatService *appchat.Service, hub *infraWS.Hub) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		hub:         hub,
	}
}

type CreateConversationRequest struct {
	Type          string    `json:"type" binding:"required"` // direct
	ParticipantID uuid.UUID `json:"participant_id" binding:"required"`
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

	conv, err := h.chatService.GetOrCreateDirectConversation(c.Request.Context(), userID, req.ParticipantID)
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
