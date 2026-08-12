package external

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestZennClient は任意のレスポンスを返すクライアントを生成する。
func newTestZennClient(fn roundTripFunc) *zennClient {
	return &zennClient{client: &http.Client{Transport: fn}}
}

func TestZennClient_FetchArticles_SinglePage(t *testing.T) {
	var requestedQuery string
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		requestedQuery = req.URL.Query().Encode()
		body := `{"articles":[{"id":1,"title":"記事A","slug":"a","emoji":"🐱","article_type":"tech","liked_count":5,"comments_count":2,"published_at":"2026-01-01T00:00:00Z"}],"next_page":null}`
		return jsonResponse(http.StatusOK, body), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	require.NoError(t, err)
	require.Len(t, articles, 1)
	assert.Equal(t, int64(1), articles[0].ZennID)
	assert.Equal(t, "記事A", articles[0].Title)
	assert.Equal(t, "a", articles[0].Slug)
	assert.Equal(t, "tech", articles[0].ArticleType)
	assert.Equal(t, 5, articles[0].LikedCount)
	assert.Equal(t, 2, articles[0].CommentsCount)
	assert.Equal(t, "order=latest&page=1&username=testuser", requestedQuery)
}

// next_page をたどって全ページ分を取得する。
func TestZennClient_FetchArticles_Pagination(t *testing.T) {
	var pages []string
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		page := req.URL.Query().Get("page")
		pages = append(pages, page)
		if page == "1" {
			return jsonResponse(http.StatusOK, `{"articles":[{"id":1,"title":"1件目"}],"next_page":2}`), nil
		}
		return jsonResponse(http.StatusOK, `{"articles":[{"id":2,"title":"2件目"}],"next_page":null}`), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	require.NoError(t, err)
	require.Len(t, articles, 2)
	assert.Equal(t, int64(1), articles[0].ZennID)
	assert.Equal(t, int64(2), articles[1].ZennID)
	assert.Equal(t, []string{"1", "2"}, pages)
}

// 記事が 1 件も無ければ空で返す。
func TestZennClient_FetchArticles_Empty(t *testing.T) {
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"articles":[],"next_page":null}`), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	require.NoError(t, err)
	assert.Empty(t, articles)
}

func TestZennClient_FetchArticles_ServerError(t *testing.T) {
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusInternalServerError, ""), nil
	})

	_, err := c.FetchArticles(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
	assert.Contains(t, domainErr.Message, "500")
}

func TestZennClient_FetchArticles_NetworkError(t *testing.T) {
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	_, err := c.FetchArticles(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestZennClient_FetchArticles_InvalidJSON(t *testing.T) {
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, "invalid json"), nil
	})

	_, err := c.FetchArticles(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

// 2 ページ目で失敗したらエラーを返す。
func TestZennClient_FetchArticles_ErrorOnSecondPage(t *testing.T) {
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("page") == "1" {
			return jsonResponse(http.StatusOK, `{"articles":[{"id":1}],"next_page":2}`), nil
		}
		return jsonResponse(http.StatusInternalServerError, ""), nil
	})

	articles, err := c.FetchArticles(context.Background(), "testuser")
	assert.Nil(t, articles)
	require.Error(t, err)
}

// ユーザー名はクエリとしてエスケープされる。
func TestZennClient_FetchArticles_EscapesUsername(t *testing.T) {
	var rawQuery, username string
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		rawQuery = req.URL.RawQuery
		username = req.URL.Query().Get("username")
		return jsonResponse(http.StatusOK, `{"articles":[],"next_page":null}`), nil
	})

	_, err := c.FetchArticles(context.Background(), "a&page=99")
	require.NoError(t, err)
	assert.Contains(t, rawQuery, "username=a%26page%3D99")
	assert.Equal(t, "a&page=99", username)
}

func TestZennClient_FetchArticles_PropagatesContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey{}, "propagated")

	var got interface{}
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		got = req.Context().Value(ctxKey{})
		return jsonResponse(http.StatusOK, `{"articles":[],"next_page":null}`), nil
	})

	_, err := c.FetchArticles(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, "propagated", got)
}

func TestZennClient_UserExists(t *testing.T) {
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
			c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
				return jsonResponse(tt.status, ""), nil
			})
			exists, err := c.UserExists(context.Background(), "testuser")
			require.NoError(t, err)
			assert.Equal(t, tt.want, exists)
		})
	}
}

func TestZennClient_UserExists_NetworkError(t *testing.T) {
	c := newTestZennClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	exists, err := c.UserExists(context.Background(), "testuser")
	assert.False(t, exists)
	assert.Error(t, err)
}
