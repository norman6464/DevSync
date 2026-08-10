package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockNotificationStatsRepo は usecase/repository.NotificationStatsRepository のモック。
type mockNotificationStatsRepo struct{ mock.Mock }

func (m *mockNotificationStatsRepo) GetNotificationStats(ctx context.Context, userID uint) (*model.NotificationStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.NotificationStats)
	return s, args.Error(1)
}

func TestGetNotificationStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockNotificationStatsRepo)
		expected := &model.NotificationStats{TotalNotifications: 5, UnreadCount: 2, NotificationsThisMonth: 4}
		stats.On("GetNotificationStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetNotificationStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockNotificationStatsRepo)
		uc := usecase.NewGetNotificationStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetNotificationStats")
	})
}
