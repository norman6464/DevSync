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

// mockPostViewRepo は usecase/repository.PostViewRepository のモック。
type mockPostViewRepo struct{ mock.Mock }

func (m *mockPostViewRepo) RecordViewIfAbsent(ctx context.Context, view *model.PostView) (bool, error) {
	args := m.Called(ctx, view)
	return args.Bool(0), args.Error(1)
}

func (m *mockPostViewRepo) GetViewCount(ctx context.Context, postID uint) (int64, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostViewRepo) GetMostViewed(ctx context.Context, limit int) ([]model.ViewCount, error) {
	args := m.Called(ctx, limit)
	vc, _ := args.Get(0).([]model.ViewCount)
	return vc, args.Error(1)
}

func TestRecordPostViewUseCase_Execute(t *testing.T) {
	t.Run("未閲覧なら記録する", func(t *testing.T) {
		views := new(mockPostViewRepo)
		views.On("RecordViewIfAbsent", mock.Anything, mock.AnythingOfType("*model.PostView")).Return(true, nil)
		uc := usecase.NewRecordPostViewUseCase(views)

		err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		views.AssertExpectations(t)
	})

	t.Run("閲覧済みでも成功する（記録はスキップ）", func(t *testing.T) {
		views := new(mockPostViewRepo)
		views.On("RecordViewIfAbsent", mock.Anything, mock.AnythingOfType("*model.PostView")).Return(false, nil)
		uc := usecase.NewRecordPostViewUseCase(views)

		err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		views.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400", func(t *testing.T) {
		views := new(mockPostViewRepo)
		uc := usecase.NewRecordPostViewUseCase(views)

		err := uc.Execute(context.Background(), 0, 5)

		assert.Error(t, err)
		views.AssertNotCalled(t, "RecordViewIfAbsent")
	})

	t.Run("postID が 0 は 400", func(t *testing.T) {
		views := new(mockPostViewRepo)
		uc := usecase.NewRecordPostViewUseCase(views)

		err := uc.Execute(context.Background(), 1, 0)

		assert.Error(t, err)
		views.AssertNotCalled(t, "RecordViewIfAbsent")
	})

	t.Run("記録のエラーは伝播する", func(t *testing.T) {
		views := new(mockPostViewRepo)
		views.On("RecordViewIfAbsent", mock.Anything, mock.AnythingOfType("*model.PostView")).Return(false, errors.New("db error"))
		uc := usecase.NewRecordPostViewUseCase(views)

		err := uc.Execute(context.Background(), 1, 5)

		assert.Error(t, err)
	})
}

func TestGetPostViewCountUseCase_Execute(t *testing.T) {
	t.Run("閲覧数を返す", func(t *testing.T) {
		views := new(mockPostViewRepo)
		views.On("GetViewCount", mock.Anything, uint(5)).Return(int64(42), nil)
		uc := usecase.NewGetPostViewCountUseCase(views)

		count, err := uc.Execute(context.Background(), 5)

		assert.NoError(t, err)
		assert.Equal(t, int64(42), count)
		views.AssertExpectations(t)
	})

	t.Run("postID が 0 は 400", func(t *testing.T) {
		views := new(mockPostViewRepo)
		uc := usecase.NewGetPostViewCountUseCase(views)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		views.AssertNotCalled(t, "GetViewCount")
	})
}

func TestGetMostViewedPostsUseCase_Execute(t *testing.T) {
	t.Run("ランキングを返す", func(t *testing.T) {
		views := new(mockPostViewRepo)
		expected := []model.ViewCount{{PostID: 1, Count: 100}, {PostID: 2, Count: 50}}
		views.On("GetMostViewed", mock.Anything, 20).Return(expected, nil)
		uc := usecase.NewGetMostViewedPostsUseCase(views)

		got, err := uc.Execute(context.Background(), 20)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		views.AssertExpectations(t)
	})

	t.Run("limit が範囲外は 400", func(t *testing.T) {
		for _, limit := range []int{0, 101} {
			views := new(mockPostViewRepo)
			uc := usecase.NewGetMostViewedPostsUseCase(views)

			_, err := uc.Execute(context.Background(), limit)

			assert.Error(t, err)
			views.AssertNotCalled(t, "GetMostViewed")
		}
	})
}
