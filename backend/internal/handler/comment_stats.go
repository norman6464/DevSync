package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// CommentStatsHandler はユーザーコメント活動統計関連の HTTP ハンドラ。
type CommentStatsHandler struct {
	getStats *usecase.GetCommentStatsUseCase
}

// NewCommentStatsHandler は CommentStatsHandler を生成する。
func NewCommentStatsHandler(getStats *usecase.GetCommentStatsUseCase) *CommentStatsHandler {
	return &CommentStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのコメント活動集計統計を返す。
func (h *CommentStatsHandler) GetStats(c *gin.Context) {
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
