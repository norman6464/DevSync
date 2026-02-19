package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// StudyCircleStatsServiceInterface はStudyCircleStatsHandlerが依存するサービスメソッドを定義する。
type StudyCircleStatsServiceInterface interface {
	GetCircleStats(circleID uint) (*model.StudyCircleStats, error)
}

// StudyCircleStatsHandler はスタディサークル統計関連のHTTPハンドラ。
type StudyCircleStatsHandler struct {
	service StudyCircleStatsServiceInterface
}

// NewStudyCircleStatsHandler は新しいStudyCircleStatsHandlerインスタンスを生成する。
func NewStudyCircleStatsHandler(s StudyCircleStatsServiceInterface) *StudyCircleStatsHandler {
	return &StudyCircleStatsHandler{service: s}
}

// GetStats は指定サークルの集計統計を返す。
func (h *StudyCircleStatsHandler) GetStats(c *gin.Context) {
	circleID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetCircleStats(circleID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
