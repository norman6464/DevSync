package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// WeeklyGoalHandler はカテゴリ別週間学習目標の HTTP ハンドラー。各操作は 1 責務の usecase に委譲する。
type WeeklyGoalHandler struct {
	setGoal     *usecase.SetWeeklyGoalUseCase
	listGoals   *usecase.ListWeeklyGoalsUseCase
	getProgress *usecase.GetWeeklyGoalProgressUseCase
}

// NewWeeklyGoalHandler は WeeklyGoalHandler を生成する。
func NewWeeklyGoalHandler(
	setGoal *usecase.SetWeeklyGoalUseCase,
	listGoals *usecase.ListWeeklyGoalsUseCase,
	getProgress *usecase.GetWeeklyGoalProgressUseCase,
) *WeeklyGoalHandler {
	return &WeeklyGoalHandler{
		setGoal:     setGoal,
		listGoals:   listGoals,
		getProgress: getProgress,
	}
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

	goal, err := h.setGoal.Execute(c.Request.Context(), userID, req.Category, req.TargetMinutes)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goal)
}

// GetGoals は認証ユーザーの全カテゴリ週間目標を取得する。
func (h *WeeklyGoalHandler) GetGoals(c *gin.Context) {
	userID := c.GetUint("userID")

	goals, err := h.listGoals.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(goals))
}

// GetProgress は認証ユーザーの全カテゴリ週間目標の達成状況を取得する。
func (h *WeeklyGoalHandler) GetProgress(c *gin.Context) {
	userID := c.GetUint("userID")

	progress, err := h.getProgress.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(progress))
}
