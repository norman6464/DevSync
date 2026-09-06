package model

import "time"

// QiitaArticle はQiitaから同期した記事データを表す。
type QiitaArticle struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	QiitaID       string    `json:"qiita_id"` // Qiita側の記事ID
	Title         string    `json:"title"`
	URL           string    `json:"url"` // 記事のURL
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
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
