package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostStatsServiceInterface はPostStatsHandlerが依存するサービスメソッドを定義する。
type PostStatsServiceInterface interface {
	GetPostStats(userID uint) (*model.PostStats, error)
}

// PostStatsHandler はユーザー投稿統計関連のHTTPハンドラ。
type PostStatsHandler struct {
	service PostStatsServiceInterface
}

// NewPostStatsHandler は新しいPostStatsHandlerインスタンスを生成する。
func NewPostStatsHandler(s PostStatsServiceInterface) *PostStatsHandler {
	return &PostStatsHandler{service: s}
}

// GetStats は指定ユーザーの投稿集計統計を返す。
func (h *PostStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetPostStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
