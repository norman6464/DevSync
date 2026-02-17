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
