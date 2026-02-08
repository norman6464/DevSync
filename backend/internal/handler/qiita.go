package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// QiitaHandler はQiita連携関連のHTTPハンドラ。
// Qiitaアカウントの接続・切断・記事同期・統計情報の取得を処理する。
type QiitaHandler struct {
	service *service.QiitaService
}

// NewQiitaHandler は新しいQiitaHandlerインスタンスを生成する。
func NewQiitaHandler(s *service.QiitaService) *QiitaHandler {
	return &QiitaHandler{service: s}
}

// Connect はQiitaユーザー名を設定し、記事を同期する。
func (h *QiitaHandler) Connect(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Qiita username"})
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect Qiita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Qiita connected successfully",
		"articles_count": count,
	})
}

// Disconnect はQiitaユーザー名を削除し、キャッシュされた記事を削除する。
func (h *QiitaHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.service.Disconnect(userID); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to disconnect Qiita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Qiita disconnected successfully"})
}

// Sync は現在のユーザーのQiita記事を再同期する。
func (h *QiitaHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.Sync(userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Qiita not connected"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync Qiita articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Qiita synced successfully",
		"articles_count": count,
	})
}

// GetArticles は指定ユーザーの全Qiita記事を返す。
func (h *QiitaHandler) GetArticles(c *gin.Context) {
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

// GetStats は指定ユーザーのQiita統計情報を返す。
func (h *QiitaHandler) GetStats(c *gin.Context) {
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
