package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// AuthServiceInterface は認証サービスの抽象インターフェース。
type AuthServiceInterface interface {
	Register(input service.RegisterInput) (*service.AuthResponse, error)
	Login(input service.LoginInput) (*service.AuthResponse, error)
	GenerateLoginState() (string, error)
	ValidateLoginState(state string) error
	GitHubLogin(ghUser *service.GitHubUserInfo, accessToken string) (*service.AuthResponse, error)
	GetMe(userID uint) (*model.User, error)
	RequestPasswordReset(email string) (string, error)
	ResetPassword(token string, newPassword string) error
	DeleteAccount(userID uint, password string) error
}

// AuthGitHubServiceInterface はAuthHandler用のGitHubサービスの抽象インターフェース。
type AuthGitHubServiceInterface interface {
	GetLoginOAuthURL(state string) string
	ExchangeCode(code string) (string, error)
	GetGitHubUser(token string) (*service.GitHubUserInfo, error)
	SyncUserData(userID uint) error
}

// AuthHandler は認証関連のHTTPハンドラ。
// ユーザー登録・ログイン・GitHub OAuth・パスワードリセット・アカウント削除を処理する。
type AuthHandler struct {
	authService   AuthServiceInterface       // 認証ビジネスロジック
	githubService AuthGitHubServiceInterface // GitHub OAuth連携
}

// NewAuthHandler は新しいAuthHandlerインスタンスを生成する。
func NewAuthHandler(authService AuthServiceInterface, githubService AuthGitHubServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		githubService: githubService,
	}
}

// cookieMaxAge はJWTトークンの有効期限と同じ72時間（秒単位）。
const cookieMaxAge = 72 * 60 * 60

// setAuthCookie はJWTトークンをhttpOnly Cookieとしてレスポンスにセットする。
// 本番環境（ENVIRONMENT=production）ではSecure属性をtrueに設定する。
func setAuthCookie(c *gin.Context, token string) {
	secure := os.Getenv("ENVIRONMENT") == "production"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, cookieMaxAge, "/", "", secure, true)
}

// clearAuthCookie は認証Cookieをクリアする（MaxAge=-1で即時削除）。
func clearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "/", "", false, true)
}

// Register はユーザー新規登録を処理する。
func (h *AuthHandler) Register(c *gin.Context) {
	req := bindJSON[dto.RegisterRequest](c)
	if req == nil {
		return
	}

	// DTOをservice層の入力に変換
	input := service.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.authService.Register(input)
	if err != nil {
		respondError(c, err)
		return
	}

	setAuthCookie(c, resp.Token)
	respondCreated(c, gin.H{"user": resp.User})
}

// Login はメール・パスワードによるログインを処理する。
func (h *AuthHandler) Login(c *gin.Context) {
	req := bindJSON[dto.LoginRequest](c)
	if req == nil {
		return
	}

	// DTOをservice層の入力に変換
	input := service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.authService.Login(input)
	if err != nil {
		respondError(c, err)
		return
	}

	setAuthCookie(c, resp.Token)
	respondOK(c, gin.H{"user": resp.User})
}

// GitHubLogin はGitHub OAuthログインのURLを生成して返す。
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	state, err := h.authService.GenerateLoginState()
	if err != nil {
		respondError(c, err)
		return
	}
	url := h.githubService.GetLoginOAuthURL(state)
	respondOK(c, dto.URLResponse{URL: url})
}

// GitHubLoginCallback はGitHub OAuthコールバックを処理する。
// 認可コードをアクセストークンに交換し、GitHubユーザー情報を取得してログインする。
// ログイン後、バックグラウンドでGitHubデータを同期する。
func (h *AuthHandler) GitHubLoginCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		respondError(c, service.ErrBadRequest)
		return
	}

	if err := h.authService.ValidateLoginState(state); err != nil {
		respondError(c, err)
		return
	}

	accessToken, err := h.githubService.ExchangeCode(code)
	if err != nil {
		respondError(c, err)
		return
	}

	ghUser, err := h.githubService.GetGitHubUser(accessToken)
	if err != nil {
		respondError(c, err)
		return
	}

	resp, err := h.authService.GitHubLogin(ghUser, accessToken)
	if err != nil {
		respondError(c, err)
		return
	}

	// バックグラウンドでGitHubデータを同期
	go h.githubService.SyncUserData(resp.User.ID)

	setAuthCookie(c, resp.Token)
	respondOK(c, gin.H{"user": resp.User})
}

// Me は認証済みユーザーの情報を返す。
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.authService.GetMe(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, user)
}

// RequestPasswordReset はパスワードリセットトークンを生成する。
// 本番環境ではメール送信を行うべきだが、デモ用にトークンを直接返す。
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	req := bindJSON[dto.PasswordResetRequest](c)
	if req == nil {
		return
	}

	token, err := h.authService.RequestPasswordReset(req.Email)
	if err != nil {
		respondError(c, err)
		return
	}

	// セキュリティのため、メールの存在有無に関わらず同一レスポンスを返す
	// 本番ではここでメール送信を行う
	// token変数はメール送信時に使用する（レスポンスには含めない）
	_ = token
	respondOK(c, dto.MessageResponse{Message: "If the email exists, a reset link has been sent"})
}

// ResetPassword はトークンを使ってパスワードをリセットする。
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	req := bindJSON[dto.ResetPasswordRequest](c)
	if req == nil {
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.MessageResponse{Message: "Password has been reset successfully"})
}

// Logout はユーザーのログアウトを処理し、認証Cookieをクリアする。
func (h *AuthHandler) Logout(c *gin.Context) {
	clearAuthCookie(c)
	respondOK(c, dto.MessageResponse{Message: "logged out successfully"})
}

// DeleteAccount はユーザーアカウントを完全に削除する。
// パスワード検証を行い、関連する全データを削除する。
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.DeleteAccountRequest](c)
	if req == nil {
		return
	}

	if err := h.authService.DeleteAccount(userID, req.Password); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.MessageResponse{Message: "Account deleted successfully"})
}
