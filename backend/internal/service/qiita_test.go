package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

// newTestQiitaService はテスト用のQiitaServiceを生成する。
func newTestQiitaService(fn roundTripFunc) *QiitaService {
	return &QiitaService{
		httpClient: &http.Client{Transport: fn},
	}
}

func TestFetchArticles_Success(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		body := `[{"id":"abc123","title":"Go入門","url":"https://qiita.com/test/items/abc123","likes_count":10,"comments_count":2,"tags":[{"name":"Go"},{"name":"初心者"}],"created_at":"2025-01-15T10:00:00+09:00"}]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
	assert.Equal(t, "abc123", articles[0].QiitaID)
	assert.Equal(t, "Go入門", articles[0].Title)
	assert.Equal(t, 10, articles[0].LikesCount)
	assert.Equal(t, "Go,初心者", articles[0].Tags)
}

func TestFetchArticles_Empty(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("[]")),
		}, nil
	})

	articles, err := svc.FetchArticles("newuser")
	assert.NoError(t, err)
	assert.Empty(t, articles)
}

func TestFetchArticles_NotFound(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	articles, err := svc.FetchArticles("unknown")
	assert.Nil(t, articles)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
}

func TestFetchArticles_ServerError(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
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

func TestFetchArticles_NetworkError(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	articles, err := svc.FetchArticles("testuser")
	assert.Nil(t, articles)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestFetchArticles_InvalidJSON(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("invalid json")),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.Nil(t, articles)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestQiitaValidateUsername_Valid(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{}")),
		}, nil
	})

	valid, err := svc.ValidateUsername("validuser")
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestQiitaValidateUsername_Invalid(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	valid, err := svc.ValidateUsername("invaliduser")
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestQiitaValidateUsername_NetworkError(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("timeout")
	})

	valid, err := svc.ValidateUsername("testuser")
	assert.Error(t, err)
	assert.False(t, valid)
}

// ============================================================
// エッジケーステスト
// ============================================================

func TestFetchArticles_MultipleArticles(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		body := `[
			{"id":"a1","title":"記事1","url":"https://qiita.com/u/items/a1","likes_count":5,"comments_count":1,"tags":[{"name":"Go"}],"created_at":"2025-01-01T00:00:00+09:00"},
			{"id":"a2","title":"記事2","url":"https://qiita.com/u/items/a2","likes_count":0,"comments_count":0,"tags":[],"created_at":"2025-02-01T00:00:00+09:00"},
			{"id":"a3","title":"記事3","url":"https://qiita.com/u/items/a3","likes_count":100,"comments_count":50,"tags":[{"name":"React"},{"name":"TypeScript"},{"name":"初心者"}],"created_at":"2025-03-01T00:00:00+09:00"}
		]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.NoError(t, err)
	assert.Len(t, articles, 3)
	assert.Equal(t, "", articles[1].Tags)
	assert.Equal(t, "React,TypeScript,初心者", articles[2].Tags)
	assert.Equal(t, 100, articles[2].LikesCount)
}

func TestFetchArticles_RequestURL(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		assert.Contains(t, req.URL.String(), "/users/myuser/items")
		assert.Contains(t, req.URL.String(), "page=1")
		assert.Contains(t, req.URL.String(), "per_page=100")

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("[]")),
		}, nil
	})

	_, err := svc.FetchArticles("myuser")
	assert.NoError(t, err)
}

func TestFetchArticles_RateLimited(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.Nil(t, articles)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
	assert.Contains(t, domainErr.Message, "429")
}

func TestFetchArticles_NoTags(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		body := `[{"id":"notag","title":"タグなし記事","url":"https://qiita.com/u/items/notag","likes_count":0,"comments_count":0,"tags":[],"created_at":"2025-01-01T00:00:00+09:00"}]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
	assert.Equal(t, "", articles[0].Tags)
	assert.Equal(t, "notag", articles[0].QiitaID)
}

func TestQiitaValidateUsername_ServerError(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	valid, err := svc.ValidateUsername("testuser")
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestQiitaValidateUsername_RequestURL(t *testing.T) {
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		assert.Contains(t, req.URL.String(), "/users/checkuser")

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{}")),
		}, nil
	})

	valid, err := svc.ValidateUsername("checkuser")
	assert.NoError(t, err)
	assert.True(t, valid)
}
