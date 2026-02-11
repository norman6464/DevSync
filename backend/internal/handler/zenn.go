package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// ZennHandler はZenn連携関連のHTTPハンドラ。
// Zennアカウントの接続・切断・記事同期・統計情報の取得を処理する。
type ZennHandler struct {
	service *service.ZennService
}

// NewZennHandler は新しいZennHandlerインスタンスを生成する。
func NewZennHandler(s *service.ZennService) *ZennHandler {
	return &ZennHandler{service: s}
}

// Connect はZennユーザー名を設定し、記事を同期する。
func (h *ZennHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	count, err := h.service.Connect(userID, req.Username)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{
		"message":        "Zenn connected successfully",
		"articles_count": count,
	})
}

// Disconnect はZennユーザー名を削除し、キャッシュされた記事を削除する。
func (h *ZennHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.service.Disconnect(userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "Zenn disconnected successfully"})
}

// Sync は現在のユーザーのZenn記事を再同期する。
func (h *ZennHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.Sync(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{
		"message":        "Zenn synced successfully",
		"articles_count": count,
	})
}

// GetArticles は指定ユーザーの全Zenn記事を返す。
func (h *ZennHandler) GetArticles(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	articles, err := h.service.GetArticles(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, articles)
}

// GetStats は指定ユーザーのZenn統計情報を返す。
func (h *ZennHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	stats, err := h.service.GetStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
