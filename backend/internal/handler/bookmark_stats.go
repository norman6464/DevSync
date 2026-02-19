package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// BookmarkStatsServiceInterface はBookmarkStatsHandlerが依存するサービスメソッドを定義する。
type BookmarkStatsServiceInterface interface {
	GetBookmarkStats(userID uint) (*model.BookmarkStats, error)
}

// BookmarkStatsHandler はユーザーブックマーク統計関連のHTTPハンドラ。
type BookmarkStatsHandler struct {
	service BookmarkStatsServiceInterface
}

// NewBookmarkStatsHandler は新しいBookmarkStatsHandlerインスタンスを生成する。
func NewBookmarkStatsHandler(s BookmarkStatsServiceInterface) *BookmarkStatsHandler {
	return &BookmarkStatsHandler{service: s}
}

// GetStats は指定ユーザーのブックマーク集計統計を返す。
func (h *BookmarkStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetBookmarkStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
