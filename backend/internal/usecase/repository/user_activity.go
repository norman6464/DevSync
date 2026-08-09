package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// UserActivityRepository はユーザーアクティビティの永続化に対する、usecase 側が要求する契約。
type UserActivityRepository interface {
	FindByUserID(ctx context.Context, userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error)
}
