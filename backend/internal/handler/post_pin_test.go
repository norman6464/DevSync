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

// mockPostPinRepo は usecase/repository.PostPinRepository のモック（ctx 付き）。
type mockPostPinRepo struct{ mock.Mock }

func (m *mockPostPinRepo) Pin(ctx context.Context, pin *model.PostPin) error {
	return m.Called(ctx, pin).Error(0)
}
func (m *mockPostPinRepo) Unpin(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}
func (m *mockPostPinRepo) GetByUserID(ctx context.Context, userID uint) ([]model.PostPin, error) {
	args := m.Called(ctx, userID)
	pins, _ := args.Get(0).([]model.PostPin)
	return pins, args.Error(1)
}
func (m *mockPostPinRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockPostPinRepo) IsPinned(ctx context.Context, userID, postID uint) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}
func (m *mockPostPinRepo) UpdateOrder(ctx context.Context, userID uint, postIDs []uint) error {
	return m.Called(ctx, userID, postIDs).Error(0)
}

// mockPostReader は usecase/repository.PostReader のモック（ctx 付き）。
type mockPostReader struct{ mock.Mock }

func (m *mockPostReader) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*model.Post)
	return p, args.Error(1)
}

// setupPostPinHandler は本物の usecase + port モックで PostPinHandler を組む。
func setupPostPinHandler() (*PostPinHandler, *mockPostPinRepo, *mockPostReader) {
	pins := new(mockPostPinRepo)
	posts := new(mockPostReader)
	h := NewPostPinHandler(
		usecase.NewPinPostUseCase(pins, posts),
		usecase.NewUnpinPostUseCase(pins),
		usecase.NewListPinnedPostsUseCase(pins),
		usecase.NewCountPinnedPostsUseCase(pins),
		usecase.NewReorderPinnedPostsUseCase(pins),
	)
	return h, pins, posts
}

// ---------- Pin ----------

func TestPostPin_Success(t *testing.T) {
	h, pins, posts := setupPostPinHandler()
	posts.On("FindByID", mock.Anything, uint(5)).Return(&model.Post{UserID: 1}, nil)
	pins.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), nil)
	pins.On("Pin", mock.Anything, mock.AnythingOfType("*model.PostPin")).Return(nil)

	r := newRouter(1)
	r.POST("/posts/:postId/pin", h.Pin)
	w := doRequest(r, http.MethodPost, "/posts/5/pin", nil)
	assertStatus(t, w, http.StatusOK)
	pins.AssertExpectations(t)
}

func TestPostPin_InvalidID(t *testing.T) {
	h, _, _ := setupPostPinHandler()
	r := newRouter(1)
	r.POST("/posts/:postId/pin", h.Pin)
	w := doRequest(r, http.MethodPost, "/posts/abc/pin", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostPin_RepositoryError(t *testing.T) {
	h, pins, posts := setupPostPinHandler()
	posts.On("FindByID", mock.Anything, uint(5)).Return(&model.Post{UserID: 1}, nil)
	pins.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.POST("/posts/:postId/pin", h.Pin)
	w := doRequest(r, http.MethodPost, "/posts/5/pin", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- Unpin ----------

func TestPostUnpin_Success(t *testing.T) {
	h, pins, _ := setupPostPinHandler()
	pins.On("IsPinned", mock.Anything, uint(1), uint(5)).Return(true, nil)
	pins.On("Unpin", mock.Anything, uint(1), uint(5)).Return(nil)

	r := newRouter(1)
	r.DELETE("/posts/:postId/pin", h.Unpin)
	w := doRequest(r, http.MethodDelete, "/posts/5/pin", nil)
	assertStatus(t, w, http.StatusOK)
	pins.AssertExpectations(t)
}

func TestPostUnpin_InvalidID(t *testing.T) {
	h, _, _ := setupPostPinHandler()
	r := newRouter(1)
	r.DELETE("/posts/:postId/pin", h.Unpin)
	w := doRequest(r, http.MethodDelete, "/posts/abc/pin", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetByUserID ----------

func TestPostPinGetByUserID_Success(t *testing.T) {
	h, pins, _ := setupPostPinHandler()
	pins.On("GetByUserID", mock.Anything, uint(1)).Return([]model.PostPin{
		{UserID: 1, PostID: 10, PinOrder: 1},
		{UserID: 1, PostID: 20, PinOrder: 2},
	}, nil)

	r := newRouter(1)
	r.GET("/users/:userId/pins", h.GetByUserID)
	w := doRequest(r, http.MethodGet, "/users/1/pins", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.NotNil(t, body["pins"])
	pins.AssertExpectations(t)
}

func TestPostPinGetByUserID_InvalidID(t *testing.T) {
	h, _, _ := setupPostPinHandler()
	r := newRouter(1)
	r.GET("/users/:userId/pins", h.GetByUserID)
	w := doRequest(r, http.MethodGet, "/users/abc/pins", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostPinGetByUserID_RepositoryError(t *testing.T) {
	h, pins, _ := setupPostPinHandler()
	pins.On("GetByUserID", mock.Anything, uint(1)).Return([]model.PostPin(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/pins", h.GetByUserID)
	w := doRequest(r, http.MethodGet, "/users/1/pins", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	pins.AssertExpectations(t)
}

// ---------- Reorder ----------

func TestPostPinReorder_Success(t *testing.T) {
	h, pins, _ := setupPostPinHandler()
	pins.On("GetByUserID", mock.Anything, uint(1)).Return([]model.PostPin{
		{PostID: 10}, {PostID: 20}, {PostID: 30},
	}, nil)
	pins.On("UpdateOrder", mock.Anything, uint(1), []uint{10, 20, 30}).Return(nil)

	r := newRouter(1)
	r.PUT("/pins/reorder", h.Reorder)
	w := doRequest(r, http.MethodPut, "/pins/reorder", map[string]interface{}{
		"post_ids": []uint{10, 20, 30},
	})
	assertStatus(t, w, http.StatusOK)
	pins.AssertExpectations(t)
}

func TestPostPinReorder_InvalidBody(t *testing.T) {
	h, _, _ := setupPostPinHandler()
	r := newRouter(1)
	r.PUT("/pins/reorder", h.Reorder)
	w := doRequestRaw(r, http.MethodPut, "/pins/reorder", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}
