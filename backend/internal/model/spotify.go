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

// SpotifyToken は Spotify のトークンエンドポイントから取得する認証情報を表す。
type SpotifyToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// SpotifyCurrentlyPlaying は現在再生中の曲情報を表す。
type SpotifyCurrentlyPlaying struct {
	IsPlaying  bool   `json:"is_playing"`
	TrackName  string `json:"track_name"`
	ArtistName string `json:"artist_name"`
	AlbumName  string `json:"album_name"`
	AlbumImage string `json:"album_image"`
	TrackURL   string `json:"track_url"`
	ProgressMs int    `json:"progress_ms"`
	DurationMs int    `json:"duration_ms"`
}

// SpotifyRecentTrackResponse は最近再生した曲のレスポンスを表す。
type SpotifyRecentTrackResponse struct {
	TrackName  string `json:"track_name"`
	ArtistName string `json:"artist_name"`
	AlbumName  string `json:"album_name"`
	AlbumImage string `json:"album_image"`
	TrackURL   string `json:"track_url"`
	PlayedAt   string `json:"played_at"`
}
