package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"not null;index"`
	User            User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title           string         `json:"title" gorm:"not null;size:200"`
	Description     string         `json:"description" gorm:"type:text"`
	TechStack       string         `json:"tech_stack" gorm:"type:text"` // JSON array of technologies
	DemoURL         string         `json:"demo_url" gorm:"size:500"`
	GithubURL       string         `json:"github_url" gorm:"size:500"`
	ImageURL        string         `json:"image_url" gorm:"size:500"`
	Role            string         `json:"role" gorm:"size:100"`            // e.g., "Lead Developer", "Contributor"
	StartDate       *time.Time     `json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
	Featured        bool           `json:"featured" gorm:"default:false"`
	GithubRepoID    *uint          `json:"github_repo_id"`                  // Link to GitHubRepository if exists
	GithubRepo      *GitHubRepository `json:"github_repo,omitempty" gorm:"foreignKey:GithubRepoID"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}
