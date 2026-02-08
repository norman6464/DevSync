package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type LearningLogHandler struct {
	logRepo *repository.LearningLogRepository
}

func NewLearningLogHandler(logRepo *repository.LearningLogRepository) *LearningLogHandler {
	return &LearningLogHandler{logRepo: logRepo}
}

// Create creates a new learning log
func (h *LearningLogHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Category string `json:"category"`
		Duration int    `json:"duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and content are required"})
		return
	}

	log := &model.LearningLog{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		Category: model.LogCategory(req.Category),
		Duration: req.Duration,
	}

	if req.Category == "" {
		log.Category = model.LogCategoryOther
	}

	if err := h.logRepo.Create(log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create log"})
		return
	}

	c.JSON(http.StatusCreated, log)
}

// Update updates a learning log
func (h *LearningLogHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	logID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	log, err := h.logRepo.FindByID(uint(logID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	if log.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req struct {
		Title    *string `json:"title"`
		Content  *string `json:"content"`
		Category *string `json:"category"`
		Duration *int    `json:"duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Title != nil {
		log.Title = *req.Title
	}
	if req.Content != nil {
		log.Content = *req.Content
	}
	if req.Category != nil {
		log.Category = model.LogCategory(*req.Category)
	}
	if req.Duration != nil {
		log.Duration = *req.Duration
	}

	if err := h.logRepo.Update(log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update log"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// Delete deletes a learning log
func (h *LearningLogHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	logID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	log, err := h.logRepo.FindByID(uint(logID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	if log.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	if err := h.logRepo.Delete(uint(logID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "log deleted"})
}

// GetByID gets a learning log by ID
func (h *LearningLogHandler) GetByID(c *gin.Context) {
	logID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	log, err := h.logRepo.FindByID(uint(logID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// GetMyLogs gets all learning logs for the current user
func (h *LearningLogHandler) GetMyLogs(c *gin.Context) {
	userID := c.GetUint("userID")

	logs, err := h.logRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetByUserID gets all learning logs for a user
func (h *LearningLogHandler) GetByUserID(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	logs, err := h.logRepo.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetStreakInfo returns streak data for a user
func (h *LearningLogHandler) GetStreakInfo(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	info, err := h.logRepo.GetStreakInfo(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get streak info"})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetCalendarData returns daily log counts for calendar visualization
func (h *LearningLogHandler) GetCalendarData(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	entries, err := h.logRepo.GetCalendarData(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get calendar data"})
		return
	}

	c.JSON(http.StatusOK, entries)
}
