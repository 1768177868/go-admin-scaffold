package sse

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// EventType defines the type of SSE event
const (
	EventTypeNotification = "notification"
	EventTypeAlert        = "alert"
	EventTypeUpdate       = "update"
)

// Event represents a server-sent event
type Event struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Time    time.Time   `json:"time"`
	UserID  string      `json:"user_id,omitempty"`  // 目标用户ID，为空表示广播
	GroupID string      `json:"group_id,omitempty"` // 目标组ID，为空表示非组消息
}

// Client represents an SSE client connection
type Client struct {
	ID       string
	Groups   map[string]bool
	Messages chan *Event
}

// Manager manages SSE connections and event broadcasting
type Manager struct {
	clients    map[string]*Client
	groups     map[string]map[string]bool // group -> userIDs
	register   chan *Client
	unregister chan *Client
	events     chan *Event
	mu         sync.RWMutex
}

// NewManager creates a new SSE manager
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[string]*Client),
		groups:     make(map[string]map[string]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		events:     make(chan *Event, 100), // 缓冲通道，避免阻塞
	}
}

// Start starts the SSE manager
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.ID] = client
			log.Printf("SSE client %s registered", client.ID)
			m.mu.Unlock()

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client.ID]; ok {
				// 从所有组中移除客户端
				for groupID := range client.Groups {
					if group, exists := m.groups[groupID]; exists {
						delete(group, client.ID)
						if len(group) == 0 {
							delete(m.groups, groupID)
						}
					}
				}
				delete(m.clients, client.ID)
				close(client.Messages)
				log.Printf("SSE client %s unregistered", client.ID)
			}
			m.mu.Unlock()

		case event := <-m.events:
			m.handleEvent(event)
		}
	}
}

// handleEvent processes and distributes events to relevant clients
func (m *Manager) handleEvent(event *Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 如果指定了用户ID，只发送给该用户
	if event.UserID != "" {
		if client, ok := m.clients[event.UserID]; ok {
			select {
			case client.Messages <- event:
				log.Printf("Event sent to user %s", event.UserID)
			default:
				log.Printf("Failed to send event to user %s: channel full", event.UserID)
			}
		}
		return
	}

	// 如果指定了组ID，发送给组内所有成员
	if event.GroupID != "" {
		if group, ok := m.groups[event.GroupID]; ok {
			for userID := range group {
				if client, ok := m.clients[userID]; ok {
					select {
					case client.Messages <- event:
						log.Printf("Event sent to group member %s", userID)
					default:
						log.Printf("Failed to send event to group member %s: channel full", userID)
					}
				}
			}
		}
		return
	}

	// 如果既没有指定用户也没有指定组，广播给所有客户端
	for _, client := range m.clients {
		select {
		case client.Messages <- event:
			log.Printf("Event broadcasted to client %s", client.ID)
		default:
			log.Printf("Failed to broadcast event to client %s: channel full", client.ID)
		}
	}
}

// Register registers a new client
func (m *Manager) Register(userID string) *Client {
	client := &Client{
		ID:       userID,
		Groups:   make(map[string]bool),
		Messages: make(chan *Event, 100),
	}
	m.register <- client
	return client
}

// Unregister removes a client
func (m *Manager) Unregister(client *Client) {
	m.unregister <- client
}

// JoinGroup adds a client to a group
func (m *Manager) JoinGroup(groupID, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.groups[groupID]; !ok {
		m.groups[groupID] = make(map[string]bool)
	}
	m.groups[groupID][userID] = true

	if client, ok := m.clients[userID]; ok {
		client.Groups[groupID] = true
	}

	// 发送加入通知
	m.events <- &Event{
		ID:      fmt.Sprintf("join_%s_%d", groupID, time.Now().UnixNano()),
		Type:    EventTypeNotification,
		Data:    fmt.Sprintf("用户 %s 加入了组 %s", userID, groupID),
		Time:    time.Now(),
		GroupID: groupID,
	}
}

// LeaveGroup removes a client from a group
func (m *Manager) LeaveGroup(groupID, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if group, ok := m.groups[groupID]; ok {
		delete(group, userID)
		if len(group) == 0 {
			delete(m.groups, groupID)
		}
	}

	if client, ok := m.clients[userID]; ok {
		delete(client.Groups, groupID)
	}

	// 发送离开通知
	m.events <- &Event{
		ID:      fmt.Sprintf("leave_%s_%d", groupID, time.Now().UnixNano()),
		Type:    EventTypeNotification,
		Data:    fmt.Sprintf("用户 %s 离开了组 %s", userID, groupID),
		Time:    time.Now(),
		GroupID: groupID,
	}
}

// SendEvent sends an event to the specified target(s)
func (m *Manager) SendEvent(event *Event) {
	if event.ID == "" {
		event.ID = fmt.Sprintf("evt_%d", time.Now().UnixNano())
	}
	if event.Time.IsZero() {
		event.Time = time.Now()
	}
	m.events <- event
}

// Broadcast sends an event to all connected clients
func (m *Manager) Broadcast(eventType string, data interface{}) {
	event := &Event{
		ID:   fmt.Sprintf("broadcast_%d", time.Now().UnixNano()),
		Type: eventType,
		Data: data,
		Time: time.Now(),
	}
	m.events <- event
}

// SendToUser sends an event to a specific user
func (m *Manager) SendToUser(userID, eventType string, data interface{}) {
	event := &Event{
		ID:     fmt.Sprintf("user_%s_%d", userID, time.Now().UnixNano()),
		Type:   eventType,
		Data:   data,
		Time:   time.Now(),
		UserID: userID,
	}
	m.events <- event
}

// SendToGroup sends an event to all members of a group
func (m *Manager) SendToGroup(groupID, eventType string, data interface{}) {
	event := &Event{
		ID:      fmt.Sprintf("group_%s_%d", groupID, time.Now().UnixNano()),
		Type:    eventType,
		Data:    data,
		Time:    time.Now(),
		GroupID: groupID,
	}
	m.events <- event
}
