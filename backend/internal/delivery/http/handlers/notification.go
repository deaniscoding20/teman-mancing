package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appnotif "teman-mancing/backend/internal/application/notification"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/middleware"
)

type NotificationHandler struct {
	notifService *appnotif.Service
}

func NewNotificationHandler(notifService *appnotif.Service) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

// GET /notifications
func (h *NotificationHandler) List(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	list, total, err := h.notifService.ListNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to retrieve notifications")
		return
	}

	unreadCount, _ := h.notifService.GetUnreadCount(c.Request.Context(), userID)

	dto.SuccessWithMeta(c, http.StatusOK, "Notifications retrieved", list, gin.H{
		"total":        total,
		"unread_count": unreadCount,
	})
}

// PATCH /notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	notifID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.notifService.MarkAsRead(c.Request.Context(), notifID, userID); err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	dto.Success(c, http.StatusOK, "Notification marked as read", nil)
}

// POST /notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	if err := h.notifService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to mark all notifications as read")
		return
	}

	dto.Success(c, http.StatusOK, "All notifications marked as read", nil)
}
