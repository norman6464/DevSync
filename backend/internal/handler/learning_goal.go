package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// LearningGoalHandler は学習目標関連のHTTPハンドラ。
// 学習目標のCRUD・統計情報の取得を処理する。
type LearningGoalHandler struct {
	service *service.LearningGoalService
}

// NewLearningGoalHandler は新しいLearningGoalHandlerインスタンスを生成する。
func NewLearningGoalHandler(s *service.LearningGoalService) *LearningGoalHandler {
	return &LearningGoalHandler{service: s}
}

// CreateGoalInput は学習目標作成のリクエストボディ。
type CreateGoalInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	TargetDate  string `json:"target_date"`
}

// Create は新しい学習目標を作成する。
func (h *LearningGoalHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[CreateGoalInput](c)
	if req == nil {
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

	// カテゴリが未指定の場合はデフォルト値を設定
	if req.Category == "" {
		goal.Category = model.GoalCategoryOther
	}

	// 目標日が指定されている場合はパースして設定
	if req.TargetDate != "" {
		targetDate, err := time.Parse("2006-01-02", req.TargetDate)
		if err == nil {
			goal.TargetDate = &targetDate
		}
	}

	if err := h.service.Create(goal); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, goal)
}

// UpdateGoalInput は学習目標更新のリクエストボディ。
type UpdateGoalInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	TargetDate  *string `json:"target_date"`
	Progress    *int    `json:"progress"`
	Status      *string `json:"status"`
}

// Update は指定された学習目標を更新する。
func (h *LearningGoalHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	goalID, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[UpdateGoalInput](c)
	if req == nil {
		return
	}

	updates := &model.LearningGoal{}
	if req.Title != nil {
		updates.Title = *req.Title
	}
	if req.Description != nil {
		updates.Description = *req.Description
	}
	if req.Category != nil {
		updates.Category = model.GoalCategory(*req.Category)
	}
	if req.TargetDate != nil {
		if *req.TargetDate == "" {
			updates.TargetDate = nil
		} else {
			targetDate, err := time.Parse("2006-01-02", *req.TargetDate)
			if err == nil {
				updates.TargetDate = &targetDate
			}
		}
	}
	if req.Progress != nil {
		updates.Progress = *req.Progress
	}
	if req.Status != nil {
		updates.Status = model.GoalStatus(*req.Status)
	}

	goal, err := h.service.Update(goalID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goal)
}

// Delete は指定された学習目標を削除する。
func (h *LearningGoalHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	goalID, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(goalID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// GetByID は指定されたIDの学習目標を取得する。
func (h *LearningGoalHandler) GetByID(c *gin.Context) {
	goalID, ok := parseID(c, "id")
	if !ok {
		return
	}

	goal, err := h.service.GetByID(goalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// GetByUserID は指定されたユーザーの学習目標一覧を取得する。
func (h *LearningGoalHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	goals, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get goals"})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// GetMyGoals は認証ユーザー自身の学習目標一覧を取得する。
func (h *LearningGoalHandler) GetMyGoals(c *gin.Context) {
	userID := c.GetUint("userID")

	goals, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get goals"})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// GetStats は指定されたユーザーの学習目標統計情報を取得する。
func (h *LearningGoalHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	stats, err := h.service.GetStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
