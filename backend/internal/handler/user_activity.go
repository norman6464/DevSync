package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// UserActivityServiceInterface はアクティビティサービスの抽象インターフェース。
type UserActivityServiceInterface interface {
	GetTimeline(userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error)
}

// UserActivityHandler はアクティビティタイムライン関連のHTTPハンドラ。
type UserActivityHandler struct {
	service UserActivityServiceInterface
}

// NewUserActivityHandler は新しいUserActivityHandlerインスタンスを生成する。
func NewUserActivityHandler(s UserActivityServiceInterface) *UserActivityHandler {
	return &UserActivityHandler{service: s}
}

// GetTimeline はユーザーのアクティビティタイムラインを取得する。
func (h *UserActivityHandler) GetTimeline(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	activityType := c.Query("type")
	limit, offset := parseLimitOffset(c)

	activities, total, err := h.service.GetTimeline(userID, activityType, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.UserActivityListResponse{
		Activities: ensureSlice(activities),
		Total:      total,
		Limit:      limit,
		Offset:     offset,
	})
}
