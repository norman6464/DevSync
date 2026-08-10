package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockQAStatsRepo は usecase/repository.QAStatsRepository のモック。
type mockQAStatsRepo struct{ mock.Mock }

func (m *mockQAStatsRepo) GetQAStats(ctx context.Context, userID uint) (*model.QAStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.QAStats)
	return s, args.Error(1)
}

func TestGetQAStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockQAStatsRepo)
		expected := &model.QAStats{TotalQuestions: 2, TotalAnswers: 3, BestAnswerCount: 1, TotalVotesReceived: 12}
		stats.On("GetQAStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetQAStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockQAStatsRepo)
		uc := usecase.NewGetQAStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetQAStats")
	})
}
