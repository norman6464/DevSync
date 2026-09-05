package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit はリクエストボディのサイズを制限するミドルウェアを返す。
// maxBytes を超えるボディを持つリクエストは 413 Request Entity Too Large で拒否される。
// GET/HEAD/OPTIONS などボディを持たないメソッドはスキップする。
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil || c.Request.Method == "GET" ||
			c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Content-Lengthヘッダーによる早期チェック
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "リクエストボディが大きすぎます",
			})
			return
		}

		// ボディを制限付きリーダーで置き換え
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
