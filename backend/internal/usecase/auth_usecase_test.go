package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const testAuthSecret = "test-secret-for-auth-usecase"

// mockAuthUsers は usecase/repository.AuthUserRepository のモック。
type mockAuthUsers struct{ mock.Mock }

func (m *mockAuthUsers) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthUsers) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthUsers) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthUsers) FindByGitHubID(ctx context.Context, githubID int64) (*model.User, error) {
	args := m.Called(ctx, githubID)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthUsers) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil && user.ID == 0 {
		user.ID = 1
	}
	return args.Error(0)
}

func (m *mockAuthUsers) Update(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *mockAuthUsers) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	return m.Called(ctx, userID, hashedPassword).Error(0)
}

func (m *mockAuthUsers) DeleteWithRelatedData(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// mockResetTokens は usecase/repository.PasswordResetTokenRepository のモック。
type mockResetTokens struct{ mock.Mock }

func (m *mockResetTokens) Create(ctx context.Context, token *model.PasswordResetToken) error {
	return m.Called(ctx, token).Error(0)
}

func (m *mockResetTokens) FindByToken(ctx context.Context, hashedToken string) (*model.PasswordResetToken, error) {
	args := m.Called(ctx, hashedToken)
	t, _ := args.Get(0).(*model.PasswordResetToken)
	return t, args.Error(1)
}

func (m *mockResetTokens) MarkAsUsed(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockResetTokens) InvalidateUserTokens(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

func authHashed(plain string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(h)
}

// ============================================================
// 登録・ログイン
// ============================================================

func TestRegisterUserUseCase(t *testing.T) {
	ctx := context.Background()
	input := usecase.AuthUserInput{Name: "Test", Username: "tester", Email: "test@example.com", Password: "password123"}

	t.Run("ユーザーを作成してトークンを発行する", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByEmail", mock.Anything, input.Email).Return(nil, nil)
		users.On("FindByUsername", mock.Anything, input.Username).Return(nil, nil)
		users.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			// パスワードは bcrypt ハッシュとして保存する
			return u.Password != input.Password &&
				bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)) == nil
		})).Return(nil)

		got, err := usecase.NewRegisterUserUseCase(users, testAuthSecret).Execute(ctx, input)
		require.NoError(t, err)
		assert.NotEmpty(t, got.Token)
		assert.Equal(t, input.Email, got.User.Email)
		users.AssertExpectations(t)

		// 発行したトークンは同じ鍵で検証できる
		userID, err := usecase.NewValidateAuthTokenUseCase(testAuthSecret).Execute(got.Token)
		require.NoError(t, err)
		assert.Equal(t, got.User.ID, userID)
	})

	t.Run("メールアドレスが登録済みなら 409", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByEmail", mock.Anything, input.Email).Return(&model.User{ID: 9}, nil)

		_, err := usecase.NewRegisterUserUseCase(users, testAuthSecret).Execute(ctx, input)
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeAlreadyExists, domainErr.Code)
		users.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("パスワードが短ければ検証エラー", func(t *testing.T) {
		users := new(mockAuthUsers)

		short := input
		short.Password = "short"
		_, err := usecase.NewRegisterUserUseCase(users, testAuthSecret).Execute(ctx, short)
		assert.Error(t, err)
		users.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
	})
}

func TestLoginUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("パスワードが一致すればトークンを返す", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByEmail", mock.Anything, "test@example.com").
			Return(&model.User{ID: 3, Email: "test@example.com", Password: authHashed("password123")}, nil)

		got, err := usecase.NewLoginUseCase(users, testAuthSecret).
			Execute(ctx, usecase.LoginInput{Email: "test@example.com", Password: "password123"})
		require.NoError(t, err)
		assert.NotEmpty(t, got.Token)
		assert.Equal(t, uint(3), got.User.ID)
	})

	t.Run("ユーザーが存在しなければ 401", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByEmail", mock.Anything, "none@example.com").Return(nil, nil)

		_, err := usecase.NewLoginUseCase(users, testAuthSecret).
			Execute(ctx, usecase.LoginInput{Email: "none@example.com", Password: "password123"})
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("パスワードが違えば 401", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByEmail", mock.Anything, "test@example.com").
			Return(&model.User{ID: 3, Password: authHashed("password123")}, nil)

		_, err := usecase.NewLoginUseCase(users, testAuthSecret).
			Execute(ctx, usecase.LoginInput{Email: "test@example.com", Password: "wrong"})
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})
}

// ============================================================
// トークンと state
// ============================================================

func TestValidateAuthTokenUseCase(t *testing.T) {
	users := new(mockAuthUsers)
	users.On("FindByEmail", mock.Anything, mock.Anything).
		Return(&model.User{ID: 7, Password: authHashed("password123")}, nil)

	got, err := usecase.NewLoginUseCase(users, testAuthSecret).
		Execute(context.Background(), usecase.LoginInput{Email: "a@example.com", Password: "password123"})
	require.NoError(t, err)

	t.Run("自分が発行したトークンは検証できる", func(t *testing.T) {
		userID, err := usecase.NewValidateAuthTokenUseCase(testAuthSecret).Execute(got.Token)
		require.NoError(t, err)
		assert.Equal(t, uint(7), userID)
	})

	t.Run("鍵が違えば 401", func(t *testing.T) {
		_, err := usecase.NewValidateAuthTokenUseCase("another-secret").Execute(got.Token)
		assert.Error(t, err)
	})

	t.Run("壊れたトークンは 401", func(t *testing.T) {
		_, err := usecase.NewValidateAuthTokenUseCase(testAuthSecret).Execute("not-a-token")
		assert.Error(t, err)
	})
}

func TestOAuthStateUseCase(t *testing.T) {
	state := usecase.NewOAuthStateUseCase(testAuthSecret)

	t.Run("発行した state からユーザー ID を取り出せる", func(t *testing.T) {
		token, err := state.Generate(12)
		require.NoError(t, err)

		userID, err := state.Validate(token)
		require.NoError(t, err)
		assert.Equal(t, uint(12), userID)
	})

	t.Run("ログイン用 state は連携用として使えない", func(t *testing.T) {
		loginState, err := usecase.NewGitHubLoginStateUseCase(testAuthSecret).Generate()
		require.NoError(t, err)

		_, err = state.Validate(loginState)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("連携用 state はログイン用として使えない", func(t *testing.T) {
		oauthState, err := state.Generate(1)
		require.NoError(t, err)

		err = usecase.NewGitHubLoginStateUseCase(testAuthSecret).Validate(oauthState)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})
}

// ============================================================
// GitHub ログイン
// ============================================================

func TestGitHubLoginUseCase(t *testing.T) {
	ctx := context.Background()
	ghUser := &model.GitHubUserInfo{ID: 42, Login: "dev", Email: "dev@example.com", Name: "Dev", AvatarURL: "a.png"}

	t.Run("GitHub ID で見つかればトークンを更新してログインする", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByGitHubID", mock.Anything, int64(42)).Return(&model.User{ID: 5}, nil)
		users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.GitHubToken == "token" && u.GitHubUsername == "dev" && u.AvatarURL == "a.png"
		})).Return(nil)

		got, err := usecase.NewGitHubLoginUseCase(users, testAuthSecret).Execute(ctx, ghUser, "token")
		require.NoError(t, err)
		assert.Equal(t, uint(5), got.User.ID)
		users.AssertExpectations(t)
	})

	t.Run("メールアドレスで見つかれば連携してログインする", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByGitHubID", mock.Anything, int64(42)).Return(nil, nil)
		users.On("FindByEmail", mock.Anything, "dev@example.com").Return(&model.User{ID: 6}, nil)
		users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.GitHubID == 42 && u.GitHubConnected
		})).Return(nil)

		got, err := usecase.NewGitHubLoginUseCase(users, testAuthSecret).Execute(ctx, ghUser, "token")
		require.NoError(t, err)
		assert.Equal(t, uint(6), got.User.ID)
	})

	t.Run("見つからなければ新規作成する", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByGitHubID", mock.Anything, int64(42)).Return(nil, nil)
		users.On("FindByEmail", mock.Anything, "dev@example.com").Return(nil, nil)
		users.On("FindByUsername", mock.Anything, "dev").Return(nil, nil)
		users.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Username == "dev" && u.GitHubConnected && u.Email == "dev@example.com"
		})).Return(nil)

		got, err := usecase.NewGitHubLoginUseCase(users, testAuthSecret).Execute(ctx, ghUser, "token")
		require.NoError(t, err)
		assert.NotEmpty(t, got.Token)
		users.AssertExpectations(t)
	})

	t.Run("メールアドレスが無ければプレースホルダを使う", func(t *testing.T) {
		users := new(mockAuthUsers)
		noEmail := &model.GitHubUserInfo{ID: 43, Login: "nomail"}
		users.On("FindByGitHubID", mock.Anything, int64(43)).Return(nil, nil)
		users.On("FindByUsername", mock.Anything, "nomail").Return(nil, nil)
		users.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			// 名前が無い場合はログイン名で代用する
			return u.Email == "nomail@github.local" && u.Name == "nomail"
		})).Return(nil)

		_, err := usecase.NewGitHubLoginUseCase(users, testAuthSecret).Execute(ctx, noEmail, "token")
		require.NoError(t, err)
		users.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
	})

	t.Run("ユーザー名が使用済みなら連番を付ける", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByGitHubID", mock.Anything, int64(42)).Return(nil, nil)
		users.On("FindByEmail", mock.Anything, "dev@example.com").Return(nil, nil)
		users.On("FindByUsername", mock.Anything, "dev").Return(&model.User{ID: 1}, nil)
		users.On("FindByUsername", mock.Anything, "dev2").Return(&model.User{ID: 2}, nil)
		users.On("FindByUsername", mock.Anything, "dev3").Return(nil, nil)
		users.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Username == "dev3"
		})).Return(nil)

		_, err := usecase.NewGitHubLoginUseCase(users, testAuthSecret).Execute(ctx, ghUser, "token")
		require.NoError(t, err)
		users.AssertExpectations(t)
	})
}

// ============================================================
// 退会・パスワードリセット
// ============================================================

func TestDeleteAccountUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("パスワードが一致すれば削除する", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Password: authHashed("password123")}, nil)
		users.On("DeleteWithRelatedData", mock.Anything, uint(1)).Return(nil)

		require.NoError(t, usecase.NewDeleteAccountUseCase(users).Execute(ctx, 1, "password123"))
		users.AssertExpectations(t)
	})

	// GitHub のみで登録したユーザーはパスワードを持たないため検証を省略する。
	t.Run("パスワード未設定のユーザーは検証をスキップする", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
		users.On("DeleteWithRelatedData", mock.Anything, uint(1)).Return(nil)

		require.NoError(t, usecase.NewDeleteAccountUseCase(users).Execute(ctx, 1, ""))
		users.AssertExpectations(t)
	})

	t.Run("パスワードが違えば 403", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Password: authHashed("password123")}, nil)

		err := usecase.NewDeleteAccountUseCase(users).Execute(ctx, 1, "wrong")
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeForbidden, domainErr.Code)
		users.AssertNotCalled(t, "DeleteWithRelatedData", mock.Anything, mock.Anything)
	})

	t.Run("パスワード設定済みで未入力なら 400", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Password: authHashed("password123")}, nil)

		err := usecase.NewDeleteAccountUseCase(users).Execute(ctx, 1, "")
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	})

	t.Run("ユーザーが存在しなければ 404", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		err := usecase.NewDeleteAccountUseCase(users).Execute(ctx, 1, "password123")
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	})
}

func TestPasswordResetUseCases(t *testing.T) {
	ctx := context.Background()

	t.Run("トークンはハッシュ化して保存し、平文を返す", func(t *testing.T) {
		users := new(mockAuthUsers)
		tokens := new(mockResetTokens)
		users.On("FindByEmail", mock.Anything, "a@example.com").Return(&model.User{ID: 1}, nil)
		tokens.On("InvalidateUserTokens", mock.Anything, uint(1)).Return(nil)

		var saved *model.PasswordResetToken
		tokens.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			saved = args.Get(1).(*model.PasswordResetToken)
		})

		plain, err := usecase.NewRequestPasswordResetUseCase(users, tokens).Execute(ctx, "a@example.com")
		require.NoError(t, err)
		assert.NotEmpty(t, plain)
		require.NotNil(t, saved)
		assert.NotEqual(t, plain, saved.Token, "平文をそのまま保存しない")
		assert.Len(t, saved.Token, 64, "SHA-256 の 16 進表現で保存する")
		assert.True(t, saved.ExpiresAt.After(time.Now()))
	})

	t.Run("未登録のメールアドレスでもエラーにしない", func(t *testing.T) {
		users := new(mockAuthUsers)
		tokens := new(mockResetTokens)
		users.On("FindByEmail", mock.Anything, "none@example.com").Return(nil, nil)

		plain, err := usecase.NewRequestPasswordResetUseCase(users, tokens).Execute(ctx, "none@example.com")
		require.NoError(t, err)
		assert.Empty(t, plain)
		tokens.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("有効なトークンでパスワードを更新する", func(t *testing.T) {
		users := new(mockAuthUsers)
		tokens := new(mockResetTokens)
		tokens.On("FindByToken", mock.Anything, mock.AnythingOfType("string")).
			Return(&model.PasswordResetToken{ID: 2, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
		users.On("UpdatePassword", mock.Anything, uint(1), mock.MatchedBy(func(hashed string) bool {
			return bcrypt.CompareHashAndPassword([]byte(hashed), []byte("newpassword123")) == nil
		})).Return(nil)
		tokens.On("MarkAsUsed", mock.Anything, uint(2)).Return(nil)

		require.NoError(t, usecase.NewResetPasswordUseCase(users, tokens).Execute(ctx, "raw", "newpassword123"))
		users.AssertExpectations(t)
		tokens.AssertExpectations(t)
	})

	t.Run("トークンが無ければ 400", func(t *testing.T) {
		users := new(mockAuthUsers)
		tokens := new(mockResetTokens)
		tokens.On("FindByToken", mock.Anything, mock.Anything).Return(nil, nil)

		err := usecase.NewResetPasswordUseCase(users, tokens).Execute(ctx, "raw", "newpassword123")
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
		users.AssertNotCalled(t, "UpdatePassword", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("期限切れのトークンは 400", func(t *testing.T) {
		users := new(mockAuthUsers)
		tokens := new(mockResetTokens)
		tokens.On("FindByToken", mock.Anything, mock.Anything).
			Return(&model.PasswordResetToken{ID: 2, UserID: 1, ExpiresAt: time.Now().Add(-time.Hour)}, nil)

		err := usecase.NewResetPasswordUseCase(users, tokens).Execute(ctx, "raw", "newpassword123")
		assert.Error(t, err)
		users.AssertNotCalled(t, "UpdatePassword", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("使用済みのトークンは 400", func(t *testing.T) {
		users := new(mockAuthUsers)
		tokens := new(mockResetTokens)
		tokens.On("FindByToken", mock.Anything, mock.Anything).
			Return(&model.PasswordResetToken{ID: 2, UserID: 1, Used: true, ExpiresAt: time.Now().Add(time.Hour)}, nil)

		err := usecase.NewResetPasswordUseCase(users, tokens).Execute(ctx, "raw", "newpassword123")
		assert.Error(t, err)
		users.AssertNotCalled(t, "UpdatePassword", mock.Anything, mock.Anything, mock.Anything)
	})
}

// ============================================================
// Me
// ============================================================

func TestGetMeUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("ユーザーを返す", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Name: "Me"}, nil)

		got, err := usecase.NewGetMeUseCase(users).Execute(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, "Me", got.Name)
	})

	// 移行前と同じく DomainError ではないエラー（handler では 500）を返す。
	t.Run("存在しなければ DomainError ではないエラー", func(t *testing.T) {
		users := new(mockAuthUsers)
		users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		_, err := usecase.NewGetMeUseCase(users).Execute(ctx, 1)
		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})

	t.Run("取得エラーはそのまま返す", func(t *testing.T) {
		users := new(mockAuthUsers)
		dbErr := errors.New("db error")
		users.On("FindByID", mock.Anything, uint(1)).Return(nil, dbErr)

		_, err := usecase.NewGetMeUseCase(users).Execute(ctx, 1)
		assert.ErrorIs(t, err, dbErr)
	})
}
