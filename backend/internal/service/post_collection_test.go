package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestPostCollectionService はPostCollectionServiceのテスト用インスタンスを生成するヘルパー。
func newTestPostCollectionService() (*PostCollectionService, *MockPostCollectionRepository) {
	repo := new(MockPostCollectionRepository)
	svc := NewPostCollectionService(repo)
	return svc, repo
}

// ============================================================
// コレクション作成テスト
// ============================================================

func TestPostCollectionCreate_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{
		UserID:      1,
		Title:       "Goベストプラクティス",
		Description: "Go言語のベストプラクティス集",
		IsPublic:    true,
	}

	repo.On("Create", collection).Run(func(args mock.Arguments) {
		c := args.Get(0).(*model.PostCollection)
		c.ID = 10
	}).Return(nil)

	result, err := svc.Create(collection)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(10), result.ID)
	repo.AssertExpectations(t)
}

func TestPostCollectionCreate_EmptyTitle(t *testing.T) {
	svc, _ := newTestPostCollectionService()

	collection := &model.PostCollection{
		UserID: 1,
		Title:  "",
	}

	result, err := svc.Create(collection)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPostCollectionCreate_RepoError(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{
		UserID: 1,
		Title:  "Test",
	}

	repo.On("Create", collection).Return(assert.AnError)

	result, err := svc.Create(collection)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestPostCollectionGetByID_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{Title: "Go Tips", UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)

	result, err := svc.GetByID(10)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Go Tips", result.Title)
	repo.AssertExpectations(t)
}

func TestPostCollectionGetByID_NotFound(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	repo.On("FindByID", uint(999)).Return(nil, assert.AnError)

	result, err := svc.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestPostCollectionGetByUserID_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collections := []model.PostCollection{
		{Title: "Go Tips", UserID: 1},
		{Title: "React集", UserID: 1},
	}
	repo.On("FindByUserID", uint(1)).Return(collections, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestPostCollectionGetByUserID_Empty(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	repo.On("FindByUserID", uint(999)).Return([]model.PostCollection{}, nil)

	result, err := svc.GetByUserID(999)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetPublicByUserID テスト
// ============================================================

func TestPostCollectionGetPublicByUserID_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collections := []model.PostCollection{
		{Title: "Public集", UserID: 2, IsPublic: true},
	}
	repo.On("FindPublicByUserID", uint(2)).Return(collections, nil)

	result, err := svc.GetPublicByUserID(2)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.True(t, result[0].IsPublic)
	repo.AssertExpectations(t)
}

// ============================================================
// Update テスト（所有権チェック）
// ============================================================

func TestPostCollectionUpdate_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{Title: "Old", UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)
	repo.On("Update", collection).Return(nil)

	result, err := svc.Update(10, 1, "New Title", "New Desc", true)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Desc", result.Description)
	assert.True(t, result.IsPublic)
	repo.AssertExpectations(t)
}

func TestPostCollectionUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{Title: "Test", UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)

	result, err := svc.Update(10, 999, "New", "", false)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestPostCollectionUpdate_NotFound(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	repo.On("FindByID", uint(999)).Return(nil, assert.AnError)

	result, err := svc.Update(999, 1, "New", "", false)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestPostCollectionUpdate_RepoError(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{Title: "Test", UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)
	repo.On("Update", collection).Return(assert.AnError)

	result, err := svc.Update(10, 1, "New", "", false)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete テスト（所有権チェック）
// ============================================================

func TestPostCollectionDelete_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)
	repo.On("Delete", uint(10)).Return(nil)

	err := svc.Delete(10, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostCollectionDelete_Forbidden(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)

	err := svc.Delete(10, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestPostCollectionDelete_NotFound(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	repo.On("FindByID", uint(999)).Return(nil, assert.AnError)

	err := svc.Delete(999, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// AddPost テスト（所有権チェック）
// ============================================================

func TestPostCollectionAddPost_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)
	repo.On("AddPost", mock.MatchedBy(func(item *model.PostCollectionItem) bool {
		return item.CollectionID == 10 && item.PostID == 5
	})).Return(nil)

	err := svc.AddPost(10, 1, 5, "良い記事")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostCollectionAddPost_Forbidden(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)

	err := svc.AddPost(10, 999, 5, "")
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestPostCollectionAddPost_NotFound(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	repo.On("FindByID", uint(999)).Return(nil, assert.AnError)

	err := svc.AddPost(999, 1, 5, "")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// RemovePost テスト（所有権チェック）
// ============================================================

func TestPostCollectionRemovePost_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)
	repo.On("RemovePost", uint(10), uint(5)).Return(nil)

	err := svc.RemovePost(10, 1, 5)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostCollectionRemovePost_Forbidden(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	collection := &model.PostCollection{UserID: 1}
	collection.ID = 10
	repo.On("FindByID", uint(10)).Return(collection, nil)

	err := svc.RemovePost(10, 999, 5)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// GetPosts テスト
// ============================================================

func TestPostCollectionGetPosts_Success(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	items := []model.PostCollectionItem{
		{CollectionID: 10, PostID: 1, OrderIndex: 0},
		{CollectionID: 10, PostID: 2, OrderIndex: 1},
	}
	repo.On("GetPostsByCollectionID", uint(10)).Return(items, nil)

	result, err := svc.GetPosts(10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestPostCollectionGetPosts_Empty(t *testing.T) {
	svc, repo := newTestPostCollectionService()

	repo.On("GetPostsByCollectionID", uint(10)).Return([]model.PostCollectionItem{}, nil)

	result, err := svc.GetPosts(10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}
