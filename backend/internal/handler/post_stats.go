package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// PostStatsHandler はユーザー投稿統計関連の HTTP ハンドラ。
type PostStatsHandler struct {
	getStats *usecase.GetPostStatsUseCase
}

// NewPostStatsHandler は PostStatsHandler を生成する。
func NewPostStatsHandler(getStats *usecase.GetPostStatsUseCase) *PostStatsHandler {
	return &PostStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーの投稿集計統計を返す。
func (h *PostStatsHandler) GetStats(c *gin.Context) {
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
