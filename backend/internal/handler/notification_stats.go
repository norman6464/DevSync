package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NotificationStatsServiceInterface はNotificationStatsHandlerが依存するサービスメソッドを定義する。
type NotificationStatsServiceInterface interface {
	GetNotificationStats(userID uint) (*model.NotificationStats, error)
}

// NotificationStatsHandler はユーザー通知統計関連のHTTPハンドラ。
type NotificationStatsHandler struct {
	service NotificationStatsServiceInterface
}

// NewNotificationStatsHandler は新しいNotificationStatsHandlerインスタンスを生成する。
func NewNotificationStatsHandler(s NotificationStatsServiceInterface) *NotificationStatsHandler {
	return &NotificationStatsHandler{service: s}
}

// GetStats は指定ユーザーの通知集計統計を返す。
func (h *NotificationStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetNotificationStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
