package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitStore_AllowsWithinLimit(t *testing.T) {
	store := NewRateLimitStore()
	for i := 0; i < 5; i++ {
		assert.True(t, store.IsAllowed("ip:127.0.0.1", 5, 60))
	}
}

func TestRateLimitStore_DeniesOverLimit(t *testing.T) {
	store := NewRateLimitStore()
	for i := 0; i < 5; i++ {
		store.IsAllowed("ip:127.0.0.1", 5, 60)
	}
	assert.False(t, store.IsAllowed("ip:127.0.0.1", 5, 60))
}

func TestRateLimitStore_DifferentKeysIndependent(t *testing.T) {
	store := NewRateLimitStore()
	store.IsAllowed("ip:1.1.1.1", 1, 60)
	assert.False(t, store.IsAllowed("ip:1.1.1.1", 1, 60))
	assert.True(t, store.IsAllowed("ip:2.2.2.2", 1, 60))
}

func TestRateLimitStore_ResetsAfterWindow(t *testing.T) {
	store := NewRateLimitStore()
	assert.True(t, store.IsAllowed("test", 1, 1))
	assert.False(t, store.IsAllowed("test", 1, 1))

	time.Sleep(1100 * time.Millisecond)
	assert.True(t, store.IsAllowed("test", 1, 1))
}

func TestRateLimitStore_Cleanup(t *testing.T) {
	store := NewRateLimitStore()
	store.IsAllowed("old", 10, 1)

	time.Sleep(1100 * time.Millisecond)
	store.Cleanup(1)

	// クリーンアップ後はキーが削除されている
	store.mu.RLock()
	_, exists := store.requests["old"]
	store.mu.RUnlock()
	assert.False(t, exists)
}

func TestRateLimit_AllowsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := NewRateLimitStore()

	r := gin.New()
	r.Use(RateLimit(store, 5, 60))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimit_BlocksExceedingRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := NewRateLimitStore()

	r := gin.New()
	r.Use(RateLimit(store, 2, 60))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 3回目はブロック
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "リクエストが多すぎます")
}
