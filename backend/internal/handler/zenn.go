package handler

import (
	"errors"
	"net/http"
	"strconv"

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
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Zenn username"})
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect Zenn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Zenn connected successfully",
		"articles_count": count,
	})
}

// Disconnect はZennユーザー名を削除し、キャッシュされた記事を削除する。
func (h *ZennHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.service.Disconnect(userID); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to disconnect Zenn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Zenn disconnected successfully"})
}

// Sync は現在のユーザーのZenn記事を再同期する。
func (h *ZennHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.Sync(userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Zenn not connected"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync Zenn articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Zenn synced successfully",
		"articles_count": count,
	})
}

// GetArticles は指定ユーザーの全Zenn記事を返す。
func (h *ZennHandler) GetArticles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	articles, err := h.service.GetArticles(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// GetStats は指定ユーザーのZenn統計情報を返す。
func (h *ZennHandler) GetStats(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	stats, err := h.service.GetStats(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
