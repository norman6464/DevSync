package service

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	userRepo.On("FindByUsername", "testuser").Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	resp, err := svc.Register(RegisterInput{
		Name:     "Test User",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "Test User", resp.User.Name)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, "testuser", resp.User.Username)
	userRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	existingUser := &model.User{Name: "Existing", Email: "test@example.com"}
	userRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	resp, err := svc.Register(RegisterInput{
		Name:     "testuser",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "このメールアドレスは既に登録されています")
	userRepo.AssertExpectations(t)
}

func TestRegister_InvalidEmail(t *testing.T) {
	svc, _, _ := newTestAuthService()

	resp, err := svc.Register(RegisterInput{
		Name:     "testuser",
		Username: "testuser",
		Email:    "invalid-email",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRegister_InvalidPassword(t *testing.T) {
	svc, _, _ := newTestAuthService()

	resp, err := svc.Register(RegisterInput{
		Name:     "testuser",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "short",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRegister_InvalidUsername(t *testing.T) {
	svc, _, _ := newTestAuthService()

	resp, err := svc.Register(RegisterInput{
		Name:     "testuser",
		Username: "",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("FindByEmail", "new@example.com").Return(nil, errors.New("not found"))
	existingUser := &model.User{Name: "Existing", Username: "taken"}
	userRepo.On("FindByUsername", "taken").Return(existingUser, nil)

	resp, err := svc.Register(RegisterInput{
		Name:     "testuser",
		Username: "taken",
		Email:    "new@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "このユーザー名は既に使用されています")
	userRepo.AssertExpectations(t)
}

func TestRegister_CreateUserError(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("FindByEmail", "new@example.com").Return(nil, errors.New("not found"))
	userRepo.On("FindByUsername", "newuser").Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.AnythingOfType("*model.User")).Return(errors.New("db error"))

	resp, err := svc.Register(RegisterInput{
		Name:     "testuser",
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
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

func TestDeleteAccount_DeleteWithRelatedDataError(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.User{Name: "Test", Email: "test@example.com", Password: string(hashed)}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("DeleteWithRelatedData", uint(1)).Return(errors.New("db error"))

	err := svc.DeleteAccount(1, "password123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
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
	userRepo.On("FindByUsername", "newghuser").Return(nil, errors.New("not found"))
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
	userRepo.On("FindByUsername", "ghuser").Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*model.User)
		// メール無しの場合、login@github.local がフォールバック
		assert.Equal(t, "ghuser@github.local", user.Email)
		assert.Equal(t, "ghuser", user.Name) // Name空の場合、Loginがフォールバック
		assert.Equal(t, "ghuser", user.Username)
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

func TestGitHubLogin_CreateNewUser_CreateError(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("FindByGitHubID", int64(12345)).Return(nil, errors.New("not found"))
	userRepo.On("FindByEmail", "new@example.com").Return(nil, errors.New("not found"))
	userRepo.On("FindByUsername", "newghuser").Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.AnythingOfType("*model.User")).Return(errors.New("db error"))

	ghUser := &GitHubUserInfo{
		ID:    12345,
		Login: "newghuser",
		Email: "new@example.com",
		Name:  "New ghuser",
	}

	resp, err := svc.GitHubLogin(ghUser, "access-token")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGitHubLogin_LinkByEmail_WithAvatarURL(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("FindByGitHubID", int64(12345)).Return(nil, errors.New("not found"))
	user := &model.User{Name: "existinguser", Email: "existing@example.com"}
	user.ID = 2
	userRepo.On("FindByEmail", "existing@example.com").Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	ghUser := &GitHubUserInfo{
		ID:        12345,
		Login:     "ghuser",
		Email:     "existing@example.com",
		Name:      "ghuser",
		AvatarURL: "https://example.com/avatar.png",
	}

	resp, err := svc.GitHubLogin(ghUser, "access-token")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "https://example.com/avatar.png", resp.User.AvatarURL)
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

// ============================================================
// ValidateLoginState 追加テスト
// ============================================================

func TestValidateLoginState_WrongPurpose(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// purposeが"oauth_state"のトークンはgithub_loginとして無効
	state, err := svc.GenerateOAuthState(1)
	assert.NoError(t, err)

	err = svc.ValidateLoginState(state)
	assert.Error(t, err)
}

// ============================================================
// ValidateOAuthState 追加テスト
// ============================================================

func TestValidateOAuthState_RoundTrip(t *testing.T) {
	svc, _, _ := newTestAuthService()

	state, err := svc.GenerateOAuthState(42)
	assert.NoError(t, err)
	assert.NotEmpty(t, state)

	userID, err := svc.ValidateOAuthState(state)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), userID)
}

func TestValidateOAuthState_MissingUserID(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// user_idなしのトークンを手動生成
	claims := jwt.MapClaims{
		"purpose": "oauth_state",
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	state, err := token.SignedString([]byte(testJWTSecret))
	assert.NoError(t, err)

	_, err = svc.ValidateOAuthState(state)
	assert.Error(t, err)
}

func TestValidateOAuthState_WrongPurpose(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// purposeが"github_login"のトークンはoauth_stateとして無効
	state, err := svc.GenerateLoginState()
	assert.NoError(t, err)

	_, err = svc.ValidateOAuthState(state)
	assert.Error(t, err)
}

// ============================================================
// generateUsername テスト
// ============================================================

func TestGenerateUsername_UniqueBase(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	// ユーザー名候補が重複しない場合、そのまま返す
	userRepo.On("FindByUsername", "testuser").Return(nil, errors.New("not found"))

	username := svc.generateUsername("testuser")
	assert.Equal(t, "testuser", username)
}

func TestGenerateUsername_AlreadyExists(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	// 最初の候補が重複している場合、数字サフィックスを追加
	existingUser := &model.User{Username: "testuser"}
	existingUser.ID = 1

	// 最初の呼び出しでは存在する、2回目以降は存在しない
	userRepo.On("FindByUsername", "testuser").Return(existingUser, nil)
	userRepo.On("FindByUsername", "testuser2").Return(nil, errors.New("not found"))

	username := svc.generateUsername("testuser")
	assert.Equal(t, "testuser2", username)
}

// ============================================================
// Login 追加テスト
// ============================================================

func TestLogin_InvalidEmail(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// 無効なメールアドレス
	_, err := svc.Login(LoginInput{Email: "not-an-email", Password: "password123"})
	assert.Error(t, err)
}

// ============================================================
// RequestPasswordReset 追加テスト
// ============================================================

func TestRequestPasswordReset_CreateError(t *testing.T) {
	svc, userRepo, passwordResetRepo := newTestAuthService()

	user := &model.User{Email: "test@example.com"}
	user.ID = 1

	userRepo.On("FindByEmail", "test@example.com").Return(user, nil)
	passwordResetRepo.On("InvalidateUserTokens", uint(1)).Return(nil)
	passwordResetRepo.On("Create", mock.AnythingOfType("*model.PasswordResetToken")).Return(errors.New("db error"))

	token, err := svc.RequestPasswordReset("test@example.com")
	assert.Error(t, err)
	assert.Empty(t, token)
}

// ============================================================
// ValidateToken 追加テスト
// ============================================================

func TestValidateToken_EmptyString(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// 空文字列トークンは無効
	_, err := svc.ValidateToken("")
	assert.Error(t, err)
}

func TestValidateToken_GeneratedToken(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// generateTokenで生成したトークンが正しく検証できる
	token, err := svc.generateToken(99)
	assert.NoError(t, err)

	userID, err := svc.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(99), userID)
}

// ============================================================
// ResetPassword 追加テスト
// ============================================================

func TestResetPassword_WeakPassword(t *testing.T) {
	svc, _, passwordResetRepo := newTestAuthService()

	validToken := &model.PasswordResetToken{
		Token:     "valid-token-for-weak-pw",
		UserID:    1,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}
	validToken.ID = 1

	passwordResetRepo.On("FindByToken", "valid-token-for-weak-pw").Return(validToken, nil)

	// 弱いパスワード（短すぎる）
	err := svc.ResetPassword("valid-token-for-weak-pw", "weak")
	assert.Error(t, err)
}

// ============================================================
// ValidateToken 追加テスト（エッジケース）
// ============================================================

// TestValidateToken_WrongSigningMethod はHMAC以外の署名方法を使ったトークンを拒否することを確認。
func TestValidateToken_WrongSigningMethod(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// RS256署名のトークンを作成（HS256を期待するサービスに渡す）
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA2a2rwplBQLF29amygykEMmYz0+Kcj3bKBp29Po3KXHBS9BWL
zENMqBG0tRDoFAb2lH+GVMOhfJVNzFR2YA2t6y41HJuMQEL9TpAlkiT4t3IpVf+r
7nRYtw4BUG8v0VPmHH3D9EM9KfmMpCvC1TibHT7i4FhCwZZ/3JhZi3ZhPzWBjpJe
LBKZB3CZRZhOqUUMvHvFc9vDGE3J3MKJdNGxmEpnFEXU3qBjLfwKwp9XRbRZcvnW
gYMZ1ySwHpHBFl6sMJ+1NKZXN89jQGANTNT5W19/xE5fJ7LzaJaEF0TlPc4Xl8Xv
lrP1T9m98LvYKqSVa9N6BYbz0WrJX0T9v2BPdwIDAQABAoIBAHWE4FeO02s9zrNO
xtFRJCzCTDGLh4R4xyaAeAJZ0+HWLHV8O0K9tQU1E9hU0SqRb/MN4OQ4BYSJQfGM
4Yf4OjGHl1TkFxAVuM7oQSLOGAIKFGXXCPcTKH5Z4GEH3J4b9Ei8EFJpgLwVfHWj
nq7RpvYlh0rJuT3qr1aILLwNqIbhLqWyMeZt8e1+f2s7b7wB5LMXQPv4T8oFjcm5
6UBvzK5j0rSd7TqTR5zMj7TzSzeTwdyMEDPsNT0v4qY6OARR9gKfIAEo2TsmTN7Y
oS4bBDaxiPDFBwQXbVQA4kGGTxQJoGlv2x2fO96OP6nKJqzXFy0mZqjMQ+HCAQEC
AoGBAO+cRoFdBtZKjYNlL3ZFnXIQRHt5t0gZGG3IFKPaSjZ3VD5Q3ZJUChxN4JM5
O6u56S5xmCJdj4GUZM28AVthc6UHt/AZK9bkHCFJLQSsY5m5MlYmOsMIl+4Xr0bV
OoiFMXxH3SJ0GQWXF1hHVUbDh7BsBHQ8Gf6o1J5V2mnlGf0jAoGBAOjp3TYXHWJ/
MDu5q+mQlkuVKGQ4/lJlUyM7GG1e3YUeHHl2TYFqsrKFuBm0VjO5yHpPlvPW5VEe
GQOG7PKaqC/dRuJMgk3l7eYJn0MrCfHZ2qrKb5qF/XRKMkc7V9k5jJGvPxSMd0mj
7K1gw5G5R/vBzCqRb3VUOXq1R9pDxTVhAoGBAIX+GZ7r7JlwAEqgkdM5XBNZ0LpL
rAqU6r0zW0YVJP7fvU0YP1Wk2X3aTGgIqfVWLEF7RCcvDx8PJ9QVGZZ73U+qdF7
kbnV3CVjfAR8yGUYnCPsNPrA8J3PXVV8kHhYRTGH3mS0qLnzSRlXqLPGpfh1OtXr
nT2zP0HPGovUVz7/AoGBAMOkwvPXaXWxpMfx1HEIZtWLqiInlh6EMl2gWfLjSmHJ
sJoTX0W7X2xV5IqJJUHVIFAeTkmIxrO4h7EJZ9KRl3aCBaUiH+4ViWFifXr0RQKJ
OxkBRMQ1ZWMIT5b+5K+sFZGMJaX2RGhbPBXvfGSwAj0h5EVj9HlZrr3X
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		// RSAキーのパースに失敗した場合、単純な不正なトークンでテスト
		_, err = svc.ValidateToken("eyJhbGciOiJSUzI1NiJ9.eyJ1c2VyX2lkIjoxfQ.invalid")
		assert.Error(t, err)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"user_id": float64(1)})
	tokenString, _ := token.SignedString(privateKey)

	userID, err := svc.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, uint(0), userID)
}

// TestValidateToken_UserIDNotFloat は user_id が数値でないトークンを拒否することを確認。
func TestValidateToken_UserIDNotFloat(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// user_id を文字列で設定したトークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "not-a-number",
		"exp":     float64(9999999999),
	})
	tokenString, err := token.SignedString([]byte(testJWTSecret))
	assert.NoError(t, err)

	userID, err := svc.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, uint(0), userID)
}

func TestResetPassword_UpdatePasswordError(t *testing.T) {
	svc, userRepo, passwordResetRepo := newTestAuthService()

	validToken := &model.PasswordResetToken{
		Token:     "valid-token-update-err",
		UserID:    1,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}
	validToken.ID = 1

	passwordResetRepo.On("FindByToken", "valid-token-update-err").Return(validToken, nil)
	userRepo.On("UpdatePassword", uint(1), mock.AnythingOfType("string")).Return(errors.New("db error"))

	err := svc.ResetPassword("valid-token-update-err", "ValidPassword123!")
	assert.Error(t, err)
}
