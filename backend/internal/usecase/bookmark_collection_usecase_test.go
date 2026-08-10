package usecase_test

import (
	"context"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockBookmarkCollectionRepo は usecase/repository.BookmarkCollectionRepository のモック。
type mockBookmarkCollectionRepo struct{ mock.Mock }

func (m *mockBookmarkCollectionRepo) Create(ctx context.Context, collection *model.BookmarkCollection) error {
	return m.Called(ctx, collection).Error(0)
}

func (m *mockBookmarkCollectionRepo) FindByID(ctx context.Context, id uint) (*model.BookmarkCollection, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.BookmarkCollection)
	return c, args.Error(1)
}

func (m *mockBookmarkCollectionRepo) FindByUserID(ctx context.Context, userID uint) ([]model.BookmarkCollection, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]model.BookmarkCollection)
	return c, args.Error(1)
}

func (m *mockBookmarkCollectionRepo) Update(ctx context.Context, collection *model.BookmarkCollection) error {
	return m.Called(ctx, collection).Error(0)
}

func (m *mockBookmarkCollectionRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockBookmarkCollectionRepo) AddPost(ctx context.Context, item *model.BookmarkCollectionItem) error {
	return m.Called(ctx, item).Error(0)
}

func (m *mockBookmarkCollectionRepo) RemovePost(ctx context.Context, collectionID, postID uint) error {
	return m.Called(ctx, collectionID, postID).Error(0)
}

func (m *mockBookmarkCollectionRepo) GetPosts(ctx context.Context, collectionID uint, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(ctx, collectionID, limit, offset)
	p, _ := args.Get(0).([]model.Post)
	return p, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookmarkCollectionRepo) HasPost(ctx context.Context, collectionID, postID uint) (bool, error) {
	args := m.Called(ctx, collectionID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *mockBookmarkCollectionRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateBookmarkCollectionUseCase_Execute(t *testing.T) {
	t.Run("前後の空白を落として作成する", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(c *model.BookmarkCollection) bool {
			return c.Name == "名前" && c.Description == "説明" && c.Color == "blue"
		})).Return(nil)
		uc := usecase.NewCreateBookmarkCollectionUseCase(repo)

		err := uc.Execute(context.Background(), &model.BookmarkCollection{
			Name: "  名前  ", Description: "  説明  ", Color: "  blue  ",
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("説明とカラーは空でも作成できる", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.BookmarkCollection")).Return(nil)
		uc := usecase.NewCreateBookmarkCollectionUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), &model.BookmarkCollection{Name: "名前"}))
		repo.AssertExpectations(t)
	})

	for name, in := range map[string]*model.BookmarkCollection{
		"名前が空":        {Name: ""},
		"名前が 100 文字超": {Name: strings.Repeat("a", 101)},
		"説明が 500 文字超": {Name: "名前", Description: strings.Repeat("a", 501)},
		"カラーが 20 文字超": {Name: "名前", Color: strings.Repeat("a", 21)},
	} {
		t.Run(name+"は 400（作成しない）", func(t *testing.T) {
			repo := new(mockBookmarkCollectionRepo)
			uc := usecase.NewCreateBookmarkCollectionUseCase(repo)

			assert.Error(t, uc.Execute(context.Background(), in))
			repo.AssertNotCalled(t, "Create")
		})
	}
}

func TestUpdateBookmarkCollectionUseCase_Execute(t *testing.T) {
	t.Run("トリム後に空でないフィールドだけ更新する", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		existing := &model.BookmarkCollection{Name: "元名", Description: "元説明", Color: "red", UserID: 1}
		repo.On("FindByID", mock.Anything, uint(5)).Return(existing, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(c *model.BookmarkCollection) bool {
			return c.Name == "新名" && c.Description == "元説明" && c.Color == "red"
		})).Return(nil)
		uc := usecase.NewUpdateBookmarkCollectionUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1, &model.BookmarkCollection{Name: "  新名  ", Description: "   "})

		assert.NoError(t, err)
		assert.Equal(t, "新名", got.Name)
		assert.Equal(t, "元説明", got.Description)
		repo.AssertExpectations(t)
	})

	t.Run("他人のコレクションは 403（更新しない）", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.BookmarkCollection{UserID: 99}, nil)
		uc := usecase.NewUpdateBookmarkCollectionUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, &model.BookmarkCollection{Name: "新名"})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
	})
}

func TestDeleteBookmarkCollectionUseCase_Execute(t *testing.T) {
	t.Run("他人のコレクションは 403（削除しない）", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.BookmarkCollection{UserID: 99}, nil)
		uc := usecase.NewDeleteBookmarkCollectionUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 1))
		repo.AssertNotCalled(t, "Delete")
	})
}

func TestAddPostToBookmarkCollectionUseCase_Execute(t *testing.T) {
	t.Run("未追加なら追加する", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.BookmarkCollection{UserID: 1}, nil)
		repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(false, nil)
		repo.On("AddPost", mock.Anything, mock.MatchedBy(func(i *model.BookmarkCollectionItem) bool {
			return i.CollectionID == 5 && i.PostID == 10
		})).Return(nil)
		uc := usecase.NewAddPostToBookmarkCollectionUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 5, 10, 1))
		repo.AssertExpectations(t)
	})

	// post_collection は 400 だが、こちらは 409（Conflict）で異なる
	t.Run("追加済みなら Conflict（追加しない）", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.BookmarkCollection{UserID: 1}, nil)
		repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(true, nil)
		uc := usecase.NewAddPostToBookmarkCollectionUseCase(repo)

		err := uc.Execute(context.Background(), 5, 10, 1)

		assert.Error(t, err)
		var de *domain.DomainError
		if assert.ErrorAs(t, err, &de) {
			assert.Equal(t, domain.ErrCodeConflict, de.Code)
		}
		repo.AssertNotCalled(t, "AddPost")
	})

	t.Run("他人のコレクションは 403（存在確認もしない）", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.BookmarkCollection{UserID: 99}, nil)
		uc := usecase.NewAddPostToBookmarkCollectionUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 10, 1))
		repo.AssertNotCalled(t, "HasPost")
		repo.AssertNotCalled(t, "AddPost")
	})
}

func TestRemovePostFromBookmarkCollectionUseCase_Execute(t *testing.T) {
	t.Run("他人のコレクションは 403（取り除かない）", func(t *testing.T) {
		repo := new(mockBookmarkCollectionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.BookmarkCollection{UserID: 99}, nil)
		uc := usecase.NewRemovePostFromBookmarkCollectionUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 5, 10, 1))
		repo.AssertNotCalled(t, "RemovePost")
	})
}
