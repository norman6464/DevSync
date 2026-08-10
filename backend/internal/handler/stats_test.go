package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// 汎用Stats モック & テストヘルパー
// ============================================================

// statsTestCase は全statsハンドラー共通のテストケース定義。
type statsTestCase struct {
	name       string
	setupRoute func(h interface{}) http.Handler
	method     string
	path       string
	setupMock  func(m mock.Mock)
	wantStatus int
}

// PostStats のハンドラーテストは post_stats_test.go（DIP 版）へ移設。

// NoteStats のハンドラーテストは note_stats_test.go（DIP 版）へ移設。

// BookReviewStats のハンドラーテストは book_review_stats_test.go（DIP 版）へ移設。

// BookmarkStats のハンドラーテストは bookmark_stats_test.go（DIP 版）へ移設。

// CodeSnippetStats のハンドラーテストは code_snippet_stats_test.go（DIP 版）へ移設。

// CommentStats のハンドラーテストは comment_stats_test.go（DIP 版）へ移設。

// FollowStats のハンドラーテストは follow_stats_test.go（DIP 版）へ移設。

// LearningLogStats のハンドラーテストは learning_log_stats_test.go（DIP 版）へ移設。

// LearningResourceStats のハンドラーテストは learning_resource_stats_test.go（DIP 版）へ移設。

// MentionStats のハンドラーテストは mention_stats_test.go（DIP 版）へ移設。

// MessageStats のハンドラーテストは message_stats_test.go（DIP 版）へ移設。

// NotificationStats のハンドラーテストは notification_stats_test.go（DIP 版）へ移設。

// ProjectStats のハンドラーテストは project_stats_test.go（DIP 版）へ移設。

// QAStats のハンドラーテストは qa_stats_test.go（DIP 版）へ移設。

// ---------- ReactionStats ----------

type MockReactionStatsService struct{ mock.Mock }

func (m *MockReactionStatsService) GetReactionStats(userID uint) (*model.ReactionStats, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.ReactionStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReactionStatsService) GetReactionSummary(userID uint) (*model.ReactionSummary, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.ReactionSummary), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestReactionStats_GetStats_Success(t *testing.T) {
	svc := new(MockReactionStatsService)
	h := NewReactionStatsHandler(svc)
	svc.On("GetReactionStats", uint(5)).Return(&model.ReactionStats{}, nil)
	r := newRouter(1)
	r.GET("/users/:id/stats/reactions", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/reactions", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestReactionStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockReactionStatsService)
	h := NewReactionStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/stats/reactions", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/stats/reactions", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestReactionStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockReactionStatsService)
	h := NewReactionStatsHandler(svc)
	svc.On("GetReactionStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/stats/reactions", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/reactions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestReactionStats_GetSummary_Success(t *testing.T) {
	svc := new(MockReactionStatsService)
	h := NewReactionStatsHandler(svc)
	svc.On("GetReactionSummary", uint(5)).Return(&model.ReactionSummary{
		TotalReactions: 10,
		EmojiCounts:    []model.ReactionCount{{Emoji: "👍", Count: 5}},
		TopPosts:       []model.TopReactedPost{{ID: 1, Title: "Test", ReactionCount: 5}},
	}, nil)
	r := newRouter(1)
	r.GET("/users/:id/reaction-summary", h.GetSummary)
	w := doRequest(r, http.MethodGet, "/users/5/reaction-summary", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestReactionStats_GetSummary_InvalidID(t *testing.T) {
	svc := new(MockReactionStatsService)
	h := NewReactionStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/reaction-summary", h.GetSummary)
	w := doRequest(r, http.MethodGet, "/users/abc/reaction-summary", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestReactionStats_GetSummary_ServiceError(t *testing.T) {
	svc := new(MockReactionStatsService)
	h := NewReactionStatsHandler(svc)
	svc.On("GetReactionSummary", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/reaction-summary", h.GetSummary)
	w := doRequest(r, http.MethodGet, "/users/5/reaction-summary", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// RoadmapStats のハンドラーテストは roadmap_stats_test.go（DIP 版）へ移設。

// StudyCircleStats のハンドラーテストは study_circle_stats_test.go（DIP 版）へ移設。
