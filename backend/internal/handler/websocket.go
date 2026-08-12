package handler

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/norman6464/devsync/backend/internal/infra/ws"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// WebSocketHandler はWebSocket関連のHTTPハンドラ。
// リアルタイム通信のためのWebSocket接続確立を処理する。
type WebSocketHandler struct {
	hub            *ws.Hub
	validateToken  *usecase.ValidateAuthTokenUseCase
	allowedOrigins map[string]bool // CORS設定と連動した許可オリジン
	upgrader       websocket.Upgrader
}

// NewWebSocketHandler は新しいWebSocketHandlerインスタンスを生成する。
// allowedOriginsにはCORS設定のオリジン一覧を渡す。
func NewWebSocketHandler(hub *ws.Hub, validateToken *usecase.ValidateAuthTokenUseCase, allowedOrigins []string) *WebSocketHandler {
	originsMap := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originsMap[o] = true
	}
	h := &WebSocketHandler{
		hub:            hub,
		validateToken:  validateToken,
		allowedOrigins: originsMap,
	}
	h.upgrader = websocket.Upgrader{
		CheckOrigin: h.checkOrigin,
	}
	return h
}

// checkOrigin はWebSocket接続のOriginヘッダーを検証する。
// CORS設定で許可されたオリジンのみ接続を許可する。
func (h *WebSocketHandler) checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	// オリジンをパースしてスキーム+ホスト部分で比較
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	originBase := parsed.Scheme + "://" + parsed.Host
	return h.allowedOrigins[originBase]
}

// HandleWebSocket はWebSocket接続を確立し、クライアントをハブに登録する。
// 認証はhttpOnly Cookie内のJWTトークンで行う（URLにトークンを含めない）。
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// httpOnly Cookieからトークンを取得（URLクエリパラメータには含めない）
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		respondUnauthorized(c, "authentication required")
		return
	}

	// トークンを検証してユーザーIDを取得
	userID, err := h.validateToken.Execute(token)
	if err != nil {
		respondUnauthorized(c, "invalid token")
		return
	}

	// HTTP接続をWebSocketにアップグレード
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// クライアントを作成してハブに登録
	client := &ws.Client{
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
