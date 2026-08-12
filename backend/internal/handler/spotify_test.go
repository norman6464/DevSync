package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockSpotifyRepo は usecase/repository.SpotifyRepository のモック。
type mockSpotifyRepo struct{ mock.Mock }

func (m *mockSpotifyRepo) DeleteUserData(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

// mockSpotifyAPIClient は usecase/repository.SpotifyAPIClient のモック。
type mockSpotifyAPIClient struct{ mock.Mock }

func (m *mockSpotifyAPIClient) AuthorizeURL(state string) string {
	return m.Called(state).String(0)
}

func (m *mockSpotifyAPIClient) ExchangeCode(ctx context.Context, code string) (*model.SpotifyToken, error) {
	args := m.Called(ctx, code)
	t, _ := args.Get(0).(*model.SpotifyToken)
	return t, args.Error(1)
}

func (m *mockSpotifyAPIClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*model.SpotifyToken, error) {
	args := m.Called(ctx, refreshToken)
	t, _ := args.Get(0).(*model.SpotifyToken)
	return t, args.Error(1)
}

func (m *mockSpotifyAPIClient) FetchCurrentlyPlaying(ctx context.Context, token string) (*model.SpotifyCurrentlyPlaying, error) {
	args := m.Called(ctx, token)
	p, _ := args.Get(0).(*model.SpotifyCurrentlyPlaying)
	return p, args.Error(1)
}

func (m *mockSpotifyAPIClient) FetchRecentlyPlayed(ctx context.Context, token string) ([]model.SpotifyRecentTrackResponse, error) {
	args := m.Called(ctx, token)
	tracks, _ := args.Get(0).([]model.SpotifyRecentTrackResponse)
	return tracks, args.Error(1)
}

// spotifyPorts は Spotify 連携の usecase に注入した port モックをまとめる。
type spotifyPorts struct {
	Users  *mockUserPort
	Repo   *mockSpotifyRepo
	Client *mockSpotifyAPIClient
}

// setupSpotifyHandler は本物の usecase に port モックを注入した SpotifyHandler を生成する。
func setupSpotifyHandler() (*SpotifyHandler, *spotifyPorts, *usecase.OAuthStateUseCase) {
	ports := &spotifyPorts{
		Users:  new(mockUserPort),
		Repo:   new(mockSpotifyRepo),
		Client: new(mockSpotifyAPIClient),
	}
	oauthState := usecase.NewOAuthStateUseCase(testJWTSecret)
	h := NewSpotifyHandler(SpotifyUseCases{
		OAuthURL:         usecase.NewGetSpotifyOAuthURLUseCase(ports.Client),
		Connect:          usecase.NewConnectSpotifyUseCase(ports.Users, ports.Client),
		Disconnect:       usecase.NewDisconnectSpotifyUseCase(ports.Users, ports.Repo),
		CurrentlyPlaying: usecase.NewGetSpotifyCurrentlyPlayingUseCase(ports.Users, ports.Client),
		RecentlyPlayed:   usecase.NewGetSpotifyRecentlyPlayedUseCase(ports.Users, ports.Client),
	}, oauthState)
	return h, ports, oauthState
}

// spotifyLinkedUser は Spotify 連携済みのユーザーを返す。
func spotifyLinkedUser(expiry time.Time) *model.User {
	return &model.User{
		ID:                  1,
		SpotifyConnected:    true,
		SpotifyToken:        "access-token",
		SpotifyRefreshToken: "refresh-token",
		SpotifyTokenExpiry:  expiry,
	}
}

// ---------- Connect ----------

func TestSpotifyConnect_Success(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Client.On("AuthorizeURL", mock.AnythingOfType("string")).Return("https://accounts.spotify.com/authorize?state=test-state")

	r := newRouter(1)
	r.GET("/spotify/connect", h.Connect)
	w := doRequest(r, http.MethodGet, "/spotify/connect", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "accounts.spotify.com")
	ports.Client.AssertExpectations(t)
}

// ---------- Callback ----------

func TestSpotifyCallback_Success(t *testing.T) {
	h, ports, authSvc := setupSpotifyHandler()
	state, err := authSvc.Generate(1)
	require.NoError(t, err)
	ports.Client.On("ExchangeCode", mock.Anything, "test-code").Return(&model.SpotifyToken{
		AccessToken: "new-access", RefreshToken: "new-refresh", ExpiresIn: 3600,
	}, nil)
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	ports.Users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.SpotifyConnected && u.SpotifyToken == "new-access" &&
			u.SpotifyRefreshToken == "new-refresh" && u.SpotifyTokenExpiry.After(time.Now())
	})).Return(nil)

	r := newRouter(1)
	r.GET("/spotify/callback", h.Callback)
	w := doRequest(r, http.MethodGet, "/spotify/callback?code=test-code&state="+state+"", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Users.AssertExpectations(t)
	ports.Client.AssertExpectations(t)
}

func TestSpotifyCallback_MissingParams(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()

	r := newRouter(1)
	r.GET("/spotify/callback", h.Callback)
	w := doRequest(r, http.MethodGet, "/spotify/callback", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Client.AssertNotCalled(t, "ExchangeCode", mock.Anything, mock.Anything)
}

func TestSpotifyCallback_InvalidState(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()

	r := newRouter(1)
	r.GET("/spotify/callback", h.Callback)
	w := doRequest(r, http.MethodGet, "/spotify/callback?code=test-code&state=bad-state", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Client.AssertNotCalled(t, "ExchangeCode", mock.Anything, mock.Anything)
}

func TestSpotifyCallback_ExchangeError(t *testing.T) {
	h, ports, authSvc := setupSpotifyHandler()
	state, err := authSvc.Generate(1)
	require.NoError(t, err)
	ports.Client.On("ExchangeCode", mock.Anything, "test-code").Return(nil, service.ErrBadRequest)

	r := newRouter(1)
	r.GET("/spotify/callback", h.Callback)
	w := doRequest(r, http.MethodGet, "/spotify/callback?code=test-code&state="+state+"", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// ユーザー情報の保存に失敗したら 500 を返す。
func TestSpotifyCallback_UpdateError(t *testing.T) {
	h, ports, authSvc := setupSpotifyHandler()
	state, err := authSvc.Generate(1)
	require.NoError(t, err)
	ports.Client.On("ExchangeCode", mock.Anything, "test-code").
		Return(&model.SpotifyToken{AccessToken: "a", RefreshToken: "r", ExpiresIn: 3600}, nil)
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	ports.Users.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

	r := newRouter(1)
	r.GET("/spotify/callback", h.Callback)
	w := doRequest(r, http.MethodGet, "/spotify/callback?code=test-code&state="+state+"", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	ports.Users.AssertExpectations(t)
}

// ---------- 再生情報 ----------

// 有効期限に余裕があればトークンを更新せずに使う。
func TestSpotifyGetCurrentlyPlaying_UsesValidToken(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(spotifyLinkedUser(time.Now().Add(time.Hour)), nil)
	ports.Client.On("FetchCurrentlyPlaying", mock.Anything, "access-token").
		Return(&model.SpotifyCurrentlyPlaying{IsPlaying: true, TrackName: "Song"}, nil)

	r := newRouter(1)
	r.GET("/spotify/:userId/currently-playing", h.GetCurrentlyPlaying)
	w := doRequest(r, http.MethodGet, "/spotify/1/currently-playing", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"track_name":"Song"`)
	ports.Client.AssertNotCalled(t, "RefreshAccessToken", mock.Anything, mock.Anything)
}

// 期限が近ければリフレッシュして保存し、新しいトークンで API を叩く。
func TestSpotifyGetCurrentlyPlaying_RefreshesExpiredToken(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(spotifyLinkedUser(time.Now().Add(time.Minute)), nil)
	ports.Client.On("RefreshAccessToken", mock.Anything, "refresh-token").
		Return(&model.SpotifyToken{AccessToken: "refreshed", ExpiresIn: 3600}, nil)
	ports.Users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		// リフレッシュトークンが返らない場合は既存の値を維持する
		return u.SpotifyToken == "refreshed" && u.SpotifyRefreshToken == "refresh-token"
	})).Return(nil)
	ports.Client.On("FetchCurrentlyPlaying", mock.Anything, "refreshed").Return(nil, nil)

	r := newRouter(1)
	r.GET("/spotify/:userId/currently-playing", h.GetCurrentlyPlaying)
	w := doRequest(r, http.MethodGet, "/spotify/1/currently-playing", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Users.AssertExpectations(t)
	ports.Client.AssertExpectations(t)
}

// 連携していないユーザーは 400。
func TestSpotifyGetCurrentlyPlaying_NotConnected(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)

	r := newRouter(1)
	r.GET("/spotify/:userId/currently-playing", h.GetCurrentlyPlaying)
	w := doRequest(r, http.MethodGet, "/spotify/1/currently-playing", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Client.AssertNotCalled(t, "FetchCurrentlyPlaying", mock.Anything, mock.Anything)
}

func TestSpotifyGetCurrentlyPlaying_UserNotFound(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.GET("/spotify/:userId/currently-playing", h.GetCurrentlyPlaying)
	w := doRequest(r, http.MethodGet, "/spotify/1/currently-playing", nil)

	assertStatus(t, w, http.StatusNotFound)
}

func TestSpotifyGetCurrentlyPlaying_InvalidID(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()

	r := newRouter(1)
	r.GET("/spotify/:userId/currently-playing", h.GetCurrentlyPlaying)
	w := doRequest(r, http.MethodGet, "/spotify/abc/currently-playing", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestSpotifyGetRecentlyPlayed_Success(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(spotifyLinkedUser(time.Now().Add(time.Hour)), nil)
	ports.Client.On("FetchRecentlyPlayed", mock.Anything, "access-token").
		Return([]model.SpotifyRecentTrackResponse{{TrackName: "Song", PlayedAt: "2024-01-01T00:00:00Z"}}, nil)

	r := newRouter(1)
	r.GET("/spotify/:userId/recently-played", h.GetRecentlyPlayed)
	w := doRequest(r, http.MethodGet, "/spotify/1/recently-played", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"track_name":"Song"`)
	ports.Client.AssertExpectations(t)
}

// 履歴が無ければ null ではなく空配列を返す。
func TestSpotifyGetRecentlyPlayed_Empty(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(spotifyLinkedUser(time.Now().Add(time.Hour)), nil)
	ports.Client.On("FetchRecentlyPlayed", mock.Anything, "access-token").Return(nil, nil)

	r := newRouter(1)
	r.GET("/spotify/:userId/recently-played", h.GetRecentlyPlayed)
	w := doRequest(r, http.MethodGet, "/spotify/1/recently-played", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestSpotifyGetRecentlyPlayed_APIError(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(spotifyLinkedUser(time.Now().Add(time.Hour)), nil)
	ports.Client.On("FetchRecentlyPlayed", mock.Anything, "access-token").Return(nil, errors.New("api error"))

	r := newRouter(1)
	r.GET("/spotify/:userId/recently-played", h.GetRecentlyPlayed)
	w := doRequest(r, http.MethodGet, "/spotify/1/recently-played", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	ports.Client.AssertExpectations(t)
}

// ---------- Disconnect ----------

func TestSpotifyDisconnect_Success(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(spotifyLinkedUser(time.Now().Add(time.Hour)), nil)
	ports.Users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return !u.SpotifyConnected && u.SpotifyToken == "" && u.SpotifyRefreshToken == "" &&
			u.SpotifyTokenExpiry.IsZero()
	})).Return(nil)
	ports.Repo.On("DeleteUserData", mock.Anything, uint(1)).Return(nil)

	r := newRouter(1)
	r.POST("/spotify/disconnect", h.Disconnect)
	w := doRequest(r, http.MethodPost, "/spotify/disconnect", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Users.AssertExpectations(t)
	ports.Repo.AssertExpectations(t)
}

// 再生履歴の削除に失敗しても連携解除は成功扱いにする（移行前と同じ）。
func TestSpotifyDisconnect_DeleteDataFails(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(spotifyLinkedUser(time.Now().Add(time.Hour)), nil)
	ports.Users.On("Update", mock.Anything, mock.Anything).Return(nil)
	ports.Repo.On("DeleteUserData", mock.Anything, uint(1)).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/spotify/disconnect", h.Disconnect)
	w := doRequest(r, http.MethodPost, "/spotify/disconnect", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Repo.AssertExpectations(t)
}

func TestSpotifyDisconnect_UpdateError(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(spotifyLinkedUser(time.Now().Add(time.Hour)), nil)
	ports.Users.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/spotify/disconnect", h.Disconnect)
	w := doRequest(r, http.MethodPost, "/spotify/disconnect", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	ports.Repo.AssertNotCalled(t, "DeleteUserData", mock.Anything, mock.Anything)
}

func TestSpotifyDisconnect_UserNotFound(t *testing.T) {
	h, ports, _ := setupSpotifyHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/spotify/disconnect", h.Disconnect)
	w := doRequest(r, http.MethodPost, "/spotify/disconnect", nil)

	assertStatus(t, w, http.StatusNotFound)
	ports.Users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}
