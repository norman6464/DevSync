package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// UserDashboardService はユーザーダッシュボード統計のビジネスロジックを提供する。
type UserDashboardService struct {
	repo repository.UserDashboardRepositoryInterface
}

// NewUserDashboardService は新しいUserDashboardServiceインスタンスを生成する。
func NewUserDashboardService(repo repository.UserDashboardRepositoryInterface) *UserDashboardService {
	return &UserDashboardService{repo: repo}
}

// GetStats は指定ユーザーのダッシュボード統計情報を取得する。
// userID=0 の場合はBadRequestエラーを返す。
func (s *UserDashboardService) GetStats(userID uint) (*model.UserDashboardStats, error) {
	if userID == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	return s.repo.GetDashboardStats(userID)
}
