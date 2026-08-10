package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// UserDashboardRepository はユーザーダッシュボード統計の取得に対する、usecase 側が要求する契約。
type UserDashboardRepository interface {
	GetDashboardStats(ctx context.Context, userID uint) (*model.UserDashboardStats, error)
}
