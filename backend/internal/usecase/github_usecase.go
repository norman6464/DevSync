package usecase

import (
	"context"
	"log"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// githubLanguageRepoLimit は言語バイト数を取得するリポジトリ数の上限。
// GitHub API のレート制限を避けるため、更新の新しい順に上位のみを対象にする。
const githubLanguageRepoLimit = 20

// GetGitHubOAuthURLUseCase は GitHub 連携用の認可 URL を組み立てる。
type GetGitHubOAuthURLUseCase struct {
	client repository.GitHubAPIClient
}

// NewGetGitHubOAuthURLUseCase は GetGitHubOAuthURLUseCase を生成する。
func NewGetGitHubOAuthURLUseCase(client repository.GitHubAPIClient) *GetGitHubOAuthURLUseCase {
	return &GetGitHubOAuthURLUseCase{client: client}
}

// Execute は連携用の認可 URL を返す。URL の組み立てだけで外部通信を行わないため ctx は取らない。
func (uc *GetGitHubOAuthURLUseCase) Execute(state string) string {
	return uc.client.ConnectAuthorizeURL(state)
}

// GetGitHubLoginURLUseCase は GitHub ログイン用の認可 URL を組み立てる。
type GetGitHubLoginURLUseCase struct {
	client repository.GitHubAPIClient
}

// NewGetGitHubLoginURLUseCase は GetGitHubLoginURLUseCase を生成する。
func NewGetGitHubLoginURLUseCase(client repository.GitHubAPIClient) *GetGitHubLoginURLUseCase {
	return &GetGitHubLoginURLUseCase{client: client}
}

// Execute はログイン用の認可 URL を返す。URL の組み立てだけで外部通信を行わないため ctx は取らない。
func (uc *GetGitHubLoginURLUseCase) Execute(state string) string {
	return uc.client.LoginAuthorizeURL(state)
}

// ExchangeGitHubCodeUseCase は OAuth 認可コードをアクセストークンに交換する。
type ExchangeGitHubCodeUseCase struct {
	client repository.GitHubAPIClient
}

// NewExchangeGitHubCodeUseCase は ExchangeGitHubCodeUseCase を生成する。
func NewExchangeGitHubCodeUseCase(client repository.GitHubAPIClient) *ExchangeGitHubCodeUseCase {
	return &ExchangeGitHubCodeUseCase{client: client}
}

// Execute は認可コードをアクセストークンに交換する。
func (uc *ExchangeGitHubCodeUseCase) Execute(ctx context.Context, code string) (string, error) {
	return uc.client.ExchangeCode(ctx, code)
}

// GetGitHubUserUseCase はアクセストークンから GitHub ユーザー情報を取得する。
type GetGitHubUserUseCase struct {
	client repository.GitHubAPIClient
}

// NewGetGitHubUserUseCase は GetGitHubUserUseCase を生成する。
func NewGetGitHubUserUseCase(client repository.GitHubAPIClient) *GetGitHubUserUseCase {
	return &GetGitHubUserUseCase{client: client}
}

// Execute は GitHub ユーザー情報を返す。
func (uc *GetGitHubUserUseCase) Execute(ctx context.Context, token string) (*model.GitHubUserInfo, error) {
	return uc.client.GetUser(ctx, token)
}

// SyncGitHubDataUseCase は GitHub のコントリビューション・リポジトリ・言語統計を同期する。
type SyncGitHubDataUseCase struct {
	users  repository.UserSkillsReader
	github repository.GitHubRepository
	client repository.GitHubAPIClient
}

// NewSyncGitHubDataUseCase は SyncGitHubDataUseCase を生成する。
func NewSyncGitHubDataUseCase(
	users repository.UserSkillsReader,
	github repository.GitHubRepository,
	client repository.GitHubAPIClient,
) *SyncGitHubDataUseCase {
	return &SyncGitHubDataUseCase{users: users, github: github, client: client}
}

// Execute は指定ユーザーの GitHub データを同期する。
func (uc *SyncGitHubDataUseCase) Execute(ctx context.Context, userID uint) error {
	user, err := findLinkedUser(ctx, uc.users, userID)
	if err != nil {
		return err
	}
	return uc.SyncUser(ctx, user)
}

// SyncUser は取得済みのユーザーに対して同期を実行する。連携直後の同期で使う。
func (uc *SyncGitHubDataUseCase) SyncUser(ctx context.Context, user *model.User) error {
	if user.GitHubToken == "" {
		return domain.NewError(domain.ErrCodeBadRequest, "GitHubが連携されていません", nil)
	}

	if err := uc.syncContributions(ctx, user); err != nil {
		return domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubコントリビューションの同期に失敗", err)
	}
	if err := uc.syncReposAndLanguages(ctx, user); err != nil {
		return domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubリポジトリの同期に失敗", err)
	}
	return nil
}

// SyncInBackground は同期をバックグラウンドで実行する。レスポンスを待たせないための入り口で、
// リクエストの終了で打ち切られないよう ctx のキャンセルは切り離す。
func (uc *SyncGitHubDataUseCase) SyncInBackground(ctx context.Context, userID uint) {
	detached := context.WithoutCancel(ctx)
	go func() {
		if err := uc.Execute(detached, userID); err != nil {
			log.Printf("[ERROR] GitHubデータ同期に失敗 (userID=%d): %v", userID, err)
		}
	}()
}

// syncContributions は日別コントリビューションを取得して保存する。
func (uc *SyncGitHubDataUseCase) syncContributions(ctx context.Context, user *model.User) error {
	days, err := uc.client.FetchContributions(ctx, user.GitHubToken, user.GitHubUsername)
	if err != nil {
		return err
	}

	contributions := make([]model.GitHubContribution, 0, len(days))
	for _, day := range days {
		contributions = append(contributions, model.GitHubContribution{
			UserID: user.ID,
			Date:   day.Date,
			Count:  day.Count,
		})
	}
	return uc.github.UpsertContributions(ctx, contributions)
}

// syncReposAndLanguages はリポジトリ一覧と言語統計を取得して保存する。
// 言語バイト数の取得はレート制限を考慮して上位 20 リポジトリに限定し、
// 個々の取得に失敗しても同期全体は継続する。
func (uc *SyncGitHubDataUseCase) syncReposAndLanguages(ctx context.Context, user *model.User) error {
	repos, err := uc.client.FetchRepos(ctx, user.GitHubToken)
	if err != nil {
		return err
	}

	modelRepos := make([]model.GitHubRepository, 0, len(repos))
	langMap := make(map[string]*model.GitHubLanguageStat)

	for _, r := range repos {
		modelRepos = append(modelRepos, model.GitHubRepository{
			UserID:       user.ID,
			GitHubRepoID: r.ID,
			Name:         r.Name,
			FullName:     r.FullName,
			Description:  r.Description,
			Language:     r.Language,
			Stars:        r.Stars,
			Forks:        r.Forks,
			IsPrivate:    r.Private,
		})

		if r.Language == "" {
			continue
		}
		if stat, ok := langMap[r.Language]; ok {
			stat.RepoCount++
		} else {
			langMap[r.Language] = &model.GitHubLanguageStat{
				UserID:    user.ID,
				Language:  r.Language,
				RepoCount: 1,
			}
		}
	}

	if err := uc.github.UpsertRepos(ctx, modelRepos); err != nil {
		return domain.NewError(domain.ErrCodeDatabase, "リポジトリデータの保存に失敗しました", err)
	}

	limit := githubLanguageRepoLimit
	if len(repos) < limit {
		limit = len(repos)
	}
	for _, r := range repos[:limit] {
		langs, err := uc.client.FetchRepoLanguages(ctx, user.GitHubToken, r.FullName)
		if err != nil {
			continue
		}
		for lang, bytes := range langs {
			if stat, ok := langMap[lang]; ok {
				stat.Bytes += bytes
			} else {
				langMap[lang] = &model.GitHubLanguageStat{
					UserID:   user.ID,
					Language: lang,
					Bytes:    bytes,
				}
			}
		}
	}

	langStats := make([]model.GitHubLanguageStat, 0, len(langMap))
	for _, stat := range langMap {
		langStats = append(langStats, *stat)
	}
	if err := uc.github.UpsertLanguageStats(ctx, langStats); err != nil {
		return domain.NewError(domain.ErrCodeDatabase, "言語統計の保存に失敗しました", err)
	}
	return nil
}

// ConnectGitHubUseCase は OAuth コールバック後に GitHub アカウントを連携する。
type ConnectGitHubUseCase struct {
	users  repository.ExternalAccountLinker
	client repository.GitHubAPIClient
	sync   *SyncGitHubDataUseCase
}

// NewConnectGitHubUseCase は ConnectGitHubUseCase を生成する。
func NewConnectGitHubUseCase(
	users repository.ExternalAccountLinker,
	client repository.GitHubAPIClient,
	sync *SyncGitHubDataUseCase,
) *ConnectGitHubUseCase {
	return &ConnectGitHubUseCase{users: users, client: client, sync: sync}
}

// Execute は認可コードからトークンとユーザー情報を取得し、連携情報を保存する。
// 保存後はバックグラウンドでデータ同期を開始する。
func (uc *ConnectGitHubUseCase) Execute(ctx context.Context, userID uint, code string) error {
	accessToken, err := uc.client.ExchangeCode(ctx, code)
	if err != nil {
		return err
	}

	ghUser, err := uc.client.GetUser(ctx, accessToken)
	if err != nil {
		return err
	}

	user, err := findLinkedUser(ctx, uc.users, userID)
	if err != nil {
		return err
	}

	user.GitHubToken = accessToken
	user.GitHubID = ghUser.ID
	user.GitHubUsername = ghUser.Login
	user.GitHubConnected = true
	if ghUser.AvatarURL != "" {
		user.AvatarURL = ghUser.AvatarURL
	}

	if err := uc.users.Update(ctx, user); err != nil {
		return domain.NewError(domain.ErrCodeInternal, "ユーザー情報の更新に失敗しました", err)
	}

	uc.sync.SyncInBackground(ctx, user.ID)
	return nil
}

// DisconnectGitHubUseCase は GitHub 連携を解除する。
type DisconnectGitHubUseCase struct {
	users  repository.ExternalAccountLinker
	github repository.GitHubRepository
}

// NewDisconnectGitHubUseCase は DisconnectGitHubUseCase を生成する。
func NewDisconnectGitHubUseCase(
	users repository.ExternalAccountLinker,
	github repository.GitHubRepository,
) *DisconnectGitHubUseCase {
	return &DisconnectGitHubUseCase{users: users, github: github}
}

// Execute は連携情報を消し、同期済みの GitHub データを削除する。
// データ削除の失敗は連携解除自体を失敗させない（解除は完了しているため）。
func (uc *DisconnectGitHubUseCase) Execute(ctx context.Context, userID uint) error {
	user, err := findLinkedUser(ctx, uc.users, userID)
	if err != nil {
		return err
	}

	user.GitHubToken = ""
	user.GitHubUsername = ""
	user.GitHubConnected = false

	if err := uc.users.Update(ctx, user); err != nil {
		return domain.NewError(domain.ErrCodeInternal, "ユーザー情報の更新に失敗しました", err)
	}

	if err := uc.github.DeleteUserData(ctx, userID); err != nil {
		log.Printf("[ERROR] GitHub連携データの削除に失敗 (userID=%d): %v", userID, err)
	}
	return nil
}

// GetGitHubContributionsUseCase は同期済みのコントリビューションを取得する。
type GetGitHubContributionsUseCase struct {
	github repository.GitHubRepository
}

// NewGetGitHubContributionsUseCase は GetGitHubContributionsUseCase を生成する。
func NewGetGitHubContributionsUseCase(github repository.GitHubRepository) *GetGitHubContributionsUseCase {
	return &GetGitHubContributionsUseCase{github: github}
}

// Execute は日付の昇順でコントリビューションを返す。
func (uc *GetGitHubContributionsUseCase) Execute(ctx context.Context, userID uint) ([]model.GitHubContribution, error) {
	return uc.github.GetContributions(ctx, userID)
}

// GetGitHubLanguagesUseCase は同期済みの言語統計を取得する。
type GetGitHubLanguagesUseCase struct {
	github repository.GitHubRepository
}

// NewGetGitHubLanguagesUseCase は GetGitHubLanguagesUseCase を生成する。
func NewGetGitHubLanguagesUseCase(github repository.GitHubRepository) *GetGitHubLanguagesUseCase {
	return &GetGitHubLanguagesUseCase{github: github}
}

// Execute はバイト数の降順で言語統計を返す。
func (uc *GetGitHubLanguagesUseCase) Execute(ctx context.Context, userID uint) ([]model.GitHubLanguageStat, error) {
	return uc.github.GetLanguageStats(ctx, userID)
}

// GetGitHubReposUseCase は同期済みのリポジトリ一覧を取得する。
type GetGitHubReposUseCase struct {
	github repository.GitHubRepository
}

// NewGetGitHubReposUseCase は GetGitHubReposUseCase を生成する。
func NewGetGitHubReposUseCase(github repository.GitHubRepository) *GetGitHubReposUseCase {
	return &GetGitHubReposUseCase{github: github}
}

// Execute はスター数の降順でリポジトリを返す。
func (uc *GetGitHubReposUseCase) Execute(ctx context.Context, userID uint) ([]model.GitHubRepository, error) {
	return uc.github.GetRepos(ctx, userID)
}
