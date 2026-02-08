package model

import "time"

// QiitaArticle はQiitaから同期した記事データを表す。
type QiitaArticle struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	QiitaID       string    `json:"qiita_id" gorm:"not null;uniqueIndex"` // Qiita側の記事ID
	Title         string    `json:"title" gorm:"not null"`
	URL           string    `json:"url" gorm:"not null"` // 記事のURL
	LikesCount    int       `json:"likes_count" gorm:"default:0"`
	CommentsCount int       `json:"comments_count" gorm:"default:0"`
	Tags          string    `json:"tags"` // カンマ区切りのタグ名
	PublishedAt   time.Time `json:"published_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// QiitaStats はユーザーのQiita記事統計情報を表す。
// DBテーブルには対応せず、集計結果を格納するDTO。
type QiitaStats struct {
	TotalArticles int `json:"total_articles"`
	TotalLikes    int `json:"total_likes"`
	TotalComments int `json:"total_comments"`
}
