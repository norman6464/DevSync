package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCommentStatsRepo は usecase/repository.CommentStatsRepository のモック。
type mockCommentStatsRepo struct{ mock.Mock }

func (m *mockCommentStatsRepo) GetCommentStats(ctx context.Context, userID uint) (*model.CommentStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.CommentStats)
	return s, args.Error(1)
}

func TestGetCommentStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockCommentStatsRepo)
		expected := &model.CommentStats{TotalComments: 3, TotalReplies: 2, CommentsReceived: 4, CommentsThisMonth: 5}
		stats.On("GetCommentStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetCommentStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockCommentStatsRepo)
		uc := usecase.NewGetCommentStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetCommentStats")
	})
}
