package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// CodeSnippetStatsHandler はユーザーコードスニペット統計関連の HTTP ハンドラ。
type CodeSnippetStatsHandler struct {
	getStats *usecase.GetCodeSnippetStatsUseCase
}

// NewCodeSnippetStatsHandler は CodeSnippetStatsHandler を生成する。
func NewCodeSnippetStatsHandler(getStats *usecase.GetCodeSnippetStatsUseCase) *CodeSnippetStatsHandler {
	return &CodeSnippetStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのコードスニペット活動集計統計を返す。
func (h *CodeSnippetStatsHandler) GetStats(c *gin.Context) {
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
