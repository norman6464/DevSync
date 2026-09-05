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

// mockLearningResourceRepo は usecase/repository.LearningResourceRepository のモック。
type mockLearningResourceRepo struct{ mock.Mock }

func (m *mockLearningResourceRepo) Create(ctx context.Context, r *model.LearningResource) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockLearningResourceRepo) Update(ctx context.Context, r *model.LearningResource) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockLearningResourceRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockLearningResourceRepo) FindByID(ctx context.Context, id uint) (*model.LearningResource, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.LearningResource)
	return r, args.Error(1)
}

func (m *mockLearningResourceRepo) FindByUserID(ctx context.Context, userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, userID, includePrivate, limit, offset)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockLearningResourceRepo) FindPublic(ctx context.Context, limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, limit, offset, category, difficulty)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockLearningResourceRepo) FindByDifficulty(ctx context.Context, difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, difficulty, limit, offset)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockLearningResourceRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockLearningResourceRepo) FindSavedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockLearningResourceRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockLearningResourceRepo) Like(ctx context.Context, userID, resourceID uint) error {
	return m.Called(ctx, userID, resourceID).Error(0)
}

func (m *mockLearningResourceRepo) Unlike(ctx context.Context, userID, resourceID uint) error {
	return m.Called(ctx, userID, resourceID).Error(0)
}

func (m *mockLearningResourceRepo) HasLiked(ctx context.Context, userID, resourceID uint) (bool, error) {
	args := m.Called(ctx, userID, resourceID)
	return args.Bool(0), args.Error(1)
}

func (m *mockLearningResourceRepo) Save(ctx context.Context, userID, resourceID uint) error {
	return m.Called(ctx, userID, resourceID).Error(0)
}

func (m *mockLearningResourceRepo) Unsave(ctx context.Context, userID, resourceID uint) error {
	return m.Called(ctx, userID, resourceID).Error(0)
}

func (m *mockLearningResourceRepo) HasSaved(ctx context.Context, userID, resourceID uint) (bool, error) {
	args := m.Called(ctx, userID, resourceID)
	return args.Bool(0), args.Error(1)
}

// assertResourceCode は err が期待の HTTP ステータスに対応する DomainError であることを検証する。
func assertResourceCode(t *testing.T, err error, want domain.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	domainErr := domain.GetDomainError(err)
	require.NotNil(t, domainErr, "DomainError であること")
	assert.Equal(t, want, domainErr.Code)
}

// validResource は検証を通る学習リソースを返す。
func validResource() *model.LearningResource {
	return &model.LearningResource{
		UserID: 1, Title: "Go 入門", Description: "説明文です", URL: "https://example.com",
		Category: model.ResourceCategoryArticle, Difficulty: model.ResourceDifficultyBeginner,
	}
}

func TestCreateLearningResourceUseCase_Execute(t *testing.T) {
	t.Run("検証を通れば作成する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.LearningResource")).Return(nil)
		uc := usecase.NewCreateLearningResourceUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), validResource()))
		repo.AssertExpectations(t)
	})

	t.Run("タイトルが空なら 400 で作成しない", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		uc := usecase.NewCreateLearningResourceUseCase(repo)

		r := validResource()
		r.Title = ""
		assertResourceCode(t, uc.Execute(context.Background(), r), domain.ErrCodeValidation)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("タグが 1001 文字なら 400", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		uc := usecase.NewCreateLearningResourceUseCase(repo)

		r := validResource()
		r.Tags = strings.Repeat("a", 1001)
		assertResourceCode(t, uc.Execute(context.Background(), r), domain.ErrCodeValidation)
	})
}

func TestGetLearningResourceUseCase_Execute(t *testing.T) {
	t.Run("公開リソースは他人でも取得できる", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.LearningResource{ID: 1, UserID: 99, IsPublic: true}, nil)
		uc := usecase.NewGetLearningResourceUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.Equal(t, uint(1), got.ID)
	})

	t.Run("自分の非公開リソースは取得できる", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.LearningResource{ID: 1, UserID: 1, IsPublic: false}, nil)
		uc := usecase.NewGetLearningResourceUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		assert.NoError(t, err)
	})

	t.Run("他人の非公開リソースは 403", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.LearningResource{ID: 1, UserID: 99, IsPublic: false}, nil)
		uc := usecase.NewGetLearningResourceUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		assertResourceCode(t, err, domain.ErrCodeForbidden)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetLearningResourceUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetLearningResourceUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})
}

func TestListLearningResourcesByUserUseCase_Execute(t *testing.T) {
	t.Run("自分の一覧は非公開も含める", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByUserID", mock.Anything, uint(1), true, 20, 0).
			Return([]model.LearningResource{{ID: 1}}, int64(1), nil)
		uc := usecase.NewListLearningResourcesByUserUseCase(repo)

		_, _, err := uc.Execute(context.Background(), 1, 1, 20, 0)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("他人の一覧は公開のみ", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByUserID", mock.Anything, uint(99), false, 20, 0).
			Return([]model.LearningResource{}, int64(0), nil)
		uc := usecase.NewListLearningResourcesByUserUseCase(repo)

		_, _, err := uc.Execute(context.Background(), 99, 1, 20, 0)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestListLearningResourcesByDifficultyUseCase_Execute(t *testing.T) {
	t.Run("未知の難易度は 400 で repo を引かない", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		uc := usecase.NewListLearningResourcesByDifficultyUseCase(repo)

		_, _, err := uc.Execute(context.Background(), "expert", 20, 0)
		assertResourceCode(t, err, domain.ErrCodeBadRequest)
		assert.Equal(t, "無効な難易度です", domain.GetDomainError(err).Message)
		repo.AssertNotCalled(t, "FindByDifficulty", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("有効な難易度は repo に委譲する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByDifficulty", mock.Anything, "beginner", 20, 0).
			Return([]model.LearningResource{{ID: 1}}, int64(1), nil)
		uc := usecase.NewListLearningResourcesByDifficultyUseCase(repo)

		got, total, err := uc.Execute(context.Background(), "beginner", 20, 0)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
	})
}

func TestUpdateLearningResourceUseCase_Execute(t *testing.T) {
	owned := func() *model.LearningResource {
		return &model.LearningResource{
			ID: 1, UserID: 1, Title: "旧題", Description: "旧説明", URL: "https://old.example.com",
			Category: model.ResourceCategoryArticle, Difficulty: model.ResourceDifficultyBeginner,
			Tags: "旧タグ",
		}
	}

	t.Run("指定した項目だけ更新する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.LearningResource")).Return(nil)
		uc := usecase.NewUpdateLearningResourceUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1, &model.LearningResource{Title: "新題"})
		require.NoError(t, err)
		assert.Equal(t, "新題", got.Title)
		assert.Equal(t, "旧説明", got.Description, "未指定の項目は据え置き")
		assert.Equal(t, "旧タグ", got.Tags)
	})

	t.Run("空白のみの入力は変更なしとして扱う", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.LearningResource")).Return(nil)
		uc := usecase.NewUpdateLearningResourceUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1, &model.LearningResource{
			Title: "   ", Description: "   ", Tags: "   ",
		})
		require.NoError(t, err)
		assert.Equal(t, "旧題", got.Title)
		assert.Equal(t, "旧説明", got.Description)
		assert.Equal(t, "旧タグ", got.Tags)
	})

	t.Run("所有者でなければ 403", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.LearningResource{ID: 1, UserID: 99}, nil)
		uc := usecase.NewUpdateLearningResourceUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &model.LearningResource{Title: "新題"})
		assertResourceCode(t, err, domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("タグが 1001 文字なら 400", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned(), nil)
		uc := usecase.NewUpdateLearningResourceUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &model.LearningResource{
			Tags: strings.Repeat("a", 1001),
		})
		assertResourceCode(t, err, domain.ErrCodeValidation)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("画像URLが 2001 文字なら 400", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned(), nil)
		uc := usecase.NewUpdateLearningResourceUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &model.LearningResource{
			ImageURL: strings.Repeat("a", 2001),
		})
		assertResourceCode(t, err, domain.ErrCodeValidation)
	})
}

func TestLearningResourceLikeAndSave(t *testing.T) {
	others := &model.LearningResource{ID: 1, UserID: 99}
	mine := &model.LearningResource{ID: 1, UserID: 1}

	t.Run("他人のリソースにはいいねできる", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(others, nil)
		repo.On("Like", mock.Anything, uint(1), uint(1)).Return(nil)
		uc := usecase.NewLikeLearningResourceUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 1))
		repo.AssertExpectations(t)
	})

	t.Run("自分のリソースへのいいねは 403", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(mine, nil)
		uc := usecase.NewLikeLearningResourceUseCase(repo)

		assertResourceCode(t, uc.Execute(context.Background(), 1, 1), domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Like", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("不在のリソースへのいいねは 404", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewLikeLearningResourceUseCase(repo)

		assertResourceCode(t, uc.Execute(context.Background(), 1, 1), domain.ErrCodeNotFound)
	})

	t.Run("取得の DB 障害も 404 に潰れる", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewLikeLearningResourceUseCase(repo)

		assertResourceCode(t, uc.Execute(context.Background(), 1, 1), domain.ErrCodeNotFound)
	})

	t.Run("いいね取消も自分のリソースなら 403", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(mine, nil)
		uc := usecase.NewUnlikeLearningResourceUseCase(repo)

		assertResourceCode(t, uc.Execute(context.Background(), 1, 1), domain.ErrCodeForbidden)
	})

	t.Run("自分のリソースの保存は 403", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(mine, nil)
		uc := usecase.NewSaveLearningResourceUseCase(repo)

		assertResourceCode(t, uc.Execute(context.Background(), 1, 1), domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("他人のリソースは保存できる", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(others, nil)
		repo.On("Save", mock.Anything, uint(1), uint(1)).Return(nil)
		uc := usecase.NewSaveLearningResourceUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 1))
		repo.AssertExpectations(t)
	})

	t.Run("保存取消も自分のリソースなら 403", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(mine, nil)
		uc := usecase.NewUnsaveLearningResourceUseCase(repo)

		assertResourceCode(t, uc.Execute(context.Background(), 1, 1), domain.ErrCodeForbidden)
	})
}

func TestLearningResourcePassThroughUseCases(t *testing.T) {
	ctx := context.Background()

	t.Run("公開一覧は絞り込み条件をそのまま渡す", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindPublic", mock.Anything, 20, 0, "article", "beginner").
			Return([]model.LearningResource{{ID: 1}}, int64(1), nil)
		got, total, err := usecase.NewListPublicLearningResourcesUseCase(repo).
			Execute(ctx, 20, 0, "article", "beginner")
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
		repo.AssertExpectations(t)
	})

	t.Run("検索は repo に委譲する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("Search", mock.Anything, "go", 20, 0).Return([]model.LearningResource{}, int64(0), nil)
		_, total, err := usecase.NewSearchLearningResourcesUseCase(repo).Execute(ctx, "go", 20, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("保存一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindSavedByUserID", mock.Anything, uint(1), 20, 0).
			Return([]model.LearningResource{{ID: 1}}, int64(1), nil)
		got, _, err := usecase.NewListSavedLearningResourcesUseCase(repo).Execute(ctx, 1, 20, 0)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("いいね済み判定は repo に委譲する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("HasLiked", mock.Anything, uint(1), uint(2)).Return(true, nil)
		got, err := usecase.NewHasLikedLearningResourceUseCase(repo).Execute(ctx, 1, 2)
		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("保存済み判定は repo に委譲する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("HasSaved", mock.Anything, uint(1), uint(2)).Return(false, nil)
		got, err := usecase.NewHasSavedLearningResourceUseCase(repo).Execute(ctx, 1, 2)
		assert.NoError(t, err)
		assert.False(t, got)
	})

	t.Run("件数は repo に委譲する", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(4), nil)
		got, err := usecase.NewCountLearningResourcesUseCase(repo).Execute(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), got)
	})

	t.Run("公開切替は所有者のみ", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.LearningResource{ID: 1, UserID: 99}, nil)
		_, err := usecase.NewUpdateLearningResourceVisibilityUseCase(repo).Execute(ctx, 1, 1, true)
		assertResourceCode(t, err, domain.ErrCodeForbidden)
	})

	t.Run("削除は所有者のみ", func(t *testing.T) {
		repo := new(mockLearningResourceRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.LearningResource{ID: 1, UserID: 1}, nil)
		repo.On("Delete", mock.Anything, uint(1)).Return(nil)
		assert.NoError(t, usecase.NewDeleteLearningResourceUseCase(repo).Execute(ctx, 1, 1))
		repo.AssertExpectations(t)
	})
}
