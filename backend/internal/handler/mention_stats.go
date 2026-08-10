package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// MentionStatsHandler はユーザーメンション統計関連の HTTP ハンドラ。
type MentionStatsHandler struct {
	getStats *usecase.GetMentionStatsUseCase
}

// NewMentionStatsHandler は MentionStatsHandler を生成する。
func NewMentionStatsHandler(getStats *usecase.GetMentionStatsUseCase) *MentionStatsHandler {
	return &MentionStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのメンション集計統計を返す。
func (h *MentionStatsHandler) GetStats(c *gin.Context) {
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
