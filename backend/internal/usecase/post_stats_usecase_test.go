package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockPostStatsRepo は usecase/repository.PostStatsRepository のモック。
type mockPostStatsRepo struct{ mock.Mock }

func (m *mockPostStatsRepo) GetPostStats(ctx context.Context, userID uint) (*model.PostStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.PostStats)
	return s, args.Error(1)
}

func TestGetPostStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockPostStatsRepo)
		expected := &model.PostStats{
			TotalPosts: 5, PublishedPosts: 4, DraftPosts: 1,
			TotalLikesReceived: 12, TotalComments: 7, PostsThisMonth: 2,
		}
		stats.On("GetPostStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetPostStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockPostStatsRepo)
		uc := usecase.NewGetPostStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetPostStats")
	})
}
