package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// RoadmapStatsHandler はユーザーロードマップ統計関連の HTTP ハンドラ。
type RoadmapStatsHandler struct {
	getStats *usecase.GetRoadmapStatsUseCase
}

// NewRoadmapStatsHandler は RoadmapStatsHandler を生成する。
func NewRoadmapStatsHandler(getStats *usecase.GetRoadmapStatsUseCase) *RoadmapStatsHandler {
	return &RoadmapStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのロードマップ統計を返す。
func (h *RoadmapStatsHandler) GetStats(c *gin.Context) {
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
