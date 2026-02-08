// Package handler はDevSyncアプリケーションのHTTPハンドラ層を提供する。
// Ginフレームワークを使用し、リクエストの受付・バリデーション・レスポンス返却を担当する。
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck はヘルスチェックエンドポイントのハンドラ。
// サーバーの稼働状態を確認するために使用する。
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
