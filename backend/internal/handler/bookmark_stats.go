package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// BookmarkStatsHandler はユーザーブックマーク統計関連の HTTP ハンドラ。
type BookmarkStatsHandler struct {
	getStats *usecase.GetBookmarkStatsUseCase
}

// NewBookmarkStatsHandler は BookmarkStatsHandler を生成する。
func NewBookmarkStatsHandler(getStats *usecase.GetBookmarkStatsUseCase) *BookmarkStatsHandler {
	return &BookmarkStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのブックマーク集計統計を返す。
func (h *BookmarkStatsHandler) GetStats(c *gin.Context) {
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
