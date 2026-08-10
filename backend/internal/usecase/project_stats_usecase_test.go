package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockProjectStatsRepo は usecase/repository.ProjectStatsRepository のモック。
type mockProjectStatsRepo struct{ mock.Mock }

func (m *mockProjectStatsRepo) GetProjectStats(ctx context.Context, userID uint) (*model.ProjectStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ProjectStats)
	return s, args.Error(1)
}

func TestGetProjectStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockProjectStatsRepo)
		expected := &model.ProjectStats{TotalProjects: 4, FeaturedProjects: 1, OngoingProjects: 3, CompletedProjects: 1}
		stats.On("GetProjectStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetProjectStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockProjectStatsRepo)
		uc := usecase.NewGetProjectStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetProjectStats")
	})
}
