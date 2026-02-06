package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
)

type QiitaHandler struct {
	qiitaRepo    *repository.QiitaRepository
	userRepo     *repository.UserRepository
	qiitaService *service.QiitaService
}

func NewQiitaHandler(qiitaRepo *repository.QiitaRepository, userRepo *repository.UserRepository, qiitaService *service.QiitaService) *QiitaHandler {
	return &QiitaHandler{
		qiitaRepo:    qiitaRepo,
		userRepo:     userRepo,
		qiitaService: qiitaService,
	}
}

// Connect sets the Qiita username and syncs articles
func (h *QiitaHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	// Validate username exists on Qiita
	valid, err := h.qiitaService.ValidateUsername(req.Username)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Qiita username"})
		return
	}

	// Update user's Qiita username
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.QiitaUsername = req.Username
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Sync articles
	articles, err := h.qiitaService.FetchArticles(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Qiita articles"})
		return
	}

	// Set updated time
	now := time.Now()
	for i := range articles {
		articles[i].UpdatedAt = now
	}

	if err := h.qiitaRepo.UpsertArticles(userID, articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Qiita connected successfully",
		"articles_count": len(articles),
	})
}

// Disconnect removes the Qiita username and deletes cached articles
func (h *QiitaHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.QiitaUsername = ""
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Delete cached articles
	if err := h.qiitaRepo.DeleteUserArticles(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Qiita disconnected successfully"})
}

// Sync refreshes the Qiita articles for the current user
func (h *QiitaHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.QiitaUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Qiita not connected"})
		return
	}

	articles, err := h.qiitaService.FetchArticles(user.QiitaUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Qiita articles"})
		return
	}

	now := time.Now()
	for i := range articles {
		articles[i].UpdatedAt = now
	}

	if err := h.qiitaRepo.UpsertArticles(userID, articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Qiita synced successfully",
		"articles_count": len(articles),
	})
}

// GetArticles returns all Qiita articles for a user
func (h *QiitaHandler) GetArticles(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	articles, err := h.qiitaRepo.GetArticles(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// GetStats returns Qiita statistics for a user
func (h *QiitaHandler) GetStats(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	stats, err := h.qiitaRepo.GetStats(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
