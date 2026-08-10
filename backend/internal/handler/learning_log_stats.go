package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LearningLogStatsHandler はユーザー学習ログ統計関連の HTTP ハンドラ。
type LearningLogStatsHandler struct {
	getStats *usecase.GetLearningLogStatsUseCase
}

// NewLearningLogStatsHandler は LearningLogStatsHandler を生成する。
func NewLearningLogStatsHandler(getStats *usecase.GetLearningLogStatsUseCase) *LearningLogStatsHandler {
	return &LearningLogStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーの学習ログ集計統計を返す。
func (h *LearningLogStatsHandler) GetStats(c *gin.Context) {
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
