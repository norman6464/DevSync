package handler

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// linkedUser は GitHub 連携済みのユーザーを返す。
func linkedUser(id uint) *model.User {
	return &model.User{
		ID:              id,
		Username:        "dev",
		GitHubUsername:  "dev",
		GitHubToken:     "token",
		GitHubConnected: true,
	}
}

func TestGitHubConnect_Success(t *testing.T) {
	h, ports, authSvc := setupGitHubHandlerMock()
	authSvc.On("GenerateOAuthState", uint(1)).Return("test-state", nil)
	ports.Client.On("ConnectAuthorizeURL", "test-state").
		Return("https://github.com/login/oauth/authorize?state=test-state")

	r := newRouter(1)
	r.GET("/github/connect", h.Connect)
	w := doRequest(r, "GET", "/github/connect", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.NotEmpty(t, body["url"])
	authSvc.AssertExpectations(t)
	ports.Client.AssertExpectations(t)
}

func TestGitHubConnect_StateError(t *testing.T) {
	h, ports, authSvc := setupGitHubHandlerMock()
	authSvc.On("GenerateOAuthState", uint(1)).Return("", service.ErrBadRequest)

	r := newRouter(1)
	r.GET("/github/connect", h.Connect)
	w := doRequest(r, "GET", "/github/connect", nil)

	assertStatus(t, w, http.StatusBadRequest)
	authSvc.AssertExpectations(t)
	ports.Client.AssertNotCalled(t, "ConnectAuthorizeURL", mock.Anything)
}

// 連携はトークン交換 → ユーザー取得 → 連携情報の保存まで通す。
func TestGitHubCallback_Success(t *testing.T) {
	h, ports, authSvc := setupGitHubHandlerMock()
	authSvc.On("ValidateOAuthState", "valid-state").Return(1, nil)
	ports.Client.On("ExchangeCode", mock.Anything, "test-code").Return("access-token", nil)
	ports.Client.On("GetUser", mock.Anything, "access-token").
		Return(&model.GitHubUserInfo{ID: 42, Login: "dev", AvatarURL: "https://example.com/a.png"}, nil)
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	ports.Users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.GitHubToken == "access-token" && u.GitHubID == 42 &&
			u.GitHubUsername == "dev" && u.GitHubConnected && u.AvatarURL == "https://example.com/a.png"
	})).Return(nil)
	// 連携直後の同期はバックグラウンドで走るため、呼ばれても呼ばれなくてもよい
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil).Maybe()
	ports.Client.On("FetchContributions", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("skip")).Maybe()

	r := newRouter(1)
	r.GET("/github/callback", h.Callback)
	w := doRequest(r, "GET", "/github/callback?code=test-code&state=valid-state", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.NotEmpty(t, body["message"])
	authSvc.AssertExpectations(t)
}

// トークン交換に失敗したらユーザー情報は更新しない。
func TestGitHubCallback_ExchangeError(t *testing.T) {
	h, ports, authSvc := setupGitHubHandlerMock()
	authSvc.On("ValidateOAuthState", "valid-state").Return(1, nil)
	ports.Client.On("ExchangeCode", mock.Anything, "test-code").Return("", service.ErrBadRequest)

	r := newRouter(1)
	r.GET("/github/callback", h.Callback)
	w := doRequest(r, "GET", "/github/callback?code=test-code&state=valid-state", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// アバターが空なら既存のアバターを維持する。
func TestGitHubCallback_KeepsAvatarWhenEmpty(t *testing.T) {
	h, ports, authSvc := setupGitHubHandlerMock()
	authSvc.On("ValidateOAuthState", "valid-state").Return(1, nil)
	ports.Client.On("ExchangeCode", mock.Anything, "test-code").Return("access-token", nil)
	ports.Client.On("GetUser", mock.Anything, "access-token").Return(&model.GitHubUserInfo{ID: 42, Login: "dev"}, nil)
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, AvatarURL: "existing.png"}, nil)
	ports.Users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.AvatarURL == "existing.png"
	})).Return(nil)
	ports.Client.On("FetchContributions", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("skip")).Maybe()

	r := newRouter(1)
	r.GET("/github/callback", h.Callback)
	w := doRequest(r, "GET", "/github/callback?code=test-code&state=valid-state", nil)

	assertStatus(t, w, http.StatusOK)
}

func TestGitHubCallback_MissingParams(t *testing.T) {
	h, _, _ := setupGitHubHandlerMock()

	r := newRouter(1)
	r.GET("/github/callback", h.Callback)
	w := doRequest(r, "GET", "/github/callback", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestGitHubCallback_InvalidState(t *testing.T) {
	h, ports, authSvc := setupGitHubHandlerMock()
	authSvc.On("ValidateOAuthState", "bad-state").Return(0, service.ErrBadRequest)

	r := newRouter(1)
	r.GET("/github/callback", h.Callback)
	w := doRequest(r, "GET", "/github/callback?code=test-code&state=bad-state", nil)

	assertStatus(t, w, http.StatusBadRequest)
	authSvc.AssertExpectations(t)
	ports.Client.AssertNotCalled(t, "ExchangeCode", mock.Anything, mock.Anything)
}

// ---------- 参照系 ----------

func TestGitHubGetContributions_Success(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Repo.On("GetContributions", mock.Anything, uint(1)).Return([]model.GitHubContribution{
		{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Count: 5},
	}, nil)

	r := newRouter(1)
	r.GET("/github/:userId/contributions", h.GetContributions)
	w := doRequest(r, "GET", "/github/1/contributions", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Repo.AssertExpectations(t)
}

// 同期データが無ければ null ではなく空配列を返す。
func TestGitHubGetContributions_Empty(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Repo.On("GetContributions", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.GET("/github/:userId/contributions", h.GetContributions)
	w := doRequest(r, "GET", "/github/1/contributions", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestGitHubGetContributions_RepositoryError(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Repo.On("GetContributions", mock.Anything, uint(1)).Return(nil, service.ErrNotFound)

	r := newRouter(1)
	r.GET("/github/:userId/contributions", h.GetContributions)
	w := doRequest(r, "GET", "/github/1/contributions", nil)

	assertStatus(t, w, http.StatusNotFound)
	ports.Repo.AssertExpectations(t)
}

func TestGitHubGetContributions_InvalidID(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	r := newRouter(1)
	r.GET("/github/:userId/contributions", h.GetContributions)
	w := doRequest(r, "GET", "/github/abc/contributions", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Repo.AssertNotCalled(t, "GetContributions", mock.Anything, mock.Anything)
}

func TestGitHubGetLanguages_Success(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Repo.On("GetLanguageStats", mock.Anything, uint(1)).
		Return([]model.GitHubLanguageStat{{Language: "Go", Bytes: 50000}}, nil)

	r := newRouter(1)
	r.GET("/github/:userId/languages", h.GetLanguages)
	w := doRequest(r, "GET", "/github/1/languages", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Repo.AssertExpectations(t)
}

func TestGitHubGetLanguages_RepositoryError(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Repo.On("GetLanguageStats", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/github/:userId/languages", h.GetLanguages)
	w := doRequest(r, "GET", "/github/1/languages", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	ports.Repo.AssertExpectations(t)
}

func TestGitHubGetLanguages_InvalidID(t *testing.T) {
	h, _, _ := setupGitHubHandlerMock()
	r := newRouter(1)
	r.GET("/github/:userId/languages", h.GetLanguages)
	w := doRequest(r, "GET", "/github/abc/languages", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestGitHubGetRepos_Success(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Repo.On("GetRepos", mock.Anything, uint(1)).
		Return([]model.GitHubRepository{{Name: "test-repo"}}, nil)

	r := newRouter(1)
	r.GET("/github/:userId/repos", h.GetRepos)
	w := doRequest(r, "GET", "/github/1/repos", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Repo.AssertExpectations(t)
}

func TestGitHubGetRepos_RepositoryError(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Repo.On("GetRepos", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/github/:userId/repos", h.GetRepos)
	w := doRequest(r, "GET", "/github/1/repos", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	ports.Repo.AssertExpectations(t)
}

func TestGitHubGetRepos_InvalidID(t *testing.T) {
	h, _, _ := setupGitHubHandlerMock()
	r := newRouter(1)
	r.GET("/github/:userId/repos", h.GetRepos)
	w := doRequest(r, "GET", "/github/abc/repos", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- 同期 ----------

func TestGitHubSync_Success(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(linkedUser(1), nil)
	ports.Client.On("FetchContributions", mock.Anything, "token", "dev").
		Return([]model.GitHubContributionDay{{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Count: 3}}, nil)
	ports.Repo.On("UpsertContributions", mock.Anything, mock.MatchedBy(func(c []model.GitHubContribution) bool {
		return len(c) == 1 && c[0].UserID == 1 && c[0].Count == 3
	})).Return(nil)
	ports.Client.On("FetchRepos", mock.Anything, "token").
		Return([]model.GitHubRepoSummary{{ID: 10, Name: "repo", FullName: "dev/repo", Language: "Go", Stars: 3}}, nil)
	ports.Repo.On("UpsertRepos", mock.Anything, mock.MatchedBy(func(r []model.GitHubRepository) bool {
		return len(r) == 1 && r[0].GitHubRepoID == 10 && r[0].UserID == 1
	})).Return(nil)
	ports.Client.On("FetchRepoLanguages", mock.Anything, "token", "dev/repo").Return(map[string]int64{"Go": 1200}, nil)
	ports.Repo.On("UpsertLanguageStats", mock.Anything, mock.MatchedBy(func(s []model.GitHubLanguageStat) bool {
		return len(s) == 1 && s[0].Language == "Go" && s[0].Bytes == 1200 && s[0].RepoCount == 1
	})).Return(nil)

	r := newRouter(1)
	r.POST("/github/sync", h.Sync)
	w := doRequest(r, "POST", "/github/sync", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Repo.AssertExpectations(t)
	ports.Client.AssertExpectations(t)
}

// 言語バイト数の取得に失敗しても同期は続行する（リポジトリ数だけ記録される）。
func TestGitHubSync_ContinuesWhenLanguageFetchFails(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(linkedUser(1), nil)
	ports.Client.On("FetchContributions", mock.Anything, "token", "dev").Return(nil, nil)
	ports.Repo.On("UpsertContributions", mock.Anything, mock.Anything).Return(nil)
	ports.Client.On("FetchRepos", mock.Anything, "token").
		Return([]model.GitHubRepoSummary{{ID: 10, FullName: "dev/repo", Language: "Go"}}, nil)
	ports.Repo.On("UpsertRepos", mock.Anything, mock.Anything).Return(nil)
	ports.Client.On("FetchRepoLanguages", mock.Anything, "token", "dev/repo").Return(nil, errors.New("rate limited"))
	ports.Repo.On("UpsertLanguageStats", mock.Anything, mock.MatchedBy(func(s []model.GitHubLanguageStat) bool {
		return len(s) == 1 && s[0].Language == "Go" && s[0].Bytes == 0 && s[0].RepoCount == 1
	})).Return(nil)

	r := newRouter(1)
	r.POST("/github/sync", h.Sync)
	w := doRequest(r, "POST", "/github/sync", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Repo.AssertExpectations(t)
}

// 連携していないユーザーの同期は 400 になる。
func TestGitHubSync_NotConnected(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)

	r := newRouter(1)
	r.POST("/github/sync", h.Sync)
	w := doRequest(r, "POST", "/github/sync", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Client.AssertNotCalled(t, "FetchContributions", mock.Anything, mock.Anything, mock.Anything)
}

func TestGitHubSync_UserNotFound(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/github/sync", h.Sync)
	w := doRequest(r, "POST", "/github/sync", nil)

	assertStatus(t, w, http.StatusNotFound)
	ports.Users.AssertExpectations(t)
}

// GitHub API の失敗は 503 になる。
func TestGitHubSync_APIError(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(linkedUser(1), nil)
	ports.Client.On("FetchContributions", mock.Anything, "token", "dev").Return(nil, errors.New("api error"))

	r := newRouter(1)
	r.POST("/github/sync", h.Sync)
	w := doRequest(r, "POST", "/github/sync", nil)

	assertStatus(t, w, http.StatusServiceUnavailable)
	ports.Repo.AssertNotCalled(t, "UpsertContributions", mock.Anything, mock.Anything)
}

// ---------- 連携解除 ----------

func TestGitHubDisconnect_Success(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(linkedUser(1), nil)
	ports.Users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.GitHubToken == "" && u.GitHubUsername == "" && !u.GitHubConnected
	})).Return(nil)
	ports.Repo.On("DeleteUserData", mock.Anything, uint(1)).Return(nil)

	r := newRouter(1)
	r.POST("/github/disconnect", h.Disconnect)
	w := doRequest(r, "POST", "/github/disconnect", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.NotEmpty(t, body["message"])
	ports.Users.AssertExpectations(t)
	ports.Repo.AssertExpectations(t)
}

// 同期データの削除に失敗しても連携解除は成功扱いにする（移行前と同じ）。
func TestGitHubDisconnect_DeleteDataFails(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(linkedUser(1), nil)
	ports.Users.On("Update", mock.Anything, mock.Anything).Return(nil)
	ports.Repo.On("DeleteUserData", mock.Anything, uint(1)).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/github/disconnect", h.Disconnect)
	w := doRequest(r, "POST", "/github/disconnect", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Repo.AssertExpectations(t)
}

func TestGitHubDisconnect_UserNotFound(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/github/disconnect", h.Disconnect)
	w := doRequest(r, "POST", "/github/disconnect", nil)

	assertStatus(t, w, http.StatusNotFound)
	ports.Users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	ports.Repo.AssertNotCalled(t, "DeleteUserData", mock.Anything, mock.Anything)
}

// ユーザー更新に失敗したらデータ削除は行わない。
func TestGitHubDisconnect_UpdateError(t *testing.T) {
	h, ports, _ := setupGitHubHandlerMock()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(linkedUser(1), nil)
	ports.Users.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/github/disconnect", h.Disconnect)
	w := doRequest(r, "POST", "/github/disconnect", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	ports.Repo.AssertNotCalled(t, "DeleteUserData", mock.Anything, mock.Anything)
}
