package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// CommentStatsServiceInterface はCommentStatsHandlerが依存するサービスメソッドを定義する。
type CommentStatsServiceInterface interface {
	GetCommentStats(userID uint) (*model.CommentStats, error)
}

// CommentStatsHandler はユーザーコメント活動統計関連のHTTPハンドラ。
type CommentStatsHandler struct {
	service CommentStatsServiceInterface
}

// NewCommentStatsHandler は新しいCommentStatsHandlerインスタンスを生成する。
func NewCommentStatsHandler(s CommentStatsServiceInterface) *CommentStatsHandler {
	return &CommentStatsHandler{service: s}
}

// GetStats は指定ユーザーのコメント活動集計統計を返す。
func (h *CommentStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetCommentStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
