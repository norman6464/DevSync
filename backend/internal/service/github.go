package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// GitHubService はGitHub連携のビジネスロジックを提供する。
// OAuth認証フロー、データ同期（コントリビューション・リポジトリ・言語統計）を担当する。
type GitHubService struct {
	cfg        *config.Config
	userRepo   repository.UserRepositoryInterface
	githubRepo repository.GitHubRepositoryInterface
}

// NewGitHubService は新しいGitHubServiceインスタンスを生成する。
func NewGitHubService(cfg *config.Config, userRepo repository.UserRepositoryInterface, githubRepo repository.GitHubRepositoryInterface) *GitHubService {
	return &GitHubService{cfg: cfg, userRepo: userRepo, githubRepo: githubRepo}
}

// GetOAuthURL はGitHub連携用のOAuth認可URLを生成する。
func (s *GitHubService) GetOAuthURL(state string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user,repo&state=%s",
		s.cfg.GitHubClientID, s.cfg.GitHubRedirectURL, state,
	)
}

// GetLoginOAuthURL はGitHubログイン用のOAuth認可URLを生成する（メール取得スコープ付き）。
func (s *GitHubService) GetLoginOAuthURL(state string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user,user:email,repo&state=%s",
		s.cfg.GitHubClientID, s.cfg.GitHubRedirectURL, state,
	)
}

// ExchangeCode はOAuth認可コードをアクセストークンに交換する。
func (s *GitHubService) ExchangeCode(code string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"client_id":     s.cfg.GitHubClientID,
		"client_secret": s.cfg.GitHubClientSecret,
		"code":          code,
	})

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", fmt.Errorf("github oauth error: %s", result.Error)
	}
	return result.AccessToken, nil
}

// GitHubUserInfo はGitHub APIから取得するユーザー情報を表す。
type GitHubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// GetGitHubUser はアクセストークンを使ってGitHubユーザー情報を取得する。
// メールアドレスが非公開の場合は /user/emails APIから取得を試みる。
func (s *GitHubService) GetGitHubUser(token string) (*GitHubUserInfo, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// メールアドレスが非公開の場合、/user/emails APIからプライマリメールを取得
	if user.Email == "" {
		user.Email = s.fetchPrimaryEmail(token)
	}

	return &user, nil
}

// fetchPrimaryEmail はGitHub /user/emails APIから認証済みプライマリメールアドレスを取得する。
func (s *GitHubService) fetchPrimaryEmail(token string) string {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(func() []byte { b, _ := io.ReadAll(resp.Body); return b }(), &emails); err != nil {
		return ""
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	return ""
}

// SyncData はGitHubのコントリビューション、リポジトリ、言語統計を同期する。
func (s *GitHubService) SyncData(user *model.User) error {
	if user.GitHubToken == "" {
		return fmt.Errorf("github not connected")
	}

	// コントリビューションを同期
	if err := s.syncContributions(user); err != nil {
		return fmt.Errorf("sync contributions: %w", err)
	}

	// リポジトリと言語統計を同期
	if err := s.syncReposAndLanguages(user); err != nil {
		return fmt.Errorf("sync repos: %w", err)
	}

	return nil
}

// syncContributions はGitHub GraphQL APIを使って過去1年のコントリビューションデータを同期する。
func (s *GitHubService) syncContributions(user *model.User) error {
	now := time.Now()
	from := now.AddDate(-1, 0, 0).Format("2006-01-02T15:04:05Z")
	to := now.Format("2006-01-02T15:04:05Z")

	query := fmt.Sprintf(`query {
		user(login: "%s") {
			contributionsCollection(from: "%s", to: "%s") {
				contributionCalendar {
					weeks {
						contributionDays {
							date
							contributionCount
						}
					}
				}
			}
		}
	}`, user.GitHubUsername, from, to)

	body, _ := json.Marshal(map[string]string{"query": query})
	req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+user.GitHubToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			User struct {
				ContributionsCollection struct {
					ContributionCalendar struct {
						Weeks []struct {
							ContributionDays []struct {
								Date              string `json:"date"`
								ContributionCount int    `json:"contributionCount"`
							} `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	var contributions []model.GitHubContribution
	for _, week := range result.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			date, _ := time.Parse("2006-01-02", day.Date)
			contributions = append(contributions, model.GitHubContribution{
				UserID: user.ID,
				Date:   date,
				Count:  day.ContributionCount,
			})
		}
	}

	return s.githubRepo.UpsertContributions(contributions)
}

// syncReposAndLanguages はGitHub REST APIを使ってリポジトリと言語統計を同期する。
// レート制限を考慮し、言語バイト数の取得は上位20リポジトリに限定する。
func (s *GitHubService) syncReposAndLanguages(user *model.User) error {
	var allRepos []struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		Language    string `json:"language"`
		Stars       int    `json:"stargazers_count"`
		Forks       int    `json:"forks_count"`
		Private     bool   `json:"private"`
	}

	page := 1
	for {
		url := fmt.Sprintf("https://api.github.com/user/repos?per_page=100&page=%d&sort=updated", page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+user.GitHubToken)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var repos []struct {
			ID          int64  `json:"id"`
			Name        string `json:"name"`
			FullName    string `json:"full_name"`
			Description string `json:"description"`
			Language    string `json:"language"`
			Stars       int    `json:"stargazers_count"`
			Forks       int    `json:"forks_count"`
			Private     bool   `json:"private"`
		}
		if err := json.Unmarshal(data, &repos); err != nil {
			return err
		}
		if len(repos) == 0 {
			break
		}
		allRepos = append(allRepos, repos...)
		if len(repos) < 100 {
			break
		}
		page++
	}

	// リポジトリデータをモデルに変換して保存
	var modelRepos []model.GitHubRepository
	langMap := make(map[string]*model.GitHubLanguageStat)

	for _, r := range allRepos {
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

		if r.Language != "" {
			if stat, ok := langMap[r.Language]; ok {
				stat.RepoCount++
			} else {
				langMap[r.Language] = &model.GitHubLanguageStat{
					UserID:    user.ID,
					Language:  r.Language,
					Bytes:     0,
					RepoCount: 1,
				}
			}
		}
	}

	if err := s.githubRepo.UpsertRepos(modelRepos); err != nil {
		return err
	}

	// 上位20リポジトリの言語バイト数を取得（レート制限対策）
	limit := 20
	if len(allRepos) < limit {
		limit = len(allRepos)
	}
	for i := 0; i < limit; i++ {
		r := allRepos[i]
		url := fmt.Sprintf("https://api.github.com/repos/%s/languages", r.FullName)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+user.GitHubToken)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		var langs map[string]int64
		json.NewDecoder(resp.Body).Decode(&langs)
		resp.Body.Close()

		for lang, b := range langs {
			if stat, ok := langMap[lang]; ok {
				stat.Bytes += b
			} else {
				langMap[lang] = &model.GitHubLanguageStat{
					UserID:    user.ID,
					Language:  lang,
					Bytes:     b,
					RepoCount: 0,
				}
			}
		}
	}

	var langStats []model.GitHubLanguageStat
	for _, stat := range langMap {
		langStats = append(langStats, *stat)
	}

	return s.githubRepo.UpsertLanguageStats(langStats)
}

// ConnectGitHub はOAuthコールバック後にユーザーのGitHubアカウントを連携する。
// 連携完了後にバックグラウンドでデータ同期を開始する。
func (s *GitHubService) ConnectGitHub(userID uint, code, state string) error {
	accessToken, err := s.ExchangeCode(code)
	if err != nil {
		return err
	}

	ghUser, err := s.GetGitHubUser(accessToken)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound
	}

	user.GitHubToken = accessToken
	user.GitHubID = ghUser.ID
	user.GitHubUsername = ghUser.Login
	user.GitHubConnected = true
	if ghUser.AvatarURL != "" {
		user.AvatarURL = ghUser.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	go s.SyncData(user)
	return nil
}

// DisconnectGitHub はユーザーのGitHub連携を解除し、関連データを削除する。
func (s *GitHubService) DisconnectGitHub(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound
	}

	user.GitHubToken = ""
	user.GitHubUsername = ""
	user.GitHubConnected = false

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	s.githubRepo.DeleteUserData(userID)
	return nil
}

// SyncUserData は指定ユーザーIDのGitHubデータを同期する。
func (s *GitHubService) SyncUserData(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound
	}
	return s.SyncData(user)
}

// GetContributions は指定ユーザーのGitHubコントリビューションデータを取得する。
func (s *GitHubService) GetContributions(userID uint) ([]model.GitHubContribution, error) {
	return s.githubRepo.GetContributions(userID)
}

// GetLanguages は指定ユーザーのGitHub言語統計を取得する。
func (s *GitHubService) GetLanguages(userID uint) ([]model.GitHubLanguageStat, error) {
	return s.githubRepo.GetLanguageStats(userID)
}

// GetRepos は指定ユーザーのGitHubリポジトリ一覧を取得する。
func (s *GitHubService) GetRepos(userID uint) ([]model.GitHubRepository, error) {
	return s.githubRepo.GetRepos(userID)
}
