package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockUserDashboardRepo は usecase/repository.UserDashboardRepository のモック。
type mockUserDashboardRepo struct{ mock.Mock }

func (m *mockUserDashboardRepo) GetDashboardStats(ctx context.Context, userID uint) (*model.UserDashboardStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.UserDashboardStats)
	return s, args.Error(1)
}

func TestGetUserDashboardStatsUseCase_Execute(t *testing.T) {
	t.Run("ダッシュボード統計を返す", func(t *testing.T) {
		repo := new(mockUserDashboardRepo)
		expected := &model.UserDashboardStats{
			PostCount: 10, LikesReceived: 50, CommentsReceived: 7,
			ViewsReceived: 120, FollowerCount: 20, FollowingCount: 3,
		}
		repo.On("GetDashboardStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetUserDashboardStatsUseCase(repo)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		repo := new(mockUserDashboardRepo)
		uc := usecase.NewGetUserDashboardStatsUseCase(repo)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "GetDashboardStats")
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockUserDashboardRepo)
		repo.On("GetDashboardStats", mock.Anything, uint(10)).
			Return((*model.UserDashboardStats)(nil), errors.New("db error"))
		uc := usecase.NewGetUserDashboardStatsUseCase(repo)

		_, err := uc.Execute(context.Background(), 10)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
