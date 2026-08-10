package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockFollowStatsRepo は usecase/repository.FollowStatsRepository のモック。
type mockFollowStatsRepo struct{ mock.Mock }

func (m *mockFollowStatsRepo) GetFollowStats(ctx context.Context, userID uint) (*model.FollowStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.FollowStats)
	return s, args.Error(1)
}

func TestGetFollowStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockFollowStatsRepo)
		expected := &model.FollowStats{FollowerCount: 3, FollowingCount: 2}
		stats.On("GetFollowStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetFollowStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockFollowStatsRepo)
		uc := usecase.NewGetFollowStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetFollowStats")
	})
}
