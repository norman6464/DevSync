package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

type LearningLogHandler struct {
	service *service.LearningLogService
}

func NewLearningLogHandler(s *service.LearningLogService) *LearningLogHandler {
	return &LearningLogHandler{service: s}
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

	if err := h.service.Create(log); err != nil {
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

	updates := &model.LearningLog{}
	if req.Title != nil {
		updates.Title = *req.Title
	}
	if req.Content != nil {
		updates.Content = *req.Content
	}
	if req.Category != nil {
		updates.Category = model.LogCategory(*req.Category)
	}
	if req.Duration != nil {
		updates.Duration = *req.Duration
	}

	log, err := h.service.Update(uint(logID), userID, updates)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
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

	if err := h.service.Delete(uint(logID), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
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

	log, err := h.service.GetByID(uint(logID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// GetMyLogs gets all learning logs for the current user
func (h *LearningLogHandler) GetMyLogs(c *gin.Context) {
	userID := c.GetUint("userID")

	logs, err := h.service.GetByUserID(userID)
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

	logs, err := h.service.GetByUserID(uint(userID))
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

	info, err := h.service.GetStreakInfo(uint(userID))
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

	entries, err := h.service.GetCalendarData(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get calendar data"})
		return
	}

	c.JSON(http.StatusOK, entries)
}
