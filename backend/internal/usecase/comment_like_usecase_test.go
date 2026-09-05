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

// mockCommentLikeRepo は usecase/repository.CommentLikeRepository のモック。
type mockCommentLikeRepo struct{ mock.Mock }

func (m *mockCommentLikeRepo) Like(ctx context.Context, userID, commentID uint) error {
	return m.Called(ctx, userID, commentID).Error(0)
}

func (m *mockCommentLikeRepo) Unlike(ctx context.Context, userID, commentID uint) error {
	return m.Called(ctx, userID, commentID).Error(0)
}

func (m *mockCommentLikeRepo) HasLiked(ctx context.Context, userID, commentID uint) (bool, error) {
	args := m.Called(ctx, userID, commentID)
	return args.Bool(0), args.Error(1)
}

func (m *mockCommentLikeRepo) CountByCommentID(ctx context.Context, commentID uint) (int64, error) {
	args := m.Called(ctx, commentID)
	return args.Get(0).(int64), args.Error(1)
}

// mockCommentReader は usecase/repository.CommentReader のモック。
type mockCommentReader struct{ mock.Mock }

func (m *mockCommentReader) FindCommentByID(ctx context.Context, id uint) (*model.Comment, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.Comment)
	return c, args.Error(1)
}

func TestLikeCommentUseCase_Execute(t *testing.T) {
	t.Run("他人のコメントに未いいねならいいねできる", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return(&model.Comment{UserID: 2}, nil)
		likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(false, nil)
		likes.On("Like", mock.Anything, uint(1), uint(5)).Return(nil)
		uc := usecase.NewLikeCommentUseCase(likes, reader)

		err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		likes.AssertExpectations(t)
		reader.AssertExpectations(t)
	})

	t.Run("コメントが存在しなければ ErrNotFound", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return((*model.Comment)(nil), errors.New("not found"))
		uc := usecase.NewLikeCommentUseCase(likes, reader)

		err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		likes.AssertNotCalled(t, "Like")
	})

	t.Run("自分のコメントには ErrForbidden", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return(&model.Comment{UserID: 1}, nil)
		uc := usecase.NewLikeCommentUseCase(likes, reader)

		err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		likes.AssertNotCalled(t, "Like")
	})

	t.Run("すでにいいね済みならエラー", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return(&model.Comment{UserID: 2}, nil)
		likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(true, nil)
		uc := usecase.NewLikeCommentUseCase(likes, reader)

		err := uc.Execute(context.Background(), 1, 5)

		assert.Error(t, err)
		likes.AssertNotCalled(t, "Like")
	})
}

func TestUnlikeCommentUseCase_Execute(t *testing.T) {
	t.Run("いいね済みなら解除できる", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return(&model.Comment{UserID: 2}, nil)
		likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(true, nil)
		likes.On("Unlike", mock.Anything, uint(1), uint(5)).Return(nil)
		uc := usecase.NewUnlikeCommentUseCase(likes, reader)

		err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		likes.AssertExpectations(t)
	})

	t.Run("いいねしていなければエラー", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return(&model.Comment{UserID: 2}, nil)
		likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(false, nil)
		uc := usecase.NewUnlikeCommentUseCase(likes, reader)

		err := uc.Execute(context.Background(), 1, 5)

		assert.Error(t, err)
		likes.AssertNotCalled(t, "Unlike")
	})

	t.Run("自分のコメントには ErrForbidden", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return(&model.Comment{UserID: 1}, nil)
		uc := usecase.NewUnlikeCommentUseCase(likes, reader)

		err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		likes.AssertNotCalled(t, "Unlike")
	})
}

func TestGetCommentLikeStatusUseCase_Execute(t *testing.T) {
	t.Run("いいね状態と件数を返す", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return(&model.Comment{UserID: 2}, nil)
		likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(true, nil)
		likes.On("CountByCommentID", mock.Anything, uint(5)).Return(int64(3), nil)
		uc := usecase.NewGetCommentLikeStatusUseCase(likes, reader)

		liked, count, err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		assert.True(t, liked)
		assert.Equal(t, int64(3), count)
		likes.AssertExpectations(t)
	})

	t.Run("コメントが存在しなければ ErrNotFound", func(t *testing.T) {
		likes := new(mockCommentLikeRepo)
		reader := new(mockCommentReader)
		reader.On("FindCommentByID", mock.Anything, uint(5)).Return((*model.Comment)(nil), errors.New("not found"))
		uc := usecase.NewGetCommentLikeStatusUseCase(likes, reader)

		_, _, err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		likes.AssertNotCalled(t, "HasLiked")
	})
}
