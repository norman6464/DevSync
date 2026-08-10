package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// QAStatsHandler はユーザー Q&A 統計関連の HTTP ハンドラ。
type QAStatsHandler struct {
	getStats *usecase.GetQAStatsUseCase
}

// NewQAStatsHandler は QAStatsHandler を生成する。
func NewQAStatsHandler(getStats *usecase.GetQAStatsUseCase) *QAStatsHandler {
	return &QAStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーの Q&A 活動集計統計を返す。
func (h *QAStatsHandler) GetStats(c *gin.Context) {
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
