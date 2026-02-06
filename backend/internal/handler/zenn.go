package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
)

type ZennHandler struct {
	zennRepo    *repository.ZennRepository
	userRepo    *repository.UserRepository
	zennService *service.ZennService
}

func NewZennHandler(zennRepo *repository.ZennRepository, userRepo *repository.UserRepository, zennService *service.ZennService) *ZennHandler {
	return &ZennHandler{
		zennRepo:    zennRepo,
		userRepo:    userRepo,
		zennService: zennService,
	}
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

	// Validate username exists on Zenn
	valid, err := h.zennService.ValidateUsername(req.Username)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Zenn username"})
		return
	}

	// Update user's Zenn username
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.ZennUsername = req.Username
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Sync articles
	articles, err := h.zennService.FetchArticles(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Zenn articles"})
		return
	}

	// Set updated time
	now := time.Now()
	for i := range articles {
		articles[i].UpdatedAt = now
	}

	if err := h.zennRepo.UpsertArticles(userID, articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Zenn connected successfully",
		"articles_count": len(articles),
	})
}

// Disconnect removes the Zenn username and deletes cached articles
func (h *ZennHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.ZennUsername = ""
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Delete cached articles
	if err := h.zennRepo.DeleteUserArticles(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Zenn disconnected successfully"})
}

// Sync refreshes the Zenn articles for the current user
func (h *ZennHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.ZennUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zenn not connected"})
		return
	}

	articles, err := h.zennService.FetchArticles(user.ZennUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Zenn articles"})
		return
	}

	now := time.Now()
	for i := range articles {
		articles[i].UpdatedAt = now
	}

	if err := h.zennRepo.UpsertArticles(userID, articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Zenn synced successfully",
		"articles_count": len(articles),
	})
}

// GetArticles returns all Zenn articles for a user
func (h *ZennHandler) GetArticles(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	articles, err := h.zennRepo.GetArticles(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// GetStats returns Zenn statistics for a user
func (h *ZennHandler) GetStats(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	stats, err := h.zennRepo.GetStats(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
