package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockPostPinRepo は usecase/repository.PostPinRepository のモック。
type mockPostPinRepo struct{ mock.Mock }

func (m *mockPostPinRepo) Pin(ctx context.Context, pin *model.PostPin) error {
	return m.Called(ctx, pin).Error(0)
}
func (m *mockPostPinRepo) Unpin(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}
func (m *mockPostPinRepo) GetByUserID(ctx context.Context, userID uint) ([]model.PostPin, error) {
	args := m.Called(ctx, userID)
	pins, _ := args.Get(0).([]model.PostPin)
	return pins, args.Error(1)
}
func (m *mockPostPinRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockPostPinRepo) IsPinned(ctx context.Context, userID, postID uint) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}
func (m *mockPostPinRepo) UpdateOrder(ctx context.Context, userID uint, postIDs []uint) error {
	return m.Called(ctx, userID, postIDs).Error(0)
}

// mockPostReader は usecase/repository.PostReader のモック。
type mockPostReader struct{ mock.Mock }

func (m *mockPostReader) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*model.Post)
	return p, args.Error(1)
}

func TestPinPostUseCase_Execute(t *testing.T) {
	t.Run("自分の投稿を上限内でピン留めできる", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		posts := new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(&model.Post{UserID: 1}, nil)
		pins.On("CountByUserID", mock.Anything, uint(1)).Return(int64(1), nil)
		pins.On("Pin", mock.Anything, mock.AnythingOfType("*model.PostPin")).Return(nil)
		uc := usecase.NewPinPostUseCase(pins, posts)

		err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		pins.AssertExpectations(t)
	})

	t.Run("存在しない投稿は ErrNotFound", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		posts := new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return((*model.Post)(nil), errors.New("not found"))
		uc := usecase.NewPinPostUseCase(pins, posts)

		err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		pins.AssertNotCalled(t, "Pin")
	})

	t.Run("他人の投稿は Forbidden", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		posts := new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(&model.Post{UserID: 2}, nil)
		uc := usecase.NewPinPostUseCase(pins, posts)

		err := uc.Execute(context.Background(), 1, 5)

		assert.Error(t, err)
		pins.AssertNotCalled(t, "Pin")
	})

	t.Run("上限に達していると BadRequest", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		posts := new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(&model.Post{UserID: 1}, nil)
		pins.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)
		uc := usecase.NewPinPostUseCase(pins, posts)

		err := uc.Execute(context.Background(), 1, 5)

		assert.Error(t, err)
		pins.AssertNotCalled(t, "Pin")
	})
}

func TestUnpinPostUseCase_Execute(t *testing.T) {
	t.Run("ピン留め済みなら解除できる", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		pins.On("IsPinned", mock.Anything, uint(1), uint(5)).Return(true, nil)
		pins.On("Unpin", mock.Anything, uint(1), uint(5)).Return(nil)
		uc := usecase.NewUnpinPostUseCase(pins)

		err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		pins.AssertExpectations(t)
	})

	t.Run("未ピン留めなら NotFound", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		pins.On("IsPinned", mock.Anything, uint(1), uint(5)).Return(false, nil)
		uc := usecase.NewUnpinPostUseCase(pins)

		err := uc.Execute(context.Background(), 1, 5)

		assert.Error(t, err)
		pins.AssertNotCalled(t, "Unpin")
	})
}

func TestListPinnedPostsUseCase_Execute(t *testing.T) {
	pins := new(mockPostPinRepo)
	expected := []model.PostPin{{PostID: 10}, {PostID: 20}}
	pins.On("GetByUserID", mock.Anything, uint(1)).Return(expected, nil)
	uc := usecase.NewListPinnedPostsUseCase(pins)

	got, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	pins.AssertExpectations(t)
}

func TestReorderPinnedPostsUseCase_Execute(t *testing.T) {
	t.Run("自分のピン留め投稿を並べ替えできる", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		pins.On("GetByUserID", mock.Anything, uint(1)).
			Return([]model.PostPin{{PostID: 10}, {PostID: 20}}, nil)
		pins.On("UpdateOrder", mock.Anything, uint(1), []uint{20, 10}).Return(nil)
		uc := usecase.NewReorderPinnedPostsUseCase(pins)

		err := uc.Execute(context.Background(), 1, []uint{20, 10})

		assert.NoError(t, err)
		pins.AssertExpectations(t)
	})

	t.Run("上限を超える並べ替えは BadRequest", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		uc := usecase.NewReorderPinnedPostsUseCase(pins)

		err := uc.Execute(context.Background(), 1, []uint{1, 2, 3, 4})

		assert.Error(t, err)
		pins.AssertNotCalled(t, "UpdateOrder")
	})

	t.Run("自分のピン留めでない投稿を含むと Forbidden", func(t *testing.T) {
		pins := new(mockPostPinRepo)
		pins.On("GetByUserID", mock.Anything, uint(1)).
			Return([]model.PostPin{{PostID: 10}}, nil)
		uc := usecase.NewReorderPinnedPostsUseCase(pins)

		err := uc.Execute(context.Background(), 1, []uint{10, 99})

		assert.Error(t, err)
		pins.AssertNotCalled(t, "UpdateOrder")
	})
}

func TestCountPinnedPostsUseCase_Execute(t *testing.T) {
	pins := new(mockPostPinRepo)
	pins.On("CountByUserID", mock.Anything, uint(1)).Return(int64(2), nil)
	uc := usecase.NewCountPinnedPostsUseCase(pins)

	count, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
	pins.AssertExpectations(t)
}
