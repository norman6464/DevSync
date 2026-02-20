package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ---------- Search ----------

func TestYouTubeSearch_Success(t *testing.T) {
	h, svc := setupYouTubeHandler()
	videos := []model.YouTubeVideo{{VideoID: "abc123", Title: "Go Tutorial"}}
	svc.On("Search", "golang", "en").Return(videos, false, nil)

	r := newRouter(1)
	r.GET("/youtube/search", h.Search)
	w := doRequest(r, http.MethodGet, "/youtube/search?q=golang&lang=en", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "query", "golang")
	assertJSONEqual(t, body, "total", float64(1))
	svc.AssertExpectations(t)
}

func TestYouTubeSearch_DefaultLanguage(t *testing.T) {
	h, svc := setupYouTubeHandler()
	svc.On("Search", "react", "ja").Return([]model.YouTubeVideo{}, false, nil)

	r := newRouter(1)
	r.GET("/youtube/search", h.Search)
	w := doRequest(r, http.MethodGet, "/youtube/search?q=react", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestYouTubeSearch_MissingQuery(t *testing.T) {
	h, _ := setupYouTubeHandler()

	r := newRouter(1)
	r.GET("/youtube/search", h.Search)
	w := doRequest(r, http.MethodGet, "/youtube/search", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestYouTubeSearch_ServiceError(t *testing.T) {
	h, svc := setupYouTubeHandler()
	svc.On("Search", "test", "ja").Return(nil, false, errors.New("api error"))

	r := newRouter(1)
	r.GET("/youtube/search", h.Search)
	w := doRequest(r, http.MethodGet, "/youtube/search?q=test", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestYouTubeSearch_NilVideos(t *testing.T) {
	h, svc := setupYouTubeHandler()
	svc.On("Search", "empty", "ja").Return(nil, true, nil)

	r := newRouter(1)
	r.GET("/youtube/search", h.Search)
	w := doRequest(r, http.MethodGet, "/youtube/search?q=empty", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "total", float64(0))
	svc.AssertExpectations(t)
}

// ---------- Recommend ----------

func TestYouTubeRecommend_Success(t *testing.T) {
	h, svc := setupYouTubeHandler()
	videos := []model.YouTubeVideo{{VideoID: "xyz789", Title: "TypeScript Tips"}}
	skills := []string{"Go", "TypeScript"}
	svc.On("GetRecommendations", uint(1)).Return(videos, skills, nil)
	svc.On("IsAvailable").Return(true)

	r := newRouter(1)
	r.GET("/youtube/recommend", h.Recommend)
	w := doRequest(r, http.MethodGet, "/youtube/recommend", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "available", true)
	svc.AssertExpectations(t)
}

func TestYouTubeRecommend_ServiceError(t *testing.T) {
	h, svc := setupYouTubeHandler()
	svc.On("GetRecommendations", uint(1)).Return(nil, nil, errors.New("fail"))

	r := newRouter(1)
	r.GET("/youtube/recommend", h.Recommend)
	w := doRequest(r, http.MethodGet, "/youtube/recommend", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestYouTubeRecommend_NilResults(t *testing.T) {
	h, svc := setupYouTubeHandler()
	svc.On("GetRecommendations", uint(1)).Return(nil, nil, nil)
	svc.On("IsAvailable").Return(false)

	r := newRouter(1)
	r.GET("/youtube/recommend", h.Recommend)
	w := doRequest(r, http.MethodGet, "/youtube/recommend", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "available", false)
	svc.AssertExpectations(t)
}

// ---------- Status ----------

func TestYouTubeStatus_Available(t *testing.T) {
	h, svc := setupYouTubeHandler()
	svc.On("IsAvailable").Return(true)

	r := newRouter(1)
	r.GET("/youtube/status", h.Status)
	w := doRequest(r, http.MethodGet, "/youtube/status", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "available", true)
	svc.AssertExpectations(t)
}

func TestYouTubeStatus_Unavailable(t *testing.T) {
	h, svc := setupYouTubeHandler()
	svc.On("IsAvailable").Return(false)

	r := newRouter(1)
	r.GET("/youtube/status", h.Status)
	w := doRequest(r, http.MethodGet, "/youtube/status", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "available", false)
	svc.AssertExpectations(t)
}

// ---------- ヘルパー ----------

func assertJSONEqual(t *testing.T, body map[string]interface{}, key string, expected interface{}) {
	t.Helper()
	if body[key] != expected {
		t.Errorf("expected %s=%v, got %v", key, expected, body[key])
	}
}
