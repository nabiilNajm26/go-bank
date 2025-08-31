package http

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type NotificationMessage struct {
	Type    string      `json:"type"`
	UserID  uuid.UUID   `json:"user_id"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type WebSocketHandler struct {
	connections map[uuid.UUID]*websocket.Conn
	mutex       sync.RWMutex
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		connections: make(map[uuid.UUID]*websocket.Conn),
	}
}

func (h *WebSocketHandler) HandleConnection(c *websocket.Conn) {
	userID := c.Locals("userID").(uuid.UUID)
	
	h.mutex.Lock()
	h.connections[userID] = c
	h.mutex.Unlock()

	defer func() {
		h.mutex.Lock()
		delete(h.connections, userID)
		h.mutex.Unlock()
		c.Close()
	}()

	// Send welcome message
	welcomeMsg := NotificationMessage{
		Type:    "welcome",
		UserID:  userID,
		Message: "Connected to real-time notifications",
	}
	h.sendMessage(c, welcomeMsg)

	// Listen for messages
	for {
		var msg map[string]interface{}
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Echo back or handle specific message types
		response := NotificationMessage{
			Type:    "echo",
			UserID:  userID,
			Message: "Message received",
			Data:    msg,
		}
		h.sendMessage(c, response)
	}
}

func (h *WebSocketHandler) NotifyUser(userID uuid.UUID, message NotificationMessage) {
	h.mutex.RLock()
	conn, exists := h.connections[userID]
	h.mutex.RUnlock()

	if exists {
		h.sendMessage(conn, message)
	}
}

func (h *WebSocketHandler) NotifyTransfer(fromUserID, toUserID uuid.UUID, transaction interface{}) {
	// Notify sender
	senderMsg := NotificationMessage{
		Type:    "transfer_sent",
		UserID:  fromUserID,
		Message: "Transfer sent successfully",
		Data:    transaction,
	}
	h.NotifyUser(fromUserID, senderMsg)

	// Notify receiver
	receiverMsg := NotificationMessage{
		Type:    "transfer_received",
		UserID:  toUserID,
		Message: "You received a transfer",
		Data:    transaction,
	}
	h.NotifyUser(toUserID, receiverMsg)
}

func (h *WebSocketHandler) sendMessage(conn *websocket.Conn, message NotificationMessage) {
	err := conn.WriteJSON(message)
	if err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

func (h *WebSocketHandler) GetConnectedUsers() []uuid.UUID {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	users := make([]uuid.UUID, 0, len(h.connections))
	for userID := range h.connections {
		users = append(users, userID)
	}
	return users
}