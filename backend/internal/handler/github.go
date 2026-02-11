package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// GitHubHandler はGitHub連携関連のHTTPハンドラ。
// GitHub OAuth認証・データ同期・コントリビューション取得を処理する。
type GitHubHandler struct {
	githubService *service.GitHubService
	authService   *service.AuthService
}

// NewGitHubHandler は新しいGitHubHandlerインスタンスを生成する。
func NewGitHubHandler(
	githubService *service.GitHubService,
	authService *service.AuthService,
) *GitHubHandler {
	return &GitHubHandler{
		githubService: githubService,
		authService:   authService,
	}
}

// Connect はGitHub OAuth認証のURLを生成して返す。
func (h *GitHubHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")
	state, err := h.authService.GenerateOAuthState(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}
	url := h.githubService.GetOAuthURL(state)
	respondOK(c, gin.H{"url": url})
}

// Callback はGitHub OAuthコールバックを処理してアカウントを連携する。
func (h *GitHubHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code or state"})
		return
	}

	userID, err := h.authService.ValidateOAuthState(state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	if err := h.githubService.ConnectGitHub(userID, code, state); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "github connected"})
}

// GetContributions は指定ユーザーのGitHubコントリビューション情報を取得する。
func (h *GitHubHandler) GetContributions(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	contributions, err := h.githubService.GetContributions(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, contributions)
}

// GetLanguages は指定ユーザーのGitHubリポジトリの使用言語統計を取得する。
func (h *GitHubHandler) GetLanguages(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	stats, err := h.githubService.GetLanguages(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, stats)
}

// GetRepos は指定ユーザーのGitHubリポジトリ一覧を取得する。
func (h *GitHubHandler) GetRepos(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	repos, err := h.githubService.GetRepos(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, repos)
}

// Sync は現在のユーザーのGitHubデータを手動で同期する。
func (h *GitHubHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.githubService.SyncUserData(userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "sync complete"})
}

// Disconnect は現在のユーザーのGitHub連携を解除する。
func (h *GitHubHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.githubService.DisconnectGitHub(userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "github disconnected"})
}
