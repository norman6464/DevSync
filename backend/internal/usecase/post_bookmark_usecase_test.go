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
	"github.com/stretchr/testify/require"
)

// mockPostBookmarkRepo は usecase/repository.PostBookmarkRepository のモック。
type mockPostBookmarkRepo struct{ mock.Mock }

func (m *mockPostBookmarkRepo) Bookmark(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}

func (m *mockPostBookmarkRepo) Unbookmark(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}

func (m *mockPostBookmarkRepo) HasBookmarked(ctx context.Context, userID, postID uint) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPostBookmarkRepo) FindBookmarkedByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Post, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostBookmarkRepo) CountBookmarkedByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// mockPostAuthorReader は usecase/repository.PostAuthorReader のモック。
type mockPostAuthorReader struct{ mock.Mock }

func (m *mockPostAuthorReader) FindAuthorID(ctx context.Context, postID uint) (uint, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(uint), args.Error(1)
}

// ============================================================
// 追加・解除
// ============================================================

func TestBookmarkPostUseCase(t *testing.T) {
	t.Run("他人の投稿はブックマークできる", func(t *testing.T) {
		bookmarks := new(mockPostBookmarkRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
		bookmarks.On("Bookmark", mock.Anything, uint(1), uint(5)).Return(nil)

		require.NoError(t, usecase.NewBookmarkPostUseCase(bookmarks, authors).Execute(context.Background(), 1, 5))
		bookmarks.AssertExpectations(t)
		authors.AssertExpectations(t)
	})

	t.Run("自分の投稿は 403", func(t *testing.T) {
		bookmarks := new(mockPostBookmarkRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

		err := usecase.NewBookmarkPostUseCase(bookmarks, authors).Execute(context.Background(), 1, 5)
		assert.ErrorIs(t, err, domain.ErrForbidden)
		bookmarks.AssertNotCalled(t, "Bookmark", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("投稿が存在しなければ 404", func(t *testing.T) {
		bookmarks := new(mockPostBookmarkRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(0), nil)

		err := usecase.NewBookmarkPostUseCase(bookmarks, authors).Execute(context.Background(), 1, 5)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		bookmarks.AssertNotCalled(t, "Bookmark", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("投稿者の取得に失敗したら 404", func(t *testing.T) {
		bookmarks := new(mockPostBookmarkRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(0), errors.New("db error"))

		err := usecase.NewBookmarkPostUseCase(bookmarks, authors).Execute(context.Background(), 1, 5)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("保存エラーはそのまま返す", func(t *testing.T) {
		bookmarks := new(mockPostBookmarkRepo)
		authors := new(mockPostAuthorReader)
		saveErr := errors.New("db error")
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
		bookmarks.On("Bookmark", mock.Anything, uint(1), uint(5)).Return(saveErr)

		assert.ErrorIs(t, usecase.NewBookmarkPostUseCase(bookmarks, authors).Execute(context.Background(), 1, 5), saveErr)
	})
}

func TestUnbookmarkPostUseCase(t *testing.T) {
	t.Run("他人の投稿のブックマークは解除できる", func(t *testing.T) {
		bookmarks := new(mockPostBookmarkRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
		bookmarks.On("Unbookmark", mock.Anything, uint(1), uint(5)).Return(nil)

		require.NoError(t, usecase.NewUnbookmarkPostUseCase(bookmarks, authors).Execute(context.Background(), 1, 5))
		bookmarks.AssertExpectations(t)
		authors.AssertExpectations(t)
	})

	t.Run("自分の投稿は 403", func(t *testing.T) {
		bookmarks := new(mockPostBookmarkRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

		err := usecase.NewUnbookmarkPostUseCase(bookmarks, authors).Execute(context.Background(), 1, 5)
		assert.ErrorIs(t, err, domain.ErrForbidden)
		bookmarks.AssertNotCalled(t, "Unbookmark", mock.Anything, mock.Anything, mock.Anything)
	})
}

// ============================================================
// 参照
// ============================================================

func TestBookmarkQueryUseCases(t *testing.T) {
	bookmarks := new(mockPostBookmarkRepo)

	bookmarks.On("HasBookmarked", mock.Anything, uint(1), uint(5)).Return(true, nil)
	has, err := usecase.NewHasBookmarkedPostUseCase(bookmarks).Execute(context.Background(), 1, 5)
	require.NoError(t, err)
	assert.True(t, has)

	bookmarks.On("FindBookmarkedByUserID", mock.Anything, uint(1), 2, 5).
		Return([]model.Post{{ID: 3}}, int64(7), nil)
	posts, total, err := usecase.NewListBookmarkedPostsUseCase(bookmarks).Execute(context.Background(), 1, 2, 5)
	require.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, int64(7), total)

	bookmarks.On("CountBookmarkedByUserID", mock.Anything, uint(1)).Return(int64(7), nil)
	count, err := usecase.NewCountBookmarkedPostsUseCase(bookmarks).Execute(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, int64(7), count)

	bookmarks.AssertExpectations(t)
}

func TestBookmarkQueryUseCases_Errors(t *testing.T) {
	bookmarks := new(mockPostBookmarkRepo)
	dbErr := errors.New("db error")

	bookmarks.On("HasBookmarked", mock.Anything, uint(1), uint(5)).Return(false, dbErr)
	_, err := usecase.NewHasBookmarkedPostUseCase(bookmarks).Execute(context.Background(), 1, 5)
	assert.ErrorIs(t, err, dbErr)

	bookmarks.On("FindBookmarkedByUserID", mock.Anything, uint(1), 1, 20).Return(nil, int64(0), dbErr)
	_, _, err = usecase.NewListBookmarkedPostsUseCase(bookmarks).Execute(context.Background(), 1, 1, 20)
	assert.ErrorIs(t, err, dbErr)

	bookmarks.On("CountBookmarkedByUserID", mock.Anything, uint(1)).Return(int64(0), dbErr)
	_, err = usecase.NewCountBookmarkedPostsUseCase(bookmarks).Execute(context.Background(), 1)
	assert.ErrorIs(t, err, dbErr)

	bookmarks.AssertExpectations(t)
}
