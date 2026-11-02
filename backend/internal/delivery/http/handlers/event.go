package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appevent "teman-mancing/backend/internal/application/event"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/middleware"
	domainevent "teman-mancing/backend/internal/domain/event"
)

type EventHandler struct {
	eventService *appevent.Service
}

func NewEventHandler(eventService *appevent.Service) *EventHandler {
	return &EventHandler{eventService: eventService}
}

type CreateEventRequest struct {
	Name                 string     `json:"name" binding:"required"`
	Description          string     `json:"description" binding:"required"`
	FishingSpotID        *uuid.UUID `json:"fishing_spot_id"`
	DestinationName      string     `json:"destination_name" binding:"required"`
	DestinationAddress   string     `json:"destination_address" binding:"required"`
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
	EventDate            string     `json:"event_date" binding:"required"` // YYYY-MM-DD
	StartTime            string     `json:"start_time" binding:"required"` // HH:MM:SS
	EndTime              *string    `json:"end_time"`
	MeetingPoint         string     `json:"meeting_point" binding:"required"`
	Difficulty           string     `json:"difficulty"`
	HasParkingFee        bool       `json:"has_parking_fee"`
	ParkingFee           *float64   `json:"parking_fee"`
	RequiresBoat         bool       `json:"requires_boat"`
	BoatRentalFee        *float64   `json:"boat_rental_fee"`
	AdditionalCost       *string    `json:"additional_cost"`
	Visibility           string     `json:"visibility"`
	MaxParticipants      *int       `json:"max_participants"`
	RegistrationDeadline *time.Time `json:"registration_deadline"`
}

type UpdateEventRequest struct {
	Name               *string  `json:"name"`
	Description        *string  `json:"description"`
	DestinationName    *string  `json:"destination_name"`
	DestinationAddress *string  `json:"destination_address"`
	Latitude           *float64 `json:"latitude"`
	Longitude          *float64 `json:"longitude"`
	MeetingPoint       *string  `json:"meeting_point"`
	Difficulty         *string  `json:"difficulty"`
	HasParkingFee      *bool    `json:"has_parking_fee"`
	ParkingFee         *float64 `json:"parking_fee"`
	RequiresBoat       *bool    `json:"requires_boat"`
	BoatRentalFee      *float64 `json:"boat_rental_fee"`
	AdditionalCost     *string  `json:"additional_cost"`
	MaxParticipants    *int     `json:"max_participants"`
	Status             *string  `json:"status"`
}

// GET /events
func (h *EventHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	query := c.Query("q")
	date := c.Query("date")
	difficulty := c.Query("difficulty")
	visibility := c.Query("visibility")
	status := c.Query("status")

	events, total, err := h.eventService.ListEvents(c.Request.Context(), appevent.ListEventsFilter{
		Page:       page,
		Limit:      limit,
		Query:      query,
		Date:       date,
		Difficulty: difficulty,
		Visibility: visibility,
		Status:     status,
	})
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	dto.SuccessWithMeta(c, http.StatusOK, "Events retrieved", events, gin.H{
		"page":        page,
		"limit":       limit,
		"total_items": total,
		"total_pages": totalPages,
	})
}

// GET /events/:id
func (h *EventHandler) GetByID(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	ev, err := h.eventService.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		if errors.Is(err, appevent.ErrEventNotFound) {
			dto.Error(c, http.StatusNotFound, "Event not found")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to retrieve event")
		return
	}

	participants, _ := h.eventService.ListParticipants(c.Request.Context(), eventID, "")

	dto.Success(c, http.StatusOK, "Event details retrieved", gin.H{
		"event":        ev,
		"participants": participants,
	})
}

// POST /events
func (h *EventHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event payload", dto.FieldError{Field: "validation", Message: err.Error()})
		return
	}

	ev, err := h.eventService.CreateEvent(c.Request.Context(), appevent.CreateEventInput{
		CreatorID:            userID,
		Name:                 req.Name,
		Description:          req.Description,
		FishingSpotID:        req.FishingSpotID,
		DestinationName:      req.DestinationName,
		DestinationAddress:   req.DestinationAddress,
		Latitude:             req.Latitude,
		Longitude:            req.Longitude,
		EventDate:            req.EventDate,
		StartTime:            req.StartTime,
		EndTime:              req.EndTime,
		MeetingPoint:         req.MeetingPoint,
		Difficulty:           req.Difficulty,
		HasParkingFee:        req.HasParkingFee,
		ParkingFee:           req.ParkingFee,
		RequiresBoat:         req.RequiresBoat,
		BoatRentalFee:        req.BoatRentalFee,
		AdditionalCost:       req.AdditionalCost,
		Visibility:           req.Visibility,
		MaxParticipants:      req.MaxParticipants,
		RegistrationDeadline: req.RegistrationDeadline,
	})
	if err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusCreated, "Event created successfully", ev)
}

// PATCH /events/:id
func (h *EventHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid update payload")
		return
	}

	ev, err := h.eventService.UpdateEvent(c.Request.Context(), eventID, userID, appevent.UpdateEventInput{
		Name:               req.Name,
		Description:        req.Description,
		DestinationName:    req.DestinationName,
		DestinationAddress: req.DestinationAddress,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		MeetingPoint:       req.MeetingPoint,
		Difficulty:         req.Difficulty,
		HasParkingFee:      req.HasParkingFee,
		ParkingFee:         req.ParkingFee,
		RequiresBoat:       req.RequiresBoat,
		BoatRentalFee:      req.BoatRentalFee,
		AdditionalCost:     req.AdditionalCost,
		MaxParticipants:    req.MaxParticipants,
		Status:             req.Status,
	})
	if err != nil {
		if errors.Is(err, appevent.ErrEventNotFound) {
			dto.Error(c, http.StatusNotFound, "Event not found")
			return
		}
		if errors.Is(err, appevent.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "You do not have permission to edit this event")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to update event")
		return
	}

	dto.Success(c, http.StatusOK, "Event updated successfully", ev)
}

// DELETE /events/:id
func (h *EventHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	if err := h.eventService.CancelEvent(c.Request.Context(), eventID, userID); err != nil {
		if errors.Is(err, appevent.ErrEventNotFound) {
			dto.Error(c, http.StatusNotFound, "Event not found")
			return
		}
		if errors.Is(err, appevent.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "You do not have permission to cancel this event")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to cancel event")
		return
	}

	dto.Success(c, http.StatusOK, "Event cancelled successfully", nil)
}

// POST /events/:id/join
func (h *EventHandler) Join(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	participant, err := h.eventService.JoinEvent(c.Request.Context(), eventID, userID)
	if err != nil {
		if errors.Is(err, domainevent.ErrAlreadyParticipant) {
			dto.Error(c, http.StatusConflict, "You have already joined or requested to join this event")
			return
		}
		if errors.Is(err, domainevent.ErrEventCapacityFull) {
			dto.Error(c, http.StatusUnprocessableEntity, "Event participant capacity has been reached")
			return
		}
		if errors.Is(err, domainevent.ErrEventCancelled) || errors.Is(err, domainevent.ErrEventClosed) {
			dto.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	msg := "Successfully joined event"
	if participant.Status == domainevent.ParticipantStatusPending {
		msg = "Join request sent to event creator. Waiting for approval."
	}

	dto.Success(c, http.StatusOK, msg, participant)
}

// POST /events/:id/leave
func (h *EventHandler) Leave(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	if err := h.eventService.LeaveEvent(c.Request.Context(), eventID, userID); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusOK, "You have left the event", nil)
}

// GET /events/:id/requests
func (h *EventHandler) ListRequests(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	ev, err := h.eventService.GetEvent(c.Request.Context(), eventID)
	if err != nil || ev.CreatorID != userID {
		dto.Error(c, http.StatusForbidden, "You do not own this event")
		return
	}

	requests, err := h.eventService.ListParticipants(c.Request.Context(), eventID, domainevent.ParticipantStatusPending)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to list requests")
		return
	}

	dto.Success(c, http.StatusOK, "Join requests retrieved", requests)
}

// POST /events/:id/requests/:userId/approve
func (h *EventHandler) ApproveRequest(c *gin.Context) {
	actorID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	if err := h.eventService.ApproveJoinRequest(c.Request.Context(), eventID, targetUserID, actorID); err != nil {
		if errors.Is(err, domainevent.ErrEventCapacityFull) {
			dto.Error(c, http.StatusUnprocessableEntity, "Event capacity is full")
			return
		}
		if errors.Is(err, appevent.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "Only creator can approve join requests")
			return
		}
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusOK, "Join request approved successfully", nil)
}

// POST /events/:id/requests/:userId/reject
func (h *EventHandler) RejectRequest(c *gin.Context) {
	actorID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid event ID")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	if err := h.eventService.RejectJoinRequest(c.Request.Context(), eventID, targetUserID, actorID); err != nil {
		if errors.Is(err, appevent.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "Only creator can reject join requests")
			return
		}
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusOK, "Join request rejected", nil)
}
