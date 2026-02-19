package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningResourceStatsServiceInterface はLearningResourceStatsHandlerが依存するサービスメソッドを定義する。
type LearningResourceStatsServiceInterface interface {
	GetLearningResourceStats(userID uint) (*model.LearningResourceStats, error)
}

// LearningResourceStatsHandler はユーザー学習リソース統計関連のHTTPハンドラ。
type LearningResourceStatsHandler struct {
	service LearningResourceStatsServiceInterface
}

// NewLearningResourceStatsHandler は新しいLearningResourceStatsHandlerインスタンスを生成する。
func NewLearningResourceStatsHandler(s LearningResourceStatsServiceInterface) *LearningResourceStatsHandler {
	return &LearningResourceStatsHandler{service: s}
}

// GetStats は指定ユーザーの学習リソース活動集計統計を返す。
func (h *LearningResourceStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetLearningResourceStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
