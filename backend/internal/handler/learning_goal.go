package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type LearningGoalHandler struct {
	goalRepo *repository.LearningGoalRepository
}

func NewLearningGoalHandler(goalRepo *repository.LearningGoalRepository) *LearningGoalHandler {
	return &LearningGoalHandler{goalRepo: goalRepo}
}

// Create creates a new learning goal
func (h *LearningGoalHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category"`
		TargetDate  string `json:"target_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	goal := &model.LearningGoal{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    model.GoalCategory(req.Category),
		Status:      model.GoalStatusActive,
		Progress:    0,
	}

	if req.Category == "" {
		goal.Category = model.GoalCategoryOther
	}

	if req.TargetDate != "" {
		targetDate, err := time.Parse("2006-01-02", req.TargetDate)
		if err == nil {
			goal.TargetDate = &targetDate
		}
	}

	if err := h.goalRepo.Create(goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create goal"})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// Update updates a learning goal
func (h *LearningGoalHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}

	goal, err := h.goalRepo.FindByID(uint(goalID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}

	if goal.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Category    *string `json:"category"`
		TargetDate  *string `json:"target_date"`
		Progress    *int    `json:"progress"`
		Status      *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Title != nil {
		goal.Title = *req.Title
	}
	if req.Description != nil {
		goal.Description = *req.Description
	}
	if req.Category != nil {
		goal.Category = model.GoalCategory(*req.Category)
	}
	if req.TargetDate != nil {
		if *req.TargetDate == "" {
			goal.TargetDate = nil
		} else {
			targetDate, err := time.Parse("2006-01-02", *req.TargetDate)
			if err == nil {
				goal.TargetDate = &targetDate
			}
		}
	}
	if req.Progress != nil {
		progress := *req.Progress
		if progress < 0 {
			progress = 0
		}
		if progress > 100 {
			progress = 100
		}
		goal.Progress = progress

		// Auto-complete if progress reaches 100
		if progress == 100 && goal.Status == model.GoalStatusActive {
			goal.Status = model.GoalStatusCompleted
			now := time.Now()
			goal.CompletedAt = &now
		}
	}
	if req.Status != nil {
		goal.Status = model.GoalStatus(*req.Status)
		if goal.Status == model.GoalStatusCompleted && goal.CompletedAt == nil {
			now := time.Now()
			goal.CompletedAt = &now
		}
	}

	if err := h.goalRepo.Update(goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update goal"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// Delete deletes a learning goal
func (h *LearningGoalHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}

	goal, err := h.goalRepo.FindByID(uint(goalID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}

	if goal.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	if err := h.goalRepo.Delete(uint(goalID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete goal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "goal deleted"})
}

// GetByID gets a learning goal by ID
func (h *LearningGoalHandler) GetByID(c *gin.Context) {
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}

	goal, err := h.goalRepo.FindByID(uint(goalID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// GetByUserID gets all learning goals for a user
func (h *LearningGoalHandler) GetByUserID(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	goals, err := h.goalRepo.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get goals"})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// GetMyGoals gets all learning goals for the current user
func (h *LearningGoalHandler) GetMyGoals(c *gin.Context) {
	userID := c.GetUint("userID")

	goals, err := h.goalRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get goals"})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// GetStats gets learning goal statistics for a user
func (h *LearningGoalHandler) GetStats(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	stats, err := h.goalRepo.GetStats(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
