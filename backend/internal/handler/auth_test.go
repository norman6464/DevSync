// Package handler は認証ハンドラーのテストを提供する。
package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// testJWTSecret はテスト用の JWT 署名鍵。
const testJWTSecret = "test-secret-key-for-handler-tests"

// mockAuthUserRepo は usecase/repository.AuthUserRepository のモック。
type mockAuthUserRepo struct{ mock.Mock }

func (m *mockAuthUserRepo) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthUserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthUserRepo) FindByGitHubID(ctx context.Context, githubID int64) (*model.User, error) {
	args := m.Called(ctx, githubID)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthUserRepo) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil && user.ID == 0 {
		user.ID = 1 // 作成後に ID が入る挙動を模す
	}
	return args.Error(0)
}

func (m *mockAuthUserRepo) Update(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *mockAuthUserRepo) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	return m.Called(ctx, userID, hashedPassword).Error(0)
}

func (m *mockAuthUserRepo) DeleteWithRelatedData(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// mockPasswordResetRepo は usecase/repository.PasswordResetTokenRepository のモック。
type mockPasswordResetRepo struct{ mock.Mock }

func (m *mockPasswordResetRepo) Create(ctx context.Context, token *model.PasswordResetToken) error {
	return m.Called(ctx, token).Error(0)
}

func (m *mockPasswordResetRepo) FindByToken(ctx context.Context, hashedToken string) (*model.PasswordResetToken, error) {
	args := m.Called(ctx, hashedToken)
	t, _ := args.Get(0).(*model.PasswordResetToken)
	return t, args.Error(1)
}

func (m *mockPasswordResetRepo) MarkAsUsed(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPasswordResetRepo) InvalidateUserTokens(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

// authPorts は認証の usecase に注入した port モックをまとめる。
type authPorts struct {
	Users  *mockAuthUserRepo
	Tokens *mockPasswordResetRepo
	GitHub *githubPorts
}

// setupAuthHandler は本物の usecase に port モックを注入した AuthHandler とルーターを生成する。
func setupAuthHandler() (*gin.Engine, *AuthHandler, *authPorts) {
	gin.SetMode(gin.TestMode)
	ports := &authPorts{
		Users:  new(mockAuthUserRepo),
		Tokens: new(mockPasswordResetRepo),
		GitHub: newGitHubPorts(),
	}

	h := NewAuthHandler(AuthUseCases{
		Register:             usecase.NewRegisterUserUseCase(ports.Users, testJWTSecret),
		Login:                usecase.NewLoginUseCase(ports.Users, testJWTSecret),
		GitHubLogin:          usecase.NewGitHubLoginUseCase(ports.Users, testJWTSecret),
		LoginState:           usecase.NewGitHubLoginStateUseCase(testJWTSecret),
		GetMe:                usecase.NewGetMeUseCase(ports.Users),
		RequestPasswordReset: usecase.NewRequestPasswordResetUseCase(ports.Users, ports.Tokens),
		ResetPassword:        usecase.NewResetPasswordUseCase(ports.Users, ports.Tokens),
		DeleteAccount:        usecase.NewDeleteAccountUseCase(ports.Users),
	}, AuthGitHubUseCases{
		LoginURL:     usecase.NewGetGitHubLoginURLUseCase(ports.GitHub.Client),
		ExchangeCode: usecase.NewExchangeGitHubCodeUseCase(ports.GitHub.Client),
		GetUser:      usecase.NewGetGitHubUserUseCase(ports.GitHub.Client),
		Sync:         usecase.NewSyncGitHubDataUseCase(ports.GitHub.Users, ports.GitHub.Repo, ports.GitHub.Client),
	})

	r := gin.New()
	r.POST("/auth/login", h.Login)
	r.POST("/auth/register", h.Register)
	r.POST("/auth/password-reset", h.RequestPasswordReset)
	r.POST("/auth/password-reset/confirm", h.ResetPassword)
	r.GET("/auth/github", h.GitHubLogin)
	r.GET("/auth/github/callback", h.GitHubLoginCallback)

	authorized := r.Group("/auth")
	authorized.Use(func(c *gin.Context) { c.Set("userID", uint(1)); c.Next() })
	authorized.GET("/me", h.Me)
	authorized.POST("/logout", h.Logout)
	authorized.DELETE("/account", h.DeleteAccount)

	return r, h, ports
}

// tokenCookie はレスポンスから token Cookie を取り出す。
func tokenCookie(w *http.Response) *http.Cookie {
	for _, c := range w.Cookies() {
		if c.Name == "token" {
			return c
		}
	}
	return nil
}

// hashedPassword はテスト用に bcrypt ハッシュを作る。
func hashedPassword(plain string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(h)
}

// ---------- ログイン ----------

func TestLogin_SetsCookie(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "test@example.com").Return(&model.User{
		ID: 1, Email: "test@example.com", Password: hashedPassword("password123"), Name: "Test User",
	}, nil)

	w := doRequest(r, http.MethodPost, "/auth/login", map[string]string{
		"email": "test@example.com", "password": "password123",
	})
	assertStatus(t, w, http.StatusOK)

	cookie := tokenCookie(w.Result())
	if assert.NotNil(t, cookie, "token Cookie がセットされるべき") {
		assert.True(t, cookie.HttpOnly, "httpOnly であるべき")
		assert.Equal(t, "/", cookie.Path)
		assert.Positive(t, cookie.MaxAge)
	}
	assert.Contains(t, w.Body.String(), `"user"`)
	ports.Users.AssertExpectations(t)
}

// パスワードが違えば 401。
func TestLogin_WrongPassword(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "test@example.com").Return(&model.User{
		ID: 1, Email: "test@example.com", Password: hashedPassword("password123"),
	}, nil)

	w := doRequest(r, http.MethodPost, "/auth/login", map[string]string{
		"email": "test@example.com", "password": "wrong-password",
	})
	assertStatus(t, w, http.StatusUnauthorized)
	assert.Nil(t, tokenCookie(w.Result()))
}

// 未登録のメールアドレスも 401（存在有無を漏らさない）。
func TestLogin_UserNotFound(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "unknown@example.com").Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/auth/login", map[string]string{
		"email": "unknown@example.com", "password": "password123",
	})
	assertStatus(t, w, http.StatusUnauthorized)
	ports.Users.AssertExpectations(t)
}

// ---------- 登録 ----------

func TestRegister_SetsCookie(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, nil)
	ports.Users.On("FindByUsername", mock.Anything, "newuser").Return(nil, nil)
	ports.Users.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		// パスワードはハッシュ化して保存する
		return u.Email == "new@example.com" && u.Password != "password123" && u.Password != ""
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/auth/register", map[string]string{
		"name": "NewUser", "username": "newuser", "email": "new@example.com",
		"password": "password123", "confirm_password": "password123",
	})
	assertStatus(t, w, http.StatusCreated)

	cookie := tokenCookie(w.Result())
	if assert.NotNil(t, cookie) {
		assert.True(t, cookie.HttpOnly)
	}
	ports.Users.AssertExpectations(t)
}

// メールアドレスが登録済みなら 409。
func TestRegister_DuplicateEmail(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "dup@example.com").Return(&model.User{ID: 2}, nil)

	w := doRequest(r, http.MethodPost, "/auth/register", map[string]string{
		"name": "Dup", "username": "dupuser", "email": "dup@example.com",
		"password": "password123", "confirm_password": "password123",
	})
	assertStatus(t, w, http.StatusConflict)
	ports.Users.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// ユーザー名が使用済みなら 409。
func TestRegister_DuplicateUsername(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, nil)
	ports.Users.On("FindByUsername", mock.Anything, "taken").Return(&model.User{ID: 3}, nil)

	w := doRequest(r, http.MethodPost, "/auth/register", map[string]string{
		"name": "New", "username": "taken", "email": "new@example.com",
		"password": "password123", "confirm_password": "password123",
	})
	assertStatus(t, w, http.StatusConflict)
	ports.Users.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegister_InvalidBody(t *testing.T) {
	r, _, ports := setupAuthHandler()

	w := doRequestRaw(r, http.MethodPost, "/auth/register", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
	ports.Users.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegister_CreateError(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, nil)
	ports.Users.On("FindByUsername", mock.Anything, "newuser").Return(nil, nil)
	ports.Users.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/auth/register", map[string]string{
		"name": "NewUser", "username": "newuser", "email": "new@example.com",
		"password": "password123", "confirm_password": "password123",
	})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- ログアウト ----------

func TestLogout_ClearsCookie(t *testing.T) {
	r, _, _ := setupAuthHandler()

	w := doRequest(r, http.MethodPost, "/auth/logout", nil)
	assertStatus(t, w, http.StatusOK)

	cookie := tokenCookie(w.Result())
	if assert.NotNil(t, cookie, "削除用の token Cookie が設定されるべき") {
		assert.Negative(t, cookie.MaxAge, "MaxAge が負で Cookie を削除する")
	}
	assert.Contains(t, w.Body.String(), "message")
}

// ---------- Me ----------

func TestMe_Success(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, Name: "Test User", Email: "test@example.com"}, nil)

	w := doRequest(r, http.MethodGet, "/auth/me", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "test@example.com")
	ports.Users.AssertExpectations(t)
}

// ユーザーが存在しない場合は移行前と同じく内部エラーになる。
func TestMe_UserNotFound(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/auth/me", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Users.AssertExpectations(t)
}

// ---------- 退会 ----------

func TestDeleteAccount_Success(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, Password: hashedPassword("password123")}, nil)
	ports.Users.On("DeleteWithRelatedData", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/auth/account", map[string]string{"password": "password123"})
	assertStatus(t, w, http.StatusOK)
	ports.Users.AssertExpectations(t)
}

func TestDeleteAccount_WrongPassword(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, Password: hashedPassword("password123")}, nil)

	w := doRequest(r, http.MethodDelete, "/auth/account", map[string]string{"password": "wrong"})
	assertStatus(t, w, http.StatusForbidden)
	ports.Users.AssertNotCalled(t, "DeleteWithRelatedData", mock.Anything, mock.Anything)
}

// パスワード未入力は DTO の binding で弾かれる（usecase まで届かない）。
// GitHub のみで登録したユーザーのパスワード検証スキップは usecase テストで検証する。
func TestDeleteAccount_EmptyPassword(t *testing.T) {
	r, _, ports := setupAuthHandler()

	w := doRequest(r, http.MethodDelete, "/auth/account", map[string]string{"password": ""})
	assertStatus(t, w, http.StatusBadRequest)
	ports.Users.AssertNotCalled(t, "DeleteWithRelatedData", mock.Anything, mock.Anything)
}

func TestDeleteAccount_UserNotFound(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/auth/account", map[string]string{"password": "password123"})
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- パスワードリセット ----------

func TestRequestPasswordReset_Success(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "test@example.com").Return(&model.User{ID: 1}, nil)
	ports.Tokens.On("InvalidateUserTokens", mock.Anything, uint(1)).Return(nil)
	ports.Tokens.On("Create", mock.Anything, mock.MatchedBy(func(t *model.PasswordResetToken) bool {
		// 平文ではなくハッシュを保存する（SHA-256 の 64 桁）
		return t.UserID == 1 && len(t.Token) == 64 && t.ExpiresAt.After(t.CreatedAt)
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/auth/password-reset", map[string]string{"email": "test@example.com"})
	assertStatus(t, w, http.StatusOK)
	ports.Tokens.AssertExpectations(t)
}

// 未登録のメールアドレスでも 200 を返し、存在有無を漏らさない。
func TestRequestPasswordReset_UnknownEmail(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "unknown@example.com").Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/auth/password-reset", map[string]string{"email": "unknown@example.com"})
	assertStatus(t, w, http.StatusOK)
	ports.Tokens.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 検索の失敗（DB 障害）は成功に見せず 500 を返す（従来は 200 で障害が埋もれていた）。
func TestRequestPasswordReset_SearchFailure(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Users.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("db down"))

	w := doRequest(r, http.MethodPost, "/auth/password-reset", map[string]string{"email": "test@example.com"})
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Tokens.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestResetPassword_Success(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Tokens.On("FindByToken", mock.Anything, mock.AnythingOfType("string")).
		Return(&model.PasswordResetToken{ID: 5, UserID: 1, ExpiresAt: futureTime()}, nil)
	ports.Users.On("UpdatePassword", mock.Anything, uint(1), mock.MatchedBy(func(hashed string) bool {
		return hashed != "newpassword123"
	})).Return(nil)
	ports.Tokens.On("MarkAsUsed", mock.Anything, uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/auth/password-reset/confirm", map[string]string{
		"token": "raw-token", "new_password": "newpassword123",
	})
	assertStatus(t, w, http.StatusOK)
	ports.Users.AssertExpectations(t)
	ports.Tokens.AssertExpectations(t)
}

// トークンが見つからなければ 400。
func TestResetPassword_InvalidToken(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Tokens.On("FindByToken", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/auth/password-reset/confirm", map[string]string{
		"token": "bad-token", "new_password": "newpassword123",
	})
	assertStatus(t, w, http.StatusBadRequest)
	ports.Users.AssertNotCalled(t, "UpdatePassword", mock.Anything, mock.Anything, mock.Anything)
}

// 期限切れのトークンは 400。
func TestResetPassword_ExpiredToken(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.Tokens.On("FindByToken", mock.Anything, mock.AnythingOfType("string")).
		Return(&model.PasswordResetToken{ID: 5, UserID: 1, ExpiresAt: pastTime()}, nil)

	w := doRequest(r, http.MethodPost, "/auth/password-reset/confirm", map[string]string{
		"token": "raw-token", "new_password": "newpassword123",
	})
	assertStatus(t, w, http.StatusBadRequest)
	ports.Users.AssertNotCalled(t, "UpdatePassword", mock.Anything, mock.Anything, mock.Anything)
}

func TestResetPassword_ValidationError(t *testing.T) {
	r, _, ports := setupAuthHandler()

	w := doRequest(r, http.MethodPost, "/auth/password-reset/confirm", map[string]string{
		"token": "raw-token", "new_password": "short",
	})
	assertStatus(t, w, http.StatusBadRequest)
	ports.Users.AssertNotCalled(t, "UpdatePassword", mock.Anything, mock.Anything, mock.Anything)
}

// ---------- GitHub ログイン ----------

func TestAuthGitHubLogin_ReturnsURL(t *testing.T) {
	r, _, ports := setupAuthHandler()
	ports.GitHub.Client.On("LoginAuthorizeURL", mock.AnythingOfType("string")).
		Return("https://github.com/login/oauth/authorize?scope=user:email")

	w := doRequest(r, http.MethodGet, "/auth/github", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "github.com/login/oauth/authorize")
	ports.GitHub.Client.AssertExpectations(t)
}

func TestGitHubLoginCallback_MissingParams(t *testing.T) {
	r, _, ports := setupAuthHandler()

	w := doRequest(r, http.MethodGet, "/auth/github/callback", nil)
	assertStatus(t, w, http.StatusBadRequest)
	ports.GitHub.Client.AssertNotCalled(t, "ExchangeCode", mock.Anything, mock.Anything)
}

func TestGitHubLoginCallback_InvalidState(t *testing.T) {
	r, _, ports := setupAuthHandler()

	w := doRequest(r, http.MethodGet, "/auth/github/callback?code=valid-code&state=invalid-state", nil)
	assertStatus(t, w, http.StatusUnauthorized)
	ports.GitHub.Client.AssertNotCalled(t, "ExchangeCode", mock.Anything, mock.Anything)
}

// state を発行してからコールバックを通し、既存ユーザーとしてログインできることを確認する。
func TestGitHubLoginCallback_Success(t *testing.T) {
	r, _, ports := setupAuthHandler()
	state, err := usecase.NewGitHubLoginStateUseCase(testJWTSecret).Generate()
	assert.NoError(t, err)

	ports.GitHub.Client.On("ExchangeCode", mock.Anything, "valid-code").Return("access-token", nil)
	ports.GitHub.Client.On("GetUser", mock.Anything, "access-token").
		Return(&model.GitHubUserInfo{ID: 42, Login: "dev", Email: "dev@example.com"}, nil)
	ports.Users.On("FindByGitHubID", mock.Anything, int64(42)).Return(&model.User{ID: 1, Email: "dev@example.com"}, nil)
	ports.Users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.GitHubToken == "access-token" && u.GitHubUsername == "dev"
	})).Return(nil)
	// ログイン後のデータ同期はバックグラウンドで走るため、呼ばれても呼ばれなくてもよい
	ports.GitHub.Users.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	w := doRequest(r, http.MethodGet, "/auth/github/callback?code=valid-code&state="+state, nil)
	assertStatus(t, w, http.StatusOK)
	assert.NotNil(t, tokenCookie(w.Result()))
	ports.Users.AssertExpectations(t)
}

func TestGitHubLoginCallback_ExchangeCodeError(t *testing.T) {
	r, _, ports := setupAuthHandler()
	state, _ := usecase.NewGitHubLoginStateUseCase(testJWTSecret).Generate()

	ports.GitHub.Client.On("ExchangeCode", mock.Anything, "bad-code").Return("", errors.New("api error"))

	w := doRequest(r, http.MethodGet, "/auth/github/callback?code=bad-code&state="+state, nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Users.AssertNotCalled(t, "FindByGitHubID", mock.Anything, mock.Anything)
}

func TestGitHubLoginCallback_GetUserError(t *testing.T) {
	r, _, ports := setupAuthHandler()
	state, _ := usecase.NewGitHubLoginStateUseCase(testJWTSecret).Generate()

	ports.GitHub.Client.On("ExchangeCode", mock.Anything, "valid-code").Return("access-token", nil)
	ports.GitHub.Client.On("GetUser", mock.Anything, "access-token").Return(nil, errors.New("api error"))

	w := doRequest(r, http.MethodGet, "/auth/github/callback?code=valid-code&state="+state, nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Users.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 既存ユーザーがいなければ新規作成してログインする。
func TestGitHubLoginCallback_CreatesUser(t *testing.T) {
	r, _, ports := setupAuthHandler()
	state, _ := usecase.NewGitHubLoginStateUseCase(testJWTSecret).Generate()

	ports.GitHub.Client.On("ExchangeCode", mock.Anything, "valid-code").Return("access-token", nil)
	ports.GitHub.Client.On("GetUser", mock.Anything, "access-token").
		Return(&model.GitHubUserInfo{ID: 42, Login: "newdev"}, nil)
	ports.Users.On("FindByGitHubID", mock.Anything, int64(42)).Return(nil, nil)
	ports.Users.On("FindByUsername", mock.Anything, "newdev").Return(nil, nil)
	ports.Users.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		// メールアドレスが無い場合はプレースホルダを使う
		return u.Username == "newdev" && u.Email == "newdev@github.local" && u.GitHubConnected
	})).Return(nil)
	ports.GitHub.Users.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	w := doRequest(r, http.MethodGet, "/auth/github/callback?code=valid-code&state="+state, nil)
	assertStatus(t, w, http.StatusOK)
	ports.Users.AssertExpectations(t)
}

// ユーザー名が使用済みなら連番を付けて作成する。
func TestGitHubLoginCallback_UniqueUsername(t *testing.T) {
	r, _, ports := setupAuthHandler()
	state, _ := usecase.NewGitHubLoginStateUseCase(testJWTSecret).Generate()

	ports.GitHub.Client.On("ExchangeCode", mock.Anything, "valid-code").Return("access-token", nil)
	ports.GitHub.Client.On("GetUser", mock.Anything, "access-token").
		Return(&model.GitHubUserInfo{ID: 42, Login: "dev"}, nil)
	ports.Users.On("FindByGitHubID", mock.Anything, int64(42)).Return(nil, nil)
	ports.Users.On("FindByUsername", mock.Anything, "dev").Return(&model.User{ID: 9}, nil)
	ports.Users.On("FindByUsername", mock.Anything, "dev2").Return(nil, nil)
	ports.Users.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.Username == "dev2"
	})).Return(nil)
	ports.GitHub.Users.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	w := doRequest(r, http.MethodGet, "/auth/github/callback?code=valid-code&state="+state, nil)
	assertStatus(t, w, http.StatusOK)
	ports.Users.AssertExpectations(t)
}

// futureTime は未来の時刻を返す（有効なトークン用）。
func futureTime() time.Time {
	return time.Now().Add(1 * time.Hour)
}

// pastTime は過去の時刻を返す（期限切れトークン用）。
func pastTime() time.Time {
	return time.Now().Add(-1 * time.Hour)
}
