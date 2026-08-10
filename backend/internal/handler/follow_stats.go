package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// FollowStatsHandler はユーザーフォロー統計関連の HTTP ハンドラ。
type FollowStatsHandler struct {
	getStats *usecase.GetFollowStatsUseCase
}

// NewFollowStatsHandler は FollowStatsHandler を生成する。
func NewFollowStatsHandler(getStats *usecase.GetFollowStatsUseCase) *FollowStatsHandler {
	return &FollowStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのフォロー関係集計統計を返す。
func (h *FollowStatsHandler) GetStats(c *gin.Context) {
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
