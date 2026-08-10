package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// BookReviewStatsHandler はユーザー書籍レビュー統計関連の HTTP ハンドラ。
type BookReviewStatsHandler struct {
	getStats *usecase.GetBookReviewStatsUseCase
}

// NewBookReviewStatsHandler は BookReviewStatsHandler を生成する。
func NewBookReviewStatsHandler(getStats *usecase.GetBookReviewStatsUseCase) *BookReviewStatsHandler {
	return &BookReviewStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーの書籍レビュー集計統計を返す。
func (h *BookReviewStatsHandler) GetStats(c *gin.Context) {
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
