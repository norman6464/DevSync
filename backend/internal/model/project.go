package model

import (
	"time"
)

// Project はユーザーのプロジェクトショーケース情報を表す。
// GitHubリポジトリとの紐付けも可能。
type Project struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"user_id"`
	User         User              `json:"user,omitempty"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	TechStack    string            `json:"tech_stack"`     // JSON配列形式の使用技術
	DemoURL      string            `json:"demo_url"`       // デモサイトのURL
	GithubURL    string            `json:"github_url"`     // GitHubリポジトリのURL
	ImageURL     string            `json:"image_url"`      // サムネイル画像URL
	Role         string            `json:"role"`           // 担当役割（例: "Lead Developer"）
	StartDate    *time.Time        `json:"start_date"`     // プロジェクト開始日
	EndDate      *time.Time        `json:"end_date"`       // 終了日（進行中はnil）
	Featured     bool              `json:"featured"`       // 注目表示フラグ
	IsArchived   bool              `json:"is_archived"`    // アーカイブフラグ
	GithubRepoID *uint             `json:"github_repo_id"` // GitHubRepository との紐付けID
	GithubRepo   *GitHubRepository `json:"github_repo,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
