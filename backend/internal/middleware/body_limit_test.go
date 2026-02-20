package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupBodyLimitRouter(limit int64) *gin.Engine {
	r := gin.New()
	r.Use(BodyLimit(limit))
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func TestBodyLimit_AllowsSmallBody(t *testing.T) {
	r := setupBodyLimitRouter(1024) // 1KB

	body := bytes.NewBufferString(`{"name":"test"}`)
	req := httptest.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBodyLimit_RejectsLargeBody(t *testing.T) {
	r := setupBodyLimitRouter(100) // 100バイト制限

	// 200バイトのボディを送信
	body := bytes.NewBufferString(strings.Repeat("a", 200))
	req := httptest.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestBodyLimit_AllowsExactLimitSize(t *testing.T) {
	r := setupBodyLimitRouter(50)

	body := bytes.NewBufferString(strings.Repeat("a", 50))
	req := httptest.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBodyLimit_AllowsGetRequests(t *testing.T) {
	r := gin.New()
	r.Use(BodyLimit(10)) // 極小制限
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBodyLimit_AllowsEmptyBody(t *testing.T) {
	r := setupBodyLimitRouter(100)

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBodyLimit_RejectsBasedOnContentLength(t *testing.T) {
	r := setupBodyLimitRouter(100)

	body := bytes.NewBufferString(`{"data":"test"}`)
	req := httptest.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = 200 // 実際のボディより大きいContent-Lengthヘッダー
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}
