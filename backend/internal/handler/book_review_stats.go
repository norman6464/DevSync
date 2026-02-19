package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// BookReviewStatsServiceInterface はBookReviewStatsHandlerが依存するサービスメソッドを定義する。
type BookReviewStatsServiceInterface interface {
	GetBookReviewStats(userID uint) (*model.BookReviewStats, error)
}

// BookReviewStatsHandler はユーザー書籍レビュー統計関連のHTTPハンドラ。
type BookReviewStatsHandler struct {
	service BookReviewStatsServiceInterface
}

// NewBookReviewStatsHandler は新しいBookReviewStatsHandlerインスタンスを生成する。
func NewBookReviewStatsHandler(s BookReviewStatsServiceInterface) *BookReviewStatsHandler {
	return &BookReviewStatsHandler{service: s}
}

// GetStats は指定ユーザーの書籍レビュー集計統計を返す。
func (h *BookReviewStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetBookReviewStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
