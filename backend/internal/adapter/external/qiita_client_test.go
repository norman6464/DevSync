package external

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestQiitaClient は任意のレスポンスを返すクライアントを生成する。
func newTestQiitaClient(fn roundTripFunc) *qiitaClient {
	return &qiitaClient{client: &http.Client{Transport: fn}}
}

// qiitaArticlesJSON は指定件数の記事を持つレスポンス本文を作る。
func qiitaArticlesJSON(count int) string {
	items := make([]string, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, fmt.Sprintf(`{"id":"item-%d","title":"記事%d","url":"https://qiita.com/items/%d","likes_count":%d,"comments_count":0,"tags":[{"name":"go"}],"created_at":"2026-01-01T00:00:00Z"}`, i, i, i, i))
	}
	return "[" + strings.Join(items, ",") + "]"
}

func TestQiitaClient_FetchArticles_SinglePage(t *testing.T) {
	var requestedPath, requestedQuery string
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		requestedPath = req.URL.EscapedPath()
		requestedQuery = req.URL.Query().Encode()
		body := `[{"id":"abc","title":"記事A","url":"https://qiita.com/items/abc","likes_count":7,"comments_count":2,"tags":[{"name":"go"},{"name":"api"}],"created_at":"2026-01-01T00:00:00Z"}]`
		return jsonResponse(http.StatusOK, body), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	require.NoError(t, err)
	require.Len(t, articles, 1)
	assert.Equal(t, "abc", articles[0].QiitaID)
	assert.Equal(t, "記事A", articles[0].Title)
	assert.Equal(t, "https://qiita.com/items/abc", articles[0].URL)
	assert.Equal(t, 7, articles[0].LikesCount)
	assert.Equal(t, 2, articles[0].CommentsCount)
	// タグ名はカンマ区切りに連結する。
	assert.Equal(t, "go,api", articles[0].Tags)
	assert.Equal(t, "/api/v2/users/testuser/items", requestedPath)
	assert.Equal(t, "page=1&per_page=100", requestedQuery)
}

// 1 ページ分ちょうど取得できた場合は次ページも取りに行く。
func TestQiitaClient_FetchArticles_Pagination(t *testing.T) {
	var pages []string
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		page := req.URL.Query().Get("page")
		pages = append(pages, page)
		if page == "1" {
			return jsonResponse(http.StatusOK, qiitaArticlesJSON(100)), nil
		}
		return jsonResponse(http.StatusOK, qiitaArticlesJSON(3)), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	require.NoError(t, err)
	assert.Len(t, articles, 103)
	assert.Equal(t, []string{"1", "2"}, pages)
}

// タグが無い記事は空文字になる。
func TestQiitaClient_FetchArticles_NoTags(t *testing.T) {
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `[{"id":"abc","title":"タグなし","tags":[]}]`), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	require.NoError(t, err)
	require.Len(t, articles, 1)
	assert.Empty(t, articles[0].Tags)
}

func TestQiitaClient_FetchArticles_Empty(t *testing.T) {
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, "[]"), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	require.NoError(t, err)
	assert.Empty(t, articles)
}

func TestQiitaClient_FetchArticles_NotFound(t *testing.T) {
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusNotFound, ""), nil
	})

	_, err := c.FetchArticles(context.Background(), "ghost")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
}

func TestQiitaClient_FetchArticles_ServerError(t *testing.T) {
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusInternalServerError, ""), nil
	})

	_, err := c.FetchArticles(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
	assert.Contains(t, domainErr.Message, "500")
}

func TestQiitaClient_FetchArticles_NetworkError(t *testing.T) {
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	_, err := c.FetchArticles(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestQiitaClient_FetchArticles_InvalidJSON(t *testing.T) {
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, "invalid json"), nil
	})

	_, err := c.FetchArticles(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

// 2 ページ目で失敗したらエラーを返す。
func TestQiitaClient_FetchArticles_ErrorOnSecondPage(t *testing.T) {
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("page") == "1" {
			return jsonResponse(http.StatusOK, qiitaArticlesJSON(100)), nil
		}
		return jsonResponse(http.StatusInternalServerError, ""), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	assert.Nil(t, articles)
	require.Error(t, err)
}

// ユーザー名はパスに埋め込む前にエスケープする。
func TestQiitaClient_FetchArticles_EscapesUsername(t *testing.T) {
	var requestedPath string
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		requestedPath = req.URL.EscapedPath()
		return jsonResponse(http.StatusOK, "[]"), nil
	})

	_, err := c.FetchArticles(context.Background(), "a/../b")
	require.NoError(t, err)
	assert.Equal(t, "/api/v2/users/a%2F..%2Fb/items", requestedPath)
}

func TestQiitaClient_FetchArticles_PropagatesContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey{}, "propagated")

	var got interface{}
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		got = req.Context().Value(ctxKey{})
		return jsonResponse(http.StatusOK, "[]"), nil
	})

	_, err := c.FetchArticles(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, "propagated", got)
}

func TestQiitaClient_UserExists(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"200 なら存在する", http.StatusOK, true},
		{"404 なら存在しない", http.StatusNotFound, false},
		{"500 なら存在しない扱い", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestedPath string
			c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
				requestedPath = req.URL.EscapedPath()
				return jsonResponse(tt.status, ""), nil
			})
			exists, err := c.UserExists(context.Background(), "testuser")
			require.NoError(t, err)
			assert.Equal(t, tt.want, exists)
			// 記事一覧ではなくユーザーのエンドポイントを見る。
			assert.Equal(t, "/api/v2/users/testuser", requestedPath)
		})
	}
}

func TestQiitaClient_UserExists_NetworkError(t *testing.T) {
	c := newTestQiitaClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	exists, err := c.UserExists(context.Background(), "testuser")
	assert.False(t, exists)
	assert.Error(t, err)
}
