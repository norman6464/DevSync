package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

type ZennHandler struct {
	service *service.ZennService
}

func NewZennHandler(s *service.ZennService) *ZennHandler {
	return &ZennHandler{service: s}
}

// Connect sets the Zenn username and syncs articles
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

// Disconnect removes the Zenn username and deletes cached articles
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

// Sync refreshes the Zenn articles for the current user
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

// GetArticles returns all Zenn articles for a user
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

// GetStats returns Zenn statistics for a user
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
