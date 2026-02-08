package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type QiitaService struct {
	httpClient *http.Client
	userRepo   repository.UserRepositoryInterface
	qiitaRepo  repository.QiitaRepositoryInterface
}

func NewQiitaService(userRepo repository.UserRepositoryInterface, qiitaRepo repository.QiitaRepositoryInterface) *QiitaService {
	return &QiitaService{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userRepo:   userRepo,
		qiitaRepo:  qiitaRepo,
	}
}

// QiitaAPIArticle represents an article from Qiita API
type QiitaAPIArticle struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	URL           string          `json:"url"`
	LikesCount    int             `json:"likes_count"`
	CommentsCount int             `json:"comments_count"`
	Tags          []QiitaAPITag   `json:"tags"`
	CreatedAt     time.Time       `json:"created_at"`
}

type QiitaAPITag struct {
	Name string `json:"name"`
}

// FetchArticles fetches all articles for a Qiita user
func (s *QiitaService) FetchArticles(username string) ([]model.QiitaArticle, error) {
	var allArticles []model.QiitaArticle
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://qiita.com/api/v2/users/%s/items?page=%d&per_page=%d", username, page, perPage)

		resp, err := s.httpClient.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Qiita articles: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("Qiita user not found")
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Qiita API returned status %d", resp.StatusCode)
		}

		var apiArticles []QiitaAPIArticle
		if err := json.NewDecoder(resp.Body).Decode(&apiArticles); err != nil {
			return nil, fmt.Errorf("failed to decode Qiita response: %w", err)
		}

		for _, article := range apiArticles {
			// Extract tag names
			tagNames := make([]string, len(article.Tags))
			for i, tag := range article.Tags {
				tagNames[i] = tag.Name
			}

			allArticles = append(allArticles, model.QiitaArticle{
				QiitaID:       article.ID,
				Title:         article.Title,
				URL:           article.URL,
				LikesCount:    article.LikesCount,
				CommentsCount: article.CommentsCount,
				Tags:          strings.Join(tagNames, ","),
				PublishedAt:   article.CreatedAt,
			})
		}

		// Check if there are more pages
		if len(apiArticles) < perPage {
			break
		}
		page++
	}

	return allArticles, nil
}

// ValidateUsername checks if a Qiita username exists
func (s *QiitaService) ValidateUsername(username string) (bool, error) {
	url := fmt.Sprintf("https://qiita.com/api/v2/users/%s", username)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// Connect sets the Qiita username and syncs articles.
func (s *QiitaService) Connect(userID uint, username string) (int, error) {
	valid, err := s.ValidateUsername(username)
	if err != nil || !valid {
		return 0, ErrBadRequest
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, ErrNotFound
	}

	user.QiitaUsername = username
	if err := s.userRepo.Update(user); err != nil {
		return 0, err
	}

	articles, err := s.FetchArticles(username)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	for i := range articles {
		articles[i].UpdatedAt = now
	}

	if err := s.qiitaRepo.UpsertArticles(userID, articles); err != nil {
		return 0, err
	}

	return len(articles), nil
}

// Disconnect removes the Qiita username and deletes cached articles.
func (s *QiitaService) Disconnect(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound
	}

	user.QiitaUsername = ""
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return s.qiitaRepo.DeleteUserArticles(userID)
}

// Sync refreshes the Qiita articles for the user.
func (s *QiitaService) Sync(userID uint) (int, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, ErrNotFound
	}

	if user.QiitaUsername == "" {
		return 0, ErrBadRequest
	}

	articles, err := s.FetchArticles(user.QiitaUsername)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	for i := range articles {
		articles[i].UpdatedAt = now
	}

	if err := s.qiitaRepo.UpsertArticles(userID, articles); err != nil {
		return 0, err
	}

	return len(articles), nil
}

// GetArticles returns all Qiita articles for a user.
func (s *QiitaService) GetArticles(userID uint) ([]model.QiitaArticle, error) {
	return s.qiitaRepo.GetArticles(userID)
}

// GetStats returns Qiita statistics for a user.
func (s *QiitaService) GetStats(userID uint) (*model.QiitaStats, error) {
	return s.qiitaRepo.GetStats(userID)
}
