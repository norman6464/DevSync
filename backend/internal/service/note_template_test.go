package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNoteTemplateRepository はNoteTemplateRepositoryのモック実装。
type MockNoteTemplateRepository struct {
	mock.Mock
}

func (m *MockNoteTemplateRepository) Create(template *model.NoteTemplate) error {
	return m.Called(template).Error(0)
}

func (m *MockNoteTemplateRepository) FindByID(id uint) (*model.NoteTemplate, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NoteTemplate), args.Error(1)
}

func (m *MockNoteTemplateRepository) FindByUserID(userID uint) ([]model.NoteTemplate, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.NoteTemplate), args.Error(1)
}

func (m *MockNoteTemplateRepository) FindDefaultByUserID(userID uint) (*model.NoteTemplate, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NoteTemplate), args.Error(1)
}

func (m *MockNoteTemplateRepository) Update(template *model.NoteTemplate) error {
	return m.Called(template).Error(0)
}

func (m *MockNoteTemplateRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockNoteTemplateRepository) ClearDefaultFlag(userID uint) error {
	return m.Called(userID).Error(0)
}

// newTestNoteTemplateService はテスト用のNoteTemplateServiceを生成する。
func newTestNoteTemplateService() (*NoteTemplateService, *MockNoteTemplateRepository) {
	repo := new(MockNoteTemplateRepository)
	svc := NewNoteTemplateService(repo)
	return svc, repo
}

// ============================================================
// Create テスト
// ============================================================

func TestNoteTemplateService_Create(t *testing.T) {
	svc, repo := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		UserID:          1,
		Name:            "週報テンプレート",
		Description:     "週報用のテンプレート",
		DefaultTitle:    "週報 - {{date}}",
		ContentTemplate: "# 週報\n\n## 今週やったこと",
		DefaultTags:     "週報",
		IsDefault:       false,
	}

	repo.On("Create", template).Return(nil)

	err := svc.Create(template)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Create_WithDefaultFlag(t *testing.T) {
	svc, repo := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		UserID:          1,
		Name:            "デフォルトテンプレート",
		ContentTemplate: "本文",
		IsDefault:       true,
	}

	repo.On("ClearDefaultFlag", uint(1)).Return(nil)
	repo.On("Create", template).Return(nil)

	err := svc.Create(template)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestNoteTemplateService_GetByID(t *testing.T) {
	svc, repo := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		Name:            "テンプレート1",
		ContentTemplate: "本文",
	}

	repo.On("FindByID", uint(1)).Return(template, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, template, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestNoteTemplateService_GetByUserID(t *testing.T) {
	svc, repo := newTestNoteTemplateService()

	templates := []model.NoteTemplate{
		{ID: 1, UserID: 1, Name: "テンプレート1", ContentTemplate: "本文1"},
		{ID: 2, UserID: 1, Name: "テンプレート2", ContentTemplate: "本文2"},
	}

	repo.On("FindByUserID", uint(1)).Return(templates, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	repo.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestNoteTemplateService_Update(t *testing.T) {
	svc, repo := newTestNoteTemplateService()

	existing := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		Name:            "元テンプレート",
		ContentTemplate: "元本文",
		IsDefault:       false,
	}

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*model.NoteTemplate")).Return(nil)

	result, err := svc.Update(1, 1, "更新後テンプレート", "", "", "更新後本文", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, "更新後テンプレート", result.Name)
	assert.Equal(t, "更新後本文", result.ContentTemplate)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Update_Forbidden(t *testing.T) {
	svc, repo := newTestNoteTemplateService()

	existing := &model.NoteTemplate{
		ID:     1,
		UserID: 1,
	}

	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 999, "", "", "", "", "", nil)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteTemplateService_Delete(t *testing.T) {
	svc, repo := newTestNoteTemplateService()

	existing := &model.NoteTemplate{ID: 1, UserID: 1}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Delete_Forbidden(t *testing.T) {
	svc, repo := newTestNoteTemplateService()

	existing := &model.NoteTemplate{ID: 1, UserID: 1}
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}
