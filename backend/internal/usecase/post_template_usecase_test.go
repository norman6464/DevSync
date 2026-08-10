package usecase_test

import (
	"context"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockPostTemplateRepo は usecase/repository.PostTemplateRepository のモック。
type mockPostTemplateRepo struct{ mock.Mock }

func (m *mockPostTemplateRepo) Create(ctx context.Context, template *model.PostTemplate) error {
	return m.Called(ctx, template).Error(0)
}

func (m *mockPostTemplateRepo) FindByID(ctx context.Context, id uint) (*model.PostTemplate, error) {
	args := m.Called(ctx, id)
	t, _ := args.Get(0).(*model.PostTemplate)
	return t, args.Error(1)
}

func (m *mockPostTemplateRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostTemplate, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	t, _ := args.Get(0).([]model.PostTemplate)
	return t, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostTemplateRepo) Update(ctx context.Context, template *model.PostTemplate) error {
	return m.Called(ctx, template).Error(0)
}

func (m *mockPostTemplateRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func TestCreatePostTemplateUseCase_Execute(t *testing.T) {
	t.Run("前後の空白を落として作成する", func(t *testing.T) {
		repo := new(mockPostTemplateRepo)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(tm *model.PostTemplate) bool {
			return tm.Name == "名前" && tm.ContentTemplate == "本文" && tm.TitleTemplate == "題"
		})).Return(nil)
		uc := usecase.NewCreatePostTemplateUseCase(repo)

		err := uc.Execute(context.Background(), &model.PostTemplate{
			Name: "  名前  ", ContentTemplate: "  本文  ", TitleTemplate: "  題  ",
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("タイトルテンプレートは空でも作成できる", func(t *testing.T) {
		repo := new(mockPostTemplateRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.PostTemplate")).Return(nil)
		uc := usecase.NewCreatePostTemplateUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), &model.PostTemplate{Name: "名前", ContentTemplate: "本文"}))
		repo.AssertExpectations(t)
	})

	for name, in := range map[string]*model.PostTemplate{
		"名前が空":                {Name: "", ContentTemplate: "本文"},
		"名前が 100 文字超":         {Name: strings.Repeat("a", 101), ContentTemplate: "本文"},
		"内容が空":                {Name: "名前", ContentTemplate: ""},
		"内容が 50000 文字超":       {Name: "名前", ContentTemplate: strings.Repeat("a", 50001)},
		"タイトルテンプレートが 200 文字超": {Name: "名前", ContentTemplate: "本文", TitleTemplate: strings.Repeat("a", 201)},
	} {
		t.Run(name+"は 400（作成しない）", func(t *testing.T) {
			repo := new(mockPostTemplateRepo)
			uc := usecase.NewCreatePostTemplateUseCase(repo)

			assert.Error(t, uc.Execute(context.Background(), in))
			repo.AssertNotCalled(t, "Create")
		})
	}
}

func TestGetPostTemplateUseCase_Execute(t *testing.T) {
	t.Run("所有者なら取得できる", func(t *testing.T) {
		repo := new(mockPostTemplateRepo)
		expected := &model.PostTemplate{UserID: 1, Name: "名前"}
		repo.On("FindByID", mock.Anything, uint(5)).Return(expected, nil)
		uc := usecase.NewGetPostTemplateUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	// 取得も所有権チェックの対象である点が他スライスと異なる
	t.Run("他人のテンプレートは 403", func(t *testing.T) {
		repo := new(mockPostTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostTemplate{UserID: 99}, nil)
		uc := usecase.NewGetPostTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUpdatePostTemplateUseCase_Execute(t *testing.T) {
	t.Run("トリム後に空でないフィールドだけ更新する", func(t *testing.T) {
		repo := new(mockPostTemplateRepo)
		existing := &model.PostTemplate{UserID: 1, Name: "元名", ContentTemplate: "元本文", TitleTemplate: "元題"}
		repo.On("FindByID", mock.Anything, uint(5)).Return(existing, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(tm *model.PostTemplate) bool {
			return tm.Name == "新名" && tm.ContentTemplate == "元本文" && tm.TitleTemplate == "元題"
		})).Return(nil)
		uc := usecase.NewUpdatePostTemplateUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1, &model.PostTemplate{Name: "  新名  ", ContentTemplate: "   "})

		assert.NoError(t, err)
		assert.Equal(t, "新名", got.Name)
		assert.Equal(t, "元本文", got.ContentTemplate)
		repo.AssertExpectations(t)
	})

	t.Run("他人のテンプレートは 403（更新しない）", func(t *testing.T) {
		repo := new(mockPostTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostTemplate{UserID: 99}, nil)
		uc := usecase.NewUpdatePostTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, &model.PostTemplate{Name: "新名"})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
		repo.AssertExpectations(t)
	})

	t.Run("名前が 100 文字超なら 400（更新しない）", func(t *testing.T) {
		repo := new(mockPostTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostTemplate{UserID: 1}, nil)
		uc := usecase.NewUpdatePostTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, &model.PostTemplate{Name: strings.Repeat("a", 101)})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
	})
}

func TestDeletePostTemplateUseCase_Execute(t *testing.T) {
	t.Run("他人のテンプレートは 403（削除しない）", func(t *testing.T) {
		repo := new(mockPostTemplateRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostTemplate{UserID: 99}, nil)
		uc := usecase.NewDeletePostTemplateUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 1))
		repo.AssertNotCalled(t, "Delete")
		repo.AssertExpectations(t)
	})
}
