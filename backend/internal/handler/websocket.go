package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/norman6464/devsync/backend/internal/service"
)

// upgrader はHTTP接続をWebSocket接続にアップグレードするための設定。
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler はWebSocket関連のHTTPハンドラ。
// リアルタイム通信のためのWebSocket接続確立を処理する。
type WebSocketHandler struct {
	hub         *service.Hub
	authService *service.AuthService
}

// NewWebSocketHandler は新しいWebSocketHandlerインスタンスを生成する。
func NewWebSocketHandler(hub *service.Hub, authService *service.AuthService) *WebSocketHandler {
	return &WebSocketHandler{hub: hub, authService: authService}
}

// HandleWebSocket はWebSocket接続を確立し、クライアントをハブに登録する。
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	// トークンを検証してユーザーIDを取得
	userID, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// HTTP接続をWebSocketにアップグレード
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// クライアントを作成してハブに登録
	client := &service.Client{
		Hub:    h.hub,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.hub.Register(client)

	// 読み書きのgoroutineを起動
	go client.WritePump()
	go client.ReadPump()
}
