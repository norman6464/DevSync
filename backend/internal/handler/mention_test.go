package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockMentionPort は usecase/repository.MentionRepository のモック。
type mockMentionPort struct{ mock.Mock }

func (m *mockMentionPort) Create(ctx context.Context, mention *model.Mention) (bool, error) {
	args := m.Called(ctx, mention)
	return args.Bool(0), args.Error(1)
}

func (m *mockMentionPort) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Mention, error) {
	args := m.Called(ctx, userID, page, limit)
	ms, _ := args.Get(0).([]model.Mention)
	return ms, args.Error(1)
}

func (m *mockMentionPort) FindByPostID(ctx context.Context, postID uint) ([]model.Mention, error) {
	args := m.Called(ctx, postID)
	ms, _ := args.Get(0).([]model.Mention)
	return ms, args.Error(1)
}

func (m *mockMentionPort) FindByCommentID(ctx context.Context, commentID uint) ([]model.Mention, error) {
	args := m.Called(ctx, commentID)
	ms, _ := args.Get(0).([]model.Mention)
	return ms, args.Error(1)
}

func (m *mockMentionPort) DeleteByPostID(ctx context.Context, postID uint) error {
	return m.Called(ctx, postID).Error(0)
}

func (m *mockMentionPort) DeleteByCommentID(ctx context.Context, commentID uint) error {
	return m.Called(ctx, commentID).Error(0)
}

// newTestMentionHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestMentionHandler() (*MentionHandler, *mockMentionPort) {
	mentions := new(mockMentionPort)
	h := NewMentionHandler(
		usecase.NewListUserMentionsUseCase(mentions),
		usecase.NewListPostMentionsUseCase(mentions),
	)
	return h, mentions
}

// ---------- GetMyMentions ----------

func TestMentionGetMyMentions_Success(t *testing.T) {
	h, mentions := newTestMentionHandler()
	r := newRouter(1)
	r.GET("/mentions", h.GetMyMentions)

	mentions.On("FindByUserID", mock.Anything, uint(1), 1, 20).
		Return([]model.Mention{{ID: 1, UserID: 1, ActorID: 2}}, nil)

	w := doRequest(r, http.MethodGet, "/mentions", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"actor_id":2`)
	mentions.AssertExpectations(t)
}

func TestMentionGetMyMentions_WithPagination(t *testing.T) {
	h, mentions := newTestMentionHandler()
	r := newRouter(1)
	r.GET("/mentions", h.GetMyMentions)

	mentions.On("FindByUserID", mock.Anything, uint(1), 2, 5).Return([]model.Mention{}, nil)

	w := doRequest(r, http.MethodGet, "/mentions?page=2&limit=5", nil)
	assertStatus(t, w, http.StatusOK)
	mentions.AssertExpectations(t)
}

// メンションが無ければ空配列を返す。
func TestMentionGetMyMentions_Empty(t *testing.T) {
	h, mentions := newTestMentionHandler()
	r := newRouter(1)
	r.GET("/mentions", h.GetMyMentions)

	mentions.On("FindByUserID", mock.Anything, uint(1), 1, 20).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/mentions", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestMentionGetMyMentions_RepositoryError(t *testing.T) {
	h, mentions := newTestMentionHandler()
	r := newRouter(1)
	r.GET("/mentions", h.GetMyMentions)

	mentions.On("FindByUserID", mock.Anything, uint(1), 1, 20).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/mentions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	mentions.AssertExpectations(t)
}

// ---------- GetPostMentions ----------

func TestMentionGetPostMentions_Success(t *testing.T) {
	h, mentions := newTestMentionHandler()
	r := newRouter(1)
	r.GET("/posts/:postId/mentions", h.GetPostMentions)

	mentions.On("FindByPostID", mock.Anything, uint(7)).Return([]model.Mention{{ID: 1, UserID: 3}}, nil)

	w := doRequest(r, http.MethodGet, "/posts/7/mentions", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"user_id":3`)
	mentions.AssertExpectations(t)
}

func TestMentionGetPostMentions_InvalidID(t *testing.T) {
	h, mentions := newTestMentionHandler()
	r := newRouter(1)
	r.GET("/posts/:postId/mentions", h.GetPostMentions)

	w := doRequest(r, http.MethodGet, "/posts/abc/mentions", nil)
	assertStatus(t, w, http.StatusBadRequest)
	mentions.AssertNotCalled(t, "FindByPostID", mock.Anything, mock.Anything)
}

func TestMentionGetPostMentions_RepositoryError(t *testing.T) {
	h, mentions := newTestMentionHandler()
	r := newRouter(1)
	r.GET("/posts/:postId/mentions", h.GetPostMentions)

	mentions.On("FindByPostID", mock.Anything, uint(7)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/7/mentions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	mentions.AssertExpectations(t)
}
