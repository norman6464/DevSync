package service

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// SpotifyService はSpotify連携のビジネスロジックを提供する。
type SpotifyService struct {
	cfg         *config.Config
	userRepo    repository.UserRepositoryInterface
	spotifyRepo repository.SpotifyRepositoryInterface
}

// NewSpotifyService は新しいSpotifyServiceインスタンスを生成する。
func NewSpotifyService(
	cfg *config.Config,
	userRepo repository.UserRepositoryInterface,
	spotifyRepo repository.SpotifyRepositoryInterface,
) *SpotifyService {
	return &SpotifyService{cfg: cfg, userRepo: userRepo, spotifyRepo: spotifyRepo}
}

// GetOAuthURL はSpotify連携用のOAuth認可URLを生成する。
func (s *SpotifyService) GetOAuthURL(state string) string {
	return "https://accounts.spotify.com/authorize?" +
		"client_id=" + s.cfg.SpotifyClientID +
		"&response_type=code" +
		"&redirect_uri=" + url.QueryEscape(s.cfg.SpotifyRedirectURL) +
		"&scope=" + url.QueryEscape("user-read-currently-playing user-read-recently-played") +
		"&state=" + state
}

// SpotifyTokenResponse はSpotifyトークンエンドポイントのレスポンスを表す。
type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// ExchangeCode はOAuth認可コードをアクセストークン・リフレッシュトークンに交換する。
func (s *SpotifyService) ExchangeCode(code string) (*SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.cfg.SpotifyRedirectURL)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+s.basicAuth())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.AccessToken == "" {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Spotify OAuthエラー", nil)
	}
	return &result, nil
}

// RefreshAccessToken はリフレッシュトークンで新しいアクセストークンを取得する。
func (s *SpotifyService) RefreshAccessToken(refreshToken string) (*SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+s.basicAuth())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.AccessToken == "" {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Spotifyトークン更新エラー", nil)
	}
	return &result, nil
}

// getValidToken はユーザーの有効なアクセストークンを取得する。期限切れなら自動更新する。
func (s *SpotifyService) getValidToken(user *model.User) (string, error) {
	if user.SpotifyRefreshToken == "" {
		return "", domain.NewError(domain.ErrCodeBadRequest, "Spotifyが連携されていません", nil)
	}

	// トークンが有効期限内ならそのまま返す（5分のマージン）
	if time.Now().Before(user.SpotifyTokenExpiry.Add(-5 * time.Minute)) {
		return user.SpotifyToken, nil
	}

	// リフレッシュトークンでアクセストークンを更新
	tokenResp, err := s.RefreshAccessToken(user.SpotifyRefreshToken)
	if err != nil {
		return "", err
	}

	user.SpotifyToken = tokenResp.AccessToken
	user.SpotifyTokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	if tokenResp.RefreshToken != "" {
		user.SpotifyRefreshToken = tokenResp.RefreshToken
	}

	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}

	return user.SpotifyToken, nil
}

// ConnectSpotify はOAuthコールバック後にSpotifyアカウントを連携する。
func (s *SpotifyService) ConnectSpotify(userID uint, code string) error {
	tokenResp, err := s.ExchangeCode(code)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound
	}

	user.SpotifyConnected = true
	user.SpotifyToken = tokenResp.AccessToken
	user.SpotifyRefreshToken = tokenResp.RefreshToken
	user.SpotifyTokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return s.userRepo.Update(user)
}

// DisconnectSpotify はSpotify連携を解除する。
func (s *SpotifyService) DisconnectSpotify(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound
	}

	user.SpotifyConnected = false
	user.SpotifyToken = ""
	user.SpotifyRefreshToken = ""
	user.SpotifyTokenExpiry = time.Time{}

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	s.spotifyRepo.DeleteUserData(userID)
	return nil
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

// GetCurrentlyPlaying は指定ユーザーの現在再生中の曲を取得する。
func (s *SpotifyService) GetCurrentlyPlaying(userID uint) (*SpotifyCurrentlyPlaying, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrNotFound
	}

	token, err := s.getValidToken(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Spotify API呼び出しに失敗", err)
	}
	defer resp.Body.Close()

	// 204 No Content = 再生中でない
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Spotify APIエラー", nil)
	}

	var apiResp struct {
		IsPlaying bool `json:"is_playing"`
		Item      *struct {
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Name   string `json:"name"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
				} `json:"images"`
			} `json:"album"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			DurationMs int `json:"duration_ms"`
		} `json:"item"`
		ProgressMs int `json:"progress_ms"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Item == nil {
		return nil, nil
	}

	artists := make([]string, len(apiResp.Item.Artists))
	for i, a := range apiResp.Item.Artists {
		artists[i] = a.Name
	}

	albumImage := ""
	if len(apiResp.Item.Album.Images) > 0 {
		albumImage = apiResp.Item.Album.Images[0].URL
		for _, img := range apiResp.Item.Album.Images {
			if img.Height >= 200 && img.Height <= 400 {
				albumImage = img.URL
				break
			}
		}
	}

	return &SpotifyCurrentlyPlaying{
		IsPlaying:  apiResp.IsPlaying,
		TrackName:  apiResp.Item.Name,
		ArtistName: strings.Join(artists, ", "),
		AlbumName:  apiResp.Item.Album.Name,
		AlbumImage: albumImage,
		TrackURL:   apiResp.Item.ExternalURLs.Spotify,
		ProgressMs: apiResp.ProgressMs,
		DurationMs: apiResp.Item.DurationMs,
	}, nil
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

// GetRecentlyPlayed は指定ユーザーの最近再生した曲を取得する。
func (s *SpotifyService) GetRecentlyPlayed(userID uint) ([]SpotifyRecentTrackResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrNotFound
	}

	token, err := s.getValidToken(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/recently-played?limit=10", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Spotify API呼び出しに失敗", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Spotify APIエラー", nil)
	}

	var apiResp struct {
		Items []struct {
			Track struct {
				Name    string `json:"name"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Name   string `json:"name"`
					Images []struct {
						URL    string `json:"url"`
						Height int    `json:"height"`
					} `json:"images"`
				} `json:"album"`
				ExternalURLs struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			} `json:"track"`
			PlayedAt string `json:"played_at"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	tracks := make([]SpotifyRecentTrackResponse, 0, len(apiResp.Items))
	for _, item := range apiResp.Items {
		artists := make([]string, len(item.Track.Artists))
		for i, a := range item.Track.Artists {
			artists[i] = a.Name
		}
		albumImage := ""
		if len(item.Track.Album.Images) > 0 {
			albumImage = item.Track.Album.Images[0].URL
		}

		tracks = append(tracks, SpotifyRecentTrackResponse{
			TrackName:  item.Track.Name,
			ArtistName: strings.Join(artists, ", "),
			AlbumName:  item.Track.Album.Name,
			AlbumImage: albumImage,
			TrackURL:   item.Track.ExternalURLs.Spotify,
			PlayedAt:   item.PlayedAt,
		})
	}

	return tracks, nil
}

// basicAuth はBasic認証ヘッダー値を生成する。
func (s *SpotifyService) basicAuth() string {
	return base64.StdEncoding.EncodeToString(
		[]byte(s.cfg.SpotifyClientID + ":" + s.cfg.SpotifyClientSecret),
	)
}
