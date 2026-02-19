package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// UserDashboardServiceInterface はUserDashboardHandlerが依存するサービスインターフェース。
type UserDashboardServiceInterface interface {
	GetStats(userID uint) (*model.UserDashboardStats, error)
}

// UserDashboardHandler はユーザーダッシュボード統計のHTTPハンドラー。
type UserDashboardHandler struct {
	service UserDashboardServiceInterface
}

// NewUserDashboardHandler は新しいUserDashboardHandlerを生成する。
func NewUserDashboardHandler(service UserDashboardServiceInterface) *UserDashboardHandler {
	return &UserDashboardHandler{service: service}
}

// GetStats は指定ユーザーのダッシュボード統計情報を返す。
// GET /api/v1/users/:id/dashboard-stats
func (h *UserDashboardHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, stats)
}
