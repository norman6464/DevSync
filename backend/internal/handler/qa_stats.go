package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// QAStatsServiceInterface はQAStatsHandlerが依存するサービスメソッドを定義する。
type QAStatsServiceInterface interface {
	GetQAStats(userID uint) (*model.QAStats, error)
}

// QAStatsHandler はユーザーQ&A統計関連のHTTPハンドラ。
type QAStatsHandler struct {
	service QAStatsServiceInterface
}

// NewQAStatsHandler は新しいQAStatsHandlerインスタンスを生成する。
func NewQAStatsHandler(s QAStatsServiceInterface) *QAStatsHandler {
	return &QAStatsHandler{service: s}
}

// GetStats は指定ユーザーのQ&A活動集計統計を返す。
func (h *QAStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetQAStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
