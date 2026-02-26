package service

import (
	"errors"
	"strings"
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

// MockNoteCreator はNoteCreatorInterfaceのモック実装。
type MockNoteCreator struct {
	mock.Mock
}

func (m *MockNoteCreator) Create(note *model.Note) error {
	return m.Called(note).Error(0)
}

// newTestNoteTemplateService はテスト用のNoteTemplateServiceを生成する。
func newTestNoteTemplateService() (*NoteTemplateService, *MockNoteTemplateRepository, *MockNoteCreator) {
	repo := new(MockNoteTemplateRepository)
	noteCreator := new(MockNoteCreator)
	svc := NewNoteTemplateService(repo, noteCreator)
	return svc, repo, noteCreator
}

// ============================================================
// Create テスト
// ============================================================

func TestNoteTemplateService_Create(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

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
	svc, repo, _ := newTestNoteTemplateService()

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

func TestNoteTemplateService_Create_ValidationError(t *testing.T) {
	svc, _, _ := newTestNoteTemplateService()

	// Name空 → バリデーションエラー
	template := &model.NoteTemplate{
		UserID:          1,
		Name:            "",
		ContentTemplate: "本文",
	}

	err := svc.Create(template)
	assert.Error(t, err)
}

func TestNoteTemplateService_Create_DescriptionValidationError(t *testing.T) {
	svc, _, _ := newTestNoteTemplateService()

	// 超長い説明文 → バリデーションエラー
	longDesc := ""
	for i := 0; i < 600; i++ {
		longDesc += "あ"
	}

	template := &model.NoteTemplate{
		UserID:          1,
		Name:            "テンプレート",
		ContentTemplate: "本文",
		Description:     longDesc,
	}

	err := svc.Create(template)
	assert.Error(t, err)
}

func TestNoteTemplateService_Create_ClearDefaultFlagError(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		UserID:          1,
		Name:            "テンプレート",
		ContentTemplate: "本文",
		IsDefault:       true,
	}

	repo.On("ClearDefaultFlag", uint(1)).Return(errors.New("db error"))

	err := svc.Create(template)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestNoteTemplateService_GetByID(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		Name:            "テンプレート1",
		ContentTemplate: "本文",
	}

	repo.On("FindByID", uint(1)).Return(template, nil)

	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, template, result)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_GetByID_Forbidden(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		Name:            "テンプレート1",
		ContentTemplate: "本文",
	}

	repo.On("FindByID", uint(1)).Return(template, nil)

	result, err := svc.GetByID(1, 999)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrForbidden, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestNoteTemplateService_GetByUserID(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

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
	svc, repo, _ := newTestNoteTemplateService()

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
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{
		ID:     1,
		UserID: 1,
	}

	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 999, "", "", "", "", "", nil)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Update_WithAllFields(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		Name:            "元テンプレート",
		ContentTemplate: "元本文",
		IsDefault:       false,
	}

	isDefault := true
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("ClearDefaultFlag", uint(1)).Return(nil)
	repo.On("Update", mock.AnythingOfType("*model.NoteTemplate")).Return(nil)

	result, err := svc.Update(1, 1, "新名前", "新説明", "新タイトル", "新本文", "tag1,tag2", &isDefault)
	assert.NoError(t, err)
	assert.Equal(t, "新名前", result.Name)
	assert.Equal(t, "新説明", result.Description)
	assert.Equal(t, "新タイトル", result.DefaultTitle)
	assert.Equal(t, "新本文", result.ContentTemplate)
	assert.Equal(t, "tag1,tag2", result.DefaultTags)
	assert.True(t, result.IsDefault)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Update_RepoError(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		Name:            "テンプレ",
		ContentTemplate: "本文",
	}

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*model.NoteTemplate")).Return(errors.New("db error"))

	result, err := svc.Update(1, 1, "新名前", "", "", "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Update_ClearDefaultFlagError(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		Name:            "テンプレ",
		ContentTemplate: "本文",
	}

	isDefault := true
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("ClearDefaultFlag", uint(1)).Return(errors.New("db error"))

	result, err := svc.Update(1, 1, "", "", "", "", "", &isDefault)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Update_ValidationNameError(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{ID: 1, UserID: 1, Name: "テンプレ"}

	repo.On("FindByID", uint(1)).Return(existing, nil)

	// 101文字の名前 → バリデーションエラー
	longName := strings.Repeat("あ", 101)
	result, err := svc.Update(1, 1, longName, "", "", "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoteTemplateService_Update_ValidationContentTemplateError(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{ID: 1, UserID: 1, Name: "テンプレ"}

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 空白のみのcontentTemplate → TrimSpace後に空 → フィールドスキップ（変更なし）
	result, err := svc.Update(1, 1, "", "", "", "   ", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, "", result.ContentTemplate)
}

func TestNoteTemplateService_Update_ValidationDescriptionError(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{ID: 1, UserID: 1, Name: "テンプレ"}

	repo.On("FindByID", uint(1)).Return(existing, nil)

	// 501文字の説明 → バリデーションエラー
	longDesc := strings.Repeat("a", 501)
	result, err := svc.Update(1, 1, "", longDesc, "", "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoteTemplateService_Update_ValidationDefaultTitleError(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{ID: 1, UserID: 1, Name: "テンプレ"}

	repo.On("FindByID", uint(1)).Return(existing, nil)

	// 201文字のデフォルトタイトル → バリデーションエラー
	longTitle := strings.Repeat("a", 201)
	result, err := svc.Update(1, 1, "", "", longTitle, "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteTemplateService_Delete(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{ID: 1, UserID: 1}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Delete_Forbidden(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	existing := &model.NoteTemplate{ID: 1, UserID: 1}
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// UseTemplate テスト
// ============================================================

func TestNoteTemplateService_UseTemplate_Success(t *testing.T) {
	svc, repo, noteCreator := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		DefaultTitle:    "テンプレタイトル",
		ContentTemplate: "テンプレ内容",
		DefaultTags:     "tag1",
	}

	repo.On("FindByID", uint(1)).Return(template, nil)
	noteCreator.On("Create", mock.AnythingOfType("*model.Note")).Return(nil)

	note, err := svc.UseTemplate(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "テンプレタイトル", note.Title)
	assert.Equal(t, "テンプレ内容", note.Content)
	assert.Equal(t, "tag1", note.Tags)
	assert.Equal(t, uint(1), note.UserID)
	repo.AssertExpectations(t)
	noteCreator.AssertExpectations(t)
}

func TestNoteTemplateService_UseTemplate_DefaultTitle(t *testing.T) {
	svc, repo, noteCreator := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		DefaultTitle:    "",
		ContentTemplate: "内容",
	}

	repo.On("FindByID", uint(1)).Return(template, nil)
	noteCreator.On("Create", mock.AnythingOfType("*model.Note")).Return(nil)

	note, err := svc.UseTemplate(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "新しいノート", note.Title)
	repo.AssertExpectations(t)
	noteCreator.AssertExpectations(t)
}

func TestNoteTemplateService_UseTemplate_Forbidden(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	template := &model.NoteTemplate{ID: 1, UserID: 1}
	repo.On("FindByID", uint(1)).Return(template, nil)

	_, err := svc.UseTemplate(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_UseTemplate_NotFound(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	repo.On("FindByID", uint(99)).Return(nil, ErrNotFound)

	_, err := svc.UseTemplate(99, 1)
	assert.ErrorIs(t, err, ErrNotFound)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_UseTemplate_CreateError(t *testing.T) {
	svc, repo, noteCreator := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		DefaultTitle:    "タイトル",
		ContentTemplate: "内容",
	}

	repo.On("FindByID", uint(1)).Return(template, nil)
	noteCreator.On("Create", mock.AnythingOfType("*model.Note")).Return(errors.New("db error"))

	_, err := svc.UseTemplate(1, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
	noteCreator.AssertExpectations(t)
}

// ============================================================
// GetDefaultByUserID テスト
// ============================================================

func TestNoteTemplateService_GetDefaultByUserID_Success(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	template := &model.NoteTemplate{
		ID:              1,
		UserID:          1,
		Name:            "デフォルトテンプレート",
		ContentTemplate: "本文",
		IsDefault:       true,
	}
	repo.On("FindDefaultByUserID", uint(1)).Return(template, nil)

	result, err := svc.GetDefaultByUserID(1)
	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
	assert.Equal(t, "デフォルトテンプレート", result.Name)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_GetDefaultByUserID_NotFound(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	repo.On("FindDefaultByUserID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.GetDefaultByUserID(99)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Delete_NotFound(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.Delete(99, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestNoteTemplateService_Create_DefaultTitleValidationError(t *testing.T) {
	svc, _, _ := newTestNoteTemplateService()

	// 201文字のデフォルトタイトル → バリデーションエラー
	longTitle := strings.Repeat("あ", 201)
	template := &model.NoteTemplate{
		UserID:          1,
		Name:            "テンプレート",
		ContentTemplate: "本文",
		DefaultTitle:    longTitle,
	}

	err := svc.Create(template)
	assert.Error(t, err)
}

// ============================================================
// 空白バイパス脆弱性テスト
// ============================================================

func TestNoteTemplateUpdate_WhitespaceDescription(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()
	existing := &model.NoteTemplate{Name: "Template", ContentTemplate: "Content", Description: "Original Desc", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "", "   ", "", "", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, "Original Desc", result.Description)
}

func TestNoteTemplateUpdate_WhitespaceDefaultTitle(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()
	existing := &model.NoteTemplate{Name: "Template", ContentTemplate: "Content", DefaultTitle: "Original Title", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "", "", "   ", "", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, "Original Title", result.DefaultTitle)
}

func TestNoteTemplateUpdate_WhitespaceDefaultTags(t *testing.T) {
	svc, repo, _ := newTestNoteTemplateService()
	existing := &model.NoteTemplate{Name: "Template", ContentTemplate: "Content", DefaultTags: "go,rust", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "", "", "", "", "   ", nil)
	assert.NoError(t, err)
	assert.Equal(t, "go,rust", result.DefaultTags)
}
