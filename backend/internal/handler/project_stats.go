package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ProjectStatsHandler はユーザープロジェクト統計関連の HTTP ハンドラ。
type ProjectStatsHandler struct {
	getStats *usecase.GetProjectStatsUseCase
}

// NewProjectStatsHandler は ProjectStatsHandler を生成する。
func NewProjectStatsHandler(getStats *usecase.GetProjectStatsUseCase) *ProjectStatsHandler {
	return &ProjectStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのプロジェクト活動集計統計を返す。
func (h *ProjectStatsHandler) GetStats(c *gin.Context) {
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
