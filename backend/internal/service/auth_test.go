package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const testJWTSecret = "test-secret-key-for-unit-tests"

// newTestAuthService はAuthServiceのテスト用インスタンスを生成するヘルパー。
func newTestAuthService() (*AuthService, *MockUserRepository, *MockPasswordResetRepository) {
	userRepo := new(MockUserRepository)
	passwordResetRepo := new(MockPasswordResetRepository)
	svc := NewAuthService(userRepo, passwordResetRepo, testJWTSecret)
	return svc, userRepo, passwordResetRepo
}

// ============================================================
// 新規登録テスト
// ============================================================

func TestRegister_Success(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	resp, err := svc.Register(RegisterInput{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "testuser", resp.User.Name)
	assert.Equal(t, "test@example.com", resp.User.Email)
	userRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	existingUser := &model.User{Name: "Existing", Email: "test@example.com"}
	userRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	resp, err := svc.Register(RegisterInput{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "このメールアドレスは既に登録されています")
	userRepo.AssertExpectations(t)
}

// ============================================================
// ログインテスト
// ============================================================

func TestLogin_Success(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.User{Name: "testuser", Email: "test@example.com", Password: string(hashed)}
	user.ID = 1

	userRepo.On("FindByEmail", "test@example.com").Return(user, nil)

	resp, err := svc.Login(LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "testuser", resp.User.Name)
	userRepo.AssertExpectations(t)
}

func TestLogin_EmailNotFound(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("FindByEmail", "unknown@example.com").Return(nil, errors.New("not found"))

	resp, err := svc.Login(LoginInput{
		Email:    "unknown@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "認証")
	userRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &model.User{Name: "testuser", Email: "test@example.com", Password: string(hashed)}
	user.ID = 1

	userRepo.On("FindByEmail", "test@example.com").Return(user, nil)

	resp, err := svc.Login(LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "認証")
	userRepo.AssertExpectations(t)
}

// ============================================================
// トークン検証テスト
// ============================================================

func TestValidateToken_Success(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.User{Name: "Test", Email: "test@example.com", Password: string(hashed)}
	user.ID = 42

	userRepo.On("FindByEmail", "test@example.com").Return(user, nil)

	resp, err := svc.Login(LoginInput{Email: "test@example.com", Password: "password123"})
	assert.NoError(t, err)

	userID, err := svc.ValidateToken(resp.Token)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), userID)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	svc, _, _ := newTestAuthService()

	userID, err := svc.ValidateToken("invalid.token.string")
	assert.Error(t, err)
	assert.Equal(t, uint(0), userID)
}

// ============================================================
// アカウント削除テスト
// ============================================================

func TestDeleteAccount_Success(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.User{Name: "Test", Email: "test@example.com", Password: string(hashed)}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("DeleteWithRelatedData", uint(1)).Return(nil)

	err := svc.DeleteAccount(1, "password123")
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestDeleteAccount_WrongPassword(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.User{Name: "Test", Email: "test@example.com", Password: string(hashed)}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)

	err := svc.DeleteAccount(1, "wrongpassword")
	assert.ErrorIs(t, err, ErrForbidden)
	userRepo.AssertExpectations(t)
}

func TestDeleteAccount_GitHubOnlyUser(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	// GitHubのみユーザーはパスワードが空
	user := &model.User{Name: "ghuser", Email: "gh@github.local", Password: ""}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("DeleteWithRelatedData", uint(1)).Return(nil)

	err := svc.DeleteAccount(1, "")
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestDeleteAccount_UserNotFound(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.DeleteAccount(999, "password")
	assert.ErrorIs(t, err, ErrNotFound)
	userRepo.AssertExpectations(t)
}

func TestDeleteAccount_PasswordRequired(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.User{Name: "Test", Password: string(hashed)}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)

	err := svc.DeleteAccount(1, "")
	assert.ErrorIs(t, err, ErrBadRequest)
	userRepo.AssertExpectations(t)
}

// ============================================================
// ユーザー情報取得テスト
// ============================================================

func TestGetMe_Success(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	user := &model.User{Name: "testuser"}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetMe(1)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", result.Name)
	userRepo.AssertExpectations(t)
}

// ============================================================
// パスワードリセット要求テスト
// ============================================================

func TestRequestPasswordReset_Success(t *testing.T) {
	svc, userRepo, passwordResetRepo := newTestAuthService()

	user := &model.User{Email: "test@example.com"}
	user.ID = 1

	userRepo.On("FindByEmail", "test@example.com").Return(user, nil)
	passwordResetRepo.On("InvalidateUserTokens", uint(1)).Return(nil)
	passwordResetRepo.On("Create", mock.AnythingOfType("*model.PasswordResetToken")).Return(nil)

	token, err := svc.RequestPasswordReset("test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	userRepo.AssertExpectations(t)
	passwordResetRepo.AssertExpectations(t)
}

func TestRequestPasswordReset_EmailNotFound(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("FindByEmail", "unknown@example.com").Return(nil, errors.New("not found"))

	token, err := svc.RequestPasswordReset("unknown@example.com")
	assert.NoError(t, err) // エラーを返さない（メール存在の漏洩防止）
	assert.Empty(t, token)
	userRepo.AssertExpectations(t)
}

// ============================================================
// パスワードリセット実行テスト
// ============================================================

func TestResetPassword_Success(t *testing.T) {
	svc, userRepo, passwordResetRepo := newTestAuthService()

	resetToken := &model.PasswordResetToken{
		UserID:    1,
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}
	resetToken.ID = 1

	passwordResetRepo.On("FindByToken", "valid-token").Return(resetToken, nil)
	userRepo.On("UpdatePassword", uint(1), mock.AnythingOfType("string")).Return(nil)
	passwordResetRepo.On("MarkAsUsed", uint(1)).Return(nil)

	err := svc.ResetPassword("valid-token", "newpassword123")
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	passwordResetRepo.AssertExpectations(t)
}

func TestResetPassword_InvalidToken(t *testing.T) {
	svc, _, passwordResetRepo := newTestAuthService()

	passwordResetRepo.On("FindByToken", "invalid-token").Return(nil, errors.New("not found"))

	err := svc.ResetPassword("invalid-token", "newpassword123")
	assert.ErrorIs(t, err, ErrBadRequest)
	passwordResetRepo.AssertExpectations(t)
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	svc, _, passwordResetRepo := newTestAuthService()

	resetToken := &model.PasswordResetToken{
		UserID:    1,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // 期限切れ
		Used:      false,
	}
	resetToken.ID = 1

	passwordResetRepo.On("FindByToken", "expired-token").Return(resetToken, nil)

	err := svc.ResetPassword("expired-token", "newpassword123")
	assert.ErrorIs(t, err, ErrBadRequest)
	passwordResetRepo.AssertExpectations(t)
}

// ============================================================
// GitHubログインテスト
// ============================================================

func TestGitHubLogin_ExistingGitHubUser(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	user := &model.User{Name: "ghuser", GitHubUsername: "ghuser"}
	user.ID = 1

	userRepo.On("FindByGitHubID", int64(12345)).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	ghUser := &GitHubUserInfo{
		ID:        12345,
		Login:     "ghuser",
		Email:     "gh@example.com",
		Name:      "ghuser",
		AvatarURL: "https://example.com/avatar.png",
	}

	resp, err := svc.GitHubLogin(ghUser, "access-token")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	userRepo.AssertExpectations(t)
}

func TestGitHubLogin_LinkByEmail(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	// GitHub IDが見つからない
	userRepo.On("FindByGitHubID", int64(12345)).Return(nil, errors.New("not found"))
	// しかしメールが一致
	user := &model.User{Name: "existinguser", Email: "existing@example.com"}
	user.ID = 2
	userRepo.On("FindByEmail", "existing@example.com").Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	ghUser := &GitHubUserInfo{
		ID:    12345,
		Login: "ghuser",
		Email: "existing@example.com",
		Name:  "ghuser",
	}

	resp, err := svc.GitHubLogin(ghUser, "access-token")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	userRepo.AssertExpectations(t)
}

func TestGitHubLogin_CreateNewUser(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	// GitHub IDもメールも見つからない
	userRepo.On("FindByGitHubID", int64(12345)).Return(nil, errors.New("not found"))
	userRepo.On("FindByEmail", "new@example.com").Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	ghUser := &GitHubUserInfo{
		ID:    12345,
		Login: "newghuser",
		Email: "new@example.com",
		Name:  "New ghuser",
	}

	resp, err := svc.GitHubLogin(ghUser, "access-token")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	userRepo.AssertExpectations(t)
}

func TestGitHubLogin_NoEmailFallback(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	// GitHub IDが見つからず、メールも未提供
	userRepo.On("FindByGitHubID", int64(12345)).Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*model.User)
		// メール無しの場合、login@github.local がフォールバック
		assert.Equal(t, "ghuser@github.local", user.Email)
		assert.Equal(t, "ghuser", user.Name) // Name空の場合、Loginがフォールバック
	}).Return(nil)

	ghUser := &GitHubUserInfo{
		ID:    12345,
		Login: "ghuser",
		Email: "",
		Name:  "",
	}

	resp, err := svc.GitHubLogin(ghUser, "access-token")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	userRepo.AssertExpectations(t)
}

// ============================================================
// OAuthステート検証テスト
// ============================================================

func TestOAuthState_RoundTrip(t *testing.T) {
	svc, _, _ := newTestAuthService()

	state, err := svc.GenerateOAuthState(42)
	assert.NoError(t, err)
	assert.NotEmpty(t, state)

	userID, err := svc.ValidateOAuthState(state)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), userID)
}

func TestValidateOAuthState_InvalidState(t *testing.T) {
	svc, _, _ := newTestAuthService()

	userID, err := svc.ValidateOAuthState("invalid-state")
	assert.Error(t, err)
	assert.Equal(t, uint(0), userID)
}

// ============================================================
// ログインステート検証テスト
// ============================================================

func TestLoginState_RoundTrip(t *testing.T) {
	svc, _, _ := newTestAuthService()

	state, err := svc.GenerateLoginState()
	assert.NoError(t, err)
	assert.NotEmpty(t, state)

	err = svc.ValidateLoginState(state)
	assert.NoError(t, err)
}

func TestValidateLoginState_InvalidState(t *testing.T) {
	svc, _, _ := newTestAuthService()

	err := svc.ValidateLoginState("invalid-state")
	assert.Error(t, err)
}
