package model

import (
	"time"
)

// Project はユーザーのプロジェクトショーケース情報を表す。
// GitHubリポジトリとの紐付けも可能。
type Project struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	UserID       uint              `json:"user_id" gorm:"not null;index"`
	User         User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title        string            `json:"title" gorm:"not null;size:200"`
	Description  string            `json:"description" gorm:"type:text"`
	TechStack    string            `json:"tech_stack" gorm:"type:text"`            // JSON配列形式の使用技術
	DemoURL      string            `json:"demo_url" gorm:"size:500"`               // デモサイトのURL
	GithubURL    string            `json:"github_url" gorm:"size:500"`             // GitHubリポジトリのURL
	ImageURL     string            `json:"image_url" gorm:"size:500"`              // サムネイル画像URL
	Role         string            `json:"role" gorm:"size:100"`                   // 担当役割（例: "Lead Developer"）
	StartDate    *time.Time        `json:"start_date"`                             // プロジェクト開始日
	EndDate      *time.Time        `json:"end_date"`                               // 終了日（進行中はnil）
	Featured     bool              `json:"featured" gorm:"default:false"`          // 注目表示フラグ
	IsArchived   bool              `json:"is_archived" gorm:"default:false"`       // アーカイブフラグ
	GithubRepoID *uint             `json:"github_repo_id"`                         // GitHubRepository との紐付けID
	GithubRepo   *GitHubRepository `json:"github_repo,omitempty" gorm:"foreignKey:GithubRepoID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
