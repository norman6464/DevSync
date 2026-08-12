package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockPostSearchPort は usecase/repository.PostSearchRepository のモック。
type mockPostSearchPort struct{ mock.Mock }

func (m *mockPostSearchPort) SearchWithFilter(ctx context.Context, params model.PostSearchParams) ([]model.Post, int64, error) {
	args := m.Called(ctx, params)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Get(1).(int64), args.Error(2)
}

// setupSearchHandler は本物の usecase に port モックを注入した SearchHandler を生成する。
func setupSearchHandler() (*SearchHandler, *mockPostSearchPort, *mockStudyCircleRepo) {
	posts := new(mockPostSearchPort)
	circles := new(mockStudyCircleRepo)
	h := NewSearchHandler(
		usecase.NewSearchPostsUseCase(posts),
		usecase.NewSearchStudyCirclesUseCase(circles),
	)
	return h, posts, circles
}

// ---------- SearchPosts ----------

func TestSearchPosts_Success(t *testing.T) {
	h, posts, _ := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/posts", h.SearchPosts)

	posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
		return p.Query == "test" && p.Limit == 20 && p.Offset == 0 && p.SortBy == model.SearchSortByLatest
	})).Return([]model.Post{{ID: 1, Title: "Test Post", Content: "Test content"}}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/search/posts?q=test&limit=20&offset=0", nil)
	assertStatus(t, w, http.StatusOK)

	var response model.PostSearchResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, int64(1), response.Total)
	assert.Len(t, response.Posts, 1)
	assert.Equal(t, "Test Post", response.Posts[0].Title)
	assert.Equal(t, 20, response.Limit)
	posts.AssertExpectations(t)
}

func TestSearchPosts_EmptyQuery(t *testing.T) {
	h, posts, _ := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/posts", h.SearchPosts)

	w := doRequest(r, http.MethodGet, "/search/posts?q=", nil)
	assertStatus(t, w, http.StatusBadRequest)
	posts.AssertNotCalled(t, "SearchWithFilter", mock.Anything, mock.Anything)
}

func TestSearchPosts_WithTagFilter(t *testing.T) {
	h, posts, _ := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/posts", h.SearchPosts)

	posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
		return p.Query == "Go" && len(p.Tags) == 2 && p.Tags[0] == "golang" && p.Tags[1] == "beginner"
	})).Return([]model.Post{{Title: "Go入門"}}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/search/posts?q=Go&tags=golang,beginner", nil)
	assertStatus(t, w, http.StatusOK)
	posts.AssertExpectations(t)
}

func TestSearchPosts_WithSortBy(t *testing.T) {
	h, posts, _ := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/posts", h.SearchPosts)

	posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
		return p.Query == "記事" && p.SortBy == model.SearchSortByPopular
	})).Return([]model.Post{{Title: "人気記事", LikeCount: 100}}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/search/posts?q=%E8%A8%98%E4%BA%8B&sort_by=popular", nil)
	assertStatus(t, w, http.StatusOK)
	posts.AssertExpectations(t)
}

// 無効なソート順は handler が最新順へ正規化する（usecase まで届かない）。
func TestSearchPosts_InvalidSortByFallsBackToLatest(t *testing.T) {
	h, posts, _ := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/posts", h.SearchPosts)

	posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
		return p.SortBy == model.SearchSortByLatest
	})).Return([]model.Post{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/search/posts?q=Go&sort_by=unknown", nil)
	assertStatus(t, w, http.StatusOK)
	posts.AssertExpectations(t)
}

// limit の上限は 100 に切り詰められる。
func TestSearchPosts_LimitIsCapped(t *testing.T) {
	h, posts, _ := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/posts", h.SearchPosts)

	posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
		return p.Limit == 100
	})).Return([]model.Post{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/search/posts?q=Go&limit=1000", nil)
	assertStatus(t, w, http.StatusOK)

	var response model.PostSearchResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, 100, response.Limit)
	posts.AssertExpectations(t)
}

func TestSearchPosts_WithDateRange(t *testing.T) {
	h, posts, _ := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/posts", h.SearchPosts)

	posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
		return p.Query == "Go" && p.DateFrom != nil && p.DateTo != nil
	})).Return([]model.Post{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/search/posts?q=Go&date_from=2024-01-01&date_to=2024-12-31", nil)
	assertStatus(t, w, http.StatusOK)
	posts.AssertExpectations(t)
}

func TestSearchPosts_RepositoryError(t *testing.T) {
	h, posts, _ := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/posts", h.SearchPosts)

	posts.On("SearchWithFilter", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/search/posts?q=error", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	posts.AssertExpectations(t)
}

// ---------- SearchCircles ----------

func TestSearchCircles_Success(t *testing.T) {
	h, _, circles := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/circles", h.SearchCircles)

	circles.On("Search", mock.Anything, "golang", 20, 0).Return(
		[]model.StudyCircle{{ID: 1, Name: "Golang Study", Topic: "Programming"}},
		int64(1),
		nil,
	)

	w := doRequest(r, http.MethodGet, "/search/circles?q=golang&limit=20&offset=0", nil)
	assertStatus(t, w, http.StatusOK)

	var response []model.StudyCircle
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Len(t, response, 1)
	assert.Equal(t, "Golang Study", response[0].Name)
	circles.AssertExpectations(t)
}

func TestSearchCircles_EmptyQuery(t *testing.T) {
	h, _, circles := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/circles", h.SearchCircles)

	w := doRequest(r, http.MethodGet, "/search/circles?q=", nil)
	assertStatus(t, w, http.StatusBadRequest)
	circles.AssertNotCalled(t, "Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestSearchCircles_RepositoryError(t *testing.T) {
	h, _, circles := setupSearchHandler()
	r := newRouter(1)
	r.GET("/search/circles", h.SearchCircles)

	circles.On("Search", mock.Anything, "error", 20, 0).Return(nil, int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/search/circles?q=error", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	circles.AssertExpectations(t)
}
