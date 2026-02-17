// Package handler は認証ハンドラーのhttpOnly Cookie設定テストを提供する。
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository は認証ハンドラーテスト用のユーザーリポジトリモック。
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	if args.Error(0) == nil {
		user.ID = 1 // テスト用にIDをセット
	}
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByGitHubID(githubID int64) (*model.User, error) {
	args := m.Called(githubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	args := m.Called(userID, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteWithRelatedData(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) Search(query string) ([]model.User, error) {
	args := m.Called(query)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

// MockPasswordResetRepository はパスワードリセットリポジトリモック。
type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) Create(token *model.PasswordResetToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) FindByToken(token string) (*model.PasswordResetToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PasswordResetToken), args.Error(1)
}

func (m *MockPasswordResetRepository) MarkAsUsed(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) InvalidateUserTokens(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

// テスト用のJWTシークレットキー
const testJWTSecret = "test-secret-key-for-handler-tests"

// setupLoginTest はログインテスト用のルーターとモックをセットアップするヘルパー関数。
func setupLoginTest() (*gin.Engine, *MockUserRepository) {
	gin.SetMode(gin.TestMode)
	mockUserRepo := new(MockUserRepository)
	mockPwResetRepo := new(MockPasswordResetRepository)
	authService := service.NewAuthService(mockUserRepo, mockPwResetRepo, testJWTSecret)
	authHandler := NewAuthHandler(authService, nil)

	r := gin.New()
	// 認証不要なルート
	r.POST("/api/v1/auth/login", authHandler.Login)
	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/password-reset", authHandler.RequestPasswordReset)
	r.POST("/api/v1/auth/password-reset/confirm", authHandler.ResetPassword)

	// 認証が必要なルート（テスト用に簡易的なミドルウェアを追加）
	authorized := r.Group("/api/v1/auth")
	authorized.Use(func(c *gin.Context) {
		// テスト用の簡易認証ミドルウェア（userIDをcontextから取得または設定）
		if userID, exists := c.Get("userID"); !exists {
			c.Set("userID", uint(1)) // デフォルトでuserID=1を設定
		} else {
			c.Set("userID", userID)
		}
		c.Next()
	})
	authorized.GET("/me", authHandler.Me)
	authorized.POST("/logout", authHandler.Logout)
	authorized.DELETE("/account", authHandler.DeleteAccount)

	return r, mockUserRepo
}

// TestLogin_SetsCookie はログイン成功時にSet-Cookieヘッダーがセットされることをテストする。
func TestLogin_SetsCookie(t *testing.T) {
	r, mockUserRepo := setupLoginTest()

	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUserRepo.On("FindByEmail", "test@example.com").Return(&model.User{
		Email:    "test@example.com",
		Password: string(hashedPw),
		Name:     "Test User",
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Set-Cookieヘッダーの確認
	cookies := w.Result().Cookies()
	var tokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "token" {
			tokenCookie = c
			break
		}
	}
	if assert.NotNil(t, tokenCookie, "token Cookieがセットされるべき") {
		assert.True(t, tokenCookie.HttpOnly, "httpOnly属性がtrueであるべき")
		assert.Equal(t, "/", tokenCookie.Path, "Pathが/であるべき")
		assert.True(t, tokenCookie.MaxAge > 0, "MaxAgeが正の値であるべき")
	}
}

// TestLogin_ResponseHasUser はログイン成功時のレスポンスにuserが含まれることをテストする。
func TestLogin_ResponseHasUser(t *testing.T) {
	r, mockUserRepo := setupLoginTest()

	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUserRepo.On("FindByEmail", "test@example.com").Return(&model.User{
		Email:    "test@example.com",
		Password: string(hashedPw),
		Name:     "Test User",
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp, "user", "レスポンスにuserフィールドが含まれるべき")
}

// TestRegister_SetsCookie は登録成功時にSet-Cookieヘッダーがセットされることをテストする。
func TestRegister_SetsCookie(t *testing.T) {
	r, mockUserRepo := setupLoginTest()

	mockUserRepo.On("FindByEmail", "new@example.com").Return(nil, assert.AnError)
	mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	body, _ := json.Marshal(map[string]string{
		"name":             "New User",
		"email":            "new@example.com",
		"password":         "password123",
		"confirm_password": "password123",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Set-Cookieヘッダーの確認
	cookies := w.Result().Cookies()
	var tokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "token" {
			tokenCookie = c
			break
		}
	}
	if assert.NotNil(t, tokenCookie, "token Cookieがセットされるべき") {
		assert.True(t, tokenCookie.HttpOnly, "httpOnly属性がtrueであるべき")
	}
}

// TestLogout_ClearsCookie はログアウト時にCookieがクリアされることをテストする。
func TestLogout_ClearsCookie(t *testing.T) {
	r, _ := setupLoginTest()

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Set-Cookieヘッダーでtoken Cookieがクリアされることを確認
	cookies := w.Result().Cookies()
	var tokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "token" {
			tokenCookie = c
			break
		}
	}
	if assert.NotNil(t, tokenCookie, "token Cookieが設定されるべき（削除用）") {
		assert.True(t, tokenCookie.MaxAge < 0, "MaxAgeが負の値であるべき（Cookie削除）")
	}
}

// TestLogout_ReturnsOK はログアウト成功時に200レスポンスを返すことをテストする。
func TestLogout_ReturnsOK(t *testing.T) {
	r, _ := setupLoginTest()

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp, "message")
}

// TestMe_Success は正常なユーザー情報取得をテストする。
func TestMe_Success(t *testing.T) {
	r, mockUserRepo := setupLoginTest()

	mockUser := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	mockUser.ID = 1

	mockUserRepo.On("FindByID", uint(1)).Return(mockUser, nil)

	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()

	// userIDをcontextにセット（middlewareの代わり）
	ctx := gin.CreateTestContextOnly(w, r)
	ctx.Set("userID", uint(1))
	ctx.Request = req

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestMe_UserNotFound はユーザーが見つからない場合をテストする。
func TestMe_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserRepo := new(MockUserRepository)
	mockPwResetRepo := new(MockPasswordResetRepository)
	authService := service.NewAuthService(mockUserRepo, mockPwResetRepo, testJWTSecret)
	authHandler := NewAuthHandler(authService, nil)

	mockUserRepo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	// 独自のルーター（ミドルウェアなし）を作成
	r := gin.New()
	r.GET("/api/v1/auth/me", func(c *gin.Context) {
		c.Set("userID", uint(999))
		authHandler.Me(c)
	})

	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestDeleteAccount_Success は正常なアカウント削除をテストする。
func TestDeleteAccount_Success(t *testing.T) {
	r, mockUserRepo := setupLoginTest()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	mockUser.ID = 1

	mockUserRepo.On("FindByID", uint(1)).Return(mockUser, nil)
	mockUserRepo.On("DeleteWithRelatedData", uint(1)).Return(nil)

	body, _ := json.Marshal(map[string]string{
		"password": "password123",
	})
	req := httptest.NewRequest("DELETE", "/api/v1/auth/account", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp, "message")
}

// TestDeleteAccount_WrongPassword はパスワード不一致をテストする。
func TestDeleteAccount_WrongPassword(t *testing.T) {
	r, mockUserRepo := setupLoginTest()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	mockUser.ID = 1

	mockUserRepo.On("FindByID", uint(1)).Return(mockUser, nil)

	body, _ := json.Marshal(map[string]string{
		"password": "wrongpassword",
	})
	req := httptest.NewRequest("DELETE", "/api/v1/auth/account", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// パスワード不一致の場合、service層でErrForbiddenが返されるが、respondErrorで401に変換される
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRequestPasswordReset_Success は正常なパスワードリセット要求をテストする。
func TestRequestPasswordReset_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserRepo := new(MockUserRepository)
	mockPwResetRepo := new(MockPasswordResetRepository)
	authService := service.NewAuthService(mockUserRepo, mockPwResetRepo, testJWTSecret)
	authHandler := NewAuthHandler(authService, nil)

	mockUser := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	mockUser.ID = 1

	mockUserRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil)
	mockPwResetRepo.On("InvalidateUserTokens", uint(1)).Return(nil)
	mockPwResetRepo.On("Create", mock.AnythingOfType("*model.PasswordResetToken")).Return(nil)

	// 独自のルーターを作成
	r := gin.New()
	r.POST("/api/v1/auth/password-reset", authHandler.RequestPasswordReset)

	body, _ := json.Marshal(map[string]string{
		"email": "test@example.com",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/password-reset", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp, "message")
}

// TestResetPassword_ValidationError はバリデーションエラーをテストする。
func TestResetPassword_ValidationError(t *testing.T) {
	r, _ := setupLoginTest()

	// new_passwordが短すぎる（min=6）
	body, _ := json.Marshal(map[string]string{
		"token":        "valid-token",
		"new_password": "short",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/password-reset/confirm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------- モックベースのAuth Handlerテスト ----------

// TestAuthGitHubLogin_ReturnsURL はGitHub OAuthログインURLの生成をテストする。
func TestAuthGitHubLogin_ReturnsURL(t *testing.T) {
	h, authSvc, ghSvc := setupAuthHandlerMock()
	authSvc.On("GenerateLoginState").Return("test-state", nil)
	ghSvc.On("GetLoginOAuthURL", "test-state").Return("https://github.com/login/oauth/authorize?state=test-state")

	r := gin.New()
	r.GET("/auth/github", h.GitHubLogin)
	w := doRequest(r, "GET", "/auth/github", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.NotEmpty(t, body["url"])
	authSvc.AssertExpectations(t)
	ghSvc.AssertExpectations(t)
}

// TestAuthGitHubLogin_StateError はstate生成エラーをテストする。
func TestAuthGitHubLogin_StateError(t *testing.T) {
	h, authSvc, _ := setupAuthHandlerMock()
	authSvc.On("GenerateLoginState").Return("", service.ErrBadRequest)

	r := gin.New()
	r.GET("/auth/github", h.GitHubLogin)
	w := doRequest(r, "GET", "/auth/github", nil)

	assertStatus(t, w, http.StatusBadRequest)
	authSvc.AssertExpectations(t)
}

// TestAuthMeMock_Success はモック版のMe成功テスト。
func TestAuthMeMock_Success(t *testing.T) {
	h, authSvc, _ := setupAuthHandlerMock()
	user := &model.User{Name: "Test User", Email: "test@example.com"}
	authSvc.On("GetMe", uint(1)).Return(user, nil)

	r := newRouter(1)
	r.GET("/me", h.Me)
	w := doRequest(r, "GET", "/me", nil)

	assertStatus(t, w, http.StatusOK)
	authSvc.AssertExpectations(t)
}

// TestAuthMeMock_NotFound はモック版のMeユーザー未発見テスト。
func TestAuthMeMock_NotFound(t *testing.T) {
	h, authSvc, _ := setupAuthHandlerMock()
	authSvc.On("GetMe", uint(1)).Return(nil, service.ErrNotFound)

	r := newRouter(1)
	r.GET("/me", h.Me)
	w := doRequest(r, "GET", "/me", nil)

	assertStatus(t, w, http.StatusNotFound)
	authSvc.AssertExpectations(t)
}

// TestAuthLoginMock_Success はモック版のLogin成功テスト。
func TestAuthLoginMock_Success(t *testing.T) {
	h, authSvc, _ := setupAuthHandlerMock()
	resp := &service.AuthResponse{
		Token: "test-token",
		User:  model.User{Name: "Test User", Email: "test@example.com"},
	}
	authSvc.On("Login", mock.Anything).Return(resp, nil)

	r := gin.New()
	r.POST("/login", h.Login)
	w := doRequest(r, "POST", "/login", map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})

	assertStatus(t, w, http.StatusOK)
	authSvc.AssertExpectations(t)
}

// TestAuthLoginMock_Unauthorized はモック版のログイン失敗テスト。
func TestAuthLoginMock_Unauthorized(t *testing.T) {
	h, authSvc, _ := setupAuthHandlerMock()
	authSvc.On("Login", mock.Anything).Return(nil, service.ErrUnauthorized)

	r := gin.New()
	r.POST("/login", h.Login)
	w := doRequest(r, "POST", "/login", map[string]string{
		"email":    "test@example.com",
		"password": "wrong",
	})

	assertStatus(t, w, http.StatusUnauthorized)
	authSvc.AssertExpectations(t)
}

// TestAuthDeleteAccountMock_Success はモック版のアカウント削除成功テスト。
func TestAuthDeleteAccountMock_Success(t *testing.T) {
	h, authSvc, _ := setupAuthHandlerMock()
	authSvc.On("DeleteAccount", uint(1), "password123").Return(nil)

	r := newRouter(1)
	r.DELETE("/account", h.DeleteAccount)
	w := doRequest(r, "DELETE", "/account", map[string]string{
		"password": "password123",
	})

	assertStatus(t, w, http.StatusOK)
	authSvc.AssertExpectations(t)
}
