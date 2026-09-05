package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitStore はIPごとのリクエスト履歴をインメモリで管理する。
type RateLimitStore struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
}

// NewRateLimitStore は新しいRateLimitStoreを生成する。
func NewRateLimitStore() *RateLimitStore {
	return &RateLimitStore{
		requests: make(map[string][]time.Time),
	}
}

// IsAllowed は指定キーがレート制限内か判定する。
// windowSeconds秒以内にmaxRequests回以下のリクエストを許可する。
func (s *RateLimitStore) IsAllowed(key string, maxRequests int, windowSeconds int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(time.Duration(-windowSeconds) * time.Second)

	// ウィンドウ内のリクエストのみ保持
	var filtered []time.Time
	for _, t := range s.requests[key] {
		if t.After(windowStart) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= maxRequests {
		s.requests[key] = filtered
		return false
	}

	s.requests[key] = append(filtered, now)
	return true
}

// Cleanup は古いエントリを削除してメモリを解放する。
func (s *RateLimitStore) Cleanup(windowSeconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(time.Duration(-windowSeconds) * time.Second)
	for key, times := range s.requests {
		var filtered []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) == 0 {
			delete(s.requests, key)
		} else {
			s.requests[key] = filtered
		}
	}
}

// StartCleanup は定期的に古いエントリを削除するgoroutineを起動する。
func StartCleanup(store *RateLimitStore, windowSeconds int) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		store.Cleanup(windowSeconds)
	}
}

// RateLimit はIPごとのレート制限ミドルウェアを返す。
func RateLimit(store *RateLimitStore, maxRequests int, windowSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !store.IsAllowed(clientIP, maxRequests, windowSeconds) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "リクエストが多すぎます。しばらく待ってからお試しください",
			})
			return
		}

		c.Next()
	}
}
