package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// QiitaService はQiita連携のビジネスロジックを提供する。
// Qiita APIを通じた記事取得・ユーザー検証・連携管理を行う。
type QiitaService struct {
	httpClient *http.Client                      // Qiita API呼び出し用HTTPクライアント
	userRepo   repository.UserRepositoryInterface // ユーザーリポジトリ
	qiitaRepo  repository.QiitaRepositoryInterface // Qiita記事リポジトリ
}

// NewQiitaService は新しいQiitaServiceインスタンスを生成する。
// HTTPクライアントは30秒タイムアウトで初期化される。
func NewQiitaService(userRepo repository.UserRepositoryInterface, qiitaRepo repository.QiitaRepositoryInterface) *QiitaService {
	return &QiitaService{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userRepo:   userRepo,
		qiitaRepo:  qiitaRepo,
	}
}

// QiitaAPIArticle はQiita APIから返される記事のレスポンス構造体。
type QiitaAPIArticle struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	URL           string          `json:"url"`
	LikesCount    int             `json:"likes_count"`
	CommentsCount int             `json:"comments_count"`
	Tags          []QiitaAPITag   `json:"tags"`
	CreatedAt     time.Time       `json:"created_at"`
}

// QiitaAPITag はQiita APIのタグ情報を表す。
type QiitaAPITag struct {
	Name string `json:"name"`
}

// FetchArticles は指定ユーザーのQiita記事を全件取得する。
// ページネーションにより100件ずつ取得し、全ページを結合して返す。
func (s *QiitaService) FetchArticles(username string) ([]model.QiitaArticle, error) {
	var allArticles []model.QiitaArticle
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://qiita.com/api/v2/users/%s/items?page=%d&per_page=%d", username, page, perPage)

		resp, err := s.httpClient.Get(url)
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Qiita記事の取得に失敗", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return nil, domain.NewError(domain.ErrCodeNotFound, "Qiitaユーザーが見つかりません", nil)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, domain.NewError(domain.ErrCodeServiceUnavailable, fmt.Sprintf("Qiita APIエラー: ステータスコード %d", resp.StatusCode), nil)
		}

		var apiArticles []QiitaAPIArticle
		if err := json.NewDecoder(resp.Body).Decode(&apiArticles); err != nil {
			return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Qiitaレスポンスのデコードに失敗", err)
		}

		for _, article := range apiArticles {
			// タグ名を抽出してカンマ区切り文字列に変換
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

		// 取得件数がperPage未満なら最終ページ
		if len(apiArticles) < perPage {
			break
		}
		page++
	}

	return allArticles, nil
}

// ValidateUsername はQiitaユーザー名が存在するかを検証する。
func (s *QiitaService) ValidateUsername(username string) (bool, error) {
	url := fmt.Sprintf("https://qiita.com/api/v2/users/%s", username)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// Connect はQiitaアカウントを連携する。
// ユーザー名を検証後、プロフィールに保存し、記事をUpsertで同期する。
// 同期した記事数を返す。
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

// Disconnect はQiitaアカウント連携を解除する。
// ユーザー名をクリアし、キャッシュ済み記事を削除する。
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

// Sync はQiita記事を最新状態に同期する。
// 既存の連携ユーザー名を使用して記事を再取得し、Upsertする。
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

// GetArticles は指定ユーザーのQiita記事一覧を返す。
func (s *QiitaService) GetArticles(userID uint) ([]model.QiitaArticle, error) {
	return s.qiitaRepo.GetArticles(userID)
}

// GetStats は指定ユーザーのQiita統計情報を返す。
func (s *QiitaService) GetStats(userID uint) (*model.QiitaStats, error) {
	return s.qiitaRepo.GetStats(userID)
}
