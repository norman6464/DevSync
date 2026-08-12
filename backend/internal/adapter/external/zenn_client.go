package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// zennRequestTimeout は Zenn への 1 リクエストのタイムアウト。
const zennRequestTimeout = 30 * time.Second

// zennArticlesURL は記事一覧を返す Zenn のエンドポイント。
const zennArticlesURL = "https://zenn.dev/api/articles"

// zennAPIResponse は Zenn API のレスポンス構造。
type zennAPIResponse struct {
	Articles []zennAPIArticle `json:"articles"`
	NextPage *int             `json:"next_page"` // 次ページが存在しない場合は nil
}

// zennAPIArticle は Zenn API から取得する記事データ。
type zennAPIArticle struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Emoji         string    `json:"emoji"`
	ArticleType   string    `json:"article_type"`
	LikedCount    int       `json:"liked_count"`
	CommentsCount int       `json:"comments_count"`
	PublishedAt   time.Time `json:"published_at"`
}

// zennClient は [repository.ZennArticleFetcher] の HTTP 実装。
type zennClient struct {
	client *http.Client
}

// NewZennClient は ZennArticleFetcher の HTTP 実装を返す。
func NewZennClient() repository.ZennArticleFetcher {
	return &zennClient{client: &http.Client{Timeout: zennRequestTimeout}}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ZennArticleFetcher = (*zennClient)(nil)

// FetchArticles は指定ユーザーの記事を、次ページが無くなるまでたどって取得する。
func (c *zennClient) FetchArticles(ctx context.Context, username string) ([]model.ZennArticle, error) {
	var all []model.ZennArticle

	for page := 1; ; {
		apiResp, err := c.fetchPage(ctx, username, page)
		if err != nil {
			return nil, err
		}

		for _, article := range apiResp.Articles {
			all = append(all, model.ZennArticle{
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

		if apiResp.NextPage == nil {
			break
		}
		page = *apiResp.NextPage
	}

	return all, nil
}

// fetchPage は 1 ページ分の記事を取得する。レスポンスの本文はこの中で閉じる。
func (c *zennClient) fetchPage(ctx context.Context, username string, page int) (*zennAPIResponse, error) {
	query := url.Values{}
	query.Set("username", username)
	query.Set("order", "latest")
	query.Set("page", fmt.Sprint(page))

	resp, err := c.get(ctx, query)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Zenn記事の取得に失敗", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, fmt.Sprintf("Zenn APIエラー: ステータスコード %d", resp.StatusCode), nil)
	}

	var apiResp zennAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Zennレスポンスのデコードに失敗", err)
	}
	return &apiResp, nil
}

// UserExists は記事一覧の 1 ページ目を取得できるかどうかを返す。
func (c *zennClient) UserExists(ctx context.Context, username string) (bool, error) {
	query := url.Values{}
	query.Set("username", username)
	query.Set("page", "1")

	resp, err := c.get(ctx, query)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// get は記事一覧のエンドポイントへ GET する。
func (c *zennClient) get(ctx context.Context, query url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, zennArticlesURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}
