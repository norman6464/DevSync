package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningGoalServiceInterface はLearningGoalServiceが実装すべきインターフェース。
type LearningGoalServiceInterface interface {
	Create(goal *model.LearningGoal) error
	GetByID(id uint) (*model.LearningGoal, error)
	GetByUserID(userID uint) ([]model.LearningGoal, error)
	GetByCategory(userID uint, category string) ([]model.LearningGoal, error)
	GetByStatus(userID uint, status string) ([]model.LearningGoal, error)
	GetStats(userID uint) (*model.LearningGoalStats, error)
	Update(id, userID uint, updates *model.LearningGoal) (*model.LearningGoal, error)
	Delete(id, userID uint) error
	GetDeadlineAlerts(userID uint) ([]model.GoalDeadlineAlert, error)
}

// LearningGoalHandler は学習目標関連のHTTPハンドラ。
// 学習目標のCRUD・統計情報の取得を処理する。
type LearningGoalHandler struct {
	service LearningGoalServiceInterface
}

// NewLearningGoalHandler は新しいLearningGoalHandlerインスタンスを生成する。
func NewLearningGoalHandler(s LearningGoalServiceInterface) *LearningGoalHandler {
	return &LearningGoalHandler{service: s}
}

// Create は新しい学習目標を作成する。
func (h *LearningGoalHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.CreateGoalRequest](c)
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

// Update は指定された学習目標を更新する。
func (h *LearningGoalHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	goalID, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdateGoalRequest](c)
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
		respondNotFound(c, "goal not found")
		return
	}

	respondOK(c, goal)
}

// GetByUserID は指定されたユーザーの学習目標一覧を取得する。
func (h *LearningGoalHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	goals, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goals)
}

// GetMyGoals は認証ユーザー自身の学習目標一覧を取得する。
func (h *LearningGoalHandler) GetMyGoals(c *gin.Context) {
	userID := c.GetUint("userID")

	goals, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goals)
}

// GetDeadlineAlerts は認証ユーザーのデッドラインアラートを取得する。
func (h *LearningGoalHandler) GetDeadlineAlerts(c *gin.Context) {
	userID := c.GetUint("userID")

	alerts, err := h.service.GetDeadlineAlerts(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	if alerts == nil {
		alerts = []model.GoalDeadlineAlert{}
	}

	respondOK(c, alerts)
}

// GetByCategory は認証ユーザーの学習目標をカテゴリでフィルタリングして取得する。
func (h *LearningGoalHandler) GetByCategory(c *gin.Context) {
	userID := c.GetUint("userID")
	category := c.Param("category")

	goals, err := h.service.GetByCategory(userID, category)
	if err != nil {
		respondError(c, err)
		return
	}
	if goals == nil {
		goals = []model.LearningGoal{}
	}

	respondOK(c, goals)
}

// GetByStatus は認証ユーザーの学習目標をステータスでフィルタリングして取得する。
func (h *LearningGoalHandler) GetByStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	status := c.Param("status")

	goals, err := h.service.GetByStatus(userID, status)
	if err != nil {
		respondError(c, err)
		return
	}
	if goals == nil {
		goals = []model.LearningGoal{}
	}

	respondOK(c, goals)
}

// GetStats は指定されたユーザーの学習目標統計情報を取得する。
func (h *LearningGoalHandler) GetStats(c *gin.Context) {
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
