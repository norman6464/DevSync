package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LearningResourceStatsHandler はユーザー学習リソース統計関連の HTTP ハンドラ。
type LearningResourceStatsHandler struct {
	getStats *usecase.GetLearningResourceStatsUseCase
}

// NewLearningResourceStatsHandler は LearningResourceStatsHandler を生成する。
func NewLearningResourceStatsHandler(getStats *usecase.GetLearningResourceStatsUseCase) *LearningResourceStatsHandler {
	return &LearningResourceStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーの学習リソース活動集計統計を返す。
func (h *LearningResourceStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.getStats.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
