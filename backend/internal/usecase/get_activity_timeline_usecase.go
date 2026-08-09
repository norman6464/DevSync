package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetActivityTimelineUseCase はユーザーのアクティビティタイムラインを取得する。
type GetActivityTimelineUseCase struct {
	activities repository.UserActivityRepository
}

// NewGetActivityTimelineUseCase は GetActivityTimelineUseCase を生成する。
func NewGetActivityTimelineUseCase(activities repository.UserActivityRepository) *GetActivityTimelineUseCase {
	return &GetActivityTimelineUseCase{activities: activities}
}

// Execute はタイムライン（種別フィルタ・ページネーション対応）と総件数を返す。
func (uc *GetActivityTimelineUseCase) Execute(ctx context.Context, userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error) {
	return uc.activities.FindByUserID(ctx, userID, activityType, limit, offset)
}
