package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockLearningResourceStatsRepo は usecase/repository.LearningResourceStatsRepository のモック。
type mockLearningResourceStatsRepo struct{ mock.Mock }

func (m *mockLearningResourceStatsRepo) GetLearningResourceStats(ctx context.Context, userID uint) (*model.LearningResourceStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.LearningResourceStats)
	return s, args.Error(1)
}

func TestGetLearningResourceStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockLearningResourceStatsRepo)
		expected := &model.LearningResourceStats{TotalResources: 5, TotalLikes: 12, TotalSaves: 3, CategoryCount: 2}
		stats.On("GetLearningResourceStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetLearningResourceStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockLearningResourceStatsRepo)
		uc := usecase.NewGetLearningResourceStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetLearningResourceStats")
	})
}
