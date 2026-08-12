package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockAtCoderFetcher は usecase/repository.AtCoderRatingFetcher のモック。
type mockAtCoderFetcher struct{ mock.Mock }

func (m *mockAtCoderFetcher) FetchRatingHistory(ctx context.Context, username string) ([]model.AtCoderRatingEntry, error) {
	args := m.Called(ctx, username)
	h, _ := args.Get(0).([]model.AtCoderRatingEntry)
	return h, args.Error(1)
}

func (m *mockAtCoderFetcher) UserExists(ctx context.Context, username string) bool {
	return m.Called(ctx, username).Bool(0)
}

// newTestAtCoderHandler は本物の usecase に port モックを注入したハンドラーを生成する。
// ユーザーの読み書きは user スライスと同じ port モックを共用する。
func newTestAtCoderHandler() (*AtCoderHandler, *mockUserPort, *mockAtCoderFetcher) {
	users := new(mockUserPort)
	ratings := new(mockAtCoderFetcher)
	h := NewAtCoderHandler(
		usecase.NewGetAtCoderRatingUseCase(ratings),
		usecase.NewConnectAtCoderUseCase(users, ratings),
		usecase.NewDisconnectAtCoderUseCase(users),
	)
	return h, users, ratings
}

// ---------- GetRating ----------

func TestAtCoder_GetRating_Success(t *testing.T) {
	h, _, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.GET("/atcoder/:username/rating", h.GetRating)

	ratings.On("FetchRatingHistory", mock.Anything, "tourist").Return([]model.AtCoderRatingEntry{
		{NewRating: 2000}, {NewRating: 3000},
	}, nil)

	w := doRequest(r, http.MethodGet, "/atcoder/tourist/rating", nil)
	assertStatus(t, w, http.StatusOK)
	// 最新（末尾）のレーティングから色とランクを決める。
	assert.JSONEq(t, `{"username":"tourist","rating":3000,"color":"red","rank":"赤"}`, w.Body.String())
	ratings.AssertExpectations(t)
}

// 履歴が空ならレーティング 0（灰）として返す。
func TestAtCoder_GetRating_EmptyHistory(t *testing.T) {
	h, _, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.GET("/atcoder/:username/rating", h.GetRating)

	ratings.On("FetchRatingHistory", mock.Anything, "newbie").Return([]model.AtCoderRatingEntry{}, nil)

	w := doRequest(r, http.MethodGet, "/atcoder/newbie/rating", nil)
	assertStatus(t, w, http.StatusOK)
	assert.JSONEq(t, `{"username":"newbie","rating":0,"color":"gray","rank":"灰"}`, w.Body.String())
	ratings.AssertExpectations(t)
}

// 取得に失敗した場合は 400 でエラー文字列をそのまま返す（移行前と同じ）。
func TestAtCoder_GetRating_FetchError(t *testing.T) {
	h, _, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.GET("/atcoder/:username/rating", h.GetRating)

	ratings.On("FetchRatingHistory", mock.Anything, "ghost").
		Return(nil, domain.NewError(domain.ErrCodeNotFound, "AtCoderユーザーが見つかりません: ghost", nil))

	w := doRequest(r, http.MethodGet, "/atcoder/ghost/rating", nil)
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "AtCoderユーザーが見つかりません: ghost")
	ratings.AssertExpectations(t)
}

// ユーザー名の形式が不正なら外部 API を叩かない。
func TestAtCoder_GetRating_InvalidUsername(t *testing.T) {
	h, _, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.GET("/atcoder/:username/rating", h.GetRating)

	w := doRequest(r, http.MethodGet, "/atcoder/bad!user/rating", nil)
	assertStatus(t, w, http.StatusBadRequest)
	ratings.AssertNotCalled(t, "FetchRatingHistory", mock.Anything, mock.Anything)
}

// ---------- Connect ----------

func TestAtCoder_Connect_Success(t *testing.T) {
	h, users, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	ratings.On("UserExists", mock.Anything, "myuser").Return(true)
	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.AtCoderUsername == "myuser"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{"username": "myuser"})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"atcoder_username":"myuser"`)
	users.AssertExpectations(t)
	ratings.AssertExpectations(t)
}

func TestAtCoder_Connect_MissingUsername(t *testing.T) {
	h, users, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
	ratings.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
	users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// AtCoder 上に存在しないユーザー名は 400。
func TestAtCoder_Connect_UnknownAtCoderUser(t *testing.T) {
	h, users, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	ratings.On("UserExists", mock.Anything, "ghost").Return(false)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{"username": "ghost"})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "invalid AtCoder username")
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// 形式が不正なユーザー名は外部 API を叩かずに 400。
func TestAtCoder_Connect_InvalidUsernameFormat(t *testing.T) {
	h, users, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{"username": "bad!user"})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "invalid AtCoder username")
	ratings.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestAtCoder_Connect_UserNotFound(t *testing.T) {
	h, users, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	ratings.On("UserExists", mock.Anything, "myuser").Return(true)
	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{"username": "myuser"})
	assertStatus(t, w, http.StatusNotFound)
	assert.Contains(t, w.Body.String(), "user not found")
	users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestAtCoder_Connect_UpdateError(t *testing.T) {
	h, users, ratings := newTestAtCoderHandler()
	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	ratings.On("UserExists", mock.Anything, "myuser").Return(true)
	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	users.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{"username": "myuser"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- Disconnect ----------

func TestAtCoder_Disconnect_Success(t *testing.T) {
	h, users, _ := newTestAtCoderHandler()
	r := newRouter(1)
	r.DELETE("/atcoder/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, AtCoderUsername: "myuser"}, nil)
	users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.AtCoderUsername == ""
	})).Return(nil)

	w := doRequest(r, http.MethodDelete, "/atcoder/disconnect", nil)
	assertStatus(t, w, http.StatusOK)
	users.AssertExpectations(t)
}

func TestAtCoder_Disconnect_UserNotFound(t *testing.T) {
	h, users, _ := newTestAtCoderHandler()
	r := newRouter(1)
	r.DELETE("/atcoder/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/atcoder/disconnect", nil)
	assertStatus(t, w, http.StatusNotFound)
	users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 取得に失敗した場合も移行前と同じく 404 に潰す。
func TestAtCoder_Disconnect_RepositoryError(t *testing.T) {
	h, users, _ := newTestAtCoderHandler()
	r := newRouter(1)
	r.DELETE("/atcoder/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodDelete, "/atcoder/disconnect", nil)
	assertStatus(t, w, http.StatusNotFound)
}
