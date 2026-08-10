package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockRoadmapStatsRepo は usecase/repository.RoadmapStatsRepository のモック。
type mockRoadmapStatsRepo struct{ mock.Mock }

func (m *mockRoadmapStatsRepo) GetRoadmapStats(ctx context.Context, userID uint) (*model.RoadmapStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.RoadmapStats)
	return s, args.Error(1)
}

func TestGetRoadmapStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockRoadmapStatsRepo)
		expected := &model.RoadmapStats{TotalRoadmaps: 3, ActiveRoadmaps: 2, CompletedRoadmaps: 1, TotalSteps: 12, CompletedSteps: 5}
		stats.On("GetRoadmapStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetRoadmapStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockRoadmapStatsRepo)
		uc := usecase.NewGetRoadmapStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetRoadmapStats")
	})
}
