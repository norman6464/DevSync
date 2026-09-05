package model

import "time"

// ZennArticle はZennから同期した記事データを表す。
type ZennArticle struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	ZennID        int64     `json:"zenn_id" gorm:"not null;uniqueIndex"` // Zenn側の記事ID
	Title         string    `json:"title" gorm:"not null"`
	Slug          string    `json:"slug" gorm:"not null"` // URL用のスラッグ
	Emoji         string    `json:"emoji"`                // 記事のアイコン絵文字
	ArticleType   string    `json:"article_type"`         // "tech"（技術記事）または "idea"（アイデア記事）
	LikedCount    int       `json:"liked_count" gorm:"default:0"`
	CommentsCount int       `json:"comments_count" gorm:"default:0"`
	PublishedAt   time.Time `json:"published_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ZennStats はユーザーのZenn記事統計情報を表す。
// DBテーブルには対応せず、集計結果を格納するDTO。
type ZennStats struct {
	TotalArticles int `json:"total_articles"`
	TotalLikes    int `json:"total_likes"`
	TotalComments int `json:"total_comments"`
}
