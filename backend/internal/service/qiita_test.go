package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestQiitaService はテスト用のQiitaServiceを生成する。
func newTestQiitaService(fn roundTripFunc) *QiitaService {
	return &QiitaService{
		httpClient: &http.Client{Transport: fn},
	}
}

// newTestQiitaServiceWithRepos はリポジトリ付きのQiitaServiceテスト用インスタンスを生成する。
func newTestQiitaServiceWithRepos(fn roundTripFunc) (*QiitaService, *MockUserRepository, *MockQiitaRepository) {
	userRepo := new(MockUserRepository)
	qiitaRepo := new(MockQiitaRepository)
	return &QiitaService{
		httpClient: &http.Client{Transport: fn},
		userRepo:   userRepo,
		qiitaRepo:  qiitaRepo,
	}, userRepo, qiitaRepo
}

func TestNewQiitaService(t *testing.T) {
	userRepo := new(MockUserRepository)
	qiitaRepo := new(MockQiitaRepository)
	svc := NewQiitaService(userRepo, qiitaRepo)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.httpClient)
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

// ============================================================
// Disconnect テスト
// ============================================================

func TestQiitaDisconnect_Success(t *testing.T) {
	svc, userRepo, qiitaRepo := newTestQiitaServiceWithRepos(nil)

	user := &model.User{}
	user.ID = 1
	user.QiitaUsername = "testuser"

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("Update", mock.MatchedBy(func(u *model.User) bool {
		return u.QiitaUsername == ""
	})).Return(nil)
	qiitaRepo.On("DeleteUserArticles", uint(1)).Return(nil)

	err := svc.Disconnect(1)
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	qiitaRepo.AssertExpectations(t)
}

func TestQiitaDisconnect_UserNotFound(t *testing.T) {
	svc, userRepo, _ := newTestQiitaServiceWithRepos(nil)

	userRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Disconnect(999)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	userRepo.AssertExpectations(t)
}

func TestQiitaDisconnect_UpdateError(t *testing.T) {
	svc, userRepo, _ := newTestQiitaServiceWithRepos(nil)

	user := &model.User{}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("Update", mock.Anything).Return(errors.New("db error"))

	err := svc.Disconnect(1)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	userRepo.AssertExpectations(t)
}

// ============================================================
// GetArticles テスト
// ============================================================

func TestQiitaGetArticles_Success(t *testing.T) {
	svc, _, qiitaRepo := newTestQiitaServiceWithRepos(nil)

	expected := []model.QiitaArticle{
		{QiitaID: "abc123", Title: "Go入門", LikesCount: 10},
		{QiitaID: "def456", Title: "React入門", LikesCount: 5},
	}
	qiitaRepo.On("GetArticles", uint(1)).Return(expected, nil)

	articles, err := svc.GetArticles(1)
	assert.NoError(t, err)
	assert.Len(t, articles, 2)
	assert.Equal(t, "Go入門", articles[0].Title)
	qiitaRepo.AssertExpectations(t)
}

func TestQiitaGetArticles_Error(t *testing.T) {
	svc, _, qiitaRepo := newTestQiitaServiceWithRepos(nil)

	qiitaRepo.On("GetArticles", uint(1)).Return([]model.QiitaArticle(nil), errors.New("db error"))

	articles, err := svc.GetArticles(1)
	assert.Error(t, err)
	assert.Nil(t, articles)
	qiitaRepo.AssertExpectations(t)
}

// ============================================================
// GetStats テスト
// ============================================================

func TestQiitaGetStats_Success(t *testing.T) {
	svc, _, qiitaRepo := newTestQiitaServiceWithRepos(nil)

	expected := &model.QiitaStats{
		TotalArticles: 10,
		TotalLikes:    50,
		TotalComments: 20,
	}
	qiitaRepo.On("GetStats", uint(1)).Return(expected, nil)

	stats, err := svc.GetStats(1)
	assert.NoError(t, err)
	assert.Equal(t, 10, stats.TotalArticles)
	assert.Equal(t, 50, stats.TotalLikes)
	qiitaRepo.AssertExpectations(t)
}

func TestQiitaGetStats_Error(t *testing.T) {
	svc, _, qiitaRepo := newTestQiitaServiceWithRepos(nil)

	qiitaRepo.On("GetStats", uint(1)).Return(nil, errors.New("db error"))

	stats, err := svc.GetStats(1)
	assert.Error(t, err)
	assert.Nil(t, stats)
	qiitaRepo.AssertExpectations(t)
}

func TestFetchArticles_Pagination(t *testing.T) {
	callCount := 0
	svc := newTestQiitaService(func(req *http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			// page=1: 100件返す（perPage=100なのでもう1ページあると判断）
			assert.Contains(t, req.URL.String(), "page=1")
			var articles []string
			for i := 0; i < 100; i++ {
				articles = append(articles, fmt.Sprintf(`{"id":"p1_%d","title":"記事%d","url":"https://qiita.com/u/items/p1_%d","likes_count":%d,"comments_count":0,"tags":[],"created_at":"2025-01-01T00:00:00+09:00"}`, i, i, i, i))
			}
			body := "[" + strings.Join(articles, ",") + "]"
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}
		// page=2: 2件返す（100未満なので最終ページ）
		assert.Contains(t, req.URL.String(), "page=2")
		body := `[{"id":"p2_0","title":"最終記事1","url":"https://qiita.com/u/items/p2_0","likes_count":5,"comments_count":1,"tags":[{"name":"Go"}],"created_at":"2025-02-01T00:00:00+09:00"},{"id":"p2_1","title":"最終記事2","url":"https://qiita.com/u/items/p2_1","likes_count":3,"comments_count":0,"tags":[],"created_at":"2025-02-02T00:00:00+09:00"}]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	articles, err := svc.FetchArticles("testuser")
	assert.NoError(t, err)
	assert.Len(t, articles, 102)
	assert.Equal(t, "p1_0", articles[0].QiitaID)
	assert.Equal(t, "p2_1", articles[101].QiitaID)
	assert.Equal(t, 2, callCount)
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
