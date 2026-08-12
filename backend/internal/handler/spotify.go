package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// SpotifyUseCases は SpotifyHandler が依存する Spotify 連携の usecase をまとめる。
type SpotifyUseCases struct {
	OAuthURL         *usecase.GetSpotifyOAuthURLUseCase
	Connect          *usecase.ConnectSpotifyUseCase
	Disconnect       *usecase.DisconnectSpotifyUseCase
	CurrentlyPlaying *usecase.GetSpotifyCurrentlyPlayingUseCase
	RecentlyPlayed   *usecase.GetSpotifyRecentlyPlayedUseCase
}

// SpotifyHandler はSpotify連携関連のHTTPハンドラ。
type SpotifyHandler struct {
	uc         SpotifyUseCases
	oauthState *usecase.OAuthStateUseCase
}

// NewSpotifyHandler は新しいSpotifyHandlerインスタンスを生成する。
func NewSpotifyHandler(uc SpotifyUseCases, oauthState *usecase.OAuthStateUseCase) *SpotifyHandler {
	return &SpotifyHandler{uc: uc, oauthState: oauthState}
}

// Connect はSpotify OAuth認証URLを生成して返す。
func (h *SpotifyHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")
	state, err := h.oauthState.Generate(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	url := h.uc.OAuthURL.Execute(state)
	respondOK(c, dto.URLResponse{URL: url})
}

// Callback はSpotify OAuthコールバックを処理する。
func (h *SpotifyHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" || len(code) > oauthCodeMaxLen || len(state) > oauthCodeMaxLen {
		respondBadRequest(c, "codeまたはstateが不正です")
		return
	}

	userID, err := h.oauthState.Validate(state)
	if err != nil {
		respondBadRequest(c, "stateが無効です")
		return
	}

	if err := h.uc.Connect.Execute(c.Request.Context(), userID, code); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("Spotify連携が完了しました"))
}

// GetCurrentlyPlaying は指定ユーザーの現在再生中の曲を返す。
func (h *SpotifyHandler) GetCurrentlyPlaying(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	track, err := h.uc.CurrentlyPlaying.Execute(c.Request.Context(), userID)
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

	tracks, err := h.uc.RecentlyPlayed.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(tracks))
}

// Disconnect はSpotify連携を解除する。
func (h *SpotifyHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.uc.Disconnect.Execute(c.Request.Context(), userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("Spotify連携を解除しました"))
}
