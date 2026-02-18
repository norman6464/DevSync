package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockMentionService は MentionServiceInterface のモック実装。
type MockMentionService struct{ mock.Mock }

func (m *MockMentionService) ProcessMentions(actorID uint, text string, postID *uint, commentID *uint) error {
	return m.Called(actorID, text, postID, commentID).Error(0)
}
func (m *MockMentionService) GetMentionsByUserID(userID uint, page, limit int) ([]model.Mention, error) {
	args := m.Called(userID, page, limit)
	if v := args.Get(0); v != nil {
		return v.([]model.Mention), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockMentionService) GetMentionsByPostID(postID uint) ([]model.Mention, error) {
	args := m.Called(postID)
	if v := args.Get(0); v != nil {
		return v.([]model.Mention), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockMentionService) DeleteMentionsByPostID(postID uint) error {
	return m.Called(postID).Error(0)
}
func (m *MockMentionService) DeleteMentionsByCommentID(commentID uint) error {
	return m.Called(commentID).Error(0)
}

func setupMentionHandler() (*MentionHandler, *MockMentionService) {
	svc := new(MockMentionService)
	h := NewMentionHandler(svc)
	return h, svc
}

// ============================================================
// GetMyMentions テスト
// ============================================================

func TestMentionGetMyMentions_Success(t *testing.T) {
	h, svc := setupMentionHandler()
	mentions := []model.Mention{
		{UserID: 1, ActorID: 2},
	}
	svc.On("GetMentionsByUserID", uint(1), 1, 20).Return(mentions, nil)

	r := newRouter(1)
	r.GET("/mentions", h.GetMyMentions)

	w := doRequest(r, http.MethodGet, "/mentions", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestMentionGetMyMentions_Empty(t *testing.T) {
	h, svc := setupMentionHandler()
	svc.On("GetMentionsByUserID", uint(1), 1, 20).Return([]model.Mention{}, nil)

	r := newRouter(1)
	r.GET("/mentions", h.GetMyMentions)

	w := doRequest(r, http.MethodGet, "/mentions", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestMentionGetMyMentions_ServiceError(t *testing.T) {
	h, svc := setupMentionHandler()
	svc.On("GetMentionsByUserID", uint(1), 1, 20).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/mentions", h.GetMyMentions)

	w := doRequest(r, http.MethodGet, "/mentions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetPostMentions テスト
// ============================================================

func TestMentionGetPostMentions_Success(t *testing.T) {
	h, svc := setupMentionHandler()
	mentions := []model.Mention{
		{UserID: 2, ActorID: 1},
	}
	svc.On("GetMentionsByPostID", uint(5)).Return(mentions, nil)

	r := newRouter(1)
	r.GET("/posts/:postId/mentions", h.GetPostMentions)

	w := doRequest(r, http.MethodGet, "/posts/5/mentions", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestMentionGetPostMentions_InvalidID(t *testing.T) {
	h, _ := setupMentionHandler()

	r := newRouter(1)
	r.GET("/posts/:postId/mentions", h.GetPostMentions)

	w := doRequest(r, http.MethodGet, "/posts/abc/mentions", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestMentionGetPostMentions_ServiceError(t *testing.T) {
	h, svc := setupMentionHandler()
	svc.On("GetMentionsByPostID", uint(5)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/posts/:postId/mentions", h.GetPostMentions)

	w := doRequest(r, http.MethodGet, "/posts/5/mentions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
