package notification

import (
	"time"

	"github.com/google/uuid"
)

const (
	TypeSpotComment     = "spot_comment"
	TypeJoinRequest     = "join_request"
	TypeRequestApproved = "request_approved"
	TypeRequestRejected = "request_rejected"
	TypeEventUpdated    = "event_updated"
	TypeEventCancelled  = "event_cancelled"
	TypeNewMessage      = "new_message"
)

// Notification represents an in-app alert for a user.
type Notification struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Message     string     `json:"message"`
	ReferenceID *uuid.UUID `json:"reference_id,omitempty"`
	IsRead      bool       `json:"is_read"`
	CreatedAt   time.Time  `json:"created_at"`
}
