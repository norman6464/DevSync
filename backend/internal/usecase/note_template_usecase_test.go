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

// mockNoteTemplateRepo は usecase/repository.NoteTemplateRepository のモック。
type mockNoteTemplateRepo struct{ mock.Mock }

func (m *mockNoteTemplateRepo) Create(ctx context.Context, template *model.NoteTemplate) error {
	return m.Called(ctx, template).Error(0)
}

func (m *mockNoteTemplateRepo) Update(ctx context.Context, template *model.NoteTemplate) error {
	return m.Called(ctx, template).Error(0)
}

func (m *mockNoteTemplateRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockNoteTemplateRepo) FindByID(ctx context.Context, id uint) (*model.NoteTemplate, error) {
	args := m.Called(ctx, id)
	t, _ := args.Get(0).(*model.NoteTemplate)
	return t, args.Error(1)
}

func (m *mockNoteTemplateRepo) FindByUserID(ctx context.Context, userID uint) ([]model.NoteTemplate, error) {
	args := m.Called(ctx, userID)
	t, _ := args.Get(0).([]model.NoteTemplate)
	return t, args.Error(1)
}

func (m *mockNoteTemplateRepo) FindDefaultByUserID(ctx context.Context, userID uint) (*model.NoteTemplate, error) {
	args := m.Called(ctx, userID)
	t, _ := args.Get(0).(*model.NoteTemplate)
	return t, args.Error(1)
}

func (m *mockNoteTemplateRepo) ClearDefaultFlag(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *mockNoteTemplateRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// ownedNoteTemplate は指定ユーザーが所有するテンプレートを返すテスト用ヘルパー。
func ownedNoteTemplate(id, userID uint) *model.NoteTemplate {
	return &model.NoteTemplate{
		ID: id, UserID: userID,
		Name: "既存テンプレート", ContentTemplate: "## 本文", Description: "説明",
		DefaultTitle: "既定タイトル", DefaultTags: "go,test",
	}
}

// ============================================================
// Create
// ============================================================

func TestCreateNoteTemplateUseCase_Execute(t *testing.T) {
	t.Run("前後の空白を落として作成する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(tm *model.NoteTemplate) bool {
			return tm.Name == "名前" && tm.ContentTemplate == "本文" &&
				tm.Description == "説明" && tm.DefaultTitle == "題"
		})).Return(nil)
		uc := usecase.NewCreateNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.NoteTemplate{
			UserID: 1, Name: "  名前  ", ContentTemplate: " 本文 ",
			Description: " 説明 ", DefaultTitle: "  題 ",
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("説明とデフォルトタイトルは空でも作成できる", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).Return(nil)
		uc := usecase.NewCreateNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.NoteTemplate{
			UserID: 1, Name: "名前", ContentTemplate: "本文",
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("デフォルト指定つきは既存の指定を外してから作成する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		cleared := false
		repo.On("ClearDefaultFlag", mock.Anything, uint(7)).Run(func(mock.Arguments) {
			cleared = true
		}).Return(nil)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).
			Run(func(mock.Arguments) {
				assert.True(t, cleared, "ClearDefaultFlag より先に Create が呼ばれている")
			}).Return(nil)
		uc := usecase.NewCreateNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.NoteTemplate{
			UserID: 7, Name: "名前", ContentTemplate: "本文", IsDefault: true,
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("デフォルト指定の解除に失敗したら作成しない", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("ClearDefaultFlag", mock.Anything, uint(1)).Return(errors.New("db error"))
		uc := usecase.NewCreateNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.NoteTemplate{
			UserID: 1, Name: "名前", ContentTemplate: "本文", IsDefault: true,
		})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		cases := []struct {
			name string
			tmpl *model.NoteTemplate
		}{
			{"名前が空白のみ", &model.NoteTemplate{Name: "   ", ContentTemplate: "本文"}},
			{"名前が上限超過", &model.NoteTemplate{Name: strings.Repeat("あ", 101), ContentTemplate: "本文"}},
			{"本文テンプレートが空", &model.NoteTemplate{Name: "名前", ContentTemplate: ""}},
			{"本文テンプレートが上限超過", &model.NoteTemplate{Name: "名前", ContentTemplate: strings.Repeat("a", 50001)}},
			{"説明が上限超過", &model.NoteTemplate{Name: "名前", ContentTemplate: "本文", Description: strings.Repeat("あ", 501)}},
			{"デフォルトタイトルが上限超過", &model.NoteTemplate{Name: "名前", ContentTemplate: "本文", DefaultTitle: strings.Repeat("あ", 201)}},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockNoteTemplateRepo)
				uc := usecase.NewCreateNoteTemplateUseCase(repo)

				err := uc.Execute(context.Background(), c.tmpl)

				require.Error(t, err)
				assert.NotNil(t, domain.GetDomainError(err))
				repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			})
		}
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewCreateNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.NoteTemplate{Name: "名前", ContentTemplate: "本文"})

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})
}

// ============================================================
// Get / List / GetDefault / Count
// ============================================================

func TestGetNoteTemplateUseCase_Execute(t *testing.T) {
	t.Run("所有者は取得できる", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		uc := usecase.NewGetNoteTemplateUseCase(repo)

		tmpl, err := uc.Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, uint(1), tmpl.ID)
	})

	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		uc := usecase.NewGetNoteTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 9)

		assert.ErrorIs(t, err, domain.ErrForbidden)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetNoteTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetNoteTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestListNoteTemplatesUseCase_Execute(t *testing.T) {
	t.Run("ユーザーのテンプレートを返す", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByUserID", mock.Anything, uint(5)).
			Return([]model.NoteTemplate{*ownedNoteTemplate(2, 5), *ownedNoteTemplate(1, 5)}, nil)
		uc := usecase.NewListNoteTemplatesUseCase(repo)

		templates, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.Len(t, templates, 2)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByUserID", mock.Anything, uint(5)).
			Return([]model.NoteTemplate(nil), errors.New("db error"))
		uc := usecase.NewListNoteTemplatesUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestGetDefaultNoteTemplateUseCase_Execute(t *testing.T) {
	t.Run("デフォルトテンプレートを返す", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		tmpl := ownedNoteTemplate(1, 5)
		tmpl.IsDefault = true
		repo.On("FindDefaultByUserID", mock.Anything, uint(5)).Return(tmpl, nil)
		uc := usecase.NewGetDefaultNoteTemplateUseCase(repo)

		got, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.True(t, got.IsDefault)
	})

	t.Run("未設定は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindDefaultByUserID", mock.Anything, uint(5)).Return(nil, nil)
		uc := usecase.NewGetDefaultNoteTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindDefaultByUserID", mock.Anything, uint(5)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetDefaultNoteTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestCountNoteTemplatesUseCase_Execute(t *testing.T) {
	t.Run("総数を返す", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("CountByUserID", mock.Anything, uint(5)).Return(int64(3), nil)
		uc := usecase.NewCountNoteTemplatesUseCase(repo)

		count, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("CountByUserID", mock.Anything, uint(5)).Return(int64(0), errors.New("db error"))
		uc := usecase.NewCountNoteTemplatesUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

// ============================================================
// Update
// ============================================================

func TestUpdateNoteTemplateUseCase_Execute(t *testing.T) {
	t.Run("指定したフィールドだけを更新する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).Return(nil)
		uc := usecase.NewUpdateNoteTemplateUseCase(repo)

		tmpl, err := uc.Execute(context.Background(), usecase.UpdateNoteTemplateInput{
			ID: 1, UserID: 5, Name: "  新しい名前  ",
		})

		require.NoError(t, err)
		assert.Equal(t, "新しい名前", tmpl.Name)
		assert.Equal(t, "## 本文", tmpl.ContentTemplate)
		assert.Equal(t, "説明", tmpl.Description)
		assert.Equal(t, "既定タイトル", tmpl.DefaultTitle)
		assert.Equal(t, "go,test", tmpl.DefaultTags)
	})

	t.Run("全フィールドを更新できる", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		repo.On("ClearDefaultFlag", mock.Anything, uint(5)).Return(nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).Return(nil)
		uc := usecase.NewUpdateNoteTemplateUseCase(repo)

		isDefault := true
		tmpl, err := uc.Execute(context.Background(), usecase.UpdateNoteTemplateInput{
			ID: 1, UserID: 5,
			Name: "名前", Description: "新説明", DefaultTitle: "新題",
			ContentTemplate: "新本文", DefaultTags: "rust", IsDefault: &isDefault,
		})

		require.NoError(t, err)
		assert.Equal(t, "名前", tmpl.Name)
		assert.Equal(t, "新説明", tmpl.Description)
		assert.Equal(t, "新題", tmpl.DefaultTitle)
		assert.Equal(t, "新本文", tmpl.ContentTemplate)
		assert.Equal(t, "rust", tmpl.DefaultTags)
		assert.True(t, tmpl.IsDefault)
		repo.AssertExpectations(t)
	})

	t.Run("空白のみのフィールドは据え置く", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).Return(nil)
		uc := usecase.NewUpdateNoteTemplateUseCase(repo)

		tmpl, err := uc.Execute(context.Background(), usecase.UpdateNoteTemplateInput{
			ID: 1, UserID: 5, Description: "   ", DefaultTitle: "  ", DefaultTags: " ",
		})

		require.NoError(t, err)
		assert.Equal(t, "説明", tmpl.Description)
		assert.Equal(t, "既定タイトル", tmpl.DefaultTitle)
		assert.Equal(t, "go,test", tmpl.DefaultTags)
	})

	t.Run("デフォルト指定を外すときは他のテンプレートを触らない", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		existing := ownedNoteTemplate(1, 5)
		existing.IsDefault = true
		repo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).Return(nil)
		uc := usecase.NewUpdateNoteTemplateUseCase(repo)

		isDefault := false
		tmpl, err := uc.Execute(context.Background(), usecase.UpdateNoteTemplateInput{
			ID: 1, UserID: 5, IsDefault: &isDefault,
		})

		require.NoError(t, err)
		assert.False(t, tmpl.IsDefault)
		repo.AssertNotCalled(t, "ClearDefaultFlag", mock.Anything, mock.Anything)
	})

	t.Run("デフォルト指定の解除に失敗したら更新しない", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		repo.On("ClearDefaultFlag", mock.Anything, uint(5)).Return(errors.New("db error"))
		uc := usecase.NewUpdateNoteTemplateUseCase(repo)

		isDefault := true
		_, err := uc.Execute(context.Background(), usecase.UpdateNoteTemplateInput{
			ID: 1, UserID: 5, IsDefault: &isDefault,
		})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		uc := usecase.NewUpdateNoteTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteTemplateInput{ID: 1, UserID: 9, Name: "名前"})

		assert.ErrorIs(t, err, domain.ErrForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewUpdateNoteTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteTemplateInput{ID: 1, UserID: 5, Name: "名前"})

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})

	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		cases := []struct {
			name string
			in   usecase.UpdateNoteTemplateInput
		}{
			{"名前が上限超過", usecase.UpdateNoteTemplateInput{ID: 1, UserID: 5, Name: strings.Repeat("あ", 101)}},
			{"本文テンプレートが上限超過", usecase.UpdateNoteTemplateInput{ID: 1, UserID: 5, ContentTemplate: strings.Repeat("a", 50001)}},
			{"説明が上限超過", usecase.UpdateNoteTemplateInput{ID: 1, UserID: 5, Description: strings.Repeat("あ", 501)}},
			{"デフォルトタイトルが上限超過", usecase.UpdateNoteTemplateInput{ID: 1, UserID: 5, DefaultTitle: strings.Repeat("あ", 201)}},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockNoteTemplateRepo)
				repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
				uc := usecase.NewUpdateNoteTemplateUseCase(repo)

				_, err := uc.Execute(context.Background(), c.in)

				require.Error(t, err)
				assert.NotNil(t, domain.GetDomainError(err))
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			})
		}
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewUpdateNoteTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteTemplateInput{ID: 1, UserID: 5, Name: "名前"})

		assert.EqualError(t, err, "db error")
	})
}

// ============================================================
// Delete
// ============================================================

func TestDeleteNoteTemplateUseCase_Execute(t *testing.T) {
	t.Run("所有者は削除できる", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		repo.On("Delete", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewDeleteNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		uc := usecase.NewDeleteNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), 1, 9)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewDeleteNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), 1, 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockNoteTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		repo.On("Delete", mock.Anything, uint(1)).Return(errors.New("db error"))
		uc := usecase.NewDeleteNoteTemplateUseCase(repo)

		err := uc.Execute(context.Background(), 1, 5)

		assert.EqualError(t, err, "db error")
	})
}

// ============================================================
// CreateNoteFromTemplate
// ============================================================

func TestCreateNoteFromTemplateUseCase_Execute(t *testing.T) {
	newUseCase := func(templates *mockNoteTemplateRepo, notes *mockNoteRepo) *usecase.CreateNoteFromTemplateUseCase {
		return usecase.NewCreateNoteFromTemplateUseCase(templates, usecase.NewCreateNoteUseCase(notes))
	}

	t.Run("テンプレートの内容でノートを作る", func(t *testing.T) {
		templates := new(mockNoteTemplateRepo)
		notes := new(mockNoteRepo)
		templates.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		notes.On("Create", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)

		note, err := newUseCase(templates, notes).Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, uint(5), note.UserID)
		assert.Equal(t, "既定タイトル", note.Title)
		assert.Equal(t, "## 本文", note.Content)
		assert.Equal(t, "go,test", note.Tags)
		notes.AssertExpectations(t)
	})

	t.Run("デフォルトタイトルが空なら既定のノート名を使う", func(t *testing.T) {
		templates := new(mockNoteTemplateRepo)
		notes := new(mockNoteRepo)
		tmpl := ownedNoteTemplate(1, 5)
		tmpl.DefaultTitle = ""
		templates.On("FindByID", mock.Anything, uint(1)).Return(tmpl, nil)
		notes.On("Create", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)

		note, err := newUseCase(templates, notes).Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, "新しいノート", note.Title)
	})

	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		templates := new(mockNoteTemplateRepo)
		notes := new(mockNoteRepo)
		templates.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)

		_, err := newUseCase(templates, notes).Execute(context.Background(), 1, 9)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		notes.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		templates := new(mockNoteTemplateRepo)
		notes := new(mockNoteRepo)
		templates.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		_, err := newUseCase(templates, notes).Execute(context.Background(), 1, 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
		notes.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	// テンプレート由来の値もノート側の検証を通ることを保証する
	// （タグ列のサイズ上は実データで到達しないが、検証の委譲が外れたら気づけるようにしている）。
	t.Run("ノート作成の検証は CreateNoteUseCase に委ねる", func(t *testing.T) {
		templates := new(mockNoteTemplateRepo)
		notes := new(mockNoteRepo)
		tmpl := ownedNoteTemplate(1, 5)
		tmpl.DefaultTags = strings.Repeat("a", 501)
		templates.On("FindByID", mock.Anything, uint(1)).Return(tmpl, nil)

		_, err := newUseCase(templates, notes).Execute(context.Background(), 1, 5)

		require.Error(t, err)
		assert.NotNil(t, domain.GetDomainError(err))
		notes.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		templates := new(mockNoteTemplateRepo)
		notes := new(mockNoteRepo)
		templates.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteTemplate(1, 5), nil)
		notes.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

		_, err := newUseCase(templates, notes).Execute(context.Background(), 1, 5)

		assert.EqualError(t, err, "db error")
	})
}
