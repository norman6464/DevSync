package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// MessageStatsServiceInterface はMessageStatsHandlerが依存するサービスメソッドを定義する。
type MessageStatsServiceInterface interface {
	GetMessageStats(userID uint) (*model.MessageStats, error)
}

// MessageStatsHandler はユーザーメッセージ統計関連のHTTPハンドラ。
type MessageStatsHandler struct {
	service MessageStatsServiceInterface
}

// NewMessageStatsHandler は新しいMessageStatsHandlerインスタンスを生成する。
func NewMessageStatsHandler(s MessageStatsServiceInterface) *MessageStatsHandler {
	return &MessageStatsHandler{service: s}
}

// GetStats は指定ユーザーのメッセージ集計統計を返す。
func (h *MessageStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetMessageStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
