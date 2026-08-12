package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// GitHubAuthServiceInterface はGitHubHandler用の認証サービスの抽象インターフェース。
type GitHubAuthServiceInterface interface {
	GenerateOAuthState(userID uint) (string, error)
	ValidateOAuthState(state string) (uint, error)
}

// GitHubUseCases は GitHubHandler が依存する GitHub 連携の usecase をまとめる。
type GitHubUseCases struct {
	OAuthURL      *usecase.GetGitHubOAuthURLUseCase
	Connect       *usecase.ConnectGitHubUseCase
	Disconnect    *usecase.DisconnectGitHubUseCase
	Sync          *usecase.SyncGitHubDataUseCase
	Contributions *usecase.GetGitHubContributionsUseCase
	Languages     *usecase.GetGitHubLanguagesUseCase
	Repos         *usecase.GetGitHubReposUseCase
}

// GitHubHandler はGitHub連携関連のHTTPハンドラ。
// GitHub OAuth認証・データ同期・コントリビューション取得を処理する。
type GitHubHandler struct {
	uc          GitHubUseCases
	authService GitHubAuthServiceInterface
}

// NewGitHubHandler は新しいGitHubHandlerインスタンスを生成する。
func NewGitHubHandler(uc GitHubUseCases, authService GitHubAuthServiceInterface) *GitHubHandler {
	return &GitHubHandler{uc: uc, authService: authService}
}

// Connect はGitHub OAuth認証のURLを生成して返す。
func (h *GitHubHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")
	state, err := h.authService.GenerateOAuthState(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	url := h.uc.OAuthURL.Execute(state)
	respondOK(c, dto.URLResponse{URL: url})
}

// Callback はGitHub OAuthコールバックを処理してアカウントを連携する。
func (h *GitHubHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" || len(code) > oauthCodeMaxLen || len(state) > oauthCodeMaxLen {
		respondBadRequest(c, "codeまたはstateが不正です")
		return
	}

	userID, err := h.authService.ValidateOAuthState(state)
	if err != nil {
		respondBadRequest(c, "stateが無効です")
		return
	}

	if err := h.uc.Connect.Execute(c.Request.Context(), userID, code); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("GitHub連携が完了しました"))
}

// GetContributions は指定ユーザーのGitHubコントリビューション情報を取得する。
func (h *GitHubHandler) GetContributions(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	contributions, err := h.uc.Contributions.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(contributions))
}

// GetLanguages は指定ユーザーのGitHubリポジトリの使用言語統計を取得する。
func (h *GitHubHandler) GetLanguages(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	stats, err := h.uc.Languages.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(stats))
}

// GetRepos は指定ユーザーのGitHubリポジトリ一覧を取得する。
func (h *GitHubHandler) GetRepos(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	repos, err := h.uc.Repos.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(repos))
}

// Sync は現在のユーザーのGitHubデータを手動で同期する。
func (h *GitHubHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.uc.Sync.Execute(c.Request.Context(), userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("同期が完了しました"))
}

// Disconnect は現在のユーザーのGitHub連携を解除する。
func (h *GitHubHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.uc.Disconnect.Execute(c.Request.Context(), userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("GitHub連携を解除しました"))
}
