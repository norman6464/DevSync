package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

func TestGitHubConnect_Success(t *testing.T) {
	h, ghSvc, authSvc := setupGitHubHandlerMock()
	authSvc.On("GenerateOAuthState", uint(1)).Return("test-state", nil)
	ghSvc.On("GetOAuthURL", "test-state").Return("https://github.com/login/oauth/authorize?state=test-state")

	r := newRouter(1)
	r.GET("/github/connect", h.Connect)
	w := doRequest(r, "GET", "/github/connect", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	if url, ok := body["url"].(string); !ok || url == "" {
		t.Error("expected non-empty URL in response")
	}
	authSvc.AssertExpectations(t)
	ghSvc.AssertExpectations(t)
}

func TestGitHubConnect_StateError(t *testing.T) {
	h, _, authSvc := setupGitHubHandlerMock()
	authSvc.On("GenerateOAuthState", uint(1)).Return("", service.ErrBadRequest)

	r := newRouter(1)
	r.GET("/github/connect", h.Connect)
	w := doRequest(r, "GET", "/github/connect", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	authSvc.AssertExpectations(t)
}

func TestGitHubCallback_Success(t *testing.T) {
	h, ghSvc, authSvc := setupGitHubHandlerMock()
	authSvc.On("ValidateOAuthState", "valid-state").Return(1, nil)
	ghSvc.On("ConnectGitHub", uint(1), "test-code", "valid-state").Return(nil)

	r := newRouter(1)
	r.GET("/github/callback", h.Callback)
	w := doRequest(r, "GET", "/github/callback?code=test-code&state=valid-state", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	if msg, ok := body["message"].(string); !ok || msg == "" {
		t.Error("expected message in response")
	}
	authSvc.AssertExpectations(t)
	ghSvc.AssertExpectations(t)
}

func TestGitHubCallback_MissingParams(t *testing.T) {
	h, _, _ := setupGitHubHandlerMock()

	r := newRouter(1)
	r.GET("/github/callback", h.Callback)
	w := doRequest(r, "GET", "/github/callback", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestGitHubCallback_InvalidState(t *testing.T) {
	h, _, authSvc := setupGitHubHandlerMock()
	authSvc.On("ValidateOAuthState", "bad-state").Return(0, service.ErrBadRequest)

	r := newRouter(1)
	r.GET("/github/callback", h.Callback)
	w := doRequest(r, "GET", "/github/callback?code=test-code&state=bad-state", nil)

	assertStatus(t, w, http.StatusBadRequest)
	authSvc.AssertExpectations(t)
}

func TestGitHubGetContributions_Success(t *testing.T) {
	h, ghSvc, _ := setupGitHubHandlerMock()
	contributions := []model.GitHubContribution{
		{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Count: 5},
	}
	ghSvc.On("GetContributions", uint(1)).Return(contributions, nil)

	r := newRouter(1)
	r.GET("/github/:userId/contributions", h.GetContributions)
	w := doRequest(r, "GET", "/github/1/contributions", nil)

	assertStatus(t, w, http.StatusOK)
	ghSvc.AssertExpectations(t)
}

func TestGitHubGetContributions_ServiceError(t *testing.T) {
	h, ghSvc, _ := setupGitHubHandlerMock()
	ghSvc.On("GetContributions", uint(1)).Return([]model.GitHubContribution(nil), service.ErrNotFound)

	r := newRouter(1)
	r.GET("/github/:userId/contributions", h.GetContributions)
	w := doRequest(r, "GET", "/github/1/contributions", nil)

	assertStatus(t, w, http.StatusNotFound)
	ghSvc.AssertExpectations(t)
}

func TestGitHubGetLanguages_Success(t *testing.T) {
	h, ghSvc, _ := setupGitHubHandlerMock()
	languages := []model.GitHubLanguageStat{
		{Language: "Go", Bytes: 50000},
	}
	ghSvc.On("GetLanguages", uint(1)).Return(languages, nil)

	r := newRouter(1)
	r.GET("/github/:userId/languages", h.GetLanguages)
	w := doRequest(r, "GET", "/github/1/languages", nil)

	assertStatus(t, w, http.StatusOK)
	ghSvc.AssertExpectations(t)
}

func TestGitHubGetRepos_Success(t *testing.T) {
	h, ghSvc, _ := setupGitHubHandlerMock()
	repos := []model.GitHubRepository{
		{Name: "test-repo"},
	}
	ghSvc.On("GetRepos", uint(1)).Return(repos, nil)

	r := newRouter(1)
	r.GET("/github/:userId/repos", h.GetRepos)
	w := doRequest(r, "GET", "/github/1/repos", nil)

	assertStatus(t, w, http.StatusOK)
	ghSvc.AssertExpectations(t)
}

func TestGitHubSync_Success(t *testing.T) {
	h, ghSvc, _ := setupGitHubHandlerMock()
	ghSvc.On("SyncUserData", uint(1)).Return(nil)

	r := newRouter(1)
	r.POST("/github/sync", h.Sync)
	w := doRequest(r, "POST", "/github/sync", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	if msg, ok := body["message"].(string); !ok || msg == "" {
		t.Error("expected message in response")
	}
	ghSvc.AssertExpectations(t)
}

func TestGitHubSync_ServiceError(t *testing.T) {
	h, ghSvc, _ := setupGitHubHandlerMock()
	ghSvc.On("SyncUserData", uint(1)).Return(service.ErrNotFound)

	r := newRouter(1)
	r.POST("/github/sync", h.Sync)
	w := doRequest(r, "POST", "/github/sync", nil)

	assertStatus(t, w, http.StatusNotFound)
	ghSvc.AssertExpectations(t)
}

func TestGitHubDisconnect_Success(t *testing.T) {
	h, ghSvc, _ := setupGitHubHandlerMock()
	ghSvc.On("DisconnectGitHub", uint(1)).Return(nil)

	r := newRouter(1)
	r.POST("/github/disconnect", h.Disconnect)
	w := doRequest(r, "POST", "/github/disconnect", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	if msg, ok := body["message"].(string); !ok || msg == "" {
		t.Error("expected message in response")
	}
	ghSvc.AssertExpectations(t)
}

func TestGitHubDisconnect_ServiceError(t *testing.T) {
	h, ghSvc, _ := setupGitHubHandlerMock()
	ghSvc.On("DisconnectGitHub", uint(1)).Return(service.ErrNotFound)

	r := newRouter(1)
	r.POST("/github/disconnect", h.Disconnect)
	w := doRequest(r, "POST", "/github/disconnect", nil)

	assertStatus(t, w, http.StatusNotFound)
	ghSvc.AssertExpectations(t)
}
