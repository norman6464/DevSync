package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockMentionStatsRepo は usecase/repository.MentionStatsRepository のモック。
type mockMentionStatsRepo struct{ mock.Mock }

func (m *mockMentionStatsRepo) GetMentionStats(ctx context.Context, userID uint) (*model.MentionStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.MentionStats)
	return s, args.Error(1)
}

func TestGetMentionStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockMentionStatsRepo)
		expected := &model.MentionStats{MentionsReceived: 3, MentionsMade: 2, MentionsThisMonth: 2}
		stats.On("GetMentionStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetMentionStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockMentionStatsRepo)
		uc := usecase.NewGetMentionStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetMentionStats")
	})
}
