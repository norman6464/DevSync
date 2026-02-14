package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNoteFolderRepository はNoteFolderRepositoryのモック実装
type MockNoteFolderRepository struct {
	mock.Mock
}

func (m *MockNoteFolderRepository) Create(folder *model.NoteFolder) error {
	args := m.Called(folder)
	return args.Error(0)
}

func (m *MockNoteFolderRepository) FindByID(id uint) (*model.NoteFolder, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderRepository) FindByUserID(userID uint) ([]model.NoteFolder, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderRepository) FindByParentID(parentID uint) ([]model.NoteFolder, error) {
	args := m.Called(parentID)
	return args.Get(0).([]model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderRepository) GetRootFolders(userID uint) ([]model.NoteFolder, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderRepository) Update(folder *model.NoteFolder) error {
	args := m.Called(folder)
	return args.Error(0)
}

func (m *MockNoteFolderRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestNoteFolderService_Create(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	folder := &model.NoteFolder{
		UserID: 1,
		Name:   "テストフォルダ",
	}

	mockRepo.On("Create", folder).Return(nil)

	err := service.Create(folder)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_GetByID(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	expected := &model.NoteFolder{
		ID:     1,
		UserID: 1,
		Name:   "テストフォルダ",
	}

	mockRepo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := service.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_GetByUserID(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	expected := []model.NoteFolder{
		{ID: 1, UserID: 1, Name: "フォルダ1"},
		{ID: 2, UserID: 1, Name: "フォルダ2"},
	}

	mockRepo.On("FindByUserID", uint(1)).Return(expected, nil)

	result, err := service.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_GetChildren(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	expected := []model.NoteFolder{
		{ID: 2, UserID: 1, Name: "子フォルダ1"},
		{ID: 3, UserID: 1, Name: "子フォルダ2"},
	}

	mockRepo.On("FindByParentID", uint(1)).Return(expected, nil)

	result, err := service.GetChildren(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_GetRootFolders(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	expected := []model.NoteFolder{
		{ID: 1, UserID: 1, Name: "ルートフォルダ1"},
		{ID: 2, UserID: 1, Name: "ルートフォルダ2"},
	}

	mockRepo.On("GetRootFolders", uint(1)).Return(expected, nil)

	result, err := service.GetRootFolders(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Update(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	folder := &model.NoteFolder{
		ID:     1,
		UserID: 1,
		Name:   "更新後の名前",
	}

	mockRepo.On("Update", folder).Return(nil)

	err := service.Update(folder)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Delete(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
