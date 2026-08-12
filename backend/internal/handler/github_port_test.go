package handler

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// mockGitHubRepo は usecase/repository.GitHubRepository のモック。
type mockGitHubRepo struct{ mock.Mock }

func (m *mockGitHubRepo) UpsertContributions(ctx context.Context, contributions []model.GitHubContribution) error {
	return m.Called(ctx, contributions).Error(0)
}

func (m *mockGitHubRepo) GetContributions(ctx context.Context, userID uint) ([]model.GitHubContribution, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]model.GitHubContribution)
	return c, args.Error(1)
}

func (m *mockGitHubRepo) UpsertLanguageStats(ctx context.Context, stats []model.GitHubLanguageStat) error {
	return m.Called(ctx, stats).Error(0)
}

func (m *mockGitHubRepo) GetLanguageStats(ctx context.Context, userID uint) ([]model.GitHubLanguageStat, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).([]model.GitHubLanguageStat)
	return s, args.Error(1)
}

func (m *mockGitHubRepo) UpsertRepos(ctx context.Context, repos []model.GitHubRepository) error {
	return m.Called(ctx, repos).Error(0)
}

func (m *mockGitHubRepo) GetRepos(ctx context.Context, userID uint) ([]model.GitHubRepository, error) {
	args := m.Called(ctx, userID)
	r, _ := args.Get(0).([]model.GitHubRepository)
	return r, args.Error(1)
}

func (m *mockGitHubRepo) DeleteUserData(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

// mockGitHubAPIClient は usecase/repository.GitHubAPIClient のモック。
type mockGitHubAPIClient struct{ mock.Mock }

func (m *mockGitHubAPIClient) ConnectAuthorizeURL(state string) string {
	return m.Called(state).String(0)
}

func (m *mockGitHubAPIClient) LoginAuthorizeURL(state string) string {
	return m.Called(state).String(0)
}

func (m *mockGitHubAPIClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Error(1)
}

func (m *mockGitHubAPIClient) GetUser(ctx context.Context, token string) (*model.GitHubUserInfo, error) {
	args := m.Called(ctx, token)
	u, _ := args.Get(0).(*model.GitHubUserInfo)
	return u, args.Error(1)
}

func (m *mockGitHubAPIClient) FetchContributions(ctx context.Context, token, username string) ([]model.GitHubContributionDay, error) {
	args := m.Called(ctx, token, username)
	d, _ := args.Get(0).([]model.GitHubContributionDay)
	return d, args.Error(1)
}

func (m *mockGitHubAPIClient) FetchRepos(ctx context.Context, token string) ([]model.GitHubRepoSummary, error) {
	args := m.Called(ctx, token)
	r, _ := args.Get(0).([]model.GitHubRepoSummary)
	return r, args.Error(1)
}

func (m *mockGitHubAPIClient) FetchRepoLanguages(ctx context.Context, token, fullName string) (map[string]int64, error) {
	args := m.Called(ctx, token, fullName)
	l, _ := args.Get(0).(map[string]int64)
	return l, args.Error(1)
}

// githubPorts は GitHub 連携の usecase に注入した port モックをまとめる。
type githubPorts struct {
	Users  *mockUserPort
	Repo   *mockGitHubRepo
	Client *mockGitHubAPIClient
}
