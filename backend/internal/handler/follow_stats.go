package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// FollowStatsServiceInterface はFollowStatsHandlerが依存するサービスメソッドを定義する。
type FollowStatsServiceInterface interface {
	GetFollowStats(userID uint) (*model.FollowStats, error)
}

// FollowStatsHandler はユーザーフォロー統計関連のHTTPハンドラ。
type FollowStatsHandler struct {
	service FollowStatsServiceInterface
}

// NewFollowStatsHandler は新しいFollowStatsHandlerインスタンスを生成する。
func NewFollowStatsHandler(s FollowStatsServiceInterface) *FollowStatsHandler {
	return &FollowStatsHandler{service: s}
}

// GetStats は指定ユーザーのフォロー関係集計統計を返す。
func (h *FollowStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetFollowStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
