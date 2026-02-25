package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookmarkCollection_Create_Success(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	mockRepo.On("Create", mock.MatchedBy(func(c *model.BookmarkCollection) bool {
		return c.UserID == 1 && c.Name == "Go学習"
	})).Return(nil)

	collection := &model.BookmarkCollection{UserID: 1, Name: "Go学習", Color: "blue"}
	err := svc.Create(collection)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookmarkCollection_Create_EmptyName(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	collection := &model.BookmarkCollection{UserID: 1, Name: ""}
	err := svc.Create(collection)
	assert.Error(t, err)
}

func TestBookmarkCollection_Update_Success(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1, Name: "旧名"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.MatchedBy(func(c *model.BookmarkCollection) bool {
		return c.Name == "新名"
	})).Return(nil)

	updated, err := svc.Update(1, 1, &model.BookmarkCollection{Name: "新名"})
	assert.NoError(t, err)
	assert.Equal(t, "新名", updated.Name)
}

func TestBookmarkCollection_Update_Forbidden(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 2, Name: "他人のコレクション"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 1, &model.BookmarkCollection{Name: "新名"})
	assert.Error(t, err)
}

func TestBookmarkCollection_Delete_Success(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
}

func TestBookmarkCollection_Delete_Forbidden(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 2}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 1)
	assert.Error(t, err)
}

func TestBookmarkCollection_AddPost_Success(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("HasPost", uint(1), uint(10)).Return(false, nil)
	mockRepo.On("AddPost", mock.MatchedBy(func(item *model.BookmarkCollectionItem) bool {
		return item.CollectionID == 1 && item.PostID == 10
	})).Return(nil)

	err := svc.AddPost(1, 10, 1)
	assert.NoError(t, err)
}

func TestBookmarkCollection_AddPost_AlreadyExists(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("HasPost", uint(1), uint(10)).Return(true, nil)

	err := svc.AddPost(1, 10, 1)
	assert.Error(t, err)
}

func TestBookmarkCollection_GetByUserID(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	collections := []model.BookmarkCollection{
		{ID: 1, UserID: 1, Name: "Go"},
		{ID: 2, UserID: 1, Name: "React"},
	}
	mockRepo.On("FindByUserID", uint(1)).Return(collections, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// --- RemovePost テスト ---

func TestBookmarkCollection_RemovePost_Success(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("RemovePost", uint(1), uint(10)).Return(nil)

	err := svc.RemovePost(1, 10, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookmarkCollection_RemovePost_Forbidden(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 2}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.RemovePost(1, 10, 1)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "RemovePost")
}

func TestBookmarkCollection_RemovePost_NotFound(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	mockRepo.On("FindByID", uint(99)).Return(nil, ErrNotFound)

	err := svc.RemovePost(99, 10, 1)
	assert.Error(t, err)
}

func TestBookmarkCollection_RemovePost_RepoError(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("RemovePost", uint(1), uint(10)).Return(assert.AnError)

	err := svc.RemovePost(1, 10, 1)
	assert.Error(t, err)
}

// --- GetPosts テスト ---

func TestBookmarkCollection_GetPosts_Success(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	posts := []model.Post{
		{ID: 1, Title: "記事1"},
		{ID: 2, Title: "記事2"},
	}
	mockRepo.On("GetPosts", uint(1), 20, 0).Return(posts, int64(2), nil)

	result, total, err := svc.GetPosts(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
}

func TestBookmarkCollection_GetPosts_Empty(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	mockRepo.On("GetPosts", uint(1), 20, 0).Return([]model.Post{}, int64(0), nil)

	result, total, err := svc.GetPosts(1, 20, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
}

func TestBookmarkCollection_GetPosts_RepoError(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	mockRepo.On("GetPosts", uint(1), 20, 0).Return([]model.Post{}, int64(0), assert.AnError)

	_, _, err := svc.GetPosts(1, 20, 0)
	assert.Error(t, err)
}

// --- Update 追加テスト ---

func TestBookmarkCollection_Update_RepoError(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1, Name: "旧名"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.Anything).Return(assert.AnError)

	_, err := svc.Update(1, 1, &model.BookmarkCollection{Name: "新名"})
	assert.Error(t, err)
}

func TestBookmarkCollection_Update_DescriptionAndColor(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1, Name: "既存", Description: "旧説明", Color: "blue"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.MatchedBy(func(c *model.BookmarkCollection) bool {
		return c.Description == "新説明" && c.Color == "red"
	})).Return(nil)

	updated, err := svc.Update(1, 1, &model.BookmarkCollection{Description: "新説明", Color: "red"})
	assert.NoError(t, err)
	assert.Equal(t, "新説明", updated.Description)
	assert.Equal(t, "red", updated.Color)
}

// --- AddPost 追加テスト ---

func TestBookmarkCollection_AddPost_Forbidden(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 2}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.AddPost(1, 10, 1)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "HasPost")
}

func TestBookmarkCollection_AddPost_HasPostError(t *testing.T) {
	mockRepo := new(MockBookmarkCollectionRepository)
	svc := NewBookmarkCollectionService(mockRepo)

	existing := &model.BookmarkCollection{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("HasPost", uint(1), uint(10)).Return(false, assert.AnError)

	err := svc.AddPost(1, 10, 1)
	assert.Error(t, err)
}
