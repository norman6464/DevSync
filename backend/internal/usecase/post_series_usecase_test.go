package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockPostSeriesRepo は usecase/repository.PostSeriesRepository のモック。
type mockPostSeriesRepo struct{ mock.Mock }

func (m *mockPostSeriesRepo) Create(ctx context.Context, series *model.PostSeries) error {
	return m.Called(ctx, series).Error(0)
}

func (m *mockPostSeriesRepo) FindByID(ctx context.Context, id uint) (*model.PostSeries, error) {
	args := m.Called(ctx, id)
	s, _ := args.Get(0).(*model.PostSeries)
	return s, args.Error(1)
}

func (m *mockPostSeriesRepo) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]model.PostSeries, error) {
	args := m.Called(ctx, userID, offset, limit)
	s, _ := args.Get(0).([]model.PostSeries)
	return s, args.Error(1)
}

func (m *mockPostSeriesRepo) CountByUser(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostSeriesRepo) Update(ctx context.Context, series *model.PostSeries) error {
	return m.Called(ctx, series).Error(0)
}

func (m *mockPostSeriesRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPostSeriesRepo) AddPost(ctx context.Context, item *model.PostSeriesItem) error {
	return m.Called(ctx, item).Error(0)
}

func (m *mockPostSeriesRepo) RemovePost(ctx context.Context, seriesID, postID uint) error {
	return m.Called(ctx, seriesID, postID).Error(0)
}

func (m *mockPostSeriesRepo) HasPost(ctx context.Context, seriesID, postID uint) (bool, error) {
	args := m.Called(ctx, seriesID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPostSeriesRepo) GetPostsBySeriesID(ctx context.Context, seriesID uint) ([]model.PostSeriesItem, error) {
	args := m.Called(ctx, seriesID)
	i, _ := args.Get(0).([]model.PostSeriesItem)
	return i, args.Error(1)
}

func TestCreatePostSeriesUseCase_Execute(t *testing.T) {
	t.Run("前後の空白を落として作成する", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(s *model.PostSeries) bool {
			return s.Title == "タイトル" && s.Description == "説明"
		})).Return(nil)
		uc := usecase.NewCreatePostSeriesUseCase(repo)

		err := uc.Execute(context.Background(), &model.PostSeries{Title: "  タイトル  ", Description: "  説明  "})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("説明は空でも作成できる", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.PostSeries")).Return(nil)
		uc := usecase.NewCreatePostSeriesUseCase(repo)

		err := uc.Execute(context.Background(), &model.PostSeries{Title: "タイトル"})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("タイトルが空なら 400（作成しない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		uc := usecase.NewCreatePostSeriesUseCase(repo)

		err := uc.Execute(context.Background(), &model.PostSeries{Title: ""})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("タイトルが 200 文字超なら 400（作成しない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		uc := usecase.NewCreatePostSeriesUseCase(repo)

		err := uc.Execute(context.Background(), &model.PostSeries{Title: strings.Repeat("a", 201)})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("説明が 1000 文字超なら 400（作成しない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		uc := usecase.NewCreatePostSeriesUseCase(repo)

		err := uc.Execute(context.Background(), &model.PostSeries{
			Title: "タイトル", Description: strings.Repeat("a", 1001),
		})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Create")
	})
}

func TestListPostSeriesUseCase_Execute(t *testing.T) {
	t.Run("ページ番号を offset に変換する", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		// page=3 / limit=20 → offset=40
		repo.On("FindByUserID", mock.Anything, uint(1), 40, 20).Return([]model.PostSeries{}, nil)
		uc := usecase.NewListPostSeriesUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 3, 20)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUpdatePostSeriesUseCase_Execute(t *testing.T) {
	t.Run("空でないフィールドだけ更新する", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		existing := &model.PostSeries{Title: "元タイトル", Description: "元説明", UserID: 1}
		repo.On("FindByID", mock.Anything, uint(5)).Return(existing, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(s *model.PostSeries) bool {
			return s.Title == "新タイトル" && s.Description == "元説明"
		})).Return(nil)
		uc := usecase.NewUpdatePostSeriesUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1, &model.PostSeries{Title: "  新タイトル  "})

		assert.NoError(t, err)
		assert.Equal(t, "新タイトル", got.Title)
		assert.Equal(t, "元説明", got.Description)
		repo.AssertExpectations(t)
	})

	t.Run("他人のシリーズは 403（更新しない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 99}, nil)
		uc := usecase.NewUpdatePostSeriesUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, &model.PostSeries{Title: "新"})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("タイトルが 200 文字超なら 400（更新しない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 1}, nil)
		uc := usecase.NewUpdatePostSeriesUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, &model.PostSeries{Title: strings.Repeat("a", 201)})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
	})
}

func TestDeletePostSeriesUseCase_Execute(t *testing.T) {
	t.Run("所有者なら削除する", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 1}, nil)
		repo.On("Delete", mock.Anything, uint(5)).Return(nil)
		uc := usecase.NewDeletePostSeriesUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 5, 1))
		repo.AssertExpectations(t)
	})

	t.Run("他人のシリーズは 403（削除しない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 99}, nil)
		uc := usecase.NewDeletePostSeriesUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 1))
		repo.AssertNotCalled(t, "Delete")
	})
}

func TestAddPostToSeriesUseCase_Execute(t *testing.T) {
	t.Run("未追加なら追加する", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 1}, nil)
		repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(false, nil)
		repo.On("AddPost", mock.Anything, mock.MatchedBy(func(i *model.PostSeriesItem) bool {
			return i.SeriesID == 5 && i.PostID == 10 && i.OrderIndex == 3
		})).Return(nil)
		uc := usecase.NewAddPostToSeriesUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 5, 10, 3, 1))
		repo.AssertExpectations(t)
	})

	t.Run("追加済みなら 400（追加しない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 1}, nil)
		repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(true, nil)
		uc := usecase.NewAddPostToSeriesUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 10, 0, 1))
		repo.AssertNotCalled(t, "AddPost")
	})

	t.Run("他人のシリーズは 403（存在確認もしない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 99}, nil)
		uc := usecase.NewAddPostToSeriesUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 10, 0, 1))
		repo.AssertNotCalled(t, "HasPost")
		repo.AssertNotCalled(t, "AddPost")
	})

	t.Run("存在確認のエラーを伝播する", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 1}, nil)
		repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(false, errors.New("db error"))
		uc := usecase.NewAddPostToSeriesUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 10, 0, 1))
		repo.AssertNotCalled(t, "AddPost")
	})
}

func TestRemovePostFromSeriesUseCase_Execute(t *testing.T) {
	t.Run("所有者なら取り除く", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 1}, nil)
		repo.On("RemovePost", mock.Anything, uint(5), uint(10)).Return(nil)
		uc := usecase.NewRemovePostFromSeriesUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 5, 10, 1))
		repo.AssertExpectations(t)
	})

	t.Run("他人のシリーズは 403（取り除かない）", func(t *testing.T) {
		repo := new(mockPostSeriesRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{UserID: 99}, nil)
		uc := usecase.NewRemovePostFromSeriesUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 10, 1))
		repo.AssertNotCalled(t, "RemovePost")
	})
}
