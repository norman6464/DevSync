package model

import "time"

type QiitaArticle struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	QiitaID       string    `json:"qiita_id" gorm:"not null;uniqueIndex"`
	Title         string    `json:"title" gorm:"not null"`
	URL           string    `json:"url" gorm:"not null"`
	LikesCount    int       `json:"likes_count" gorm:"default:0"`
	CommentsCount int       `json:"comments_count" gorm:"default:0"`
	Tags          string    `json:"tags"` // comma-separated tag names
	PublishedAt   time.Time `json:"published_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// QiitaStats represents aggregated Qiita statistics for a user
type QiitaStats struct {
	TotalArticles int `json:"total_articles"`
	TotalLikes    int `json:"total_likes"`
	TotalComments int `json:"total_comments"`
}
