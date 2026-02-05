package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/norman6464/devsync/backend/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub         *service.Hub
	authService *service.AuthService
}

func NewWebSocketHandler(hub *service.Hub, authService *service.AuthService) *WebSocketHandler {
	return &WebSocketHandler{hub: hub, authService: authService}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	userID, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &service.Client{
		Hub:    h.hub,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
