package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/service"
)

// SpotifyServiceInterface はSpotifyサービスの抽象インターフェース。
type SpotifyServiceInterface interface {
	GetOAuthURL(state string) string
	ConnectSpotify(userID uint, code string) error
	DisconnectSpotify(userID uint) error
	GetCurrentlyPlaying(userID uint) (*service.SpotifyCurrentlyPlaying, error)
	GetRecentlyPlayed(userID uint) ([]service.SpotifyRecentTrackResponse, error)
}

// SpotifyAuthServiceInterface はSpotifyHandler用の認証サービスの抽象インターフェース。
type SpotifyAuthServiceInterface interface {
	GenerateOAuthState(userID uint) (string, error)
	ValidateOAuthState(state string) (uint, error)
}

// SpotifyHandler はSpotify連携関連のHTTPハンドラ。
type SpotifyHandler struct {
	spotifyService SpotifyServiceInterface
	authService    SpotifyAuthServiceInterface
}

// NewSpotifyHandler は新しいSpotifyHandlerインスタンスを生成する。
func NewSpotifyHandler(
	spotifyService SpotifyServiceInterface,
	authService SpotifyAuthServiceInterface,
) *SpotifyHandler {
	return &SpotifyHandler{
		spotifyService: spotifyService,
		authService:    authService,
	}
}

// Connect はSpotify OAuth認証URLを生成して返す。
func (h *SpotifyHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")
	state, err := h.authService.GenerateOAuthState(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	url := h.spotifyService.GetOAuthURL(state)
	respondOK(c, dto.URLResponse{URL: url})
}

// Callback はSpotify OAuthコールバックを処理する。
func (h *SpotifyHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" || len(code) > oauthCodeMaxLen || len(state) > oauthCodeMaxLen {
		respondBadRequest(c, "missing or invalid code/state")
		return
	}

	userID, err := h.authService.ValidateOAuthState(state)
	if err != nil {
		respondBadRequest(c, "invalid state")
		return
	}

	if err := h.spotifyService.ConnectSpotify(userID, code); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("spotify connected"))
}

// GetCurrentlyPlaying は指定ユーザーの現在再生中の曲を返す。
func (h *SpotifyHandler) GetCurrentlyPlaying(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	track, err := h.spotifyService.GetCurrentlyPlaying(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, track)
}

// GetRecentlyPlayed は指定ユーザーの最近再生した曲を返す。
func (h *SpotifyHandler) GetRecentlyPlayed(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	tracks, err := h.spotifyService.GetRecentlyPlayed(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(tracks))
}

// Disconnect はSpotify連携を解除する。
func (h *SpotifyHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.spotifyService.DisconnectSpotify(userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("spotify disconnected"))
}
