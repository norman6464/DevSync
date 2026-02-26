package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

// newTestZennService はテスト用のZennServiceを生成する。
func newTestZennService(fn roundTripFunc) *ZennService {
	return &ZennService{
		httpClient: &http.Client{Transport: fn},
	}
}

func TestZennFetchArticles_Success(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		body := `{"articles":[{"id":12345,"title":"Go入門ガイド","slug":"go-intro","emoji":"🐹","article_type":"tech","liked_count":42,"comments_count":5,"published_at":"2025-01-15T10:00:00Z"}],"next_page":null}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
	assert.Equal(t, int64(12345), articles[0].ZennID)
	assert.Equal(t, "Go入門ガイド", articles[0].Title)
	assert.Equal(t, "go-intro", articles[0].Slug)
	assert.Equal(t, "🐹", articles[0].Emoji)
	assert.Equal(t, "tech", articles[0].ArticleType)
	assert.Equal(t, 42, articles[0].LikedCount)
	assert.Equal(t, 5, articles[0].CommentsCount)
}

func TestZennFetchArticles_Pagination(t *testing.T) {
	callCount := 0
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		callCount++
		var body string
		if callCount == 1 {
			body = `{"articles":[{"id":1,"title":"記事1","slug":"a1","emoji":"📝","article_type":"tech","liked_count":10,"comments_count":1,"published_at":"2025-01-10T00:00:00Z"}],"next_page":2}`
		} else {
			body = `{"articles":[{"id":2,"title":"記事2","slug":"a2","emoji":"📖","article_type":"idea","liked_count":5,"comments_count":0,"published_at":"2025-01-11T00:00:00Z"}],"next_page":null}`
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.NoError(t, err)
	assert.Len(t, articles, 2)
	assert.Equal(t, "記事1", articles[0].Title)
	assert.Equal(t, "記事2", articles[1].Title)
	assert.Equal(t, 2, callCount)
}

func TestZennFetchArticles_Empty(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"articles":[],"next_page":null}`)),
		}, nil
	})

	articles, err := svc.FetchArticles("newuser")
	assert.NoError(t, err)
	assert.Empty(t, articles)
}

func TestZennFetchArticles_ServerError(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.Nil(t, articles)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestZennFetchArticles_NetworkError(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	articles, err := svc.FetchArticles("testuser")
	assert.Nil(t, articles)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestZennFetchArticles_InvalidJSON(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("not json")),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.Nil(t, articles)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestZennFetchArticles_StatusErrorMessage(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	_, err := svc.FetchArticles("testuser")
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Contains(t, domainErr.Message, fmt.Sprintf("%d", http.StatusTooManyRequests))
}

func TestZennValidateUsername_Valid(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"articles":[]}`)),
		}, nil
	})

	valid, err := svc.ValidateUsername("validuser")
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestZennValidateUsername_Invalid(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	valid, err := svc.ValidateUsername("invaliduser")
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestZennValidateUsername_NetworkError(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("timeout")
	})

	valid, err := svc.ValidateUsername("testuser")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestZennFetchArticles_RequestURL(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "zenn.dev", req.URL.Host)
		assert.Equal(t, "/api/articles", req.URL.Path)
		assert.Equal(t, "myuser", req.URL.Query().Get("username"))
		assert.Equal(t, "1", req.URL.Query().Get("page"))
		assert.Equal(t, "latest", req.URL.Query().Get("order"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"articles":[],"next_page":null}`)),
		}, nil
	})

	_, err := svc.FetchArticles("myuser")
	assert.NoError(t, err)
}

func TestZennFetchArticles_MultipleArticles(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		body := `{"articles":[
			{"id":1,"title":"記事A","slug":"a","emoji":"📝","article_type":"tech","liked_count":10,"comments_count":2,"published_at":"2025-01-01T00:00:00Z"},
			{"id":2,"title":"記事B","slug":"b","emoji":"📖","article_type":"idea","liked_count":20,"comments_count":4,"published_at":"2025-01-02T00:00:00Z"},
			{"id":3,"title":"記事C","slug":"c","emoji":"🔧","article_type":"tech","liked_count":30,"comments_count":6,"published_at":"2025-01-03T00:00:00Z"}
		],"next_page":null}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.NoError(t, err)
	assert.Len(t, articles, 3)
	assert.Equal(t, "記事A", articles[0].Title)
	assert.Equal(t, "記事B", articles[1].Title)
	assert.Equal(t, "記事C", articles[2].Title)
	assert.Equal(t, 30, articles[2].LikedCount)
}

func TestZennFetchArticles_NotFound(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	articles, err := svc.FetchArticles("nonexistent")
	assert.Nil(t, articles)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Contains(t, domainErr.Message, "404")
}

func TestZennValidateUsername_ServerError(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	valid, err := svc.ValidateUsername("testuser")
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestZennFetchArticles_PaginationURL(t *testing.T) {
	callCount := 0
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			assert.Equal(t, "1", req.URL.Query().Get("page"))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"articles":[],"next_page":3}`)),
			}, nil
		}
		assert.Equal(t, "3", req.URL.Query().Get("page"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"articles":[],"next_page":null}`)),
		}, nil
	})

	_, err := svc.FetchArticles("testuser")
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestZennFetchArticles_InvalidUsername(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		t.Fatal("HTTPリクエストが送信されるべきではない")
		return nil, nil
	})

	tests := []struct {
		name     string
		username string
	}{
		{"スラッシュ含む", "user/path"},
		{"アンパサンド含む", "user&param=1"},
		{"スペース含む", "user name"},
		{"空文字", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articles, err := svc.FetchArticles(tt.username)
			assert.Nil(t, articles)
			assert.Error(t, err)
		})
	}
}

func TestZennValidateUsername_InvalidFormat(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		t.Fatal("HTTPリクエストが送信されるべきではない")
		return nil, nil
	})

	valid, err := svc.ValidateUsername("user/injection")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestZennValidateUsername_RequestURL(t *testing.T) {
	svc := newTestZennService(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "zenn.dev", req.URL.Host)
		assert.Equal(t, "/api/articles", req.URL.Path)
		assert.Equal(t, "zennuser", req.URL.Query().Get("username"))
		assert.Equal(t, "1", req.URL.Query().Get("page"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"articles":[]}`)),
		}, nil
	})

	valid, err := svc.ValidateUsername("zennuser")
	assert.NoError(t, err)
	assert.True(t, valid)
}
