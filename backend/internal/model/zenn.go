package model

import "time"

type ZennArticle struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	ZennID        int64     `json:"zenn_id" gorm:"not null;uniqueIndex"`
	Title         string    `json:"title" gorm:"not null"`
	Slug          string    `json:"slug" gorm:"not null"`
	Emoji         string    `json:"emoji"`
	ArticleType   string    `json:"article_type"` // tech or idea
	LikedCount    int       `json:"liked_count" gorm:"default:0"`
	CommentsCount int       `json:"comments_count" gorm:"default:0"`
	PublishedAt   time.Time `json:"published_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ZennStats represents aggregated Zenn statistics for a user
type ZennStats struct {
	TotalArticles int `json:"total_articles"`
	TotalLikes    int `json:"total_likes"`
	TotalComments int `json:"total_comments"`
}
