package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// MessageStatsHandler はユーザーメッセージ統計関連の HTTP ハンドラ。
type MessageStatsHandler struct {
	getStats *usecase.GetMessageStatsUseCase
}

// NewMessageStatsHandler は MessageStatsHandler を生成する。
func NewMessageStatsHandler(getStats *usecase.GetMessageStatsUseCase) *MessageStatsHandler {
	return &MessageStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのメッセージ集計統計を返す。
func (h *MessageStatsHandler) GetStats(c *gin.Context) {
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
