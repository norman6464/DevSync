package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// UserActivityHandler はアクティビティタイムライン関連の HTTP ハンドラ。
type UserActivityHandler struct {
	getTimeline *usecase.GetActivityTimelineUseCase
}

// NewUserActivityHandler は UserActivityHandler を生成する。
func NewUserActivityHandler(getTimeline *usecase.GetActivityTimelineUseCase) *UserActivityHandler {
	return &UserActivityHandler{getTimeline: getTimeline}
}

// GetTimeline はユーザーのアクティビティタイムラインを取得する。
func (h *UserActivityHandler) GetTimeline(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	activityType := c.Query("type")
	limit, offset := parseLimitOffset(c)

	activities, total, err := h.getTimeline.Execute(c.Request.Context(), userID, activityType, limit, offset)
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
