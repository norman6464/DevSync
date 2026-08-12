package usecase

import (
	"context"
	"log"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// spotifyTokenRefreshMargin はアクセストークンを前倒しで更新する猶予。
// 期限ぎりぎりのトークンで API を叩いて失敗するのを避ける。
const spotifyTokenRefreshMargin = 5 * time.Minute

// GetSpotifyOAuthURLUseCase は Spotify 連携用の認可 URL を組み立てる。
type GetSpotifyOAuthURLUseCase struct {
	client repository.SpotifyAPIClient
}

// NewGetSpotifyOAuthURLUseCase は GetSpotifyOAuthURLUseCase を生成する。
func NewGetSpotifyOAuthURLUseCase(client repository.SpotifyAPIClient) *GetSpotifyOAuthURLUseCase {
	return &GetSpotifyOAuthURLUseCase{client: client}
}

// Execute は認可 URL を返す。URL の組み立てだけで外部通信を行わないため ctx は取らない。
func (uc *GetSpotifyOAuthURLUseCase) Execute(state string) string {
	return uc.client.AuthorizeURL(state)
}

// ConnectSpotifyUseCase は OAuth コールバック後に Spotify アカウントを連携する。
type ConnectSpotifyUseCase struct {
	users  repository.ExternalAccountLinker
	client repository.SpotifyAPIClient
}

// NewConnectSpotifyUseCase は ConnectSpotifyUseCase を生成する。
func NewConnectSpotifyUseCase(
	users repository.ExternalAccountLinker,
	client repository.SpotifyAPIClient,
) *ConnectSpotifyUseCase {
	return &ConnectSpotifyUseCase{users: users, client: client}
}

// Execute は認可コードをトークンに交換し、連携情報を保存する。
func (uc *ConnectSpotifyUseCase) Execute(ctx context.Context, userID uint, code string) error {
	token, err := uc.client.ExchangeCode(ctx, code)
	if err != nil {
		return err
	}

	user, err := findLinkedUser(ctx, uc.users, userID)
	if err != nil {
		return err
	}

	user.SpotifyConnected = true
	user.SpotifyToken = token.AccessToken
	user.SpotifyRefreshToken = token.RefreshToken
	user.SpotifyTokenExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return uc.users.Update(ctx, user)
}

// DisconnectSpotifyUseCase は Spotify 連携を解除する。
type DisconnectSpotifyUseCase struct {
	users   repository.ExternalAccountLinker
	spotify repository.SpotifyRepository
}

// NewDisconnectSpotifyUseCase は DisconnectSpotifyUseCase を生成する。
func NewDisconnectSpotifyUseCase(
	users repository.ExternalAccountLinker,
	spotify repository.SpotifyRepository,
) *DisconnectSpotifyUseCase {
	return &DisconnectSpotifyUseCase{users: users, spotify: spotify}
}

// Execute は連携情報を消し、保存済みの再生履歴を削除する。
// 履歴削除の失敗は連携解除自体を失敗させない（解除は完了しているため）。
func (uc *DisconnectSpotifyUseCase) Execute(ctx context.Context, userID uint) error {
	user, err := findLinkedUser(ctx, uc.users, userID)
	if err != nil {
		return err
	}

	user.SpotifyConnected = false
	user.SpotifyToken = ""
	user.SpotifyRefreshToken = ""
	user.SpotifyTokenExpiry = time.Time{}

	if err := uc.users.Update(ctx, user); err != nil {
		return err
	}

	if err := uc.spotify.DeleteUserData(ctx, userID); err != nil {
		log.Printf("[ERROR] Spotify連携データの削除に失敗 (userID=%d): %v", userID, err)
	}
	return nil
}

// GetSpotifyCurrentlyPlayingUseCase は現在再生中の曲を取得する。
type GetSpotifyCurrentlyPlayingUseCase struct {
	users  repository.ExternalAccountLinker
	client repository.SpotifyAPIClient
}

// NewGetSpotifyCurrentlyPlayingUseCase は GetSpotifyCurrentlyPlayingUseCase を生成する。
func NewGetSpotifyCurrentlyPlayingUseCase(
	users repository.ExternalAccountLinker,
	client repository.SpotifyAPIClient,
) *GetSpotifyCurrentlyPlayingUseCase {
	return &GetSpotifyCurrentlyPlayingUseCase{users: users, client: client}
}

// Execute は現在再生中の曲を返す。再生していない場合は (nil, nil) を返す。
func (uc *GetSpotifyCurrentlyPlayingUseCase) Execute(ctx context.Context, userID uint) (*model.SpotifyCurrentlyPlaying, error) {
	token, err := resolveSpotifyToken(ctx, uc.users, uc.client, userID)
	if err != nil {
		return nil, err
	}
	return uc.client.FetchCurrentlyPlaying(ctx, token)
}

// GetSpotifyRecentlyPlayedUseCase は最近再生した曲を取得する。
type GetSpotifyRecentlyPlayedUseCase struct {
	users  repository.ExternalAccountLinker
	client repository.SpotifyAPIClient
}

// NewGetSpotifyRecentlyPlayedUseCase は GetSpotifyRecentlyPlayedUseCase を生成する。
func NewGetSpotifyRecentlyPlayedUseCase(
	users repository.ExternalAccountLinker,
	client repository.SpotifyAPIClient,
) *GetSpotifyRecentlyPlayedUseCase {
	return &GetSpotifyRecentlyPlayedUseCase{users: users, client: client}
}

// Execute は最近再生した曲を新しい順に返す。
func (uc *GetSpotifyRecentlyPlayedUseCase) Execute(ctx context.Context, userID uint) ([]model.SpotifyRecentTrackResponse, error) {
	token, err := resolveSpotifyToken(ctx, uc.users, uc.client, userID)
	if err != nil {
		return nil, err
	}
	return uc.client.FetchRecentlyPlayed(ctx, token)
}

// resolveSpotifyToken は有効なアクセストークンを返す。
// 期限が近い場合はリフレッシュトークンで更新し、更新後のトークンを保存する。
func resolveSpotifyToken(
	ctx context.Context,
	users repository.ExternalAccountLinker,
	client repository.SpotifyAPIClient,
	userID uint,
) (string, error) {
	user, err := findLinkedUser(ctx, users, userID)
	if err != nil {
		return "", err
	}

	if user.SpotifyRefreshToken == "" {
		return "", domain.NewError(domain.ErrCodeBadRequest, "Spotifyが連携されていません", nil)
	}
	if time.Now().Before(user.SpotifyTokenExpiry.Add(-spotifyTokenRefreshMargin)) {
		return user.SpotifyToken, nil
	}

	token, err := client.RefreshAccessToken(ctx, user.SpotifyRefreshToken)
	if err != nil {
		return "", err
	}

	user.SpotifyToken = token.AccessToken
	user.SpotifyTokenExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	// リフレッシュトークンは再発行されないこともあるため、返ってきたときだけ差し替える
	if token.RefreshToken != "" {
		user.SpotifyRefreshToken = token.RefreshToken
	}

	if err := users.Update(ctx, user); err != nil {
		return "", err
	}
	return user.SpotifyToken, nil
}
