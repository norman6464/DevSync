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

type GitHubService struct {
	cfg        *config.Config
	userRepo   *repository.UserRepository
	githubRepo *repository.GitHubRepository
}

func NewGitHubService(cfg *config.Config, userRepo *repository.UserRepository, githubRepo *repository.GitHubRepository) *GitHubService {
	return &GitHubService{cfg: cfg, userRepo: userRepo, githubRepo: githubRepo}
}

func (s *GitHubService) GetOAuthURL(state string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user,repo&state=%s",
		s.cfg.GitHubClientID, s.cfg.GitHubRedirectURL, state,
	)
}

func (s *GitHubService) GetLoginOAuthURL(state string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user,user:email,repo&state=%s",
		s.cfg.GitHubClientID, s.cfg.GitHubRedirectURL, state,
	)
}

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

type GitHubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

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

	// GitHub may not return email from /user if it's private, try /user/emails
	if user.Email == "" {
		user.Email = s.fetchPrimaryEmail(token)
	}

	return &user, nil
}

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

func (s *GitHubService) SyncData(user *model.User) error {
	if user.GitHubToken == "" {
		return fmt.Errorf("github not connected")
	}

	// Sync contributions
	if err := s.syncContributions(user); err != nil {
		return fmt.Errorf("sync contributions: %w", err)
	}

	// Sync repos and languages
	if err := s.syncReposAndLanguages(user); err != nil {
		return fmt.Errorf("sync repos: %w", err)
	}

	return nil
}

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

	// Save repos
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

	// Fetch language bytes for top repos (limit to 20 to avoid rate limiting)
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
