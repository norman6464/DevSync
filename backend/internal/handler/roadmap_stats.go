package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// RoadmapStatsServiceInterface はRoadmapStatsHandlerが依存するサービスメソッドを定義する。
type RoadmapStatsServiceInterface interface {
	GetRoadmapStats(userID uint) (*model.RoadmapStats, error)
}

// RoadmapStatsHandler はユーザーロードマップ統計関連のHTTPハンドラ。
type RoadmapStatsHandler struct {
	service RoadmapStatsServiceInterface
}

// NewRoadmapStatsHandler は新しいRoadmapStatsHandlerインスタンスを生成する。
func NewRoadmapStatsHandler(s RoadmapStatsServiceInterface) *RoadmapStatsHandler {
	return &RoadmapStatsHandler{service: s}
}

// GetStats は指定ユーザーのロードマップ統計を返す。
func (h *RoadmapStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetRoadmapStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
