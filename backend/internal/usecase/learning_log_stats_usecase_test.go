package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockLearningLogStatsRepo は usecase/repository.LearningLogStatsRepository のモック。
type mockLearningLogStatsRepo struct{ mock.Mock }

func (m *mockLearningLogStatsRepo) GetLearningLogStats(ctx context.Context, userID uint) (*model.LearningLogStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.LearningLogStats)
	return s, args.Error(1)
}

func TestGetLearningLogStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockLearningLogStatsRepo)
		expected := &model.LearningLogStats{TotalLogs: 4, TotalDuration: 300, CategoryCount: 2, LogsThisMonth: 3}
		stats.On("GetLearningLogStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetLearningLogStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockLearningLogStatsRepo)
		uc := usecase.NewGetLearningLogStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetLearningLogStats")
	})
}
