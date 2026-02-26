package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningDashboardServiceInterface はLearningDashboardHandlerが依存するサービスのインターフェース。
type LearningDashboardServiceInterface interface {
	GetSummary(userID uint) (*model.LearningDashboardSummary, error)
}

// LearningDashboardHandler は学習ダッシュボード統合サマリーのHTTPハンドラ。
type LearningDashboardHandler struct {
	service LearningDashboardServiceInterface
}

// NewLearningDashboardHandler は新しいLearningDashboardHandlerインスタンスを生成する。
func NewLearningDashboardHandler(s LearningDashboardServiceInterface) *LearningDashboardHandler {
	return &LearningDashboardHandler{service: s}
}

// GetSummary は認証ユーザーの学習ダッシュボード統合サマリーを返す。
func (h *LearningDashboardHandler) GetSummary(c *gin.Context) {
	userID := c.GetUint("userID")

	summary, err := h.service.GetSummary(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, summary)
}
