package handler

import (
	"net/http"
	"os"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// oauthCodeMaxLen はOAuth認可コード・stateパラメータの最大許容長。
const oauthCodeMaxLen = 2048

// AuthUseCases は AuthHandler が依存する認証の usecase をまとめる。
type AuthUseCases struct {
	Register             *usecase.RegisterUserUseCase
	Login                *usecase.LoginUseCase
	GitHubLogin          *usecase.GitHubLoginUseCase
	LoginState           *usecase.GitHubLoginStateUseCase
	GetMe                *usecase.GetMeUseCase
	RequestPasswordReset *usecase.RequestPasswordResetUseCase
	ResetPassword        *usecase.ResetPasswordUseCase
	DeleteAccount        *usecase.DeleteAccountUseCase
}

// AuthGitHubUseCases は GitHub ログインで使う GitHub 連携の usecase をまとめる。
type AuthGitHubUseCases struct {
	LoginURL     *usecase.GetGitHubLoginURLUseCase
	ExchangeCode *usecase.ExchangeGitHubCodeUseCase
	GetUser      *usecase.GetGitHubUserUseCase
	Sync         *usecase.SyncGitHubDataUseCase
}

// AuthHandler は認証関連のHTTPハンドラ。
// ユーザー登録・ログイン・GitHub OAuth・パスワードリセット・アカウント削除を処理する。
type AuthHandler struct {
	uc     AuthUseCases       // 認証ビジネスロジック
	github AuthGitHubUseCases // GitHub OAuth 連携
}

// NewAuthHandler は新しいAuthHandlerインスタンスを生成する。
func NewAuthHandler(uc AuthUseCases, github AuthGitHubUseCases) *AuthHandler {
	return &AuthHandler{uc: uc, github: github}
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
	secure := os.Getenv("ENVIRONMENT") == "production"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "/", "", secure, true)
}

// userResponse はユーザー情報レスポンス（認証後のレスポンスで使用）
type userResponse struct {
	User model.User `json:"user"`
}

// registerRequest は新規ユーザー登録リクエスト
type registerRequest struct {
	Name            string `json:"name" binding:"required,max=50" validate:"required,max=50"`
	Username        string `json:"username" binding:"required,max=30" validate:"required,max=30"`
	Email           string `json:"email" binding:"required,email,max=255" validate:"required,email,max=255"`
	Password        string `json:"password" binding:"required,min=8" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required"`
}

// Register はユーザー新規登録を処理する。
func (h *AuthHandler) Register(c *gin.Context) {
	req := bindJSON[registerRequest](c)
	if req == nil {
		return
	}

	// DTO を usecase の入力に変換する
	input := usecase.AuthUserInput{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.uc.Register.Execute(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}

	setAuthCookie(c, resp.Token)
	respondCreated(c, userResponse{User: resp.User})
}

// loginRequest はログインリクエスト
type loginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255" validate:"required,email,max=255"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// Login はメール・パスワードによるログインを処理する。
func (h *AuthHandler) Login(c *gin.Context) {
	req := bindJSON[loginRequest](c)
	if req == nil {
		return
	}

	// DTO を usecase の入力に変換する
	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.uc.Login.Execute(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}

	setAuthCookie(c, resp.Token)
	respondOK(c, userResponse{User: resp.User})
}

// urlResponse はURLを返すレスポンス（OAuth等で使用）。
type urlResponse struct {
	URL string `json:"url"`
}

// GitHubLogin はGitHub OAuthログインのURLを生成して返す。
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	state, err := h.uc.LoginState.Generate()
	if err != nil {
		respondError(c, err)
		return
	}
	url := h.github.LoginURL.Execute(state)
	respondOK(c, urlResponse{URL: url})
}

// GitHubLoginCallback はGitHub OAuthコールバックを処理する。
// 認可コードをアクセストークンに交換し、GitHubユーザー情報を取得してログインする。
// ログイン後、バックグラウンドでGitHubデータを同期する。
func (h *AuthHandler) GitHubLoginCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" || len(code) > oauthCodeMaxLen || len(state) > oauthCodeMaxLen {
		respondError(c, domain.ErrBadRequest)
		return
	}

	if err := h.uc.LoginState.Validate(state); err != nil {
		respondError(c, err)
		return
	}

	accessToken, err := h.github.ExchangeCode.Execute(c.Request.Context(), code)
	if err != nil {
		respondError(c, err)
		return
	}

	ghUser, err := h.github.GetUser.Execute(c.Request.Context(), accessToken)
	if err != nil {
		respondError(c, err)
		return
	}

	resp, err := h.uc.GitHubLogin.Execute(c.Request.Context(), ghUser, accessToken)
	if err != nil {
		respondError(c, err)
		return
	}

	// バックグラウンドでGitHubデータを同期
	h.github.Sync.SyncInBackground(c.Request.Context(), resp.User.ID)

	setAuthCookie(c, resp.Token)
	respondOK(c, userResponse{User: resp.User})
}

// Me は認証済みユーザーの情報を返す。
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.uc.GetMe.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, user)
}

// passwordResetRequest はパスワードリセットリクエスト
type passwordResetRequest struct {
	Email string `json:"email" binding:"required,email,max=255" validate:"required,email,max=255"`
}

// messageResponse は汎用メッセージレスポンス
type messageResponse struct {
	Message string `json:"message"`
}

// RequestPasswordReset はパスワードリセットトークンを生成する。
// 本番環境ではメール送信を行うべきだが、デモ用にトークンを直接返す。
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	req := bindJSON[passwordResetRequest](c)
	if req == nil {
		return
	}

	token, err := h.uc.RequestPasswordReset.Execute(c.Request.Context(), req.Email)
	if err != nil {
		respondError(c, err)
		return
	}

	// セキュリティのため、メールの存在有無に関わらず同一レスポンスを返す
	// 本番ではここでメール送信を行う
	// token変数はメール送信時に使用する（レスポンスには含めない）
	_ = token
	respondOK(c, messageResponse{Message: "メールアドレスが登録されている場合、リセットリンクを送信しました"})
}

// resetPasswordRequest はトークンによるパスワードリセットリクエスト
type resetPasswordRequest struct {
	Token       string `json:"token" binding:"required" validate:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8" validate:"required,min=8"`
}

// ResetPassword はトークンを使ってパスワードをリセットする。
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	req := bindJSON[resetPasswordRequest](c)
	if req == nil {
		return
	}

	if err := h.uc.ResetPassword.Execute(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, messageResponse{Message: "パスワードをリセットしました"})
}

// Logout はユーザーのログアウトを処理し、認証Cookieをクリアする。
func (h *AuthHandler) Logout(c *gin.Context) {
	clearAuthCookie(c)
	respondOK(c, messageResponse{Message: "ログアウトしました"})
}

// deleteAccountRequest はアカウント削除リクエスト
type deleteAccountRequest struct {
	Password string `json:"password" binding:"required" validate:"required"`
}

// DeleteAccount はユーザーアカウントを完全に削除する。
// パスワード検証を行い、関連する全データを削除する。
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[deleteAccountRequest](c)
	if req == nil {
		return
	}

	if err := h.uc.DeleteAccount.Execute(c.Request.Context(), userID, req.Password); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, messageResponse{Message: "アカウントを削除しました"})
}
