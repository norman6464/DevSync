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

// mockLearningLogTemplateRepo は usecase/repository.LearningLogTemplateRepository のモック。
type mockLearningLogTemplateRepo struct{ mock.Mock }

func (m *mockLearningLogTemplateRepo) Create(ctx context.Context, template *model.LearningLogTemplate) error {
	return m.Called(ctx, template).Error(0)
}
func (m *mockLearningLogTemplateRepo) Update(ctx context.Context, template *model.LearningLogTemplate) error {
	return m.Called(ctx, template).Error(0)
}
func (m *mockLearningLogTemplateRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockLearningLogTemplateRepo) FindByID(ctx context.Context, id uint) (*model.LearningLogTemplate, error) {
	args := m.Called(ctx, id)
	t, _ := args.Get(0).(*model.LearningLogTemplate)
	return t, args.Error(1)
}
func (m *mockLearningLogTemplateRepo) FindByUserID(ctx context.Context, userID uint) ([]model.LearningLogTemplate, error) {
	args := m.Called(ctx, userID)
	t, _ := args.Get(0).([]model.LearningLogTemplate)
	return t, args.Error(1)
}
func (m *mockLearningLogTemplateRepo) FindDefaultByUserID(ctx context.Context, userID uint) (*model.LearningLogTemplate, error) {
	args := m.Called(ctx, userID)
	t, _ := args.Get(0).(*model.LearningLogTemplate)
	return t, args.Error(1)
}
func (m *mockLearningLogTemplateRepo) ClearDefaultFlag(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *mockLearningLogTemplateRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// ownedLearningLogTemplate は指定ユーザーが所有するテンプレートを返すテスト用ヘルパー。
func ownedLearningLogTemplate(id, userID uint) *model.LearningLogTemplate {
	return &model.LearningLogTemplate{
		ID: id, UserID: userID,
		Name: "既存テンプレート", DefaultTitle: "既定タイトル", DefaultContent: "既定本文",
		DefaultCategory: model.LogCategoryCoding, DefaultDuration: 60,
	}
}

// ============================================================
// Create
// ============================================================

func TestCreateLearningLogTemplateUseCase_Execute(t *testing.T) {
	t.Run("前後の空白を落として作成する", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(tm *model.LearningLogTemplate) bool {
			return tm.Name == "名前" && tm.DefaultTitle == "題" && tm.DefaultContent == "本文"
		})).Return(nil)
		uc := usecase.NewCreateLearningLogTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.LearningLogTemplate{
			UserID: 1, Name: "  名前 ", DefaultTitle: " 題  ", DefaultContent: "  本文 ",
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("カテゴリと時間は未指定でも作成できる", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateLearningLogTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.LearningLogTemplate{UserID: 1, Name: "名前"})

		assert.NoError(t, err)
	})

	t.Run("デフォルト指定つきは既存の指定を外してから作成する", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		cleared := false
		repo.On("ClearDefaultFlag", mock.Anything, uint(7)).Run(func(mock.Arguments) { cleared = true }).Return(nil)
		repo.On("Create", mock.Anything, mock.Anything).Run(func(mock.Arguments) {
			assert.True(t, cleared, "ClearDefaultFlag より先に Create が呼ばれている")
		}).Return(nil)
		uc := usecase.NewCreateLearningLogTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.LearningLogTemplate{UserID: 7, Name: "名前", IsDefault: true})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("デフォルト指定の解除に失敗したら作成しない", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("ClearDefaultFlag", mock.Anything, uint(1)).Return(errors.New("db error"))
		uc := usecase.NewCreateLearningLogTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.LearningLogTemplate{UserID: 1, Name: "名前", IsDefault: true})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		cases := []struct {
			name string
			tmpl *model.LearningLogTemplate
			msg  string
		}{
			{"名前が空白のみ", &model.LearningLogTemplate{Name: "   "}, "テンプレート名"},
			{"名前が上限超過", &model.LearningLogTemplate{Name: strings.Repeat("あ", 101)}, "テンプレート名"},
			{"デフォルトタイトルが上限超過", &model.LearningLogTemplate{Name: "名前", DefaultTitle: strings.Repeat("あ", 201)}, "デフォルトタイトル"},
			{"デフォルト本文が上限超過", &model.LearningLogTemplate{Name: "名前", DefaultContent: strings.Repeat("a", 50001)}, "デフォルト本文"},
			{"無効なカテゴリ", &model.LearningLogTemplate{Name: "名前", DefaultCategory: "unknown"}, "無効なカテゴリです"},
			{"時間が負", &model.LearningLogTemplate{Name: "名前", DefaultDuration: -1}, "デフォルト時間"},
			{"時間が上限超過", &model.LearningLogTemplate{Name: "名前", DefaultDuration: 1441}, "デフォルト時間"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockLearningLogTemplateRepo)
				uc := usecase.NewCreateLearningLogTemplateUseCase(repo)

				err := uc.Execute(context.Background(), c.tmpl)

				require.Error(t, err)
				domainErr := domain.GetDomainError(err)
				require.NotNil(t, domainErr)
				assert.Contains(t, domainErr.Message, c.msg)
				repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			})
		}
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewCreateLearningLogTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.LearningLogTemplate{UserID: 1, Name: "名前"})

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})
}

// ============================================================
// 取得 / 一覧 / デフォルト / 件数
// ============================================================

func TestGetLearningLogTemplateUseCase_Execute(t *testing.T) {
	t.Run("所有者は取得できる", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		uc := usecase.NewGetLearningLogTemplateUseCase(repo)

		tmpl, err := uc.Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, uint(1), tmpl.ID)
	})

	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		uc := usecase.NewGetLearningLogTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 9)

		assert.ErrorIs(t, err, domain.ErrForbidden)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetLearningLogTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetLearningLogTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestListLearningLogTemplatesUseCase_Execute(t *testing.T) {
	repo := new(mockLearningLogTemplateRepo)
	repo.On("FindByUserID", mock.Anything, uint(5)).
		Return([]model.LearningLogTemplate{*ownedLearningLogTemplate(2, 5), *ownedLearningLogTemplate(1, 5)}, nil)
	uc := usecase.NewListLearningLogTemplatesUseCase(repo)

	templates, err := uc.Execute(context.Background(), 5)

	require.NoError(t, err)
	assert.Len(t, templates, 2)
}

func TestGetDefaultLearningLogTemplateUseCase_Execute(t *testing.T) {
	t.Run("デフォルトテンプレートを返す", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		tmpl := ownedLearningLogTemplate(1, 5)
		tmpl.IsDefault = true
		repo.On("FindDefaultByUserID", mock.Anything, uint(5)).Return(tmpl, nil)
		uc := usecase.NewGetDefaultLearningLogTemplateUseCase(repo)

		got, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.True(t, got.IsDefault)
	})

	t.Run("未設定は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindDefaultByUserID", mock.Anything, uint(5)).Return(nil, nil)
		uc := usecase.NewGetDefaultLearningLogTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindDefaultByUserID", mock.Anything, uint(5)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetDefaultLearningLogTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestCountLearningLogTemplatesUseCase_Execute(t *testing.T) {
	repo := new(mockLearningLogTemplateRepo)
	repo.On("CountByUserID", mock.Anything, uint(5)).Return(int64(4), nil)
	uc := usecase.NewCountLearningLogTemplatesUseCase(repo)

	count, err := uc.Execute(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, int64(4), count)
}

// ============================================================
// Update
// ============================================================

func TestUpdateLearningLogTemplateUseCase_Execute(t *testing.T) {
	t.Run("指定したフィールドだけを更新する", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateLearningLogTemplateUseCase(repo)

		tmpl, err := uc.Execute(context.Background(), usecase.UpdateLearningLogTemplateInput{
			ID: 1, UserID: 5, Name: "  新しい名前  ",
		})

		require.NoError(t, err)
		assert.Equal(t, "新しい名前", tmpl.Name)
		assert.Equal(t, "既定タイトル", tmpl.DefaultTitle)
		assert.Equal(t, "既定本文", tmpl.DefaultContent)
		assert.Equal(t, model.LogCategoryCoding, tmpl.DefaultCategory)
		assert.Equal(t, 60, tmpl.DefaultDuration)
	})

	t.Run("全フィールドを更新できる", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		repo.On("ClearDefaultFlag", mock.Anything, uint(5)).Return(nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateLearningLogTemplateUseCase(repo)

		duration := 30
		isDefault := true
		tmpl, err := uc.Execute(context.Background(), usecase.UpdateLearningLogTemplateInput{
			ID: 1, UserID: 5, Name: "名前", DefaultTitle: "題", DefaultContent: "本文",
			DefaultCategory: model.LogCategoryReading, DefaultDuration: &duration, IsDefault: &isDefault,
		})

		require.NoError(t, err)
		assert.Equal(t, "題", tmpl.DefaultTitle)
		assert.Equal(t, model.LogCategoryReading, tmpl.DefaultCategory)
		assert.Equal(t, 30, tmpl.DefaultDuration)
		assert.True(t, tmpl.IsDefault)
		repo.AssertExpectations(t)
	})

	t.Run("空白のみのフィールドは据え置く", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateLearningLogTemplateUseCase(repo)

		tmpl, err := uc.Execute(context.Background(), usecase.UpdateLearningLogTemplateInput{
			ID: 1, UserID: 5, Name: "  ", DefaultTitle: " ", DefaultContent: "   ",
		})

		require.NoError(t, err)
		assert.Equal(t, "既存テンプレート", tmpl.Name)
		assert.Equal(t, "既定タイトル", tmpl.DefaultTitle)
		assert.Equal(t, "既定本文", tmpl.DefaultContent)
	})

	t.Run("時間 0 も明示指定なら反映する", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateLearningLogTemplateUseCase(repo)

		duration := 0
		tmpl, err := uc.Execute(context.Background(), usecase.UpdateLearningLogTemplateInput{
			ID: 1, UserID: 5, DefaultDuration: &duration,
		})

		require.NoError(t, err)
		assert.Equal(t, 0, tmpl.DefaultDuration)
	})

	t.Run("デフォルト指定を外すときは他のテンプレートを触らない", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		existing := ownedLearningLogTemplate(1, 5)
		existing.IsDefault = true
		repo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateLearningLogTemplateUseCase(repo)

		isDefault := false
		tmpl, err := uc.Execute(context.Background(), usecase.UpdateLearningLogTemplateInput{
			ID: 1, UserID: 5, IsDefault: &isDefault,
		})

		require.NoError(t, err)
		assert.False(t, tmpl.IsDefault)
		repo.AssertNotCalled(t, "ClearDefaultFlag", mock.Anything, mock.Anything)
	})

	t.Run("デフォルト指定の解除に失敗したら更新しない", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		repo.On("ClearDefaultFlag", mock.Anything, uint(5)).Return(errors.New("db error"))
		uc := usecase.NewUpdateLearningLogTemplateUseCase(repo)

		isDefault := true
		_, err := uc.Execute(context.Background(), usecase.UpdateLearningLogTemplateInput{
			ID: 1, UserID: 5, IsDefault: &isDefault,
		})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		uc := usecase.NewUpdateLearningLogTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateLearningLogTemplateInput{ID: 1, UserID: 9, Name: "名前"})

		assert.ErrorIs(t, err, domain.ErrForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		overLimit := 1441
		cases := []struct {
			name string
			in   usecase.UpdateLearningLogTemplateInput
		}{
			{"名前が上限超過", usecase.UpdateLearningLogTemplateInput{ID: 1, UserID: 5, Name: strings.Repeat("あ", 101)}},
			{"デフォルトタイトルが上限超過", usecase.UpdateLearningLogTemplateInput{ID: 1, UserID: 5, DefaultTitle: strings.Repeat("あ", 201)}},
			{"デフォルト本文が上限超過", usecase.UpdateLearningLogTemplateInput{ID: 1, UserID: 5, DefaultContent: strings.Repeat("a", 50001)}},
			{"無効なカテゴリ", usecase.UpdateLearningLogTemplateInput{ID: 1, UserID: 5, DefaultCategory: "unknown"}},
			{"時間が上限超過", usecase.UpdateLearningLogTemplateInput{ID: 1, UserID: 5, DefaultDuration: &overLimit}},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockLearningLogTemplateRepo)
				repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
				uc := usecase.NewUpdateLearningLogTemplateUseCase(repo)

				_, err := uc.Execute(context.Background(), c.in)

				require.Error(t, err)
				assert.NotNil(t, domain.GetDomainError(err))
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			})
		}
	})
}

// ============================================================
// Delete
// ============================================================

func TestDeleteLearningLogTemplateUseCase_Execute(t *testing.T) {
	t.Run("所有者は削除できる", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		repo.On("Delete", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewDeleteLearningLogTemplateUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
		repo.AssertExpectations(t)
	})

	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		uc := usecase.NewDeleteLearningLogTemplateUseCase(repo)

		assert.ErrorIs(t, uc.Execute(context.Background(), 1, 9), domain.ErrForbidden)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockLearningLogTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewDeleteLearningLogTemplateUseCase(repo)

		err := uc.Execute(context.Background(), 1, 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})
}

// ============================================================
// テンプレートからの学習ログ作成
// ============================================================

func TestCreateLearningLogFromTemplateUseCase_Execute(t *testing.T) {
	newUseCase := func(templates *mockLearningLogTemplateRepo, logs *mockLearningLogRepo) *usecase.CreateLearningLogFromTemplateUseCase {
		return usecase.NewCreateLearningLogFromTemplateUseCase(templates, usecase.NewCreateLearningLogUseCase(logs, nil))
	}

	t.Run("テンプレートの内容で学習ログを作る", func(t *testing.T) {
		templates := new(mockLearningLogTemplateRepo)
		logs := new(mockLearningLogRepo)
		templates.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		logs.On("Create", mock.Anything, mock.AnythingOfType("*model.LearningLog")).Return(nil)

		log, err := newUseCase(templates, logs).Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, uint(5), log.UserID)
		assert.Equal(t, "既定タイトル", log.Title)
		assert.Equal(t, "既定本文", log.Content)
		assert.Equal(t, model.LogCategoryCoding, log.Category)
		assert.Equal(t, 60, log.Duration)
		assert.Equal(t, model.LogSourceManual, log.Source)
	})

	t.Run("デフォルトタイトルが空ならテンプレート名を使う", func(t *testing.T) {
		templates := new(mockLearningLogTemplateRepo)
		logs := new(mockLearningLogRepo)
		tmpl := ownedLearningLogTemplate(1, 5)
		tmpl.DefaultTitle = ""
		templates.On("FindByID", mock.Anything, uint(1)).Return(tmpl, nil)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)

		log, err := newUseCase(templates, logs).Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, "既存テンプレート", log.Title)
	})

	t.Run("カテゴリが空ならその他を使う", func(t *testing.T) {
		templates := new(mockLearningLogTemplateRepo)
		logs := new(mockLearningLogRepo)
		tmpl := ownedLearningLogTemplate(1, 5)
		tmpl.DefaultCategory = ""
		templates.On("FindByID", mock.Anything, uint(1)).Return(tmpl, nil)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)

		log, err := newUseCase(templates, logs).Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, model.LogCategoryOther, log.Category)
	})

	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		templates := new(mockLearningLogTemplateRepo)
		logs := new(mockLearningLogRepo)
		templates.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)

		_, err := newUseCase(templates, logs).Execute(context.Background(), 1, 9)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		templates := new(mockLearningLogTemplateRepo)
		logs := new(mockLearningLogRepo)
		templates.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		_, err := newUseCase(templates, logs).Execute(context.Background(), 1, 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
		logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	// テンプレート由来の値も学習ログ側の検証を通ることを保証する
	// （テンプレート作成時の上限が学習ログ側より緩いため、実データでも到達しうる）。
	t.Run("学習ログの検証は CreateLearningLogUseCase に委ねる", func(t *testing.T) {
		templates := new(mockLearningLogTemplateRepo)
		logs := new(mockLearningLogRepo)
		tmpl := ownedLearningLogTemplate(1, 5)
		tmpl.DefaultContent = strings.Repeat("a", 10001)
		templates.On("FindByID", mock.Anything, uint(1)).Return(tmpl, nil)

		_, err := newUseCase(templates, logs).Execute(context.Background(), 1, 5)

		require.Error(t, err)
		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Contains(t, domainErr.Message, "内容")
		logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("書き込み失敗は学習ログ側のメッセージで返る", func(t *testing.T) {
		templates := new(mockLearningLogTemplateRepo)
		logs := new(mockLearningLogRepo)
		templates.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLogTemplate(1, 5), nil)
		logs.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

		_, err := newUseCase(templates, logs).Execute(context.Background(), 1, 5)

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "学習ログの作成に失敗しました", domainErr.Message)
	})
}
