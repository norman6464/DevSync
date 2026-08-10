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

// ---------- NoteStats ----------

type MockNoteStatsService struct{ mock.Mock }

func (m *MockNoteStatsService) GetNoteStats(userID uint) (*model.NoteStats, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.NoteStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestNoteStats_GetStats_Success(t *testing.T) {
	svc := new(MockNoteStatsService)
	h := NewNoteStatsHandler(svc)
	svc.On("GetNoteStats", uint(5)).Return(&model.NoteStats{}, nil)
	r := newRouter(1)
	r.GET("/users/:id/stats/notes", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/notes", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockNoteStatsService)
	h := NewNoteStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/stats/notes", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/stats/notes", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockNoteStatsService)
	h := NewNoteStatsHandler(svc)
	svc.On("GetNoteStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/stats/notes", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/notes", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// BookReviewStats のハンドラーテストは book_review_stats_test.go（DIP 版）へ移設。

// ---------- BookmarkStats ----------

type MockBookmarkStatsService struct{ mock.Mock }

func (m *MockBookmarkStatsService) GetBookmarkStats(userID uint) (*model.BookmarkStats, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.BookmarkStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestBookmarkStats_GetStats_Success(t *testing.T) {
	svc := new(MockBookmarkStatsService)
	h := NewBookmarkStatsHandler(svc)
	svc.On("GetBookmarkStats", uint(5)).Return(&model.BookmarkStats{}, nil)
	r := newRouter(1)
	r.GET("/users/:id/stats/bookmarks", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/bookmarks", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestBookmarkStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockBookmarkStatsService)
	h := NewBookmarkStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/stats/bookmarks", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/stats/bookmarks", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookmarkStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockBookmarkStatsService)
	h := NewBookmarkStatsHandler(svc)
	svc.On("GetBookmarkStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/stats/bookmarks", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/bookmarks", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// CodeSnippetStats のハンドラーテストは code_snippet_stats_test.go（DIP 版）へ移設。

// CommentStats のハンドラーテストは comment_stats_test.go（DIP 版）へ移設。

// FollowStats のハンドラーテストは follow_stats_test.go（DIP 版）へ移設。

// LearningLogStats のハンドラーテストは learning_log_stats_test.go（DIP 版）へ移設。

// LearningResourceStats のハンドラーテストは learning_resource_stats_test.go（DIP 版）へ移設。

// MentionStats のハンドラーテストは mention_stats_test.go（DIP 版）へ移設。

// ---------- MessageStats ----------

type MockMessageStatsService struct{ mock.Mock }

func (m *MockMessageStatsService) GetMessageStats(userID uint) (*model.MessageStats, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.MessageStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestMessageStats_GetStats_Success(t *testing.T) {
	svc := new(MockMessageStatsService)
	h := NewMessageStatsHandler(svc)
	svc.On("GetMessageStats", uint(5)).Return(&model.MessageStats{}, nil)
	r := newRouter(1)
	r.GET("/users/:id/stats/messages", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/messages", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestMessageStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockMessageStatsService)
	h := NewMessageStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/stats/messages", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/stats/messages", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestMessageStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockMessageStatsService)
	h := NewMessageStatsHandler(svc)
	svc.On("GetMessageStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/stats/messages", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/messages", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ---------- NotificationStats ----------

type MockNotificationStatsService struct{ mock.Mock }

func (m *MockNotificationStatsService) GetNotificationStats(userID uint) (*model.NotificationStats, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.NotificationStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestNotificationStats_GetStats_Success(t *testing.T) {
	svc := new(MockNotificationStatsService)
	h := NewNotificationStatsHandler(svc)
	svc.On("GetNotificationStats", uint(5)).Return(&model.NotificationStats{}, nil)
	r := newRouter(1)
	r.GET("/users/:id/stats/notifications", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/notifications", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNotificationStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockNotificationStatsService)
	h := NewNotificationStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/stats/notifications", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/stats/notifications", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNotificationStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockNotificationStatsService)
	h := NewNotificationStatsHandler(svc)
	svc.On("GetNotificationStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/stats/notifications", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/notifications", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ---------- ProjectStats ----------

type MockProjectStatsService struct{ mock.Mock }

func (m *MockProjectStatsService) GetProjectStats(userID uint) (*model.ProjectStats, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.ProjectStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestProjectStats_GetStats_Success(t *testing.T) {
	svc := new(MockProjectStatsService)
	h := NewProjectStatsHandler(svc)
	svc.On("GetProjectStats", uint(5)).Return(&model.ProjectStats{}, nil)
	r := newRouter(1)
	r.GET("/users/:id/stats/projects", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/projects", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestProjectStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockProjectStatsService)
	h := NewProjectStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/stats/projects", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/stats/projects", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestProjectStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockProjectStatsService)
	h := NewProjectStatsHandler(svc)
	svc.On("GetProjectStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/stats/projects", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/projects", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ---------- QAStats ----------

type MockQAStatsService struct{ mock.Mock }

func (m *MockQAStatsService) GetQAStats(userID uint) (*model.QAStats, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.QAStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestQAStats_GetStats_Success(t *testing.T) {
	svc := new(MockQAStatsService)
	h := NewQAStatsHandler(svc)
	svc.On("GetQAStats", uint(5)).Return(&model.QAStats{}, nil)
	r := newRouter(1)
	r.GET("/users/:id/stats/qa", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/qa", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestQAStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockQAStatsService)
	h := NewQAStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/stats/qa", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/stats/qa", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestQAStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockQAStatsService)
	h := NewQAStatsHandler(svc)
	svc.On("GetQAStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/stats/qa", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/qa", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

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

// ---------- RoadmapStats ----------

type MockRoadmapStatsService struct{ mock.Mock }

func (m *MockRoadmapStatsService) GetRoadmapStats(userID uint) (*model.RoadmapStats, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*model.RoadmapStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRoadmapStats_GetStats_Success(t *testing.T) {
	svc := new(MockRoadmapStatsService)
	h := NewRoadmapStatsHandler(svc)
	svc.On("GetRoadmapStats", uint(5)).Return(&model.RoadmapStats{}, nil)
	r := newRouter(1)
	r.GET("/users/:id/stats/roadmaps", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/roadmaps", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestRoadmapStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockRoadmapStatsService)
	h := NewRoadmapStatsHandler(svc)
	r := newRouter(1)
	r.GET("/users/:id/stats/roadmaps", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/stats/roadmaps", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestRoadmapStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockRoadmapStatsService)
	h := NewRoadmapStatsHandler(svc)
	svc.On("GetRoadmapStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/users/:id/stats/roadmaps", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/stats/roadmaps", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ---------- StudyCircleStats ----------

type MockStudyCircleStatsService struct{ mock.Mock }

func (m *MockStudyCircleStatsService) GetCircleStats(circleID uint) (*model.StudyCircleStats, error) {
	args := m.Called(circleID)
	if v := args.Get(0); v != nil {
		return v.(*model.StudyCircleStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestStudyCircleStats_GetStats_Success(t *testing.T) {
	svc := new(MockStudyCircleStatsService)
	h := NewStudyCircleStatsHandler(svc)
	svc.On("GetCircleStats", uint(5)).Return(&model.StudyCircleStats{}, nil)
	r := newRouter(1)
	r.GET("/circles/:id/stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/circles/5/stats", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestStudyCircleStats_GetStats_InvalidID(t *testing.T) {
	svc := new(MockStudyCircleStatsService)
	h := NewStudyCircleStatsHandler(svc)
	r := newRouter(1)
	r.GET("/circles/:id/stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/circles/abc/stats", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestStudyCircleStats_GetStats_ServiceError(t *testing.T) {
	svc := new(MockStudyCircleStatsService)
	h := NewStudyCircleStatsHandler(svc)
	svc.On("GetCircleStats", uint(5)).Return(nil, errors.New("db error"))
	r := newRouter(1)
	r.GET("/circles/:id/stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/circles/5/stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
