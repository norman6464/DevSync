package model

import "time"

// YouTubeVideo はYouTube動画のキャッシュデータを表す。
type YouTubeVideo struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	VideoID      string    `json:"video_id" gorm:"uniqueIndex;not null;size:20"`
	Title        string    `json:"title" gorm:"not null;size:500"`
	Description  string    `json:"description" gorm:"type:text"`
	ChannelID    string    `json:"channel_id" gorm:"size:50"`
	ChannelTitle string    `json:"channel_title" gorm:"size:200"`
	ThumbnailURL string    `json:"thumbnail_url" gorm:"size:500"`
	PublishedAt  time.Time `json:"published_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// YouTubeSearchCache は検索クエリごとのキャッシュを表す。
type YouTubeSearchCache struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Query        string    `json:"query" gorm:"index;not null;size:500"`
	Language     string    `json:"language" gorm:"size:10;default:'ja'"`
	VideoIDs     string    `json:"video_ids" gorm:"type:text"`
	CacheExpires time.Time `json:"cache_expires" gorm:"index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
