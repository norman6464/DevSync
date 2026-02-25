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
