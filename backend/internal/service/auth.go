package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService は認証・認可に関するビジネスロジックを提供する。
// JWT生成・検証、パスワードハッシュ化、GitHub OAuthログインを担当する。
type AuthService struct {
	userRepo          repository.UserRepositoryInterface
	passwordResetRepo repository.PasswordResetRepositoryInterface
	jwtSecret         []byte
}

// NewAuthService は新しいAuthServiceインスタンスを生成する。
func NewAuthService(userRepo repository.UserRepositoryInterface, passwordResetRepo repository.PasswordResetRepositoryInterface, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		jwtSecret:         []byte(jwtSecret),
	}
}

// RegisterInput はユーザー登録リクエストの入力値を表す。
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput はログインリクエストの入力値を表す。
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse は認証成功時のレスポンスを表す（JWTトークン + ユーザー情報）。
type AuthResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// Register は新規ユーザーを登録し、JWTトークンを発行する。
// メールアドレスの重複チェックとパスワードのbcryptハッシュ化を行う。
func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	// バリデーション
	if err := domain.ValidateEmail(input.Email); err != nil {
		return nil, err
	}
	if err := domain.ValidatePassword(input.Password); err != nil {
		return nil, err
	}
	if err := domain.ValidateUsername(input.Username); err != nil {
		return nil, err
	}

	existing, _ := s.userRepo.FindByEmail(input.Email)
	if existing != nil {
		return nil, domain.NewError(domain.ErrCodeAlreadyExists, "このメールアドレスは既に登録されています", nil)
	}

	existingUsername, _ := s.userRepo.FindByUsername(input.Username)
	if existingUsername != nil {
		return nil, domain.NewError(domain.ErrCodeAlreadyExists, "このユーザー名は既に使用されています", nil)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
		Username: input.Username,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

// Login はメールアドレスとパスワードで認証し、JWTトークンを発行する。
func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	// バリデーション
	if err := domain.ValidateEmail(input.Email); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

// ValidateToken はJWTトークンを検証し、含まれるユーザーIDを返す。
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewError(domain.ErrCodeUnauthorized, "unexpected signing method", nil)
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, domain.ErrUnauthorized
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, domain.ErrUnauthorized
	}

	return uint(userID), nil
}

// GenerateLoginState はGitHubログイン用のCSRF防止stateトークンを生成する（有効期限5分）。
func (s *AuthService) GenerateLoginState() (string, error) {
	claims := jwt.MapClaims{
		"purpose": "github_login",
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateLoginState はGitHubログイン用のstateトークンを検証する。
func (s *AuthService) ValidateLoginState(state string) error {
	token, err := jwt.Parse(state, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return domain.ErrUnauthorized
	}
	purpose, _ := claims["purpose"].(string)
	if purpose != "github_login" {
		return domain.ErrUnauthorized
	}
	return nil
}

// GitHubLogin はGitHubユーザー情報を使ってログイン/登録処理を行う。
// 1. GitHub IDで既存ユーザーを検索 → 見つかればトークン更新してログイン
// 2. メールアドレスで既存ユーザーを検索 → 見つかればGitHub連携してログイン
// 3. どちらも見つからなければ新規ユーザーを作成
func (s *AuthService) GitHubLogin(ghUser *GitHubUserInfo, accessToken string) (*AuthResponse, error) {
	// 1. GitHub IDで既存ユーザーを検索
	user, err := s.userRepo.FindByGitHubID(ghUser.ID)
	if err == nil && user != nil {
		user.GitHubToken = accessToken
		user.GitHubUsername = ghUser.Login
		if ghUser.AvatarURL != "" {
			user.AvatarURL = ghUser.AvatarURL
		}
		s.userRepo.Update(user)

		token, err := s.generateToken(user.ID)
		if err != nil {
			return nil, err
		}
		return &AuthResponse{Token: token, User: *user}, nil
	}

	// 2. メールアドレスで既存ユーザーを検索してGitHub連携
	if ghUser.Email != "" {
		user, err = s.userRepo.FindByEmail(ghUser.Email)
		if err == nil && user != nil {
			user.GitHubID = ghUser.ID
			user.GitHubToken = accessToken
			user.GitHubUsername = ghUser.Login
			user.GitHubConnected = true
			if ghUser.AvatarURL != "" {
				user.AvatarURL = ghUser.AvatarURL
			}
			s.userRepo.Update(user)

			token, err := s.generateToken(user.ID)
			if err != nil {
				return nil, err
			}
			return &AuthResponse{Token: token, User: *user}, nil
		}
	}

	// 3. 新規ユーザーを作成
	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}
	email := ghUser.Email
	if email == "" {
		email = ghUser.Login + "@github.local"
	}

	// GitHubのログイン名をベースにユニークなユーザー名を生成
	username := s.generateUsername(ghUser.Login)

	newUser := &model.User{
		Name:            name,
		Email:           email,
		Username:        username,
		GitHubID:        ghUser.ID,
		GitHubUsername:  ghUser.Login,
		GitHubToken:     accessToken,
		GitHubConnected: true,
		AvatarURL:       ghUser.AvatarURL,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	token, err := s.generateToken(newUser.ID)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, User: *newUser}, nil
}

// GenerateOAuthState はGitHub連携用のOAuth stateトークンを生成する（ユーザーID埋め込み、有効期限5分）。
func (s *AuthService) GenerateOAuthState(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"purpose": "oauth_state",
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateOAuthState はOAuth stateトークンを検証し、埋め込まれたユーザーIDを返す。
func (s *AuthService) ValidateOAuthState(state string) (uint, error) {
	token, err := jwt.Parse(state, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, domain.ErrUnauthorized
	}

	purpose, _ := claims["purpose"].(string)
	if purpose != "oauth_state" {
		return 0, domain.ErrUnauthorized
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, domain.ErrUnauthorized
	}

	return uint(userID), nil
}

// GetMe は現在のログインユーザー情報をIDで取得する。
func (s *AuthService) GetMe(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

// DeleteAccount はパスワード検証後にユーザーアカウントと関連データを完全削除する。
// GitHubのみで登録したユーザー（パスワード未設定）の場合はパスワード検証をスキップする。
func (s *AuthService) DeleteAccount(userID uint, password string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound
	}

	// パスワードが設定されている場合（GitHub専用アカウントでない場合）は検証
	if user.Password != "" {
		if password == "" {
			return ErrBadRequest
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return ErrForbidden
		}
	}

	return s.userRepo.DeleteWithRelatedData(userID)
}

// RequestPasswordReset はパスワードリセットトークンを生成する。
// セキュリティ上、メールアドレスが存在しない場合でもエラーを返さない。
func (s *AuthService) RequestPasswordReset(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// メールアドレスの存在有無を漏らさない
		return "", nil
	}

	// 既存トークンを無効化
	s.passwordResetRepo.InvalidateUserTokens(user.ID)

	// セキュアなランダムトークンを生成
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// リセットトークンを作成（有効期限1時間）
	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if err := s.passwordResetRepo.Create(resetToken); err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword は有効なリセットトークンを使ってパスワードを再設定する。
func (s *AuthService) ResetPassword(token string, newPassword string) error {
	resetToken, err := s.passwordResetRepo.FindByToken(token)
	if err != nil {
		return ErrBadRequest
	}

	if !resetToken.IsValid() {
		return ErrBadRequest
	}

	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(resetToken.UserID, string(hashedPassword)); err != nil {
		return err
	}

	s.passwordResetRepo.MarkAsUsed(resetToken.ID)
	return nil
}

// generateUsername は候補のユーザー名からユニークなユーザー名を生成する。
// 候補が既に使われている場合は末尾に連番を付与する（例: alice → alice2 → alice3）。
func (s *AuthService) generateUsername(base string) string {
	candidate := base
	for i := 2; ; i++ {
		existing, _ := s.userRepo.FindByUsername(candidate)
		if existing == nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s%d", base, i)
	}
}

// generateToken は指定ユーザーIDのJWTトークンを生成する（有効期限72時間）。
func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
