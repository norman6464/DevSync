package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockBookReviewStatsRepo は usecase/repository.BookReviewStatsRepository のモック。
type mockBookReviewStatsRepo struct{ mock.Mock }

func (m *mockBookReviewStatsRepo) GetBookReviewStats(ctx context.Context, userID uint) (*model.BookReviewStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.BookReviewStats)
	return s, args.Error(1)
}

func TestGetBookReviewStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockBookReviewStatsRepo)
		expected := &model.BookReviewStats{TotalReviews: 3, AverageRating: 4, MaxRating: 5, MinRating: 3, FiveStarCount: 1}
		stats.On("GetBookReviewStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetBookReviewStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockBookReviewStatsRepo)
		uc := usecase.NewGetBookReviewStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetBookReviewStats")
	})
}
