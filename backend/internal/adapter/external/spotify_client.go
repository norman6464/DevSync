package external

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// spotifyRequestTimeout は Spotify への 1 リクエストのタイムアウト。
const spotifyRequestTimeout = 30 * time.Second

const (
	spotifyAuthorizeURL        = "https://accounts.spotify.com/authorize"
	spotifyTokenURL            = "https://accounts.spotify.com/api/token"
	spotifyCurrentlyPlayingURL = "https://api.spotify.com/v1/me/player/currently-playing"
	spotifyRecentlyPlayedURL   = "https://api.spotify.com/v1/me/player/recently-played?limit=10"
	// spotifyScopes は連携時に要求するスコープ。
	spotifyScopes = "user-read-currently-playing user-read-recently-played"
)

// spotifyClient は [repository.SpotifyAPIClient] の HTTP 実装。
type spotifyClient struct {
	clientID     string
	clientSecret string
	redirectURL  string
	httpClient   *http.Client
	// tokenURL はトークンエンドポイント。テストで fake サーバーへ差し替えるためフィールドで持つ。
	tokenURL string
}

// NewSpotifyClient は SpotifyAPIClient の HTTP 実装を返す。
func NewSpotifyClient(clientID, clientSecret, redirectURL string) repository.SpotifyAPIClient {
	return &spotifyClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		httpClient:   &http.Client{Timeout: spotifyRequestTimeout},
		tokenURL:     spotifyTokenURL,
	}
}

var _ repository.SpotifyAPIClient = (*spotifyClient)(nil)

// AuthorizeURL は連携用の OAuth 認可 URL を返す。
// クライアント ID と state も含めてすべてのパラメータをエスケープする。
func (c *spotifyClient) AuthorizeURL(state string) string {
	q := url.Values{}
	q.Set("client_id", c.clientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", c.redirectURL)
	q.Set("scope", spotifyScopes)
	q.Set("state", state)
	return spotifyAuthorizeURL + "?" + q.Encode()
}

// ExchangeCode は認可コードをトークンに交換する。
func (c *spotifyClient) ExchangeCode(ctx context.Context, code string) (*model.SpotifyToken, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", c.redirectURL)

	return c.requestToken(ctx, form, "Spotify OAuthエラー")
}

// RefreshAccessToken はリフレッシュトークンでアクセストークンを更新する。
func (c *spotifyClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*model.SpotifyToken, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	return c.requestToken(ctx, form, "Spotifyトークン更新エラー")
}

// requestToken はトークンエンドポイントを呼び出す。
// クライアント側起因の失敗（4xx: invalid_grant / invalid_client 等）は 400 系、
// Spotify 側の障害（5xx）やネットワーク障害は 503 として返す。
func (c *spotifyClient) requestToken(ctx context.Context, form url.Values, errMessage string) (*model.SpotifyToken, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+c.basicAuth())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, errMessage, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 原因調査のため Spotify のエラー内容をログへ残す（トークン本体はログに出さない）。
		var apiErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		log.Printf("[WARN] Spotify トークンエンドポイントが %d を返却: error=%q error_description=%q",
			resp.StatusCode, apiErr.Error, apiErr.ErrorDescription)

		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, domain.NewError(domain.ErrCodeBadRequest, errMessage, nil)
		}
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, errMessage, nil)
	}

	var token model.SpotifyToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, errMessage, nil)
	}
	return &token, nil
}

// FetchCurrentlyPlaying は現在再生中の曲を返す。再生していない場合は (nil, nil) を返す。
func (c *spotifyClient) FetchCurrentlyPlaying(ctx context.Context, token string) (*model.SpotifyCurrentlyPlaying, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, spotifyCurrentlyPlayingURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
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

	// 一覧表示に使うため、200〜400px のジャケット画像を優先する
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

	return &model.SpotifyCurrentlyPlaying{
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

// FetchRecentlyPlayed は最近再生した曲を新しい順に返す。
func (c *spotifyClient) FetchRecentlyPlayed(ctx context.Context, token string) ([]model.SpotifyRecentTrackResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, spotifyRecentlyPlayedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
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

	tracks := make([]model.SpotifyRecentTrackResponse, 0, len(apiResp.Items))
	for _, item := range apiResp.Items {
		artists := make([]string, len(item.Track.Artists))
		for i, a := range item.Track.Artists {
			artists[i] = a.Name
		}
		albumImage := ""
		if len(item.Track.Album.Images) > 0 {
			albumImage = item.Track.Album.Images[0].URL
		}

		tracks = append(tracks, model.SpotifyRecentTrackResponse{
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

// basicAuth は Basic 認証ヘッダーの値を生成する。
func (c *spotifyClient) basicAuth() string {
	return base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
}
