package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// MentionStatsServiceInterface はMentionStatsHandlerが依存するサービスメソッドを定義する。
type MentionStatsServiceInterface interface {
	GetMentionStats(userID uint) (*model.MentionStats, error)
}

// MentionStatsHandler はユーザーメンション統計関連のHTTPハンドラ。
type MentionStatsHandler struct {
	service MentionStatsServiceInterface
}

// NewMentionStatsHandler は新しいMentionStatsHandlerインスタンスを生成する。
func NewMentionStatsHandler(s MentionStatsServiceInterface) *MentionStatsHandler {
	return &MentionStatsHandler{service: s}
}

// GetStats は指定ユーザーのメンション集計統計を返す。
func (h *MentionStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetMentionStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
