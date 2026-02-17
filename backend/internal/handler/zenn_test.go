package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
)

func TestZenn_Connect_Success(t *testing.T) {
	h, svc := setupZennHandler()
	svc.On("Connect", uint(1), "testuser").Return(3, nil)

	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{
		"username": "testuser",
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestZenn_Connect_MissingUsername(t *testing.T) {
	h, _ := setupZennHandler()

	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestZenn_Connect_ServiceError(t *testing.T) {
	h, svc := setupZennHandler()
	svc.On("Connect", uint(1), "testuser").Return(0, errors.New("api error"))

	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{
		"username": "testuser",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestZenn_Disconnect_Success(t *testing.T) {
	h, svc := setupZennHandler()
	svc.On("Disconnect", uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/zenn/disconnect", h.Disconnect)

	w := doRequest(r, http.MethodDelete, "/zenn/disconnect", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestZenn_Sync_Success(t *testing.T) {
	h, svc := setupZennHandler()
	svc.On("Sync", uint(1)).Return(8, nil)

	r := newRouter(1)
	r.POST("/zenn/sync", h.Sync)

	w := doRequest(r, http.MethodPost, "/zenn/sync", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestZenn_GetArticles_Success(t *testing.T) {
	h, svc := setupZennHandler()
	svc.On("GetArticles", uint(5)).Return([]model.ZennArticle{
		{Title: "Article 1"},
	}, nil)

	r := newRouter(1)
	r.GET("/zenn/:userId/articles", h.GetArticles)

	w := doRequest(r, http.MethodGet, "/zenn/5/articles", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestZenn_GetStats_Success(t *testing.T) {
	h, svc := setupZennHandler()
	svc.On("GetStats", uint(5)).Return(&model.ZennStats{
		TotalArticles: 8,
	}, nil)

	r := newRouter(1)
	r.GET("/zenn/:userId/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/zenn/5/stats", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}
