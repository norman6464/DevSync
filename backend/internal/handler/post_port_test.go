package handler

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// mockPostPort は usecase/repository.PostRepository のモック。
type mockPostPort struct{ mock.Mock }

func (m *mockPostPort) Create(ctx context.Context, post *model.Post) error {
	return m.Called(ctx, post).Error(0)
}

func (m *mockPostPort) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*model.Post)
	return p, args.Error(1)
}

func (m *mockPostPort) Update(ctx context.Context, post *model.Post) error {
	return m.Called(ctx, post).Error(0)
}

func (m *mockPostPort) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPostPort) FindAll(ctx context.Context, page, limit int) ([]model.Post, error) {
	args := m.Called(ctx, page, limit)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Error(1)
}

func (m *mockPostPort) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostPort) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostPort) FindDraftsByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	args := m.Called(ctx, userID)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Error(1)
}

func (m *mockPostPort) FindScheduledByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	args := m.Called(ctx, userID)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Error(1)
}

func (m *mockPostPort) Timeline(ctx context.Context, userID uint, page, limit int) ([]model.Post, error) {
	args := m.Called(ctx, userID, page, limit)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Error(1)
}

func (m *mockPostPort) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostPort) CountDraftsByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostPort) CountScheduledByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// mockPostLikePort は usecase/repository.PostLikeRepository のモック。
type mockPostLikePort struct{ mock.Mock }

func (m *mockPostLikePort) Like(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}

func (m *mockPostLikePort) Unlike(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}

func (m *mockPostLikePort) HasLiked(ctx context.Context, userID, postID uint) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

// mockFollowerNotifierPort は usecase/repository.FollowerNotifier のモック。
// フォロワー通知は goroutine で実行されるため、テストでは呼び出しを任意（Maybe）として設定する。
type mockFollowerNotifierPort struct{ mock.Mock }

func (m *mockFollowerNotifierPort) FindFollowerIDs(ctx context.Context, userID uint) ([]uint, error) {
	args := m.Called(ctx, userID)
	ids, _ := args.Get(0).([]uint)
	return ids, args.Error(1)
}

func (m *mockFollowerNotifierPort) CreateBatch(ctx context.Context, notifications []*model.Notification) error {
	return m.Called(ctx, notifications).Error(0)
}
