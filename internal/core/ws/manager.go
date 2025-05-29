package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType defines the type of WebSocket message
const (
	MessageTypePrivate      = 1
	MessageTypeGroup        = 2
	MessageTypeAnnouncement = 3
)

// Message represents a WebSocket message
type Message struct {
	Type      int    `json:"type"`
	From      string `json:"from"`
	To        string `json:"to"` // User ID or Group ID
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// Client represents a WebSocket client
type Client struct {
	ID      string
	Conn    *websocket.Conn
	Send    chan []byte
	Manager *Manager
	Groups  map[string]bool
	mu      sync.Mutex
}

// Manager manages WebSocket connections and message broadcasting
type Manager struct {
	Clients    map[string]*Client
	Groups     map[string]map[string]bool // group -> userIDs
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	mu         sync.RWMutex
}

// NewManager creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		Clients:    make(map[string]*Client),
		Groups:     make(map[string]map[string]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

// Start starts the WebSocket manager
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.mu.Lock()
			m.Clients[client.ID] = client
			m.mu.Unlock()

		case client := <-m.Unregister:
			if _, ok := m.Clients[client.ID]; ok {
				m.mu.Lock()
				delete(m.Clients, client.ID)
				close(client.Send)
				m.mu.Unlock()
			}

		case message := <-m.Broadcast:
			switch message.Type {
			case MessageTypePrivate:
				m.handlePrivateMessage(message)
			case MessageTypeGroup:
				m.handleGroupMessage(message)
			case MessageTypeAnnouncement:
				m.handleAnnouncement(message)
			}
		}
	}
}

func (m *Manager) handlePrivateMessage(message *Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Send to both sender and receiver
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	// Send to receiver
	if client, ok := m.Clients[message.To]; ok {
		client.Send <- data
	}

	// Send to sender if not the same as receiver
	if message.From != message.To {
		if client, ok := m.Clients[message.From]; ok {
			client.Send <- data
		}
	}
}

func (m *Manager) handleGroupMessage(message *Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if group, ok := m.Groups[message.To]; ok {
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling group message: %v", err)
			return
		}

		log.Printf("Broadcasting group message to group %s with %d members", message.To, len(group))

		// Send to all members of the group, including the sender
		for userID := range group {
			if client, ok := m.Clients[userID]; ok {
				log.Printf("Sending message to group member: %s", userID)
				select {
				case client.Send <- data:
					log.Printf("Message sent to member %s successfully", userID)
				default:
					log.Printf("Failed to send message to member %s: send channel full", userID)
				}
			} else {
				log.Printf("Member %s not found in clients map", userID)
			}
		}
	} else {
		log.Printf("Group %s not found", message.To)
	}
}

func (m *Manager) handleAnnouncement(message *Message) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, client := range m.Clients {
		client.Send <- data
	}
}

// JoinGroup adds a client to a group
func (m *Manager) JoinGroup(groupID, clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Groups[groupID]; !ok {
		m.Groups[groupID] = make(map[string]bool)
	}
	m.Groups[groupID][clientID] = true

	log.Printf("Client %s joined group %s", clientID, groupID)

	// Send a confirmation message to the group
	confirmMsg := &Message{
		Type:      MessageTypeAnnouncement,
		From:      clientID,
		Content:   fmt.Sprintf("用户 %s 加入了群组", clientID),
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(confirmMsg)
	if err != nil {
		log.Printf("Error marshaling join confirmation: %v", err)
		return
	}

	// Send to all members of the group
	for memberID := range m.Groups[groupID] {
		if client, ok := m.Clients[memberID]; ok {
			select {
			case client.Send <- data:
				log.Printf("Join confirmation sent to %s", memberID)
			default:
				log.Printf("Failed to send join confirmation to %s", memberID)
			}
		}
	}
}

// LeaveGroup removes a client from a group
func (m *Manager) LeaveGroup(groupID, clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if group, ok := m.Groups[groupID]; ok {
		delete(group, clientID)
		if len(group) == 0 {
			delete(m.Groups, groupID)
		}
	}
}

// WritePump handles writing messages to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(time.Second * 30)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				log.Printf("Send channel closed for client %s", c.ID)
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error getting writer for client %s: %v", c.ID, err)
				return
			}

			if _, err := w.Write(message); err != nil {
				log.Printf("Error writing message for client %s: %v", c.ID, err)
				return
			}

			if err := w.Close(); err != nil {
				log.Printf("Error closing writer for client %s: %v", c.ID, err)
				return
			}

			log.Printf("Message successfully written to client %s", c.ID)

		case <-ticker.C:
			// Send ping to keep connection alive
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping to client %s: %v", c.ID, err)
				return
			}
		}
	}
}

// ReadPump handles reading messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		log.Printf("Client %s disconnected", c.ID)
		c.Manager.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message from client %s: %v", c.ID, err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message from client %s: %v", c.ID, err)
			continue
		}

		// Set sender and timestamp only if not set
		if msg.From == "" {
			msg.From = c.ID
		}
		if msg.Timestamp == 0 {
			msg.Timestamp = time.Now().Unix()
		}

		// Add debug logging
		data, _ := json.Marshal(msg)
		log.Printf("Received message from client %s: %s", c.ID, string(data))

		// Validate message type
		if msg.Type < MessageTypePrivate || msg.Type > MessageTypeAnnouncement {
			log.Printf("Invalid message type from client %s: %d", c.ID, msg.Type)
			continue
		}

		// Broadcast the message
		c.Manager.Broadcast <- &msg
		log.Printf("Message from client %s broadcasted", c.ID)
	}
}
