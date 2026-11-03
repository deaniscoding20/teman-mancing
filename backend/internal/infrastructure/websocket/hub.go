package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"

	appchat "teman-mancing/backend/internal/application/chat"
	domainchat "teman-mancing/backend/internal/domain/chat"
)

type Client struct {
	UserID   uuid.UUID
	Username string
	Conn     *websocket.Conn
	SendChan chan []byte
}

type Hub struct {
	mu          sync.RWMutex
	clients     map[uuid.UUID][]*Client // userID -> active sockets
	chatService *appchat.Service
	logger      *zap.Logger
}

func NewHub(chatService *appchat.Service, logger *zap.Logger) *Hub {
	return &Hub{
		clients:     make(map[uuid.UUID][]*Client),
		chatService: chatService,
		logger:      logger,
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.UserID] = append(h.clients[c.UserID], c)
	h.logger.Info("WebSocket client registered", zap.String("user_id", c.UserID.String()))
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	list := h.clients[c.UserID]
	var updated []*Client
	for _, item := range list {
		if item != c {
			updated = append(updated, item)
		}
	}
	if len(updated) == 0 {
		delete(h.clients, c.UserID)
	} else {
		h.clients[c.UserID] = updated
	}
	close(c.SendChan)
	h.logger.Info("WebSocket client unregistered", zap.String("user_id", c.UserID.String()))
}

type IncomingWSMessage struct {
	Action         string    `json:"action"` // "send_message", "ping"
	ConversationID uuid.UUID `json:"conversation_id"`
	Content        string    `json:"content"`
}

type OutgoingWSMessage struct {
	Type           string    `json:"type"` // "new_message", "pong", "error"
	ConversationID uuid.UUID `json:"conversation_id,omitempty"`
	MessageID      uuid.UUID `json:"message_id,omitempty"`
	SenderID       uuid.UUID `json:"sender_id,omitempty"`
	SenderName     string    `json:"sender_name,omitempty"`
	Content        string    `json:"content,omitempty"`
	CreatedAt      string    `json:"created_at,omitempty"`
}

// HandleWS upgrades HTTP connection and manages message pump
func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request, userID uuid.UUID, username string) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		h.logger.Error("Failed to accept websocket connection", zap.Error(err))
		return
	}
	defer conn.CloseNow()

	client := &Client{
		UserID:   userID,
		Username: username,
		Conn:     conn,
		SendChan: make(chan []byte, 32),
	}
	h.Register(client)
	defer h.Unregister(client)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Writer goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msgBytes, ok := <-client.SendChan:
				if !ok {
					return
				}
				writeCtx, writeCancel := context.WithTimeout(ctx, 5*time.Second)
				err := conn.Write(writeCtx, websocket.MessageText, msgBytes)
				writeCancel()
				if err != nil {
					return
				}
			}
		}
	}()

	// Reader loop
	for {
		_, msgBytes, err := conn.Read(ctx)
		if err != nil {
			break
		}

		var incoming IncomingWSMessage
		if err := json.Unmarshal(msgBytes, &incoming); err != nil {
			continue
		}

		if incoming.Action == "ping" {
			res, _ := json.Marshal(OutgoingWSMessage{Type: "pong"})
			client.SendChan <- res
			continue
		}

		if incoming.Action == "send_message" {
			// Save message to persistent storage & validate membership
			saved, err := h.chatService.SendMessage(ctx, incoming.ConversationID, userID, username, incoming.Content)
			if err != nil {
				errBytes, _ := json.Marshal(OutgoingWSMessage{
					Type:    "error",
					Content: err.Error(),
				})
				client.SendChan <- errBytes
				continue
			}

			// Broadcast to all connected clients
			h.BroadcastMessage(saved)
		}
	}
}

// BroadcastMessage sends a chat message to all connected clients.
func (h *Hub) BroadcastMessage(msg *domainchat.Message) {
	if h == nil || msg == nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	outBytes, err := json.Marshal(OutgoingWSMessage{
		Type:           "new_message",
		ConversationID: msg.ConversationID,
		MessageID:      msg.ID,
		SenderID:       msg.SenderID,
		SenderName:     msg.SenderName,
		Content:        msg.Message,
		CreatedAt:      msg.CreatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return
	}

	for _, clientList := range h.clients {
		for _, client := range clientList {
			select {
			case client.SendChan <- outBytes:
			default:
			}
		}
	}
}
