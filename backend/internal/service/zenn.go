package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ZennService はZenn記事連携のビジネスロジックを提供する。
// Zenn APIからの記事取得、ユーザー連携、データ同期を担当する。
type ZennService struct {
	httpClient *http.Client
	userRepo   repository.UserRepositoryInterface
	zennRepo   repository.ZennRepositoryInterface
}

// NewZennService は新しいZennServiceインスタンスを生成する。
func NewZennService(userRepo repository.UserRepositoryInterface, zennRepo repository.ZennRepositoryInterface) *ZennService {
	return &ZennService{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userRepo:   userRepo,
		zennRepo:   zennRepo,
	}
}

// ZennAPIResponse はZenn APIのレスポンス構造を表す。
type ZennAPIResponse struct {
	Articles []ZennAPIArticle `json:"articles"`
	NextPage *int             `json:"next_page"` // 次ページが存在しない場合はnil
}

// ZennAPIArticle はZenn APIから取得する記事データを表す。
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

// FetchArticles はZenn APIから指定ユーザーの全記事を取得する（ページネーション対応）。
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

		// 次ページが存在しない場合は終了
		if apiResp.NextPage == nil {
			break
		}
		page = *apiResp.NextPage
	}

	return allArticles, nil
}

// ValidateUsername はZennユーザー名が存在するかを検証する。
func (s *ZennService) ValidateUsername(username string) (bool, error) {
	url := fmt.Sprintf("https://zenn.dev/api/articles?username=%s&page=1", username)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// Connect はZennユーザー名を設定し、記事データを同期する。
// 同期した記事数を返す。
func (s *ZennService) Connect(userID uint, username string) (int, error) {
	valid, err := s.ValidateUsername(username)
	if err != nil || !valid {
		return 0, ErrBadRequest
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, ErrNotFound
	}

	user.ZennUsername = username
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

	if err := s.zennRepo.UpsertArticles(userID, articles); err != nil {
		return 0, err
	}

	return len(articles), nil
}

// Disconnect はZenn連携を解除し、キャッシュされた記事データを削除する。
func (s *ZennService) Disconnect(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound
	}

	user.ZennUsername = ""
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return s.zennRepo.DeleteUserArticles(userID)
}

// Sync はZenn記事データを最新の状態に同期する。同期した記事数を返す。
func (s *ZennService) Sync(userID uint) (int, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, ErrNotFound
	}

	if user.ZennUsername == "" {
		return 0, ErrBadRequest
	}

	articles, err := s.FetchArticles(user.ZennUsername)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	for i := range articles {
		articles[i].UpdatedAt = now
	}

	if err := s.zennRepo.UpsertArticles(userID, articles); err != nil {
		return 0, err
	}

	return len(articles), nil
}

// GetArticles は指定ユーザーの全Zenn記事を取得する。
func (s *ZennService) GetArticles(userID uint) ([]model.ZennArticle, error) {
	return s.zennRepo.GetArticles(userID)
}

// GetStats は指定ユーザーのZenn統計情報を取得する。
func (s *ZennService) GetStats(userID uint) (*model.ZennStats, error) {
	return s.zennRepo.GetStats(userID)
}
