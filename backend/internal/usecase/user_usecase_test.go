package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockUserRepo は usecase/repository.UserRepository のモック。
type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) FindAll(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	u, _ := args.Get(0).([]model.User)
	return u, args.Error(1)
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockUserRepo) Search(ctx context.Context, query string) ([]model.User, error) {
	args := m.Called(ctx, query)
	u, _ := args.Get(0).([]model.User)
	return u, args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

// ============================================================
// 一覧 / 取得
// ============================================================

func TestListUsersUseCase_Execute(t *testing.T) {
	t.Run("キーワードが空なら全件を返す", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindAll", mock.Anything).Return([]model.User{{ID: 1}, {ID: 2}}, nil)
		uc := usecase.NewListUsersUseCase(repo)

		users, err := uc.Execute(context.Background(), "")

		require.NoError(t, err)
		assert.Len(t, users, 2)
		repo.AssertNotCalled(t, "Search", mock.Anything, mock.Anything)
	})

	t.Run("キーワードがあれば検索する", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("Search", mock.Anything, "go").Return([]model.User{{ID: 1}}, nil)
		uc := usecase.NewListUsersUseCase(repo)

		users, err := uc.Execute(context.Background(), "go")

		require.NoError(t, err)
		assert.Len(t, users, 1)
		repo.AssertNotCalled(t, "FindAll", mock.Anything)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindAll", mock.Anything).Return([]model.User(nil), errors.New("db error"))
		uc := usecase.NewListUsersUseCase(repo)

		_, err := uc.Execute(context.Background(), "")

		assert.EqualError(t, err, "db error")
	})
}

func TestGetUserUseCase_Execute(t *testing.T) {
	t.Run("ユーザーを返す", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Name: "太郎"}, nil)
		uc := usecase.NewGetUserUseCase(repo)

		user, err := uc.Execute(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, "太郎", user.Name)
	})

	// 不在も DB 障害も 404 に潰す（移行前と同じ）。
	t.Run("不在と DB 障害はどちらも 404", func(t *testing.T) {
		for _, c := range []struct {
			name string
			user *model.User
			err  error
		}{
			{"不在", nil, nil},
			{"DB 障害", nil, errors.New("db error")},
		} {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockUserRepo)
				repo.On("FindByID", mock.Anything, uint(1)).Return(c.user, c.err)
				uc := usecase.NewGetUserUseCase(repo)

				_, err := uc.Execute(context.Background(), 1)

				domainErr := domain.GetDomainError(err)
				require.NotNil(t, domainErr)
				assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
				assert.Equal(t, "ユーザーが見つかりません", domainErr.Message)
			})
		}
	})
}

func TestGetUserByUsernameUseCase_Execute(t *testing.T) {
	t.Run("ユーザーを返す", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByUsername", mock.Anything, "taro").Return(&model.User{ID: 1, Username: "taro"}, nil)
		uc := usecase.NewGetUserByUsernameUseCase(repo)

		user, err := uc.Execute(context.Background(), "taro")

		require.NoError(t, err)
		assert.Equal(t, "taro", user.Username)
	})

	t.Run("不在は 404", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByUsername", mock.Anything, "nobody").Return(nil, nil)
		uc := usecase.NewGetUserByUsernameUseCase(repo)

		_, err := uc.Execute(context.Background(), "nobody")

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	})
}

// ============================================================
// プロフィール更新
// ============================================================

func TestUpdateUserProfileUseCase_Execute(t *testing.T) {
	existing := func() *model.User {
		return &model.User{
			ID: 1, Name: "旧名前", Bio: "旧自己紹介", AvatarURL: "https://old.example.com/a.png",
			SkillsLanguages: "Go", SkillsFrameworks: "Gin",
			AtCoderUsername: "old_atcoder", PaizaRank: "B",
		}
	}

	t.Run("名前は空でなければ更新する", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(existing(), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateUserProfileUseCase(repo)

		user, err := uc.Execute(context.Background(), 1, 1, &usecase.UpdateProfileInput{Name: "  新名前  "})

		require.NoError(t, err)
		assert.Equal(t, "新名前", user.Name)
	})

	t.Run("名前が空なら据え置く", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(existing(), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateUserProfileUseCase(repo)

		user, err := uc.Execute(context.Background(), 1, 1, &usecase.UpdateProfileInput{Name: "   "})

		require.NoError(t, err)
		assert.Equal(t, "旧名前", user.Name)
	})

	// 自己紹介とアバターは渡された値で必ず上書きする（空文字なら空になる）。
	t.Run("自己紹介とアバターは空文字でも上書きする", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(existing(), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateUserProfileUseCase(repo)

		user, err := uc.Execute(context.Background(), 1, 1, &usecase.UpdateProfileInput{})

		require.NoError(t, err)
		assert.Equal(t, "", user.Bio)
		assert.Equal(t, "", user.AvatarURL)
	})

	t.Run("ポインタの項目は nil なら据え置く", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(existing(), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateUserProfileUseCase(repo)

		user, err := uc.Execute(context.Background(), 1, 1, &usecase.UpdateProfileInput{Name: "名前"})

		require.NoError(t, err)
		assert.Equal(t, "Go", user.SkillsLanguages)
		assert.Equal(t, "Gin", user.SkillsFrameworks)
		assert.Equal(t, "old_atcoder", user.AtCoderUsername)
		assert.Equal(t, "B", user.PaizaRank)
		assert.False(t, user.OnboardingCompleted)
	})

	t.Run("ポインタの項目は指定されればトリムして反映する", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(existing(), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateUserProfileUseCase(repo)

		languages, frameworks := " Rust ", ""
		atcoder, rank := " new_atcoder ", " S "
		completed := true
		user, err := uc.Execute(context.Background(), 1, 1, &usecase.UpdateProfileInput{
			Name: "名前", SkillsLanguages: &languages, SkillsFrameworks: &frameworks,
			AtCoderUsername: &atcoder, PaizaRank: &rank, OnboardingCompleted: &completed,
		})

		require.NoError(t, err)
		assert.Equal(t, "Rust", user.SkillsLanguages)
		assert.Equal(t, "", user.SkillsFrameworks)
		assert.Equal(t, "new_atcoder", user.AtCoderUsername)
		assert.Equal(t, "S", user.PaizaRank)
		assert.True(t, user.OnboardingCompleted)
	})

	t.Run("本人以外は 403", func(t *testing.T) {
		repo := new(mockUserRepo)
		uc := usecase.NewUpdateUserProfileUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 2, &usecase.UpdateProfileInput{Name: "名前"})

		assert.ErrorIs(t, err, domain.ErrForbidden)
		repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("不在のユーザーは 404", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewUpdateUserProfileUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &usecase.UpdateProfileInput{Name: "名前"})

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		long := strings.Repeat("あ", 501)
		longURL := strings.Repeat("a", 2001)
		cases := []struct {
			name string
			in   *usecase.UpdateProfileInput
		}{
			{"名前が上限超過", &usecase.UpdateProfileInput{Name: strings.Repeat("あ", 101)}},
			{"自己紹介が上限超過", &usecase.UpdateProfileInput{Name: "名前", Bio: long}},
			{"言語スキルが上限超過", &usecase.UpdateProfileInput{Name: "名前", SkillsLanguages: &long}},
			{"フレームワークスキルが上限超過", &usecase.UpdateProfileInput{Name: "名前", SkillsFrameworks: &long}},
			{"アバター URL が上限超過", &usecase.UpdateProfileInput{Name: "名前", AvatarURL: longURL}},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockUserRepo)
				repo.On("FindByID", mock.Anything, uint(1)).Return(existing(), nil)
				uc := usecase.NewUpdateUserProfileUseCase(repo)

				_, err := uc.Execute(context.Background(), 1, 1, c.in)

				require.Error(t, err)
				assert.NotNil(t, domain.GetDomainError(err))
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			})
		}
	})

	t.Run("書き込みの DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(existing(), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewUpdateUserProfileUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &usecase.UpdateProfileInput{Name: "名前"})

		assert.EqualError(t, err, "db error")
	})
}

// ============================================================
// プロフィール完成度
// ============================================================

func TestGetProfileCompletenessUseCase_Execute(t *testing.T) {
	t.Run("完成度を返す", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{
			AvatarURL: "https://example.com/a.png", Bio: "自己紹介",
			GitHubConnected: true, SkillsLanguages: "Go",
		}, nil)
		uc := usecase.NewGetProfileCompletenessUseCase(repo)

		result, err := uc.Execute(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, 100, result.Percentage)
		assert.Empty(t, result.MissingFields)
	})

	// 他の取得と違い、こちらは 404 に変換せず 500 のままにする（移行前と同じ）。
	t.Run("不在と DB 障害は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		for _, c := range []struct {
			name string
			err  error
		}{{"不在", nil}, {"DB 障害", errors.New("db error")}} {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockUserRepo)
				repo.On("FindByID", mock.Anything, uint(1)).Return(nil, c.err)
				uc := usecase.NewGetProfileCompletenessUseCase(repo)

				_, err := uc.Execute(context.Background(), 1)

				require.Error(t, err)
				assert.Nil(t, domain.GetDomainError(err))
			})
		}
	})
}

func TestCalculateProfileCompleteness(t *testing.T) {
	cases := []struct {
		name    string
		user    *model.User
		percent int
		missing []string
	}{
		{"すべて未設定", &model.User{}, 0, []string{"avatar", "bio", "github", "skills"}},
		{"アバターのみ", &model.User{AvatarURL: "https://example.com/a.png"}, 25, []string{"bio", "github", "skills"}},
		{"半分", &model.User{Bio: "自己紹介", GitHubConnected: true}, 50, []string{"avatar", "skills"}},
		{"フレームワークだけでもスキルは充足", &model.User{SkillsFrameworks: "Gin"}, 25, []string{"avatar", "bio", "github"}},
		{"すべて設定済み", &model.User{
			AvatarURL: "https://example.com/a.png", Bio: "自己紹介",
			GitHubConnected: true, SkillsLanguages: "Go",
		}, 100, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := usecase.CalculateProfileCompleteness(c.user)
			assert.Equal(t, c.percent, result.Percentage)
			assert.Equal(t, c.missing, result.MissingFields)
		})
	}
}

// ============================================================
// メール配信設定
// ============================================================

func TestGetEmailPreferencesUseCase_Execute(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, EmailWeeklyReport: true, EmailLanguage: "ja"}, nil)
	uc := usecase.NewGetEmailPreferencesUseCase(repo)

	user, err := uc.Execute(context.Background(), 1)

	require.NoError(t, err)
	assert.True(t, user.EmailWeeklyReport)
	assert.Equal(t, "ja", user.EmailLanguage)
}

func TestUpdateEmailPreferencesUseCase_Execute(t *testing.T) {
	t.Run("受信可否と言語を更新する", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.User{ID: 1, EmailWeeklyReport: true, EmailLanguage: "ja"}, nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateEmailPreferencesUseCase(repo)

		weekly, language := false, "en"
		user, err := uc.Execute(context.Background(), 1, &weekly, &language)

		require.NoError(t, err)
		assert.False(t, user.EmailWeeklyReport)
		assert.Equal(t, "en", user.EmailLanguage)
	})

	t.Run("nil の項目は据え置く", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.User{ID: 1, EmailWeeklyReport: true, EmailLanguage: "ja"}, nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateEmailPreferencesUseCase(repo)

		user, err := uc.Execute(context.Background(), 1, nil, nil)

		require.NoError(t, err)
		assert.True(t, user.EmailWeeklyReport)
		assert.Equal(t, "ja", user.EmailLanguage)
	})

	t.Run("許可された言語はすべて受け付ける", func(t *testing.T) {
		for _, language := range []string{"ja", "en", "ko", "zh-CN", "zh-TW", "es", "fr", "de", "pt", "ru"} {
			repo := new(mockUserRepo)
			repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
			repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			uc := usecase.NewUpdateEmailPreferencesUseCase(repo)

			lang := language
			user, err := uc.Execute(context.Background(), 1, nil, &lang)

			require.NoError(t, err, language)
			assert.Equal(t, language, user.EmailLanguage)
		}
	})

	t.Run("許可外の言語は 400 で書き込まない", func(t *testing.T) {
		for _, language := range []string{"xx", "JA", "", "japanese"} {
			repo := new(mockUserRepo)
			repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
			uc := usecase.NewUpdateEmailPreferencesUseCase(repo)

			lang := language
			_, err := uc.Execute(context.Background(), 1, nil, &lang)

			domainErr := domain.GetDomainError(err)
			require.NotNil(t, domainErr, language)
			assert.Equal(t, "無効なメール言語設定です", domainErr.Message)
			repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		}
	})

	t.Run("不在のユーザーは 404", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewUpdateEmailPreferencesUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, nil, nil)

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	})
}
