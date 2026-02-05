package model

import "time"

type GitHubContribution struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_date"`
	Date      time.Time `json:"date" gorm:"not null;uniqueIndex:idx_user_date"`
	Count     int       `json:"count" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GitHubLanguageStat struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_lang"`
	Language  string    `json:"language" gorm:"not null;uniqueIndex:idx_user_lang"`
	Bytes     int64     `json:"bytes" gorm:"not null;default:0"`
	RepoCount int       `json:"repo_count" gorm:"not null;default:0"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GitHubRepository struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	GitHubRepoID int64     `json:"github_repo_id" gorm:"not null;uniqueIndex"`
	Name         string    `json:"name" gorm:"not null"`
	FullName     string    `json:"full_name"`
	Description  string    `json:"description"`
	Language     string    `json:"language"`
	Stars        int       `json:"stars"`
	Forks        int       `json:"forks"`
	IsPrivate    bool      `json:"is_private"`
	UpdatedAt    time.Time `json:"updated_at"`
}
