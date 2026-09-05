package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// githubRequestTimeout は GitHub への 1 リクエストのタイムアウト。
const githubRequestTimeout = 30 * time.Second

// githubReposPerPage はリポジトリ一覧の 1 ページあたりの取得件数。
const githubReposPerPage = 100

const (
	githubAuthorizeURL   = "https://github.com/login/oauth/authorize"
	githubAccessTokenURL = "https://github.com/login/oauth/access_token"
	githubUserURL        = "https://api.github.com/user"
	githubUserEmailsURL  = "https://api.github.com/user/emails"
	githubGraphQLURL     = "https://api.github.com/graphql"
)

// githubClient は [repository.GitHubAPIClient] の HTTP 実装。
type githubClient struct {
	clientID     string
	clientSecret string
	redirectURL  string
	httpClient   *http.Client
}

// NewGitHubClient は GitHubAPIClient の HTTP 実装を返す。
func NewGitHubClient(clientID, clientSecret, redirectURL string) repository.GitHubAPIClient {
	return &githubClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		httpClient:   &http.Client{Timeout: githubRequestTimeout},
	}
}

var _ repository.GitHubAPIClient = (*githubClient)(nil)

// ConnectAuthorizeURL は連携用の OAuth 認可 URL を返す。
func (c *githubClient) ConnectAuthorizeURL(state string) string {
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=read:user,repo&state=%s",
		githubAuthorizeURL, c.clientID, c.redirectURL, state)
}

// LoginAuthorizeURL はログイン用の OAuth 認可 URL を返す（メール取得スコープ付き）。
func (c *githubClient) LoginAuthorizeURL(state string) string {
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=read:user,user:email,repo&state=%s",
		githubAuthorizeURL, c.clientID, c.redirectURL, state)
}

// ExchangeCode は認可コードをアクセストークンに交換する。
func (c *githubClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"code":          code,
	})
	if err != nil {
		return "", domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub APIリクエストの生成に失敗しました", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubAccessTokenURL, bytes.NewReader(body))
	if err != nil {
		return "", domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub APIリクエストの生成に失敗しました", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub APIへの接続に失敗しました", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub APIレスポンスの解析に失敗しました", err)
	}
	if result.Error != "" {
		log.Printf("[WARN] GitHub OAuthエラー: %s", result.Error)
		return "", domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub認証に失敗しました", nil)
	}
	return result.AccessToken, nil
}

// GetUser はアクセストークンからユーザー情報を取得する。
// メールアドレスが非公開の場合は /user/emails からプライマリメールを補完する。
func (c *githubClient) GetUser(ctx context.Context, token string) (*model.GitHubUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserURL, nil)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubユーザー情報の取得に失敗しました", err)
	}
	c.setAuthHeaders(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubユーザー情報の取得に失敗しました", err)
	}
	defer resp.Body.Close()

	var user model.GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubユーザー情報の解析に失敗しました", err)
	}

	if user.Email == "" {
		user.Email = c.fetchPrimaryEmail(ctx, token)
	}
	return &user, nil
}

// fetchPrimaryEmail は /user/emails から認証済みプライマリメールを取得する。
// 取得できない場合は空文字を返す（ユーザー情報の取得自体は失敗させない）。
func (c *githubClient) fetchPrimaryEmail(ctx context.Context, token string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserEmailsURL, nil)
	if err != nil {
		return ""
	}
	c.setAuthHeaders(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(data, &emails); err != nil {
		return ""
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	return ""
}

// FetchContributions は GraphQL API で過去 1 年の日別コントリビューションを取得する。
func (c *githubClient) FetchContributions(ctx context.Context, token, username string) ([]model.GitHubContributionDay, error) {
	now := time.Now()
	from := now.AddDate(-1, 0, 0).Format("2006-01-02T15:04:05Z")
	to := now.Format("2006-01-02T15:04:05Z")

	// GraphQL インジェクション防止: ユーザー名から引用符とバックスラッシュを除去する
	safeUsername := strings.ReplaceAll(username, `"`, "")
	safeUsername = strings.ReplaceAll(safeUsername, `\`, "")

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
	}`, safeUsername, from, to)

	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub GraphQLリクエストの生成に失敗しました", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubGraphQLURL, bytes.NewReader(body))
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub GraphQLリクエストの生成に失敗しました", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub GraphQL APIへの接続に失敗しました", err)
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
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub GraphQLレスポンスの解析に失敗しました", err)
	}

	var days []model.GitHubContributionDay
	for _, week := range result.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			date, _ := time.Parse("2006-01-02", day.Date)
			days = append(days, model.GitHubContributionDay{Date: date, Count: day.ContributionCount})
		}
	}
	return days, nil
}

// FetchRepos は認証ユーザーのリポジトリ一覧を全ページ取得する（更新の新しい順）。
func (c *githubClient) FetchRepos(ctx context.Context, token string) ([]model.GitHubRepoSummary, error) {
	var all []model.GitHubRepoSummary

	for page := 1; ; page++ {
		url := fmt.Sprintf("https://api.github.com/user/repos?per_page=%d&page=%d&sort=updated", githubReposPerPage, page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubリポジトリ一覧の取得に失敗しました", err)
		}
		c.setAuthHeaders(req, token)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubリポジトリ一覧の取得に失敗しました", err)
		}
		data, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubリポジトリ一覧の取得に失敗しました", readErr)
		}

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
			return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHubリポジトリ一覧の解析に失敗しました", err)
		}
		if len(repos) == 0 {
			break
		}
		for _, r := range repos {
			all = append(all, model.GitHubRepoSummary{
				ID: r.ID, Name: r.Name, FullName: r.FullName, Description: r.Description,
				Language: r.Language, Stars: r.Stars, Forks: r.Forks, Private: r.Private,
			})
		}
		if len(repos) < githubReposPerPage {
			break
		}
	}
	return all, nil
}

// FetchRepoLanguages は指定リポジトリの言語別バイト数を取得する。
func (c *githubClient) FetchRepoLanguages(ctx context.Context, token, fullName string) (map[string]int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("https://api.github.com/repos/%s/languages", fullName), nil)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub言語統計の取得に失敗しました", err)
	}
	c.setAuthHeaders(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub言語統計の取得に失敗しました", err)
	}
	defer resp.Body.Close()

	var langs map[string]int64
	if err := json.NewDecoder(resp.Body).Decode(&langs); err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "GitHub言語統計の解析に失敗しました", err)
	}
	return langs, nil
}

// setAuthHeaders は GitHub API 呼び出しに必要な共通ヘッダーを設定する。
func (c *githubClient) setAuthHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
}
