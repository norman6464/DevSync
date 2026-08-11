package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LearningDashboardHandler は学習ダッシュボード統合サマリーのHTTPハンドラ。
type LearningDashboardHandler struct {
	summary *usecase.GetLearningDashboardSummaryUseCase
}

// NewLearningDashboardHandler は新しいLearningDashboardHandlerインスタンスを生成する。
func NewLearningDashboardHandler(summary *usecase.GetLearningDashboardSummaryUseCase) *LearningDashboardHandler {
	return &LearningDashboardHandler{summary: summary}
}

// GetSummary は認証ユーザーの学習ダッシュボード統合サマリーを返す。
func (h *LearningDashboardHandler) GetSummary(c *gin.Context) {
	userID := c.GetUint("userID")

	summary, err := h.summary.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, summary)
}
