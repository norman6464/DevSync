package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newGitHubTestService() (*GitHubService, *MockUserRepository, *MockGitHubRepository) {
	cfg := &config.Config{
		GitHubClientID:    "test-client-id",
		GitHubClientSecret: "test-secret",
		GitHubRedirectURL: "http://localhost:5173/github/callback",
	}
	userRepo := new(MockUserRepository)
	githubRepo := new(MockGitHubRepository)
	svc := NewGitHubService(cfg, userRepo, githubRepo)
	return svc, userRepo, githubRepo
}

// --- GetOAuthURL ---

func TestGitHubService_GetOAuthURL(t *testing.T) {
	svc, _, _ := newGitHubTestService()
	url := svc.GetOAuthURL("test-state")
	assert.Contains(t, url, "client_id=test-client-id")
	assert.Contains(t, url, "state=test-state")
	assert.Contains(t, url, "scope=read:user,repo")
	assert.Contains(t, url, "redirect_uri=http://localhost:5173/github/callback")
}

// --- GetLoginOAuthURL ---

func TestGitHubService_GetLoginOAuthURL(t *testing.T) {
	svc, _, _ := newGitHubTestService()
	url := svc.GetLoginOAuthURL("login-state")
	assert.Contains(t, url, "client_id=test-client-id")
	assert.Contains(t, url, "state=login-state")
	assert.Contains(t, url, "scope=read:user,user:email,repo")
}

// --- SyncData ---

func TestGitHubService_SyncData_NoToken(t *testing.T) {
	svc, _, _ := newGitHubTestService()
	user := &model.User{GitHubToken: ""}
	err := svc.SyncData(user)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

// --- SyncUserData ---

func TestGitHubService_SyncUserData_UserNotFound(t *testing.T) {
	svc, userRepo, _ := newGitHubTestService()
	userRepo.On("FindByID", uint(999)).Return(nil, ErrNotFound)
	err := svc.SyncUserData(999)
	assert.ErrorIs(t, err, ErrNotFound)
	userRepo.AssertExpectations(t)
}

// --- DisconnectGitHub ---

func TestGitHubService_DisconnectGitHub_Success(t *testing.T) {
	svc, userRepo, githubRepo := newGitHubTestService()
	user := &model.User{
		GitHubToken:     "token123",
		GitHubUsername:  "testuser",
		GitHubConnected: true,
	}
	user.ID = 1
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("Update", mock.MatchedBy(func(u *model.User) bool {
		return u.GitHubToken == "" && u.GitHubUsername == "" && !u.GitHubConnected
	})).Return(nil)
	githubRepo.On("DeleteUserData", uint(1)).Return(nil)

	err := svc.DisconnectGitHub(1)
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	githubRepo.AssertExpectations(t)
}

func TestGitHubService_DisconnectGitHub_UserNotFound(t *testing.T) {
	svc, userRepo, _ := newGitHubTestService()
	userRepo.On("FindByID", uint(999)).Return(nil, ErrNotFound)
	err := svc.DisconnectGitHub(999)
	assert.ErrorIs(t, err, ErrNotFound)
	userRepo.AssertExpectations(t)
}

func TestGitHubService_DisconnectGitHub_UpdateError(t *testing.T) {
	svc, userRepo, _ := newGitHubTestService()
	user := &model.User{GitHubToken: "token", GitHubConnected: true}
	user.ID = 1
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("Update", mock.Anything).Return(errors.New("db error"))

	err := svc.DisconnectGitHub(1)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	userRepo.AssertExpectations(t)
}

// --- GetContributions ---

func TestGitHubService_GetContributions_Success(t *testing.T) {
	svc, _, githubRepo := newGitHubTestService()
	expected := []model.GitHubContribution{
		{UserID: 1, Date: time.Now(), Count: 5},
		{UserID: 1, Date: time.Now().AddDate(0, 0, -1), Count: 3},
	}
	githubRepo.On("GetContributions", uint(1)).Return(expected, nil)

	result, err := svc.GetContributions(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 5, result[0].Count)
	githubRepo.AssertExpectations(t)
}

func TestGitHubService_GetContributions_Error(t *testing.T) {
	svc, _, githubRepo := newGitHubTestService()
	githubRepo.On("GetContributions", uint(1)).Return([]model.GitHubContribution(nil), errors.New("db error"))

	result, err := svc.GetContributions(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	githubRepo.AssertExpectations(t)
}

// --- GetLanguages ---

func TestGitHubService_GetLanguages_Success(t *testing.T) {
	svc, _, githubRepo := newGitHubTestService()
	expected := []model.GitHubLanguageStat{
		{UserID: 1, Language: "Go", Bytes: 10000, RepoCount: 3},
		{UserID: 1, Language: "TypeScript", Bytes: 8000, RepoCount: 2},
	}
	githubRepo.On("GetLanguageStats", uint(1)).Return(expected, nil)

	result, err := svc.GetLanguages(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Go", result[0].Language)
	githubRepo.AssertExpectations(t)
}

func TestGitHubService_GetLanguages_Error(t *testing.T) {
	svc, _, githubRepo := newGitHubTestService()
	githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat(nil), errors.New("db error"))

	result, err := svc.GetLanguages(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	githubRepo.AssertExpectations(t)
}

// --- GetRepos ---

func TestGitHubService_GetRepos_Success(t *testing.T) {
	svc, _, githubRepo := newGitHubTestService()
	expected := []model.GitHubRepository{
		{UserID: 1, Name: "devsync", Language: "Go", Stars: 10},
		{UserID: 1, Name: "frontend", Language: "TypeScript", Stars: 5},
	}
	githubRepo.On("GetRepos", uint(1)).Return(expected, nil)

	result, err := svc.GetRepos(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "devsync", result[0].Name)
	githubRepo.AssertExpectations(t)
}

func TestGitHubService_GetRepos_Error(t *testing.T) {
	svc, _, githubRepo := newGitHubTestService()
	githubRepo.On("GetRepos", uint(1)).Return([]model.GitHubRepository(nil), errors.New("db error"))

	result, err := svc.GetRepos(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	githubRepo.AssertExpectations(t)
}

func TestGitHubService_GetRepos_Empty(t *testing.T) {
	svc, _, githubRepo := newGitHubTestService()
	githubRepo.On("GetRepos", uint(1)).Return([]model.GitHubRepository{}, nil)

	result, err := svc.GetRepos(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
	githubRepo.AssertExpectations(t)
}
