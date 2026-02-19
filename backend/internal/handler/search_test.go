package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostSearchService は PostSearchService のテスト用モック。
type MockPostSearchService struct {
	mock.Mock
}

func (m *MockPostSearchService) SearchPosts(params service.PostSearchParams) (*service.PostSearchResult, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.PostSearchResult), args.Error(1)
}

// MockCircleSearchService は CircleSearchService のテスト用モック。
type MockCircleSearchService struct {
	mock.Mock
}

func (m *MockCircleSearchService) SearchCircles(query string, limit, offset int) (interface{}, int64, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func TestSearchPosts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockPostSearchService)
	mockSvc.On("SearchPosts", mock.MatchedBy(func(p service.PostSearchParams) bool {
		return p.Query == "test" && p.Limit == 20 && p.Offset == 0
	})).Return(
		&service.PostSearchResult{
			Posts: []model.Post{{ID: 1, Title: "Test Post", Content: "Test content"}},
			Total: 1,
			Limit: 20,
		},
		nil,
	)

	h := &SearchHandler{searchService: mockSvc}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/posts?q=test&limit=20&offset=0", nil)
	c.Request.URL.RawQuery = "q=test&limit=20&offset=0"

	h.SearchPosts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response service.PostSearchResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.Total)
	assert.Len(t, response.Posts, 1)
	assert.Equal(t, "Test Post", response.Posts[0].Title)
}

func TestSearchPosts_EmptyQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &SearchHandler{searchService: new(MockPostSearchService)}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/posts?q=", nil)
	c.Request.URL.RawQuery = "q="

	h.SearchPosts(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchPosts_WithTagFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockPostSearchService)
	mockSvc.On("SearchPosts", mock.MatchedBy(func(p service.PostSearchParams) bool {
		return p.Query == "Go" && len(p.Tags) == 2 &&
			p.Tags[0] == "golang" && p.Tags[1] == "beginner"
	})).Return(
		&service.PostSearchResult{
			Posts: []model.Post{{Title: "Go入門"}},
			Total: 1,
			Limit: 20,
		},
		nil,
	)

	h := &SearchHandler{searchService: mockSvc}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/posts?q=Go&tags=golang,beginner", nil)
	c.Request.URL.RawQuery = "q=Go&tags=golang,beginner"

	h.SearchPosts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSearchPosts_WithSortBy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockPostSearchService)
	mockSvc.On("SearchPosts", mock.MatchedBy(func(p service.PostSearchParams) bool {
		return p.Query == "記事" && p.SortBy == service.SearchSortByPopular
	})).Return(
		&service.PostSearchResult{
			Posts: []model.Post{{Title: "人気記事", LikeCount: 100}},
			Total: 1,
			Limit: 20,
		},
		nil,
	)

	h := &SearchHandler{searchService: mockSvc}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/posts?q=記事&sort_by=popular", nil)
	c.Request.URL.RawQuery = "q=%E8%A8%98%E4%BA%8B&sort_by=popular"

	h.SearchPosts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSearchCircles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCircleSvc := new(MockCircleSearchService)
	mockCircleSvc.On("SearchCircles", "golang", 20, 0).Return(
		[]model.StudyCircle{{ID: 1, Name: "Golang Study", Topic: "Programming"}},
		int64(1),
		nil,
	)

	h := &SearchHandler{circleService: mockCircleSvc}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/circles?q=golang&limit=20&offset=0", nil)
	c.Request.URL.RawQuery = "q=golang&limit=20&offset=0"

	h.SearchCircles(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []model.StudyCircle
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "Golang Study", response[0].Name)
}

func TestSearchCircles_EmptyQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCircleSvc := new(MockCircleSearchService)
	h := &SearchHandler{circleService: mockCircleSvc}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/circles?q=", nil)
	c.Request.URL.RawQuery = "q="

	h.SearchCircles(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
