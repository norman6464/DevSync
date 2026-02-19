package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// ReactionStatsServiceInterface はReactionStatsHandlerが依存するサービスメソッドを定義する。
type ReactionStatsServiceInterface interface {
	GetReactionStats(userID uint) (*model.ReactionStats, error)
}

// ReactionStatsHandler はユーザーリアクション統計関連のHTTPハンドラ。
type ReactionStatsHandler struct {
	service ReactionStatsServiceInterface
}

// NewReactionStatsHandler は新しいReactionStatsHandlerインスタンスを生成する。
func NewReactionStatsHandler(s ReactionStatsServiceInterface) *ReactionStatsHandler {
	return &ReactionStatsHandler{service: s}
}

// GetStats は指定ユーザーのリアクション集計統計を返す。
func (h *ReactionStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetReactionStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
