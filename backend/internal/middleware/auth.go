// Package middleware はDevSyncアプリケーションのHTTPミドルウェアを提供する。
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// AuthRequired はJWT認証を必須とするミドルウェアを返す。
// AuthorizationヘッダーからBearerトークンを抽出し、検証に成功した場合は
// コンテキストにuserIDをセットして次のハンドラに処理を委譲する。
func AuthRequired(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// "Bearer <token>" 形式を検証
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		userID, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// 認証済みユーザーIDをコンテキストに格納
		c.Set("userID", userID)
		c.Next()
	}
}
