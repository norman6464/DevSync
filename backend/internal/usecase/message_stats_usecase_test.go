package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockMessageStatsRepo は usecase/repository.MessageStatsRepository のモック。
type mockMessageStatsRepo struct{ mock.Mock }

func (m *mockMessageStatsRepo) GetMessageStats(ctx context.Context, userID uint) (*model.MessageStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.MessageStats)
	return s, args.Error(1)
}

func TestGetMessageStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockMessageStatsRepo)
		expected := &model.MessageStats{TotalSent: 5, TotalReceived: 3, ConversationsCount: 2, MessagesThisMonth: 4}
		stats.On("GetMessageStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetMessageStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockMessageStatsRepo)
		uc := usecase.NewGetMessageStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetMessageStats")
	})
}
