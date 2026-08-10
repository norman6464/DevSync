package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockStudyCircleStatsRepo は usecase/repository.StudyCircleStatsRepository のモック。
type mockStudyCircleStatsRepo struct{ mock.Mock }

func (m *mockStudyCircleStatsRepo) GetCircleStats(ctx context.Context, circleID uint) (*model.StudyCircleStats, error) {
	args := m.Called(ctx, circleID)
	s, _ := args.Get(0).(*model.StudyCircleStats)
	return s, args.Error(1)
}

func TestGetStudyCircleStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockStudyCircleStatsRepo)
		expected := &model.StudyCircleStats{MemberCount: 3, CheckinCount: 5, TotalSteps: 4, CompletedSteps: 2}
		stats.On("GetCircleStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetStudyCircleStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("circleID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockStudyCircleStatsRepo)
		uc := usecase.NewGetStudyCircleStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circleID")
		stats.AssertNotCalled(t, "GetCircleStats")
	})
}
