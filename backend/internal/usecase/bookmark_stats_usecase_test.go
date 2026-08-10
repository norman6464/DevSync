package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockBookmarkStatsRepo は usecase/repository.BookmarkStatsRepository のモック。
type mockBookmarkStatsRepo struct{ mock.Mock }

func (m *mockBookmarkStatsRepo) GetBookmarkStats(ctx context.Context, userID uint) (*model.BookmarkStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.BookmarkStats)
	return s, args.Error(1)
}

func TestGetBookmarkStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockBookmarkStatsRepo)
		expected := &model.BookmarkStats{TotalBookmarksMade: 4, TotalBookmarksReceived: 2, BookmarksThisMonth: 3}
		stats.On("GetBookmarkStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetBookmarkStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockBookmarkStatsRepo)
		uc := usecase.NewGetBookmarkStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetBookmarkStats")
	})
}
