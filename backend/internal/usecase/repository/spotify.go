package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// SpotifyRepository は Spotify 連携データの永続化に対する、usecase 側が要求する契約。
type SpotifyRepository interface {
	// DeleteUserData は指定ユーザーの Spotify 連携データを削除する。
	DeleteUserData(ctx context.Context, userID uint) error
}

// SpotifyAPIClient は Spotify API 呼び出しに対する、usecase 側が要求する契約。
// 実装は adapter/external に置く。
type SpotifyAPIClient interface {
	// AuthorizeURL は連携用の OAuth 認可 URL を返す。
	AuthorizeURL(state string) string
	// ExchangeCode は認可コードをアクセストークン・リフレッシュトークンに交換する。
	ExchangeCode(ctx context.Context, code string) (*model.SpotifyToken, error)
	// RefreshAccessToken はリフレッシュトークンでアクセストークンを更新する。
	RefreshAccessToken(ctx context.Context, refreshToken string) (*model.SpotifyToken, error)
	// FetchCurrentlyPlaying は現在再生中の曲を返す。再生していない場合は (nil, nil) を返す。
	FetchCurrentlyPlaying(ctx context.Context, token string) (*model.SpotifyCurrentlyPlaying, error)
	// FetchRecentlyPlayed は最近再生した曲を新しい順に返す。
	FetchRecentlyPlayed(ctx context.Context, token string) ([]model.SpotifyRecentTrackResponse, error)
}
