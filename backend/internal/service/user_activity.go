package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// UserActivityService はユーザーアクティビティのビジネスロジックを提供する。
type UserActivityService struct {
	repo repository.UserActivityRepositoryInterface
}

// NewUserActivityService は新しいUserActivityServiceインスタンスを生成する。
func NewUserActivityService(repo repository.UserActivityRepositoryInterface) *UserActivityService {
	return &UserActivityService{repo: repo}
}

// RecordActivity はアクティビティを記録する。
func (s *UserActivityService) RecordActivity(userID uint, activityType model.ActivityType, targetType string, targetID uint, metadata string) error {
	activity := &model.UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		TargetType:   targetType,
		TargetID:     targetID,
		Metadata:     metadata,
	}
	return s.repo.Create(activity)
}

// GetTimeline はユーザーのアクティビティタイムラインを取得する。
func (s *UserActivityService) GetTimeline(userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error) {
	return s.repo.FindByUserID(userID, activityType, limit, offset)
}
