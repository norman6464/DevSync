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

// mockPostCollectionRepo は usecase/repository.PostCollectionRepository のモック。
type mockPostCollectionRepo struct{ mock.Mock }

func (m *mockPostCollectionRepo) Create(ctx context.Context, collection *model.PostCollection) error {
	return m.Called(ctx, collection).Error(0)
}

func (m *mockPostCollectionRepo) FindByID(ctx context.Context, id uint) (*model.PostCollection, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.PostCollection)
	return c, args.Error(1)
}

func (m *mockPostCollectionRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	c, _ := args.Get(0).([]model.PostCollection)
	return c, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostCollectionRepo) FindPublicByUserID(ctx context.Context, userID uint) ([]model.PostCollection, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]model.PostCollection)
	return c, args.Error(1)
}

func (m *mockPostCollectionRepo) Update(ctx context.Context, collection *model.PostCollection) error {
	return m.Called(ctx, collection).Error(0)
}

func (m *mockPostCollectionRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPostCollectionRepo) AddPost(ctx context.Context, item *model.PostCollectionItem) error {
	return m.Called(ctx, item).Error(0)
}

func (m *mockPostCollectionRepo) RemovePost(ctx context.Context, collectionID, postID uint) error {
	return m.Called(ctx, collectionID, postID).Error(0)
}

func (m *mockPostCollectionRepo) HasPost(ctx context.Context, collectionID, postID uint) (bool, error) {
	args := m.Called(ctx, collectionID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPostCollectionRepo) GetPostsByCollectionID(ctx context.Context, collectionID uint) ([]model.PostCollectionItem, error) {
	args := m.Called(ctx, collectionID)
	i, _ := args.Get(0).([]model.PostCollectionItem)
	return i, args.Error(1)
}

func (m *mockPostCollectionRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreatePostCollectionUseCase_Execute(t *testing.T) {
	t.Run("前後の空白を落として作成する", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(c *model.PostCollection) bool {
			return c.Title == "タイトル" && c.Description == "説明"
		})).Return(nil)
		uc := usecase.NewCreatePostCollectionUseCase(repo)

		got, err := uc.Execute(context.Background(), &model.PostCollection{Title: "  タイトル  ", Description: "  説明  "})

		assert.NoError(t, err)
		assert.Equal(t, "タイトル", got.Title)
		repo.AssertExpectations(t)
	})

	t.Run("タイトルが空なら 400（作成しない）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		uc := usecase.NewCreatePostCollectionUseCase(repo)

		_, err := uc.Execute(context.Background(), &model.PostCollection{Title: ""})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("説明が 1000 文字超なら 400（作成しない）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		uc := usecase.NewCreatePostCollectionUseCase(repo)

		_, err := uc.Execute(context.Background(), &model.PostCollection{
			Title: "タイトル", Description: strings.Repeat("a", 1001),
		})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Create")
	})
}

func TestListPostCollectionsForViewerUseCase_Execute(t *testing.T) {
	t.Run("自分のコレクションは全件をページネーション付きで返す", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		all := []model.PostCollection{{Title: "private", UserID: 1}, {Title: "public", UserID: 1, IsPublic: true}}
		repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).Return(all, int64(2), nil)
		uc := usecase.NewListPostCollectionsForViewerUseCase(repo)

		got, total, err := uc.Execute(context.Background(), 1, 1, 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, int64(2), total)
		repo.AssertNotCalled(t, "FindPublicByUserID")
		repo.AssertExpectations(t)
	})

	t.Run("他人のコレクションは公開のみを返し、総数は取得件数になる", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		pub := []model.PostCollection{{Title: "public", UserID: 2, IsPublic: true}}
		repo.On("FindPublicByUserID", mock.Anything, uint(2)).Return(pub, nil)
		uc := usecase.NewListPostCollectionsForViewerUseCase(repo)

		got, total, err := uc.Execute(context.Background(), 1, 2, 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
		repo.AssertNotCalled(t, "FindByUserID")
		repo.AssertExpectations(t)
	})

	t.Run("公開一覧のエラーを伝播する", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("FindPublicByUserID", mock.Anything, uint(2)).
			Return([]model.PostCollection(nil), errors.New("db error"))
		uc := usecase.NewListPostCollectionsForViewerUseCase(repo)

		_, total, err := uc.Execute(context.Background(), 1, 2, 20, 0)

		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
	})
}

func TestUpdatePostCollectionUseCase_Execute(t *testing.T) {
	t.Run("渡された値で上書きする", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		existing := &model.PostCollection{Title: "元", Description: "元説明", IsPublic: true, UserID: 1}
		repo.On("FindByID", mock.Anything, uint(5)).Return(existing, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(c *model.PostCollection) bool {
			return c.Title == "新" && c.Description == "" && !c.IsPublic
		})).Return(nil)
		uc := usecase.NewUpdatePostCollectionUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1, "  新  ", "", false)

		assert.NoError(t, err)
		assert.Equal(t, "新", got.Title)
		assert.Empty(t, got.Description)
		assert.False(t, got.IsPublic)
		repo.AssertExpectations(t)
	})

	t.Run("他人のコレクションは 403（更新しない）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostCollection{UserID: 99}, nil)
		uc := usecase.NewUpdatePostCollectionUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, "新", "", false)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("タイトルが空なら 400（更新しない）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostCollection{UserID: 1}, nil)
		uc := usecase.NewUpdatePostCollectionUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, "", "", false)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
	})
}

func TestDeletePostCollectionUseCase_Execute(t *testing.T) {
	t.Run("他人のコレクションは 403（削除しない）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostCollection{UserID: 99}, nil)
		uc := usecase.NewDeletePostCollectionUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 1))
		repo.AssertNotCalled(t, "Delete")
	})
}

func TestAddPostToCollectionUseCase_Execute(t *testing.T) {
	t.Run("メモの前後空白を落として追加する", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostCollection{UserID: 1}, nil)
		repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(false, nil)
		repo.On("AddPost", mock.Anything, mock.MatchedBy(func(i *model.PostCollectionItem) bool {
			return i.CollectionID == 5 && i.PostID == 10 && i.Note == "メモ"
		})).Return(nil)
		uc := usecase.NewAddPostToCollectionUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 5, 1, 10, "  メモ  "))
		repo.AssertExpectations(t)
	})

	t.Run("追加済みなら 400（追加しない）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostCollection{UserID: 1}, nil)
		repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(true, nil)
		uc := usecase.NewAddPostToCollectionUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 1, 10, ""))
		repo.AssertNotCalled(t, "AddPost")
	})

	t.Run("メモの検証は所有権チェックより先に行う（移行前と同じ順序）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		uc := usecase.NewAddPostToCollectionUseCase(repo)

		// 他人のコレクションであっても、メモが長すぎる場合は所有権を見る前にエラーになる
		err := uc.Execute(context.Background(), 5, 1, 10, strings.Repeat("a", 501))

		assert.Error(t, err)
		repo.AssertNotCalled(t, "FindByID")
	})

	t.Run("他人のコレクションは 403（存在確認もしない）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostCollection{UserID: 99}, nil)
		uc := usecase.NewAddPostToCollectionUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 1, 10, ""))
		repo.AssertNotCalled(t, "HasPost")
		repo.AssertNotCalled(t, "AddPost")
	})
}

func TestRemovePostFromCollectionUseCase_Execute(t *testing.T) {
	t.Run("他人のコレクションは 403（取り除かない）", func(t *testing.T) {
		repo := new(mockPostCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostCollection{UserID: 99}, nil)
		uc := usecase.NewRemovePostFromCollectionUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 1, 10))
		repo.AssertNotCalled(t, "RemovePost")
	})
}
