package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
)

type ZennService struct {
	httpClient *http.Client
}

func NewZennService() *ZennService {
	return &ZennService{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ZennAPIResponse represents the response from Zenn API
type ZennAPIResponse struct {
	Articles []ZennAPIArticle `json:"articles"`
	NextPage *int             `json:"next_page"`
}

// ZennAPIArticle represents an article from Zenn API
type ZennAPIArticle struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Emoji         string    `json:"emoji"`
	ArticleType   string    `json:"article_type"`
	LikedCount    int       `json:"liked_count"`
	CommentsCount int       `json:"comments_count"`
	PublishedAt   time.Time `json:"published_at"`
}

// FetchArticles fetches all articles for a Zenn user
func (s *ZennService) FetchArticles(username string) ([]model.ZennArticle, error) {
	var allArticles []model.ZennArticle
	page := 1

	for {
		url := fmt.Sprintf("https://zenn.dev/api/articles?username=%s&order=latest&page=%d", username, page)

		resp, err := s.httpClient.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Zenn articles: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Zenn API returned status %d", resp.StatusCode)
		}

		var apiResp ZennAPIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return nil, fmt.Errorf("failed to decode Zenn response: %w", err)
		}

		for _, article := range apiResp.Articles {
			allArticles = append(allArticles, model.ZennArticle{
				ZennID:        article.ID,
				Title:         article.Title,
				Slug:          article.Slug,
				Emoji:         article.Emoji,
				ArticleType:   article.ArticleType,
				LikedCount:    article.LikedCount,
				CommentsCount: article.CommentsCount,
				PublishedAt:   article.PublishedAt,
			})
		}

		// Check if there are more pages
		if apiResp.NextPage == nil {
			break
		}
		page = *apiResp.NextPage
	}

	return allArticles, nil
}

// ValidateUsername checks if a Zenn username exists
func (s *ZennService) ValidateUsername(username string) (bool, error) {
	url := fmt.Sprintf("https://zenn.dev/api/articles?username=%s&page=1", username)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
