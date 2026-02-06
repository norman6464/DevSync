package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
)

type QiitaService struct {
	httpClient *http.Client
}

func NewQiitaService() *QiitaService {
	return &QiitaService{
		httpClient: &http.Client{Timeout: 30 * time.Second},
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
