package model

import "time"

// SpotifyRecentTrack はSpotifyの最近再生した曲データを表す。
type SpotifyRecentTrack struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	TrackName  string    `json:"track_name" gorm:"not null"`
	ArtistName string    `json:"artist_name" gorm:"not null"`
	AlbumName  string    `json:"album_name"`
	AlbumImage string    `json:"album_image"`
	TrackURL   string    `json:"track_url"`
	PlayedAt   time.Time `json:"played_at" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}
