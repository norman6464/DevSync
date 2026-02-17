package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
)

// GitHubServiceInterface はGitHubサービスの抽象インターフェース。
type GitHubServiceInterface interface {
	GetOAuthURL(state string) string
	ConnectGitHub(userID uint, code, state string) error
	GetContributions(userID uint) ([]model.GitHubContribution, error)
	GetLanguages(userID uint) ([]model.GitHubLanguageStat, error)
	GetRepos(userID uint) ([]model.GitHubRepository, error)
	SyncUserData(userID uint) error
	DisconnectGitHub(userID uint) error
}

// GitHubAuthServiceInterface はGitHubHandler用の認証サービスの抽象インターフェース。
type GitHubAuthServiceInterface interface {
	GenerateOAuthState(userID uint) (string, error)
	ValidateOAuthState(state string) (uint, error)
}

// GitHubHandler はGitHub連携関連のHTTPハンドラ。
// GitHub OAuth認証・データ同期・コントリビューション取得を処理する。
type GitHubHandler struct {
	githubService GitHubServiceInterface
	authService   GitHubAuthServiceInterface
}

// NewGitHubHandler は新しいGitHubHandlerインスタンスを生成する。
func NewGitHubHandler(
	githubService GitHubServiceInterface,
	authService GitHubAuthServiceInterface,
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
		respondError(c, err)
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
		respondBadRequest(c, "missing code or state")
		return
	}

	userID, err := h.authService.ValidateOAuthState(state)
	if err != nil {
		respondBadRequest(c, "invalid state")
		return
	}

	if err := h.githubService.ConnectGitHub(userID, code, state); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("github connected"))
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

	respondOK(c, domain.NewMessageResponse("sync complete"))
}

// Disconnect は現在のユーザーのGitHub連携を解除する。
func (h *GitHubHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.githubService.DisconnectGitHub(userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("github disconnected"))
}
