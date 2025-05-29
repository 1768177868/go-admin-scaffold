package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType defines the type of WebSocket message
type MessageType int

const (
	MessageTypePrivate MessageType = iota + 1
	MessageTypeGroup
	MessageTypeAnnouncement
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to"` // User ID or Group ID
	Content   string      `json:"content"`
	Timestamp int64       `json:"timestamp"`
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
	if client, ok := m.Clients[message.To]; ok {
		data, err := json.Marshal(message)
		if err == nil {
			client.Send <- data
		}
	}
	m.mu.RUnlock()
}

func (m *Manager) handleGroupMessage(message *Message) {
	m.mu.RLock()
	if group, ok := m.Groups[message.To]; ok {
		data, err := json.Marshal(message)
		if err == nil {
			for userID := range group {
				if client, ok := m.Clients[userID]; ok {
					client.Send <- data
				}
			}
		}
	}
	m.mu.RUnlock()
}

func (m *Manager) handleAnnouncement(message *Message) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	m.mu.RLock()
	for _, client := range m.Clients {
		client.Send <- data
	}
	m.mu.RUnlock()
}

// JoinGroup adds a client to a group
func (m *Manager) JoinGroup(groupID, clientID string) {
	m.mu.Lock()
	if _, ok := m.Groups[groupID]; !ok {
		m.Groups[groupID] = make(map[string]bool)
	}
	m.Groups[groupID][clientID] = true
	m.mu.Unlock()
}

// LeaveGroup removes a client from a group
func (m *Manager) LeaveGroup(groupID, clientID string) {
	m.mu.Lock()
	if group, ok := m.Groups[groupID]; ok {
		delete(group, clientID)
		if len(group) == 0 {
			delete(m.Groups, groupID)
		}
	}
	m.mu.Unlock()
}

// WritePump handles writing messages to the WebSocket connection
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// ReadPump handles reading messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Log error here
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		msg.From = c.ID
		msg.Timestamp = time.Now().Unix()
		c.Manager.Broadcast <- &msg
	}
}
