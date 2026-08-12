// Package middleware はDevSyncアプリケーションのHTTPミドルウェアのテストを提供する。
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
)

// テスト用のJWTシークレットキー
const testJWTSecret = "test-secret-key-for-middleware-tests"

// generateTestToken はテスト用のJWTトークンを生成するヘルパー関数。
func generateTestToken(userID uint) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testJWTSecret))
	return tokenString
}

// setupTestRouter はテスト用の Gin ルーターとトークン検証 usecase を生成するヘルパー関数。
func setupTestRouter() (*gin.Engine, *usecase.ValidateAuthTokenUseCase) {
	gin.SetMode(gin.TestMode)
	validateToken := usecase.NewValidateAuthTokenUseCase(testJWTSecret)
	r := gin.New()
	return r, validateToken
}

// TestAuthRequired_WithCookie はCookieからトークンを抽出して認証成功することをテストする。
func TestAuthRequired_WithCookie(t *testing.T) {
	r, validateToken := setupTestRouter()
	r.GET("/protected", AuthRequired(validateToken), func(c *gin.Context) {
		userID := c.GetUint("userID")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	token := generateTestToken(123)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}

// TestAuthRequired_WithAuthorizationHeader は後方互換性のためAuthorizationヘッダーでも認証成功することをテストする。
func TestAuthRequired_WithAuthorizationHeader(t *testing.T) {
	r, validateToken := setupTestRouter()
	r.GET("/protected", AuthRequired(validateToken), func(c *gin.Context) {
		userID := c.GetUint("userID")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	token := generateTestToken(456)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "456")
}

// TestAuthRequired_CookiePriority はCookieとAuthorizationヘッダーの両方が存在する場合、Cookieを優先することをテストする。
func TestAuthRequired_CookiePriority(t *testing.T) {
	r, validateToken := setupTestRouter()
	r.GET("/protected", AuthRequired(validateToken), func(c *gin.Context) {
		userID := c.GetUint("userID")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	cookieToken := generateTestToken(100) // Cookie: user_id=100
	headerToken := generateTestToken(200) // Header: user_id=200

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: cookieToken})
	req.Header.Set("Authorization", "Bearer "+headerToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Cookieのuser_id=100が優先される
	assert.Contains(t, w.Body.String(), "100")
	assert.NotContains(t, w.Body.String(), "200")
}

// TestAuthRequired_NeitherProvided はCookieもAuthorizationヘッダーもない場合に401エラーを返すことをテストする。
func TestAuthRequired_NeitherProvided(t *testing.T) {
	r, validateToken := setupTestRouter()
	r.GET("/protected", AuthRequired(validateToken), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// TestAuthRequired_InvalidCookieToken は不正なCookieトークンで401エラーを返すことをテストする。
func TestAuthRequired_InvalidCookieToken(t *testing.T) {
	r, validateToken := setupTestRouter()
	r.GET("/protected", AuthRequired(validateToken), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid-token-value"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// TestAuthRequired_ExpiredCookieToken は期限切れのCookieトークンで401エラーを返すことをテストする。
func TestAuthRequired_ExpiredCookieToken(t *testing.T) {
	r, validateToken := setupTestRouter()
	r.GET("/protected", AuthRequired(validateToken), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	// 期限切れのトークンを生成
	claims := jwt.MapClaims{
		"user_id": 123,
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, _ := token.SignedString([]byte(testJWTSecret))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: expiredToken})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthRequired_InvalidAuthorizationHeaderFormat は不正なAuthorizationヘッダー形式で401エラーを返すことをテストする。
func TestAuthRequired_InvalidAuthorizationHeaderFormat(t *testing.T) {
	r, validateToken := setupTestRouter()
	r.GET("/protected", AuthRequired(validateToken), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
