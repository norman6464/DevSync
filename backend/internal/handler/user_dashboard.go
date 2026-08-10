package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// UserDashboardHandler はユーザーダッシュボード統計の HTTP ハンドラー。
type UserDashboardHandler struct {
	getStats *usecase.GetUserDashboardStatsUseCase
}

// NewUserDashboardHandler は UserDashboardHandler を生成する。
func NewUserDashboardHandler(getStats *usecase.GetUserDashboardStatsUseCase) *UserDashboardHandler {
	return &UserDashboardHandler{getStats: getStats}
}

// GetStats は指定ユーザーのダッシュボード統計情報を返す。
// GET /api/v1/users/:id/dashboard-stats
func (h *UserDashboardHandler) GetStats(c *gin.Context) {
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
