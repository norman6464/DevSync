package model

import "time"

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	Password        string    `json:"-"`
	AvatarURL       string    `json:"avatar_url"`
	Bio             string    `json:"bio"`
	GitHubID        int64     `json:"github_id" gorm:"uniqueIndex"`
	GitHubUsername  string    `json:"github_username"`
	GitHubToken     string    `json:"-"`
	GitHubConnected  bool      `json:"github_connected" gorm:"default:false"`
	ZennUsername     string    `json:"zenn_username"`
	SkillsLanguages  string    `json:"skills_languages"`
	SkillsFrameworks string    `json:"skills_frameworks"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
