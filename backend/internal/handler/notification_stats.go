package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NotificationStatsHandler はユーザー通知統計関連の HTTP ハンドラ。
type NotificationStatsHandler struct {
	getStats *usecase.GetNotificationStatsUseCase
}

// NewNotificationStatsHandler は NotificationStatsHandler を生成する。
func NewNotificationStatsHandler(getStats *usecase.GetNotificationStatsUseCase) *NotificationStatsHandler {
	return &NotificationStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーの通知集計統計を返す。
func (h *NotificationStatsHandler) GetStats(c *gin.Context) {
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
