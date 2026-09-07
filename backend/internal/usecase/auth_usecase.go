package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	// authTokenTTL は発行する JWT の有効期限。
	authTokenTTL = 72 * time.Hour
	// oauthStateTTL は OAuth の state トークンの有効期限。
	oauthStateTTL = 5 * time.Minute
	// passwordResetTTL はパスワードリセットトークンの有効期限。
	passwordResetTTL = 1 * time.Hour
	// githubLoginStatePurpose は GitHub ログイン用 state であることを示す用途。
	githubLoginStatePurpose = "github_login"
	// oauthStatePurpose は外部サービス連携用 state であることを示す用途。
	// 連携先ごとに "oauth_state:github" のようにサフィックスを付けて区別する。
	oauthStatePurpose = "oauth_state"
	// accessTokenPurpose はログイン後のアクセストークンであることを示す用途。
	accessTokenPurpose = "access_token"
	// maxUsernameCandidates は連番でユーザー名の空きを探す試行回数の上限。
	maxUsernameCandidates = 100
)

// OAuth の state を連携先ごとに区別するための識別子。
const (
	OAuthProviderGitHub  = "github"
	OAuthProviderSpotify = "spotify"
)

// AuthUserInput はユーザー登録の入力値。
type AuthUserInput struct {
	Name     string
	Username string
	Email    string
	Password string
}

// LoginInput はログインの入力値。
type LoginInput struct {
	Email    string
	Password string
}

// AuthResult は認証成功時の結果（JWT とユーザー情報）。
type AuthResult struct {
	Token string
	User  model.User
}

// authTokenIssuer は JWT の発行と検証をまとめる。usecase 間で共有する。
type authTokenIssuer struct {
	secret []byte
}

// newAuthTokenIssuer は署名鍵から発行者を作る。
func newAuthTokenIssuer(jwtSecret string) *authTokenIssuer {
	return &authTokenIssuer{secret: []byte(jwtSecret)}
}

// issue は指定ユーザーのアクセストークンを発行する。
// state トークンと同じ鍵で署名するため、用途を purpose クレームで明示する。
func (t *authTokenIssuer) issue(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"purpose": accessTokenPurpose,
		"exp":     time.Now().Add(authTokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.secret)
}

// parse は署名方式を検証したうえでクレームを取り出す。
func (t *authTokenIssuer) parse(tokenString, invalidMessage string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewError(domain.ErrCodeUnauthorized, "予期しない署名方式です", nil)
		}
		return t.secret, nil
	})
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeUnauthorized, invalidMessage, err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrUnauthorized
	}
	return claims, nil
}

// errAuthUserNotFound はログイン中のユーザーが見つからないときのエラー。
// DomainError ではないため handler では 500 になり、不在を素の DB エラーとして扱っていた
// 移行前の挙動と一致する。
var errAuthUserNotFound = errors.New("ユーザーが見つかりません")

// hashResetToken はリセットトークンを SHA-256 でハッシュ化する。
// DB には平文を保存せず、検証時も同じハッシュで比較する。
func hashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// RegisterUserUseCase は新規ユーザーを登録する。
type RegisterUserUseCase struct {
	users  repository.AuthUserRepository
	tokens *authTokenIssuer
}

// NewRegisterUserUseCase は RegisterUserUseCase を生成する。
func NewRegisterUserUseCase(users repository.AuthUserRepository, jwtSecret string) *RegisterUserUseCase {
	return &RegisterUserUseCase{users: users, tokens: newAuthTokenIssuer(jwtSecret)}
}

// Execute は入力を検証し、重複を弾いてユーザーを作成し、トークンを発行する。
func (uc *RegisterUserUseCase) Execute(ctx context.Context, input AuthUserInput) (*AuthResult, error) {
	if err := domain.ValidateEmail(input.Email); err != nil {
		return nil, err
	}
	if err := domain.ValidatePassword(input.Password); err != nil {
		return nil, err
	}
	if err := domain.ValidateUsername(input.Username); err != nil {
		return nil, err
	}

	if existing, _ := uc.users.FindByEmail(ctx, input.Email); existing != nil {
		return nil, domain.NewError(domain.ErrCodeAlreadyExists, "このメールアドレスは既に登録されています", nil)
	}
	if existing, _ := uc.users.FindByUsername(ctx, input.Username); existing != nil {
		return nil, domain.NewError(domain.ErrCodeAlreadyExists, "このユーザー名は既に使用されています", nil)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "パスワードのハッシュ化に失敗しました", err)
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
		Username: input.Username,
	}
	if err := uc.users.Create(ctx, user); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "ユーザーの作成に失敗しました", err)
	}

	token, err := uc.tokens.issue(user.ID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "トークンの生成に失敗しました", err)
	}
	return &AuthResult{Token: token, User: *user}, nil
}

// LoginUseCase はメールアドレスとパスワードで認証する。
type LoginUseCase struct {
	users  repository.AuthUserRepository
	tokens *authTokenIssuer
}

// NewLoginUseCase は LoginUseCase を生成する。
func NewLoginUseCase(users repository.AuthUserRepository, jwtSecret string) *LoginUseCase {
	return &LoginUseCase{users: users, tokens: newAuthTokenIssuer(jwtSecret)}
}

// Execute は認証に成功したらトークンを発行する。
// ユーザーが存在しない場合とパスワード不一致は、区別せず 401 を返す。
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*AuthResult, error) {
	if err := domain.ValidateEmail(input.Email); err != nil {
		return nil, err
	}

	user, err := uc.users.FindByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, domain.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	token, err := uc.tokens.issue(user.ID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "トークンの生成に失敗しました", err)
	}
	return &AuthResult{Token: token, User: *user}, nil
}

// ValidateAuthTokenUseCase はアクセストークンを検証する。
type ValidateAuthTokenUseCase struct {
	tokens *authTokenIssuer
}

// NewValidateAuthTokenUseCase は ValidateAuthTokenUseCase を生成する。
func NewValidateAuthTokenUseCase(jwtSecret string) *ValidateAuthTokenUseCase {
	return &ValidateAuthTokenUseCase{tokens: newAuthTokenIssuer(jwtSecret)}
}

// Execute はトークンからユーザー ID を取り出す。外部 I/O を行わないため ctx は取らない。
// OAuth の state トークンも同じ鍵で署名され user_id を含むため、purpose を検証して弾く。
// purpose を持たないトークンは、この検証を入れる前に発行した既存のアクセストークンとして受け入れる。
func (uc *ValidateAuthTokenUseCase) Execute(tokenString string) (uint, error) {
	claims, err := uc.tokens.parse(tokenString, "無効なトークンです")
	if err != nil {
		return 0, err
	}
	if purpose, ok := claims["purpose"].(string); ok && purpose != accessTokenPurpose {
		return 0, domain.ErrUnauthorized
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, domain.ErrUnauthorized
	}
	return uint(userID), nil
}

// GitHubLoginStateUseCase は GitHub ログイン用の state トークンを発行・検証する。
type GitHubLoginStateUseCase struct {
	tokens *authTokenIssuer
}

// NewGitHubLoginStateUseCase は GitHubLoginStateUseCase を生成する。
func NewGitHubLoginStateUseCase(jwtSecret string) *GitHubLoginStateUseCase {
	return &GitHubLoginStateUseCase{tokens: newAuthTokenIssuer(jwtSecret)}
}

// Generate は CSRF 防止用の state を発行する（有効期限 5 分）。
func (uc *GitHubLoginStateUseCase) Generate() (string, error) {
	claims := jwt.MapClaims{
		"purpose": githubLoginStatePurpose,
		"exp":     time.Now().Add(oauthStateTTL).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(uc.tokens.secret)
}

// Validate は state を検証する。用途が異なる場合も 401 として扱う。
func (uc *GitHubLoginStateUseCase) Validate(state string) error {
	claims, err := uc.tokens.parse(state, "無効なログインステートです")
	if err != nil {
		return err
	}
	if purpose, _ := claims["purpose"].(string); purpose != githubLoginStatePurpose {
		return domain.ErrUnauthorized
	}
	return nil
}

// OAuthStateUseCase は外部サービス連携用の state トークン（ユーザー ID 入り）を発行・検証する。
// provider ごとに別インスタンスを作り、連携先をまたいだ state の使い回しを防ぐ。
type OAuthStateUseCase struct {
	tokens   *authTokenIssuer
	provider string
}

// NewOAuthStateUseCase は指定した連携先向けの OAuthStateUseCase を生成する。
func NewOAuthStateUseCase(jwtSecret, provider string) *OAuthStateUseCase {
	return &OAuthStateUseCase{tokens: newAuthTokenIssuer(jwtSecret), provider: provider}
}

// Generate はユーザー ID を埋め込んだ state を発行する（有効期限 5 分）。
func (uc *OAuthStateUseCase) Generate(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"purpose": uc.purpose(),
		"exp":     time.Now().Add(oauthStateTTL).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(uc.tokens.secret)
}

// Validate は state を検証し、埋め込まれたユーザー ID を返す。
// 連携先が異なる state は受け付けない。
func (uc *OAuthStateUseCase) Validate(state string) (uint, error) {
	claims, err := uc.tokens.parse(state, "無効なOAuthステートです")
	if err != nil {
		return 0, err
	}
	if purpose, _ := claims["purpose"].(string); purpose != uc.purpose() {
		return 0, domain.ErrUnauthorized
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, domain.ErrUnauthorized
	}
	return uint(userID), nil
}

// purpose は連携先を含む用途文字列を返す。
func (uc *OAuthStateUseCase) purpose() string {
	if uc.provider == "" {
		return oauthStatePurpose
	}
	return oauthStatePurpose + ":" + uc.provider
}

// GitHubLoginUseCase は GitHub アカウントでログイン（必要なら登録）する。
type GitHubLoginUseCase struct {
	users  repository.AuthUserRepository
	tokens *authTokenIssuer
}

// NewGitHubLoginUseCase は GitHubLoginUseCase を生成する。
func NewGitHubLoginUseCase(users repository.AuthUserRepository, jwtSecret string) *GitHubLoginUseCase {
	return &GitHubLoginUseCase{users: users, tokens: newAuthTokenIssuer(jwtSecret)}
}

// Execute は GitHub ID → メールアドレス → 新規作成の順にユーザーを解決してログインさせる。
func (uc *GitHubLoginUseCase) Execute(ctx context.Context, ghUser *model.GitHubUserInfo, accessToken string) (*AuthResult, error) {
	// 1. GitHub ID で既存ユーザーを探す
	if user, err := uc.users.FindByGitHubID(ctx, ghUser.ID); err == nil && user != nil {
		user.GitHubToken = accessToken
		user.GitHubUsername = ghUser.Login
		if ghUser.AvatarURL != "" {
			user.AvatarURL = ghUser.AvatarURL
		}
		// 更新に失敗してもログイン自体は続行する（移行前と同じ）
		_ = uc.users.Update(ctx, user)
		return uc.issueFor(user)
	}

	// 2. メールアドレスで既存ユーザーを探して連携する
	if ghUser.Email != "" {
		if user, err := uc.users.FindByEmail(ctx, ghUser.Email); err == nil && user != nil {
			user.GitHubID = ghUser.ID
			user.GitHubToken = accessToken
			user.GitHubUsername = ghUser.Login
			user.GitHubConnected = true
			if ghUser.AvatarURL != "" {
				user.AvatarURL = ghUser.AvatarURL
			}
			_ = uc.users.Update(ctx, user)
			return uc.issueFor(user)
		}
	}

	// 3. 新規ユーザーを作成する
	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}
	email := ghUser.Email
	if email == "" {
		email = ghUser.Login + "@github.local"
	}

	newUser := &model.User{
		Name:            name,
		Email:           email,
		Username:        uc.uniqueUsername(ctx, ghUser.Login),
		GitHubID:        ghUser.ID,
		GitHubUsername:  ghUser.Login,
		GitHubToken:     accessToken,
		GitHubConnected: true,
		AvatarURL:       ghUser.AvatarURL,
	}
	if err := uc.users.Create(ctx, newUser); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "ユーザーの作成に失敗しました", err)
	}
	return uc.issueFor(newUser)
}

// issueFor は指定ユーザーのトークンを発行して結果にまとめる。
func (uc *GitHubLoginUseCase) issueFor(user *model.User) (*AuthResult, error) {
	token, err := uc.tokens.issue(user.ID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "トークンの生成に失敗しました", err)
	}
	return &AuthResult{Token: token, User: *user}, nil
}

// uniqueUsername は使われていないユーザー名を作る。
// 候補が使用済みなら末尾に連番を付ける（alice → alice2 → alice3）。
// 連番で決まらない場合はランダムなサフィックスにフォールバックし、無限ループを避ける。
func (uc *GitHubLoginUseCase) uniqueUsername(ctx context.Context, base string) string {
	candidate := base
	for i := 2; i <= maxUsernameCandidates; i++ {
		existing, _ := uc.users.FindByUsername(ctx, candidate)
		if existing == nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s%d", base, i)
	}

	suffix := make([]byte, 4)
	if _, err := rand.Read(suffix); err != nil {
		// 乱数が使えない場合でも衝突しにくい値にする
		return fmt.Sprintf("%s%d", base, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s-%s", base, hex.EncodeToString(suffix))
}

// GetMeUseCase はログイン中のユーザー情報を返す。
type GetMeUseCase struct {
	users repository.AuthUserRepository
}

// NewGetMeUseCase は GetMeUseCase を生成する。
func NewGetMeUseCase(users repository.AuthUserRepository) *GetMeUseCase {
	return &GetMeUseCase{users: users}
}

// Execute はユーザーを返す。存在しない場合は移行前と同じく内部エラーになる。
func (uc *GetMeUseCase) Execute(ctx context.Context, userID uint) (*model.User, error) {
	user, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errAuthUserNotFound
	}
	return user, nil
}

// DeleteAccountUseCase は退会（ユーザーと関連データの削除）を行う。
type DeleteAccountUseCase struct {
	users repository.AuthUserRepository
}

// NewDeleteAccountUseCase は DeleteAccountUseCase を生成する。
func NewDeleteAccountUseCase(users repository.AuthUserRepository) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{users: users}
}

// Execute はパスワードを検証したうえでアカウントを削除する。
// GitHub のみで登録したユーザー（パスワード未設定）は検証をスキップする。
func (uc *DeleteAccountUseCase) Execute(ctx context.Context, userID uint, password string) error {
	user, err := uc.users.FindByIDWithPassword(ctx, userID)
	if err != nil || user == nil {
		return domain.NewError(domain.ErrCodeNotFound, "ユーザーが見つかりません", err)
	}

	if user.Password != "" {
		if password == "" {
			return domain.NewError(domain.ErrCodeBadRequest, "パスワードの入力が必要です", nil)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return domain.NewError(domain.ErrCodeForbidden, "パスワードが正しくありません", nil)
		}
	}

	if err := uc.users.DeleteWithRelatedData(ctx, userID); err != nil {
		return domain.NewError(domain.ErrCodeInternal, "アカウント削除に失敗しました", err)
	}
	return nil
}

// RequestPasswordResetUseCase はパスワードリセットトークンを発行する。
type RequestPasswordResetUseCase struct {
	users  repository.AuthUserRepository
	tokens repository.PasswordResetTokenRepository
}

// NewRequestPasswordResetUseCase は RequestPasswordResetUseCase を生成する。
func NewRequestPasswordResetUseCase(
	users repository.AuthUserRepository,
	tokens repository.PasswordResetTokenRepository,
) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{users: users, tokens: tokens}
}

// Execute はリセットトークン（平文）を返す。
// メールアドレスの存在有無を漏らさないため、未登録でもエラーにせず空文字を返す。
func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, email string) (string, error) {
	user, err := uc.users.FindByEmail(ctx, email)
	if err != nil {
		// 検索の失敗はアカウント個別の事情ではないため、成功に見せず 500 を返す
		// （エラーにしてもメールアドレスの存在有無は漏れない）。「未登録」の成功扱いとは別物。
		return "", domain.NewError(domain.ErrCodeInternal, "パスワードリセットの処理に失敗しました", err)
	}
	if user == nil {
		return "", nil
	}

	// 既存のトークンを無効化する。失敗しても新しいトークンの発行は続ける（移行前と同じ）
	_ = uc.tokens.InvalidateUserTokens(ctx, user.ID)

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", domain.NewError(domain.ErrCodeInternal, "トークンの生成に失敗しました", err)
	}
	token := hex.EncodeToString(tokenBytes)

	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		Token:     hashResetToken(token),
		ExpiresAt: time.Now().Add(passwordResetTTL),
	}
	if err := uc.tokens.Create(ctx, resetToken); err != nil {
		return "", domain.NewError(domain.ErrCodeInternal, "リセットトークンの保存に失敗しました", err)
	}
	return token, nil
}

// ResetPasswordUseCase はリセットトークンでパスワードを再設定する。
type ResetPasswordUseCase struct {
	users  repository.AuthUserRepository
	tokens repository.PasswordResetTokenRepository
}

// NewResetPasswordUseCase は ResetPasswordUseCase を生成する。
func NewResetPasswordUseCase(
	users repository.AuthUserRepository,
	tokens repository.PasswordResetTokenRepository,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{users: users, tokens: tokens}
}

// Execute はトークンを検証し、新しいパスワードを保存してトークンを使用済みにする。
func (uc *ResetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	resetToken, err := uc.tokens.FindByToken(ctx, hashResetToken(token))
	if err != nil || resetToken == nil {
		return domain.NewError(domain.ErrCodeBadRequest, "無効なリセットトークンです", err)
	}
	if !resetToken.IsValid() {
		return domain.NewError(domain.ErrCodeBadRequest, "リセットトークンが期限切れです", nil)
	}
	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return domain.NewError(domain.ErrCodeInternal, "パスワードのハッシュ化に失敗しました", err)
	}
	if err := uc.users.UpdatePassword(ctx, resetToken.UserID, string(hashedPassword)); err != nil {
		return domain.NewError(domain.ErrCodeInternal, "パスワードの更新に失敗しました", err)
	}

	// 使用済みへの更新に失敗してもパスワード変更自体は完了している（移行前と同じ）
	_ = uc.tokens.MarkAsUsed(ctx, resetToken.ID)
	return nil
}
