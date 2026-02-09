// Package config はDevSyncアプリケーションの設定管理を提供する。
// 環境変数から設定値を読み込み、デフォルト値のフォールバックを行う。
package config

import (
	"fmt"
	"os"
)

// Config はアプリケーション全体の設定を保持する構造体。
type Config struct {
	Port               string // サーバーポート番号
	DBHost             string // PostgreSQLホスト
	DBPort             string // PostgreSQLポート
	DBUser             string // PostgreSQLユーザー名
	DBPass             string // PostgreSQLパスワード
	DBName             string // PostgreSQLデータベース名
	DBSSLMode          string // PostgreSQL SSL接続モード
	JWTSecret          string // JWT署名用シークレット
	GitHubClientID     string // GitHub OAuth クライアントID
	GitHubClientSecret string // GitHub OAuth クライアントシークレット
	GitHubRedirectURL  string // GitHub OAuth リダイレクトURL
	CORSOrigins        string // CORS許可オリジン（カンマ区切り）
	OpenAIAPIKey       string // OpenAI APIキー（LLM機能用、空なら無効）
}

// Load は環境変数から設定を読み込み、Configインスタンスを返す。
// 環境変数が未設定の場合はデフォルト値を使用する。
func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "devsync"),
		DBPass:             getEnv("DB_PASSWORD", "devsync"),
		DBName:             getEnv("DB_NAME", "devsync"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "devsync-dev-secret-change-me"),
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL: getEnv("GITHUB_REDIRECT_URL", "http://localhost:5173/github/callback"),
		CORSOrigins:       getEnv("CORS_ORIGINS", "http://localhost:5173"),
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
	}
}

// DSN はPostgreSQL接続文字列（Data Source Name）を返す。
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode,
	)
}

// getEnv は環境変数の値を取得する。未設定の場合はfallbackを返す。
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
