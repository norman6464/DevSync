package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// GitHubRepository は GitHub 連携データの永続化に対する、usecase 側が要求する契約。
type GitHubRepository interface {
	// UpsertContributions は日別コントリビューションを (user_id, date) で重複判定して保存する。
	UpsertContributions(ctx context.Context, contributions []model.GitHubContribution) error
	// GetContributions は日付の昇順でコントリビューションを返す。
	GetContributions(ctx context.Context, userID uint) ([]model.GitHubContribution, error)
	// UpsertLanguageStats は言語統計を (user_id, language) で重複判定して保存する。
	UpsertLanguageStats(ctx context.Context, stats []model.GitHubLanguageStat) error
	// GetLanguageStats はバイト数の降順で言語統計を返す。
	GetLanguageStats(ctx context.Context, userID uint) ([]model.GitHubLanguageStat, error)
	// UpsertRepos はリポジトリを GitHub 側のリポジトリ ID で重複判定して保存する。
	UpsertRepos(ctx context.Context, repos []model.GitHubRepository) error
	// GetRepos はスター数の降順でリポジトリを返す。
	GetRepos(ctx context.Context, userID uint) ([]model.GitHubRepository, error)
	// DeleteUserData は指定ユーザーの GitHub 連携データをすべて削除する。
	DeleteUserData(ctx context.Context, userID uint) error
}

// GitHubAPIClient は GitHub API 呼び出しに対する、usecase 側が要求する契約。
// 実装は adapter/external に置く。
type GitHubAPIClient interface {
	// ConnectAuthorizeURL は連携用の OAuth 認可 URL を返す。
	ConnectAuthorizeURL(state string) string
	// LoginAuthorizeURL はログイン用の OAuth 認可 URL を返す（メール取得スコープ付き）。
	LoginAuthorizeURL(state string) string
	// ExchangeCode は認可コードをアクセストークンに交換する。
	ExchangeCode(ctx context.Context, code string) (string, error)
	// GetUser はアクセストークンからユーザー情報を取得する。
	GetUser(ctx context.Context, token string) (*model.GitHubUserInfo, error)
	// FetchContributions は過去 1 年の日別コントリビューションを取得する。
	FetchContributions(ctx context.Context, token, username string) ([]model.GitHubContributionDay, error)
	// FetchRepos は認証ユーザーのリポジトリ一覧を取得する。
	FetchRepos(ctx context.Context, token string) ([]model.GitHubRepoSummary, error)
	// FetchRepoLanguages は指定リポジトリの言語別バイト数を取得する。
	FetchRepoLanguages(ctx context.Context, token, fullName string) (map[string]int64, error)
}
