package handler

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// mockPostReactionPort は usecase/repository.PostReactionRepository のモック。
type mockPostReactionPort struct{ mock.Mock }

func (m *mockPostReactionPort) AddReaction(ctx context.Context, userID, postID uint, emoji string) error {
	return m.Called(ctx, userID, postID, emoji).Error(0)
}
func (m *mockPostReactionPort) RemoveReaction(ctx context.Context, userID, postID uint, emoji string) error {
	return m.Called(ctx, userID, postID, emoji).Error(0)
}
func (m *mockPostReactionPort) GetReactionsByPostID(ctx context.Context, postID uint) ([]model.ReactionCount, error) {
	args := m.Called(ctx, postID)
	r, _ := args.Get(0).([]model.ReactionCount)
	return r, args.Error(1)
}
func (m *mockPostReactionPort) GetUserReactions(ctx context.Context, userID, postID uint) ([]string, error) {
	args := m.Called(ctx, userID, postID)
	r, _ := args.Get(0).([]string)
	return r, args.Error(1)
}
func (m *mockPostReactionPort) GetReactionsBatch(ctx context.Context, postIDs []uint) (map[uint][]model.ReactionCount, error) {
	args := m.Called(ctx, postIDs)
	r, _ := args.Get(0).(map[uint][]model.ReactionCount)
	return r, args.Error(1)
}
func (m *mockPostReactionPort) GetUserReactionsBatch(ctx context.Context, userID uint, postIDs []uint) (map[uint][]string, error) {
	args := m.Called(ctx, userID, postIDs)
	r, _ := args.Get(0).(map[uint][]string)
	return r, args.Error(1)
}

// mockPostAuthorPort は usecase/repository.PostAuthorReader のモック。
type mockPostAuthorPort struct{ mock.Mock }

func (m *mockPostAuthorPort) FindAuthorID(ctx context.Context, postID uint) (uint, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(uint), args.Error(1)
}
