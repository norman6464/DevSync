package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// qiitaRequestTimeout は Qiita への 1 リクエストのタイムアウト。
const qiitaRequestTimeout = 30 * time.Second

// qiitaPerPage は記事一覧の 1 ページあたりの取得件数。
const qiitaPerPage = 100

// Qiita のエンドポイント。
const (
	qiitaItemsURL = "https://qiita.com/api/v2/users/%s/items"
	qiitaUserURL  = "https://qiita.com/api/v2/users/%s"
)

// qiitaAPIArticle は Qiita API から返される記事データ。
type qiitaAPIArticle struct {
	ID            string        `json:"id"`
	Title         string        `json:"title"`
	URL           string        `json:"url"`
	LikesCount    int           `json:"likes_count"`
	CommentsCount int           `json:"comments_count"`
	Tags          []qiitaAPITag `json:"tags"`
	CreatedAt     time.Time     `json:"created_at"`
}

// qiitaAPITag は Qiita API のタグ情報。
type qiitaAPITag struct {
	Name string `json:"name"`
}

// qiitaClient は [repository.QiitaArticleFetcher] の HTTP 実装。
type qiitaClient struct {
	client *http.Client
}

// NewQiitaClient は QiitaArticleFetcher の HTTP 実装を返す。
func NewQiitaClient() repository.QiitaArticleFetcher {
	return &qiitaClient{client: &http.Client{Timeout: qiitaRequestTimeout}}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.QiitaArticleFetcher = (*qiitaClient)(nil)

// FetchArticles は指定ユーザーの記事を、取得件数が 1 ページ分に満たなくなるまでたどって取得する。
func (c *qiitaClient) FetchArticles(ctx context.Context, username string) ([]model.QiitaArticle, error) {
	var all []model.QiitaArticle

	for page := 1; ; page++ {
		apiArticles, err := c.fetchPage(ctx, username, page)
		if err != nil {
			return nil, err
		}

		for _, article := range apiArticles {
			all = append(all, model.QiitaArticle{
				QiitaID:       article.ID,
				Title:         article.Title,
				URL:           article.URL,
				LikesCount:    article.LikesCount,
				CommentsCount: article.CommentsCount,
				Tags:          joinQiitaTags(article.Tags),
				PublishedAt:   article.CreatedAt,
			})
		}

		if len(apiArticles) < qiitaPerPage {
			break
		}
	}

	return all, nil
}

// joinQiitaTags はタグ名をカンマ区切りに連結する。
func joinQiitaTags(tags []qiitaAPITag) string {
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	return strings.Join(names, ",")
}

// fetchPage は 1 ページ分の記事を取得する。レスポンスの本文はこの中で閉じる。
func (c *qiitaClient) fetchPage(ctx context.Context, username string, page int) ([]qiitaAPIArticle, error) {
	query := url.Values{}
	query.Set("page", fmt.Sprint(page))
	query.Set("per_page", fmt.Sprint(qiitaPerPage))
	endpoint := fmt.Sprintf(qiitaItemsURL, url.PathEscape(username)) + "?" + query.Encode()

	resp, err := c.get(ctx, endpoint)
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

	var apiArticles []qiitaAPIArticle
	if err := json.NewDecoder(resp.Body).Decode(&apiArticles); err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Qiitaレスポンスのデコードに失敗", err)
	}
	return apiArticles, nil
}

// UserExists はユーザーページを取得できるかどうかを返す。
func (c *qiitaClient) UserExists(ctx context.Context, username string) (bool, error) {
	resp, err := c.get(ctx, fmt.Sprintf(qiitaUserURL, url.PathEscape(username)))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// get は指定 URL へ GET する。
func (c *qiitaClient) get(ctx context.Context, endpoint string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}
