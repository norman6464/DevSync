package handler

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// mockPostCommentPort は usecase/repository.PostCommentRepository のモック。
type mockPostCommentPort struct{ mock.Mock }

func (m *mockPostCommentPort) FindCommentByID(ctx context.Context, id uint) (*model.Comment, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.Comment)
	return c, args.Error(1)
}

func (m *mockPostCommentPort) Create(ctx context.Context, comment *model.Comment) error {
	return m.Called(ctx, comment).Error(0)
}

func (m *mockPostCommentPort) Update(ctx context.Context, comment *model.Comment) error {
	return m.Called(ctx, comment).Error(0)
}

func (m *mockPostCommentPort) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPostCommentPort) ListByPostID(ctx context.Context, postID uint) ([]model.Comment, error) {
	args := m.Called(ctx, postID)
	cs, _ := args.Get(0).([]model.Comment)
	return cs, args.Error(1)
}

func (m *mockPostCommentPort) ListReplies(ctx context.Context, parentID uint) ([]model.Comment, error) {
	args := m.Called(ctx, parentID)
	cs, _ := args.Get(0).([]model.Comment)
	return cs, args.Error(1)
}
