package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port              string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	DBSSLMode         string
	JWTSecret         string
	GitHubClientID    string
	GitHubClientSecret string
	GitHubRedirectURL      string
	GitHubLoginRedirectURL string
	CORSOrigins            string
}

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
		GitHubRedirectURL:      getEnv("GITHUB_REDIRECT_URL", "http://localhost:5173/github/callback"),
		GitHubLoginRedirectURL: getEnv("GITHUB_LOGIN_REDIRECT_URL", "http://localhost:5173/auth/github/callback"),
		CORSOrigins:            getEnv("CORS_ORIGINS", "http://localhost:5173"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
