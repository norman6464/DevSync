package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningLogStatsServiceInterface はLearningLogStatsHandlerが依存するサービスメソッドを定義する。
type LearningLogStatsServiceInterface interface {
	GetLearningLogStats(userID uint) (*model.LearningLogStats, error)
}

// LearningLogStatsHandler はユーザー学習ログ統計関連のHTTPハンドラ。
type LearningLogStatsHandler struct {
	service LearningLogStatsServiceInterface
}

// NewLearningLogStatsHandler は新しいLearningLogStatsHandlerインスタンスを生成する。
func NewLearningLogStatsHandler(s LearningLogStatsServiceInterface) *LearningLogStatsHandler {
	return &LearningLogStatsHandler{service: s}
}

// GetStats は指定ユーザーの学習ログ集計統計を返す。
func (h *LearningLogStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetLearningLogStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
