package model

import "time"

// YouTubeVideo はYouTube動画のキャッシュデータを表す。
type YouTubeVideo struct {
	ID           uint      `json:"id"`
	VideoID      string    `json:"video_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ChannelID    string    `json:"channel_id"`
	ChannelTitle string    `json:"channel_title"`
	ThumbnailURL string    `json:"thumbnail_url"`
	PublishedAt  time.Time `json:"published_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// YouTubeSearchCache は検索クエリごとのキャッシュを表す。
type YouTubeSearchCache struct {
	ID           uint      `json:"id"`
	Query        string    `json:"query"`
	Language     string    `json:"language"`
	VideoIDs     string    `json:"video_ids"`
	CacheExpires time.Time `json:"cache_expires"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
