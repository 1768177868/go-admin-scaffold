package handlers

import (
	"app/internal/core/ws"
	"app/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type WSHandler struct {
	manager *ws.Manager
}

func NewWSHandler() *WSHandler {
	manager := ws.NewManager()
	go manager.Start() // Start the WebSocket manager
	return &WSHandler{manager: manager}
}

// HandleWebSocket handles WebSocket connections
// @Summary Connect to WebSocket
// @Description Establishes a WebSocket connection for real-time chat
// @Tags WebSocket
// @Accept  json
// @Produce  json
// @Param user_id query string true "User ID"
// @Success 101 {string} string "Switching Protocols to websocket"
// @Router /ws [get]
func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		response.ParamError(c, "user_id is required")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.ServerError(c)
		return
	}

	client := &ws.Client{
		ID:      userID,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		Manager: h.manager,
		Groups:  make(map[string]bool),
	}

	h.manager.Register <- client

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// JoinGroup handles joining a chat group
// @Summary Join Chat Group
// @Description Adds a user to a chat group
// @Tags WebSocket
// @Accept  json
// @Produce  json
// @Param user_id query string true "User ID"
// @Param group_id query string true "Group ID"
// @Success 200 {object} map[string]string
// @Router /ws/join [post]
func (h *WSHandler) JoinGroup(c *gin.Context) {
	userID := c.Query("user_id")
	groupID := c.Query("group_id")

	if userID == "" || groupID == "" {
		response.ParamError(c, "user_id and group_id are required")
		return
	}

	h.manager.JoinGroup(groupID, userID)
	response.Success(c, gin.H{"message": "Successfully joined group"})
}

// SendMessage handles sending messages
// @Summary Send Message
// @Description Sends a message (private, group, or announcement)
// @Tags WebSocket
// @Accept  json
// @Produce  json
// @Param message body ws.Message true "Message Object"
// @Success 200 {object} map[string]string
// @Router /ws/send [post]
func (h *WSHandler) SendMessage(c *gin.Context) {
	var message ws.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	message.Timestamp = time.Now().Unix()
	h.manager.Broadcast <- &message

	response.Success(c, gin.H{"message": "Message sent successfully"})
}
