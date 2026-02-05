package model

import "time"

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	Password        string    `json:"-" gorm:"not null"`
	AvatarURL       string    `json:"avatar_url"`
	Bio             string    `json:"bio"`
	GitHubUsername  string    `json:"github_username" gorm:"uniqueIndex"`
	GitHubToken     string    `json:"-"`
	GitHubConnected bool      `json:"github_connected" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
