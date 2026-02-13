package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSearchPosts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockPostRepository)
	mockRepo.On("Search", "test", 20, 0).Return(
		[]model.Post{{ID: 1, Title: "Test Post", Content: "Test content"}},
		int64(1),
		nil,
	)

	handler := &SearchHandler{postRepo: mockRepo}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/posts?q=test&limit=20&offset=0", nil)
	c.Request.URL.RawQuery = "q=test&limit=20&offset=0"

	handler.SearchPosts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []model.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "Test Post", response[0].Title)
}

func TestSearchPosts_EmptyQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockPostRepository)
	handler := &SearchHandler{postRepo: mockRepo}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/posts?q=", nil)
	c.Request.URL.RawQuery = "q="

	handler.SearchPosts(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchCircles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockStudyCircleRepository)
	mockRepo.On("Search", "golang", 20, 0).Return(
		[]model.StudyCircle{{ID: 1, Name: "Golang Study", Topic: "Programming"}},
		int64(1),
		nil,
	)

	handler := &SearchHandler{circleRepo: mockRepo}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/circles?q=golang&limit=20&offset=0", nil)
	c.Request.URL.RawQuery = "q=golang&limit=20&offset=0"

	handler.SearchCircles(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []model.StudyCircle
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "Golang Study", response[0].Name)
}

func TestSearchCircles_EmptyQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockStudyCircleRepository)
	handler := &SearchHandler{circleRepo: mockRepo}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/search/circles?q=", nil)
	c.Request.URL.RawQuery = "q="

	handler.SearchCircles(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
