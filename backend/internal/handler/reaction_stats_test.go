package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockReactionStatsRepo は usecase/repository.ReactionStatsRepository のモック（ctx 付き）。
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

func setupReactionStatsHandler() (*ReactionStatsHandler, *mockReactionStatsRepo) {
	repo := new(mockReactionStatsRepo)
	h := NewReactionStatsHandler(
		usecase.NewGetReactionStatsUseCase(repo),
		usecase.NewGetReactionSummaryUseCase(repo),
	)
	return h, repo
}

func TestReactionStats_GetStats_Success(t *testing.T) {
	h, repo := setupReactionStatsHandler()
	repo.On("GetReactionStats", mock.Anything, uint(5)).Return(
		&model.ReactionStats{TotalReactionsReceived: 7, UniqueReactors: 3, ReactionsThisMonth: 2}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/reactions", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/reactions", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_reactions_received"].(float64) != 7 {
		t.Errorf("expected total_reactions_received 7, got %v", data["total_reactions_received"])
	}
	repo.AssertExpectations(t)
}

func TestReactionStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupReactionStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/reactions", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/reactions", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestReactionStats_GetStats_ServiceError(t *testing.T) {
	h, repo := setupReactionStatsHandler()
	repo.On("GetReactionStats", mock.Anything, uint(5)).Return((*model.ReactionStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/reactions", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/reactions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestReactionStats_GetSummary_Success(t *testing.T) {
	h, repo := setupReactionStatsHandler()
	repo.On("GetEmojiBreakdown", mock.Anything, uint(5)).Return(
		[]model.ReactionCount{{Emoji: "👍", Count: 5}, {Emoji: "🎉", Count: 3}}, nil,
	)
	repo.On("GetTopReactedPosts", mock.Anything, uint(5), 5).Return(
		[]model.TopReactedPost{{ID: 1, Title: "Test", ReactionCount: 5}}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/reaction-summary", h.GetSummary)

	w := doRequest(r, http.MethodGet, "/users/5/reaction-summary", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_reactions"].(float64) != 8 {
		t.Errorf("expected total_reactions 8, got %v", data["total_reactions"])
	}
	repo.AssertExpectations(t)
}

func TestReactionStats_GetSummary_InvalidID(t *testing.T) {
	h, _ := setupReactionStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/reaction-summary", h.GetSummary)

	w := doRequest(r, http.MethodGet, "/users/abc/reaction-summary", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestReactionStats_GetSummary_ServiceError(t *testing.T) {
	h, repo := setupReactionStatsHandler()
	repo.On("GetEmojiBreakdown", mock.Anything, uint(5)).Return([]model.ReactionCount(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/reaction-summary", h.GetSummary)

	w := doRequest(r, http.MethodGet, "/users/5/reaction-summary", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}
