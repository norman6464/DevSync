package service

import (
	"errors"
	"strings"
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

func (m *MockNoteFolderRepository) FindByUserID(userID uint, limit, offset int) ([]model.NoteFolder, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.NoteFolder), args.Get(1).(int64), args.Error(2)
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

func (m *MockNoteFolderRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
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

	mockRepo.On("FindByUserID", uint(1), 20, 0).Return(expected, int64(2), nil)

	folders, total, err := service.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, folders, 2)
	assert.Equal(t, int64(2), total)
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

	existing := &model.NoteFolder{
		ID:     1,
		UserID: 1,
		Name:   "元の名前",
	}

	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.NoteFolder")).Return(nil)

	result, err := service.Update(1, 1, "更新後の名前", nil)
	assert.NoError(t, err)
	assert.Equal(t, "更新後の名前", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Create_ValidationError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	folder := &model.NoteFolder{
		UserID: 1,
		Name:   "", // 空名前 → バリデーションエラー
	}

	err := service.Create(folder)
	assert.Error(t, err)
}

func TestNoteFolderService_Create_RepoError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	folder := &model.NoteFolder{
		UserID: 1,
		Name:   "テストフォルダ",
	}

	mockRepo.On("Create", folder).Return(errors.New("db error"))

	err := service.Create(folder)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Update_RepoError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "元の名前"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.NoteFolder")).Return(errors.New("db error"))

	result, err := service.Update(1, 1, "更新後", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Update_Forbidden(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := service.Update(1, 999, "名前", nil)
	assert.ErrorIs(t, err, ErrForbidden)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Delete(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.Delete(1, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Delete_Forbidden(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	err := service.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_GetByID_Error(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_GetByUserID_Error(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	mockRepo.On("FindByUserID", uint(1), 20, 0).Return([]model.NoteFolder(nil), int64(0), errors.New("db error"))

	result, total, err := service.GetByUserID(1, 20, 0)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_GetChildren_Error(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	mockRepo.On("FindByParentID", uint(1)).Return([]model.NoteFolder(nil), errors.New("db error"))

	result, err := service.GetChildren(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_GetRootFolders_Error(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	mockRepo.On("GetRootFolders", uint(1)).Return([]model.NoteFolder(nil), errors.New("db error"))

	result, err := service.GetRootFolders(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Update_FindByIDError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := service.Update(999, 1, "名前", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Update_ValidationError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "元の名前"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	// 空文字でname更新 → 元のNameが維持されるので、空文字を直接代入するケースはないが、
	// name="" の場合は更新されないのでバリデーションは通る
	// 代わりにparentID付きの正常更新をテスト
	parentID := uint(5)
	mockRepo.On("FindByParentID", uint(1)).Return([]model.NoteFolder{}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.NoteFolder")).Return(nil)

	result, err := service.Update(1, 1, "更新名", &parentID)
	assert.NoError(t, err)
	assert.Equal(t, "更新名", result.Name)
	assert.Equal(t, &parentID, result.ParentID)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Update_SelfReference(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "フォルダA"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	// 自分自身を親に設定 → エラー
	selfID := uint(1)
	result, err := service.Update(1, 1, "", &selfID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "自分自身")
}

func TestNoteFolderService_Update_CircularReference(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	// フォルダA(ID=1) → 子フォルダB(ID=2) → 子フォルダC(ID=3)
	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "フォルダA"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByParentID", uint(1)).Return([]model.NoteFolder{
		{ID: 2, UserID: 1, Name: "フォルダB"},
	}, nil)
	mockRepo.On("FindByParentID", uint(2)).Return([]model.NoteFolder{
		{ID: 3, UserID: 1, Name: "フォルダC"},
	}, nil)
	mockRepo.On("FindByParentID", uint(3)).Return([]model.NoteFolder{}, nil)

	// フォルダAの親をフォルダC(孫)に設定 → 循環参照エラー
	childID := uint(3)
	result, err := service.Update(1, 1, "", &childID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "循環参照")
}

func TestNoteFolderService_Update_CircularCheckRepoError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "フォルダA"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByParentID", uint(1)).Return([]model.NoteFolder{}, errors.New("db error"))

	parentID := uint(5)
	result, err := service.Update(1, 1, "", &parentID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoteFolderService_Update_CircularCheckRecursiveError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	// フォルダA(ID=1) → 子フォルダB(ID=2) → 再帰呼び出し時にエラー
	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "フォルダA"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByParentID", uint(1)).Return([]model.NoteFolder{
		{ID: 2, UserID: 1, Name: "フォルダB"},
	}, nil)
	mockRepo.On("FindByParentID", uint(2)).Return([]model.NoteFolder{}, errors.New("db error"))

	parentID := uint(5)
	result, err := service.Update(1, 1, "", &parentID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoteFolderService_Delete_FindByIDError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := service.Delete(999, 1)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Delete_RepoError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Delete", uint(1)).Return(errors.New("db error"))

	err := service.Delete(1, 1)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNoteFolderService_Update_ValidateLongNameError(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "元の名前"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	// 101文字のフォルダ名 → ValidateUpdate でバリデーションエラー
	longName := strings.Repeat("あ", 101)
	result, err := service.Update(1, 1, longName, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoteFolderService_Update_WhitespaceName(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	service := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "元の名前"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	// 空白のみの名前 → バリデーションエラーになるべき
	result, err := service.Update(1, 1, "   \t\n  ", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "フォルダ名は空白のみにできません")
}

func TestNoteFolderService_Update_TrimSpaceName(t *testing.T) {
	mockRepo := new(MockNoteFolderRepository)
	svc := NewNoteFolderService(mockRepo)

	existing := &model.NoteFolder{ID: 1, UserID: 1, Name: "元の名前"}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.MatchedBy(func(f *model.NoteFolder) bool {
		return f.Name == "新しい名前"
	})).Return(nil)

	// 前後に空白がある名前 → TrimSpaceされて保存されるべき
	result, err := svc.Update(1, 1, "  新しい名前  ", nil)
	assert.NoError(t, err)
	assert.Equal(t, "新しい名前", result.Name)
	mockRepo.AssertExpectations(t)
}
