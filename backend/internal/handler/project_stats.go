package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// ProjectStatsServiceInterface はProjectStatsHandlerが依存するサービスメソッドを定義する。
type ProjectStatsServiceInterface interface {
	GetProjectStats(userID uint) (*model.ProjectStats, error)
}

// ProjectStatsHandler はユーザープロジェクト統計関連のHTTPハンドラ。
type ProjectStatsHandler struct {
	service ProjectStatsServiceInterface
}

// NewProjectStatsHandler は新しいProjectStatsHandlerインスタンスを生成する。
func NewProjectStatsHandler(s ProjectStatsServiceInterface) *ProjectStatsHandler {
	return &ProjectStatsHandler{service: s}
}

// GetStats は指定ユーザーのプロジェクト活動集計統計を返す。
func (h *ProjectStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetProjectStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
