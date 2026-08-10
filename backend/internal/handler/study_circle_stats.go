package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// StudyCircleStatsHandler はスタディサークル統計関連の HTTP ハンドラ。
type StudyCircleStatsHandler struct {
	getStats *usecase.GetStudyCircleStatsUseCase
}

// NewStudyCircleStatsHandler は StudyCircleStatsHandler を生成する。
func NewStudyCircleStatsHandler(getStats *usecase.GetStudyCircleStatsUseCase) *StudyCircleStatsHandler {
	return &StudyCircleStatsHandler{getStats: getStats}
}

// GetStats は指定サークルの集計統計を返す。
func (h *StudyCircleStatsHandler) GetStats(c *gin.Context) {
	circleID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.getStats.Execute(c.Request.Context(), circleID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
