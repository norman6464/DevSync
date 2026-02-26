package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// WeeklyGoalServiceInterface はWeeklyGoalHandlerが依存するサービスメソッドを定義する。
type WeeklyGoalServiceInterface interface {
	SetGoal(userID uint, category string, targetMinutes int) (*model.WeeklyGoal, error)
	GetGoals(userID uint) ([]model.WeeklyGoal, error)
	GetProgress(userID uint) ([]model.WeeklyGoalProgress, error)
}

// WeeklyGoalHandler はカテゴリ別週間学習目標のHTTPハンドラ。
type WeeklyGoalHandler struct {
	service WeeklyGoalServiceInterface
}

// NewWeeklyGoalHandler は新しいWeeklyGoalHandlerインスタンスを生成する。
func NewWeeklyGoalHandler(s WeeklyGoalServiceInterface) *WeeklyGoalHandler {
	return &WeeklyGoalHandler{service: s}
}

// setWeeklyGoalRequest は週間目標設定リクエスト。
type setWeeklyGoalRequest struct {
	Category      string `json:"category" binding:"required"`
	TargetMinutes int    `json:"target_minutes" binding:"required,min=0"`
}

// SetGoal はカテゴリ別の週間学習目標を設定する。
func (h *WeeklyGoalHandler) SetGoal(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[setWeeklyGoalRequest](c)
	if req == nil {
		return
	}

	goal, err := h.service.SetGoal(userID, req.Category, req.TargetMinutes)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goal)
}

// GetGoals は認証ユーザーの全カテゴリ週間目標を取得する。
func (h *WeeklyGoalHandler) GetGoals(c *gin.Context) {
	userID := c.GetUint("userID")

	goals, err := h.service.GetGoals(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(goals))
}

// GetProgress は認証ユーザーの全カテゴリ週間目標の達成状況を取得する。
func (h *WeeklyGoalHandler) GetProgress(c *gin.Context) {
	userID := c.GetUint("userID")

	progress, err := h.service.GetProgress(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(progress))
}
