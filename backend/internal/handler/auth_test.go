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
	r.POST("/api/v1/auth/login", authHandler.Login)
	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/logout", authHandler.Logout)

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
		"name":     "New User",
		"email":    "new@example.com",
		"password": "password123",
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
