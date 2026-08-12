// Package middleware はDevSyncアプリケーションのHTTPミドルウェアを提供する。
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// AuthRequired はJWT認証を必須とするミドルウェアを返す。
// 以下の優先順位でトークンを抽出する:
//  1. httpOnly Cookie（"token"）
//  2. Authorizationヘッダー（"Bearer <token>"）— 後方互換性のため維持
//
// 検証に成功した場合はコンテキストにuserIDをセットして次のハンドラに処理を委譲する。
func AuthRequired(validateToken *usecase.ValidateAuthTokenUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Cookieからトークンを取得（優先）
		if cookie, err := c.Cookie("token"); err == nil && cookie != "" {
			tokenString = cookie
		}

		// 2. Cookieにない場合、Authorizationヘッダーにフォールバック
		if tokenString == "" {
			header := c.GetHeader("Authorization")
			if header != "" {
				parts := strings.SplitN(header, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				} else {
					c.JSON(http.StatusUnauthorized, domain.NewErrorResponse("invalid authorization header format", string(domain.ErrCodeUnauthorized), nil))
					c.Abort()
					return
				}
			}
		}

		// CookieもAuthorizationヘッダーもない場合
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse("authentication required", string(domain.ErrCodeUnauthorized), nil))
			c.Abort()
			return
		}

		userID, err := validateToken.Execute(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse("invalid or expired token", string(domain.ErrCodeUnauthorized), nil))
			c.Abort()
			return
		}

		// 認証済みユーザーIDをコンテキストに格納
		c.Set("userID", userID)
		c.Next()
	}
}
