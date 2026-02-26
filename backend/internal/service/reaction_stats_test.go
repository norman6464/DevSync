package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newReactionStatsTestService() (*ReactionStatsService, *MockReactionStatsRepository) {
	repo := new(MockReactionStatsRepository)
	svc := NewReactionStatsService(repo)
	return svc, repo
}

func TestReactionStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newReactionStatsTestService()
	expected := &model.ReactionStats{
		TotalReactionsReceived: 120,
		UniqueReactors:         35,
		ReactionsThisMonth:     18,
	}
	repo.On("GetReactionStats", uint(1)).Return(expected, nil)

	result, err := svc.GetReactionStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestReactionStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newReactionStatsTestService()

	_, err := svc.GetReactionStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestReactionStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newReactionStatsTestService()
	repo.On("GetReactionStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetReactionStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// リアクションサマリーテスト
// ============================================================

func TestReactionStatsService_GetSummary_Success(t *testing.T) {
	svc, repo := newReactionStatsTestService()

	emojiCounts := []model.ReactionCount{
		{Emoji: "🔥", Count: 15},
		{Emoji: "👍", Count: 12},
	}
	topPosts := []model.TopReactedPost{
		{ID: 1, Title: "Go入門", ReactionCount: 10},
		{ID: 2, Title: "React入門", ReactionCount: 5},
	}
	repo.On("GetEmojiBreakdown", uint(1)).Return(emojiCounts, nil)
	repo.On("GetTopReactedPosts", uint(1), 5).Return(topPosts, nil)

	result, err := svc.GetReactionSummary(1)
	assert.NoError(t, err)
	assert.Len(t, result.EmojiCounts, 2)
	assert.Equal(t, "🔥", result.EmojiCounts[0].Emoji)
	assert.Len(t, result.TopPosts, 2)
	assert.Equal(t, 27, result.TotalReactions)
	repo.AssertExpectations(t)
}

func TestReactionStatsService_GetSummary_InvalidUserID(t *testing.T) {
	svc, _ := newReactionStatsTestService()

	_, err := svc.GetReactionSummary(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestReactionStatsService_GetSummary_EmojiBreakdownError(t *testing.T) {
	svc, repo := newReactionStatsTestService()

	repo.On("GetEmojiBreakdown", uint(1)).Return([]model.ReactionCount(nil), errors.New("db error"))

	_, err := svc.GetReactionSummary(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestReactionStatsService_GetSummary_TopPostsError(t *testing.T) {
	svc, repo := newReactionStatsTestService()

	emojiCounts := []model.ReactionCount{{Emoji: "👍", Count: 5}}
	repo.On("GetEmojiBreakdown", uint(1)).Return(emojiCounts, nil)
	repo.On("GetTopReactedPosts", uint(1), 5).Return([]model.TopReactedPost(nil), errors.New("db error"))

	_, err := svc.GetReactionSummary(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestReactionStatsService_GetSummary_Empty(t *testing.T) {
	svc, repo := newReactionStatsTestService()

	repo.On("GetEmojiBreakdown", uint(99)).Return([]model.ReactionCount{}, nil)
	repo.On("GetTopReactedPosts", uint(99), 5).Return([]model.TopReactedPost{}, nil)

	result, err := svc.GetReactionSummary(99)
	assert.NoError(t, err)
	assert.Empty(t, result.EmojiCounts)
	assert.Empty(t, result.TopPosts)
	assert.Equal(t, 0, result.TotalReactions)
}

func TestReactionStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newReactionStatsTestService()
	expected := &model.ReactionStats{}
	repo.On("GetReactionStats", uint(99)).Return(expected, nil)

	result, err := svc.GetReactionStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalReactionsReceived)
	assert.Equal(t, int64(0), result.UniqueReactors)
	assert.Equal(t, int64(0), result.ReactionsThisMonth)
	repo.AssertExpectations(t)
}
