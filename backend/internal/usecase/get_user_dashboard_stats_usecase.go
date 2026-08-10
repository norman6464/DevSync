package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetUserDashboardStatsUseCase は指定ユーザーのダッシュボード統計を取得する。
type GetUserDashboardStatsUseCase struct {
	dashboard repository.UserDashboardRepository
}

// NewGetUserDashboardStatsUseCase は GetUserDashboardStatsUseCase を生成する。
func NewGetUserDashboardStatsUseCase(dashboard repository.UserDashboardRepository) *GetUserDashboardStatsUseCase {
	return &GetUserDashboardStatsUseCase{dashboard: dashboard}
}

// Execute はユーザー ID を検証し、ダッシュボード統計を返す。
func (uc *GetUserDashboardStatsUseCase) Execute(ctx context.Context, userID uint) (*model.UserDashboardStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.dashboard.GetDashboardStats(ctx, userID)
}
