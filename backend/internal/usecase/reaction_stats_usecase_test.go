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

// mockReactionStatsRepo は usecase/repository.ReactionStatsRepository のモック。
type mockReactionStatsRepo struct{ mock.Mock }

func (m *mockReactionStatsRepo) GetReactionStats(ctx context.Context, userID uint) (*model.ReactionStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ReactionStats)
	return s, args.Error(1)
}

func (m *mockReactionStatsRepo) GetEmojiBreakdown(ctx context.Context, userID uint) ([]model.ReactionCount, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]model.ReactionCount)
	return c, args.Error(1)
}

func (m *mockReactionStatsRepo) GetTopReactedPosts(ctx context.Context, userID uint, limit int) ([]model.TopReactedPost, error) {
	args := m.Called(ctx, userID, limit)
	p, _ := args.Get(0).([]model.TopReactedPost)
	return p, args.Error(1)
}

func TestGetReactionStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		repo := new(mockReactionStatsRepo)
		expected := &model.ReactionStats{TotalReactionsReceived: 7, UniqueReactors: 3, ReactionsThisMonth: 2}
		repo.On("GetReactionStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetReactionStatsUseCase(repo)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		repo := new(mockReactionStatsRepo)
		uc := usecase.NewGetReactionStatsUseCase(repo)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "GetReactionStats")
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockReactionStatsRepo)
		repo.On("GetReactionStats", mock.Anything, uint(10)).Return((*model.ReactionStats)(nil), errors.New("db error"))
		uc := usecase.NewGetReactionStatsUseCase(repo)

		_, err := uc.Execute(context.Background(), 10)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetReactionSummaryUseCase_Execute(t *testing.T) {
	t.Run("絵文字別集計の合計を算出して返す", func(t *testing.T) {
		repo := new(mockReactionStatsRepo)
		emojis := []model.ReactionCount{{Emoji: "👍", Count: 5}, {Emoji: "🎉", Count: 3}}
		tops := []model.TopReactedPost{{ID: 1, Title: "t1", ReactionCount: 5}}
		repo.On("GetEmojiBreakdown", mock.Anything, uint(10)).Return(emojis, nil)
		repo.On("GetTopReactedPosts", mock.Anything, uint(10), 5).Return(tops, nil)
		uc := usecase.NewGetReactionSummaryUseCase(repo)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, emojis, got.EmojiCounts)
		assert.Equal(t, tops, got.TopPosts)
		assert.Equal(t, 8, got.TotalReactions)
		repo.AssertExpectations(t)
	})

	t.Run("リアクションが無い場合は合計 0", func(t *testing.T) {
		repo := new(mockReactionStatsRepo)
		repo.On("GetEmojiBreakdown", mock.Anything, uint(10)).Return([]model.ReactionCount{}, nil)
		repo.On("GetTopReactedPosts", mock.Anything, uint(10), 5).Return([]model.TopReactedPost{}, nil)
		uc := usecase.NewGetReactionSummaryUseCase(repo)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, 0, got.TotalReactions)
		repo.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		repo := new(mockReactionStatsRepo)
		uc := usecase.NewGetReactionSummaryUseCase(repo)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "GetEmojiBreakdown")
		repo.AssertNotCalled(t, "GetTopReactedPosts")
	})

	t.Run("絵文字別集計が失敗したらトップ投稿を取りに行かない", func(t *testing.T) {
		repo := new(mockReactionStatsRepo)
		repo.On("GetEmojiBreakdown", mock.Anything, uint(10)).Return([]model.ReactionCount(nil), errors.New("db error"))
		uc := usecase.NewGetReactionSummaryUseCase(repo)

		_, err := uc.Execute(context.Background(), 10)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "GetTopReactedPosts")
		repo.AssertExpectations(t)
	})

	t.Run("トップ投稿の取得エラーを返す", func(t *testing.T) {
		repo := new(mockReactionStatsRepo)
		repo.On("GetEmojiBreakdown", mock.Anything, uint(10)).Return([]model.ReactionCount{}, nil)
		repo.On("GetTopReactedPosts", mock.Anything, uint(10), 5).Return([]model.TopReactedPost(nil), errors.New("db error"))
		uc := usecase.NewGetReactionSummaryUseCase(repo)

		_, err := uc.Execute(context.Background(), 10)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
