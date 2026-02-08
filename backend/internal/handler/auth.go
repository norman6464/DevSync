package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// AuthHandler は認証関連のHTTPハンドラ。
// ユーザー登録・ログイン・GitHub OAuth・パスワードリセット・アカウント削除を処理する。
type AuthHandler struct {
	authService   *service.AuthService   // 認証ビジネスロジック
	githubService *service.GitHubService // GitHub OAuth連携
}

// NewAuthHandler は新しいAuthHandlerインスタンスを生成する。
func NewAuthHandler(authService *service.AuthService, githubService *service.GitHubService) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		githubService: githubService,
	}
}

// Register はユーザー新規登録を処理する。
func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login はメール・パスワードによるログインを処理する。
func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GitHubLogin はGitHub OAuthログインのURLを生成して返す。
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	state, err := h.authService.GenerateLoginState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}
	url := h.githubService.GetLoginOAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GitHubLoginCallback はGitHub OAuthコールバックを処理する。
// 認可コードをアクセストークンに交換し、GitHubユーザー情報を取得してログインする。
// ログイン後、バックグラウンドでGitHubデータを同期する。
func (h *AuthHandler) GitHubLoginCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code or state"})
		return
	}

	if err := h.authService.ValidateLoginState(state); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	accessToken, err := h.githubService.ExchangeCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code"})
		return
	}

	ghUser, err := h.githubService.GetGitHubUser(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get github user"})
		return
	}

	resp, err := h.authService.GitHubLogin(ghUser, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// バックグラウンドでGitHubデータを同期
	go h.githubService.SyncUserData(resp.User.ID)

	c.JSON(http.StatusOK, resp)
}

// Me は認証済みユーザーの情報を返す。
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.authService.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RequestPasswordReset はパスワードリセットトークンを生成する。
// 本番環境ではメール送信を行うべきだが、デモ用にトークンを直接返す。
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.RequestPasswordReset(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	if token == "" {
		// セキュリティのため、メールの存在有無を明かさない
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
		return
	}

	// 本番ではここでメール送信を行う（デモ用にトークンを返却）
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset token generated",
		"token":   token,
	})
}

// ResetPassword はトークンを使ってパスワードをリセットする。
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ResetPassword(input.Token, input.NewPassword); err != nil {
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}

// DeleteAccount はユーザーアカウントを完全に削除する。
// パスワード検証を行い、関連する全データを削除する。
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.DeleteAccount(userID, input.Password); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password required"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
