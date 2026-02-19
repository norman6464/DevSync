package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// CodeSnippetStatsServiceInterface はCodeSnippetStatsHandlerが依存するサービスメソッドを定義する。
type CodeSnippetStatsServiceInterface interface {
	GetCodeSnippetStats(userID uint) (*model.CodeSnippetStats, error)
}

// CodeSnippetStatsHandler はユーザーコードスニペット統計関連のHTTPハンドラ。
type CodeSnippetStatsHandler struct {
	service CodeSnippetStatsServiceInterface
}

// NewCodeSnippetStatsHandler は新しいCodeSnippetStatsHandlerインスタンスを生成する。
func NewCodeSnippetStatsHandler(s CodeSnippetStatsServiceInterface) *CodeSnippetStatsHandler {
	return &CodeSnippetStatsHandler{service: s}
}

// GetStats は指定ユーザーのコードスニペット活動集計統計を返す。
func (h *CodeSnippetStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetCodeSnippetStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
