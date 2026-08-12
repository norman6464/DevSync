package handler

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// mockPostBookmarkPort は usecase/repository.PostBookmarkRepository のモック。
type mockPostBookmarkPort struct{ mock.Mock }

func (m *mockPostBookmarkPort) Bookmark(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}

func (m *mockPostBookmarkPort) Unbookmark(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}

func (m *mockPostBookmarkPort) HasBookmarked(ctx context.Context, userID, postID uint) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPostBookmarkPort) FindBookmarkedByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Post, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostBookmarkPort) CountBookmarkedByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
