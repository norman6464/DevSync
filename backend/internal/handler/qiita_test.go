package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
)

func TestQiita_Connect_Success(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("Connect", uint(1), "testuser").Return(5, nil)

	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{
		"username": "testuser",
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestQiita_Connect_MissingUsername(t *testing.T) {
	h, _ := setupQiitaHandler()

	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestQiita_Connect_ServiceError(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("Connect", uint(1), "testuser").Return(0, errors.New("api error"))

	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{
		"username": "testuser",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestQiita_Disconnect_Success(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("Disconnect", uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/qiita/disconnect", h.Disconnect)

	w := doRequest(r, http.MethodDelete, "/qiita/disconnect", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestQiita_Sync_Success(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("Sync", uint(1)).Return(10, nil)

	r := newRouter(1)
	r.POST("/qiita/sync", h.Sync)

	w := doRequest(r, http.MethodPost, "/qiita/sync", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestQiita_GetArticles_Success(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("GetArticles", uint(5)).Return([]model.QiitaArticle{
		{Title: "Article 1"},
	}, nil)

	r := newRouter(1)
	r.GET("/qiita/:userId/articles", h.GetArticles)

	w := doRequest(r, http.MethodGet, "/qiita/5/articles", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestQiita_Disconnect_ServiceError(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("Disconnect", uint(1)).Return(errors.New("db error"))

	r := newRouter(1)
	r.DELETE("/qiita/disconnect", h.Disconnect)

	w := doRequest(r, http.MethodDelete, "/qiita/disconnect", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestQiita_Sync_ServiceError(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("Sync", uint(1)).Return(0, errors.New("sync error"))

	r := newRouter(1)
	r.POST("/qiita/sync", h.Sync)

	w := doRequest(r, http.MethodPost, "/qiita/sync", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestQiita_GetArticles_InvalidID(t *testing.T) {
	h, _ := setupQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/articles", h.GetArticles)

	w := doRequest(r, http.MethodGet, "/qiita/abc/articles", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestQiita_GetArticles_ServiceError(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("GetArticles", uint(5)).Return([]model.QiitaArticle(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/qiita/:userId/articles", h.GetArticles)

	w := doRequest(r, http.MethodGet, "/qiita/5/articles", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestQiita_GetStats_InvalidID(t *testing.T) {
	h, _ := setupQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/qiita/abc/stats", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestQiita_GetStats_ServiceError(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("GetStats", uint(5)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/qiita/:userId/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/qiita/5/stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestQiita_GetStats_Success(t *testing.T) {
	h, svc := setupQiitaHandler()
	svc.On("GetStats", uint(5)).Return(&model.QiitaStats{
		TotalArticles: 10,
	}, nil)

	r := newRouter(1)
	r.GET("/qiita/:userId/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/qiita/5/stats", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}
