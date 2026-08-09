package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockUserActivityRepo は usecase/repository.UserActivityRepository のモック。
type mockUserActivityRepo struct{ mock.Mock }

func (m *mockUserActivityRepo) FindByUserID(ctx context.Context, userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error) {
	args := m.Called(ctx, userID, activityType, limit, offset)
	acts, _ := args.Get(0).([]model.UserActivity)
	return acts, args.Get(1).(int64), args.Error(2)
}

func TestGetActivityTimelineUseCase_Execute(t *testing.T) {
	t.Run("種別フィルタなしで一覧を返す", func(t *testing.T) {
		activities := new(mockUserActivityRepo)
		expected := []model.UserActivity{
			{ID: 1, UserID: 10, ActivityType: model.ActivityPostCreated},
			{ID: 2, UserID: 10, ActivityType: model.ActivityCommentCreated},
		}
		activities.On("FindByUserID", mock.Anything, uint(10), "", 20, 0).Return(expected, int64(2), nil)
		uc := usecase.NewGetActivityTimelineUseCase(activities)

		list, total, err := uc.Execute(context.Background(), 10, "", 20, 0)

		assert.NoError(t, err)
		assert.Equal(t, expected, list)
		assert.Equal(t, int64(2), total)
		activities.AssertExpectations(t)
	})

	t.Run("種別フィルタで絞り込む", func(t *testing.T) {
		activities := new(mockUserActivityRepo)
		activities.On("FindByUserID", mock.Anything, uint(10), "post_created", 20, 0).
			Return([]model.UserActivity{{ID: 1}}, int64(1), nil)
		uc := usecase.NewGetActivityTimelineUseCase(activities)

		list, total, err := uc.Execute(context.Background(), 10, "post_created", 20, 0)

		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, int64(1), total)
		activities.AssertExpectations(t)
	})

	t.Run("空でも成功する", func(t *testing.T) {
		activities := new(mockUserActivityRepo)
		activities.On("FindByUserID", mock.Anything, uint(10), "", 20, 0).
			Return([]model.UserActivity{}, int64(0), nil)
		uc := usecase.NewGetActivityTimelineUseCase(activities)

		list, total, err := uc.Execute(context.Background(), 10, "", 20, 0)

		assert.NoError(t, err)
		assert.Empty(t, list)
		assert.Equal(t, int64(0), total)
		activities.AssertExpectations(t)
	})

	t.Run("リポジトリエラーは伝播する", func(t *testing.T) {
		activities := new(mockUserActivityRepo)
		activities.On("FindByUserID", mock.Anything, uint(10), "", 20, 0).
			Return([]model.UserActivity(nil), int64(0), errors.New("db error"))
		uc := usecase.NewGetActivityTimelineUseCase(activities)

		_, _, err := uc.Execute(context.Background(), 10, "", 20, 0)

		assert.Error(t, err)
		activities.AssertExpectations(t)
	})
}
