package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLearningLogTemplateRepository は学習ログテンプレートリポジトリのモック実装。
type MockLearningLogTemplateRepository struct {
	mock.Mock
}

func (m *MockLearningLogTemplateRepository) Create(template *model.LearningLogTemplate) error {
	return m.Called(template).Error(0)
}

func (m *MockLearningLogTemplateRepository) FindByID(id uint) (*model.LearningLogTemplate, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningLogTemplate), args.Error(1)
}

func (m *MockLearningLogTemplateRepository) FindByUserID(userID uint) ([]model.LearningLogTemplate, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.LearningLogTemplate), args.Error(1)
}

func (m *MockLearningLogTemplateRepository) FindDefaultByUserID(userID uint) (*model.LearningLogTemplate, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningLogTemplate), args.Error(1)
}

func (m *MockLearningLogTemplateRepository) Update(template *model.LearningLogTemplate) error {
	return m.Called(template).Error(0)
}

func (m *MockLearningLogTemplateRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockLearningLogTemplateRepository) ClearDefaultFlag(userID uint) error {
	return m.Called(userID).Error(0)
}

func (m *MockLearningLogTemplateRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// MockLearningLogCreator は学習ログ作成のモック実装。
type MockLearningLogCreator struct {
	mock.Mock
}

func (m *MockLearningLogCreator) Create(log *model.LearningLog) error {
	return m.Called(log).Error(0)
}

// newTestLogTemplateService はテスト用のLearningLogTemplateServiceを生成する。
func newTestLogTemplateService() (*LearningLogTemplateService, *MockLearningLogTemplateRepository, *MockLearningLogCreator) {
	repo := new(MockLearningLogTemplateRepository)
	logCreator := new(MockLearningLogCreator)
	svc := NewLearningLogTemplateService(repo, logCreator)
	return svc, repo, logCreator
}

// ============================================================
// Create テスト
// ============================================================

func TestLogTemplateService_Create(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:          1,
		Name:            "毎日のコーディング",
		DefaultTitle:    "コーディング練習",
		DefaultContent:  "## やったこと\n\n## 学んだこと",
		DefaultCategory: model.LogCategoryCoding,
		DefaultDuration: 60,
		IsDefault:       false,
	}

	repo.On("Create", template).Return(nil)

	err := svc.Create(template)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLogTemplateService_Create_WithDefaultFlag(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:          1,
		Name:            "デフォルトテンプレート",
		DefaultContent:  "本文",
		DefaultCategory: model.LogCategoryOther,
		IsDefault:       true,
	}

	repo.On("ClearDefaultFlag", uint(1)).Return(nil)
	repo.On("Create", template).Return(nil)

	err := svc.Create(template)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLogTemplateService_Create_EmptyName(t *testing.T) {
	svc, _, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:         1,
		Name:           "",
		DefaultContent: "本文",
	}

	err := svc.Create(template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "テンプレート名を入力してください")
}

func TestLogTemplateService_Create_NameTooLong(t *testing.T) {
	svc, _, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID: 1,
		Name:   strings.Repeat("あ", 101),
	}

	err := svc.Create(template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "100文字以下")
}

func TestLogTemplateService_Create_DefaultTitleTooLong(t *testing.T) {
	svc, _, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:       1,
		Name:         "テスト",
		DefaultTitle: strings.Repeat("あ", 201),
	}

	err := svc.Create(template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "デフォルトタイトルは200文字以下")
}

func TestLogTemplateService_Create_DefaultContentTooLong(t *testing.T) {
	svc, _, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:         1,
		Name:           "テスト",
		DefaultContent: strings.Repeat("あ", 50001),
	}

	err := svc.Create(template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "デフォルト本文は50000文字以下")
}

func TestLogTemplateService_Create_InvalidCategory(t *testing.T) {
	svc, _, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:          1,
		Name:            "テスト",
		DefaultCategory: "invalid",
	}

	err := svc.Create(template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "無効なカテゴリ")
}

func TestLogTemplateService_Create_InvalidDuration(t *testing.T) {
	svc, _, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:          1,
		Name:            "テスト",
		DefaultDuration: 1441,
	}

	err := svc.Create(template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "0〜1440分")
}

// ============================================================
// GetByID テスト
// ============================================================

func TestLogTemplateService_GetByID(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{UserID: 1, Name: "テスト"}
	template.ID = 1
	repo.On("FindByID", uint(1)).Return(template, nil)

	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "テスト", result.Name)
	repo.AssertExpectations(t)
}

func TestLogTemplateService_GetByID_NotFound(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	_, err := svc.GetByID(999, 1)
	assert.Error(t, err)
}

func TestLogTemplateService_GetByID_Forbidden(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{UserID: 2, Name: "他人のテンプレート"}
	template.ID = 1
	repo.On("FindByID", uint(1)).Return(template, nil)

	_, err := svc.GetByID(1, 1) // userID=1 が userID=2 のテンプレートにアクセス
	assert.ErrorIs(t, err, ErrForbidden)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestLogTemplateService_GetByUserID(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	templates := []model.LearningLogTemplate{
		{Name: "テンプレート1"},
		{Name: "テンプレート2"},
	}
	repo.On("FindByUserID", uint(1)).Return(templates, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestLogTemplateService_Update(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{
		UserID:          1,
		Name:            "旧名前",
		DefaultCategory: model.LogCategoryCoding,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "新名前", "", "", "", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "新名前", result.Name)
	repo.AssertExpectations(t)
}

func TestLogTemplateService_Update_Forbidden(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 2}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 1, "名前", "", "", "", nil, nil)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestLogTemplateService_Update_NameTooLong(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 1, Name: "旧名前"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 1, strings.Repeat("あ", 101), "", "", "", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "100文字以下")
}

func TestLogTemplateService_Update_DefaultTitleTooLong(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 1, Name: "テスト"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 1, "", strings.Repeat("あ", 201), "", "", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "デフォルトタイトルは200文字以下")
}

func TestLogTemplateService_Update_DefaultContentTooLong(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 1, Name: "テスト"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 1, "", "", strings.Repeat("あ", 50001), "", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "デフォルト本文は50000文字以下")
}

func TestLogTemplateService_Update_InvalidCategory(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 1, Name: "テスト"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 1, "", "", "", "invalid", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "無効なカテゴリ")
}

func TestLogTemplateService_Update_WithDefaultFlag(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 1, Name: "テスト", IsDefault: false}
	existing.ID = 1

	isDefault := true
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("ClearDefaultFlag", uint(1)).Return(nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "", "", "", "", nil, &isDefault)
	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
	repo.AssertExpectations(t)
}

func TestLogTemplateService_Update_InvalidDuration(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 1, Name: "テスト"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	dur := -1
	_, err := svc.Update(1, 1, "", "", "", "", &dur, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "0〜1440分")
}

// ============================================================
// Delete テスト
// ============================================================

func TestLogTemplateService_Delete(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLogTemplateService_Delete_Forbidden(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	existing := &model.LearningLogTemplate{UserID: 2}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}

// ============================================================
// UseTemplate テスト
// ============================================================

func TestLogTemplateService_UseTemplate(t *testing.T) {
	svc, repo, logCreator := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:          1,
		Name:            "毎日のコーディング",
		DefaultTitle:    "コーディング練習",
		DefaultContent:  "## やったこと",
		DefaultCategory: model.LogCategoryCoding,
		DefaultDuration: 60,
	}
	template.ID = 1

	repo.On("FindByID", uint(1)).Return(template, nil)
	logCreator.On("Create", mock.AnythingOfType("*model.LearningLog")).Return(nil)

	result, err := svc.UseTemplate(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "コーディング練習", result.Title)
	assert.Equal(t, model.LogCategoryCoding, result.Category)
	assert.Equal(t, 60, result.Duration)
	assert.Equal(t, model.LogSourceManual, result.Source)
	repo.AssertExpectations(t)
	logCreator.AssertExpectations(t)
}

func TestLogTemplateService_UseTemplate_FallbackTitle(t *testing.T) {
	svc, repo, logCreator := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:       1,
		Name:         "テンプレート名",
		DefaultTitle: "", // タイトル空
	}
	template.ID = 1

	repo.On("FindByID", uint(1)).Return(template, nil)
	logCreator.On("Create", mock.MatchedBy(func(log *model.LearningLog) bool {
		return log.Title == "テンプレート名" // Name がフォールバックとして使われる
	})).Return(nil)

	result, err := svc.UseTemplate(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "テンプレート名", result.Title)
}

func TestLogTemplateService_UseTemplate_FallbackCategory(t *testing.T) {
	svc, repo, logCreator := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:          1,
		Name:            "テスト",
		DefaultTitle:    "タイトル",
		DefaultCategory: "", // カテゴリ空
	}
	template.ID = 1

	repo.On("FindByID", uint(1)).Return(template, nil)
	logCreator.On("Create", mock.MatchedBy(func(log *model.LearningLog) bool {
		return log.Category == model.LogCategoryOther // other がデフォルト
	})).Return(nil)

	result, err := svc.UseTemplate(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, model.LogCategoryOther, result.Category)
}

func TestLogTemplateService_UseTemplate_Forbidden(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	template := &model.LearningLogTemplate{UserID: 2}
	template.ID = 1
	repo.On("FindByID", uint(1)).Return(template, nil)

	_, err := svc.UseTemplate(1, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestLogTemplateService_UseTemplate_CreateError(t *testing.T) {
	svc, repo, logCreator := newTestLogTemplateService()

	template := &model.LearningLogTemplate{
		UserID:       1,
		Name:         "テスト",
		DefaultTitle: "タイトル",
	}
	template.ID = 1

	repo.On("FindByID", uint(1)).Return(template, nil)
	logCreator.On("Create", mock.AnythingOfType("*model.LearningLog")).Return(errors.New("db error"))

	_, err := svc.UseTemplate(1, 1)
	assert.Error(t, err)
}

// ============================================================
// CountByUserID テスト
// ============================================================

func TestLogTemplateService_CountByUserID_Success(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	repo.On("CountByUserID", uint(1)).Return(int64(5), nil)

	count, err := svc.CountByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	repo.AssertExpectations(t)
}

func TestLogTemplateService_CountByUserID_Error(t *testing.T) {
	svc, repo, _ := newTestLogTemplateService()

	repo.On("CountByUserID", uint(1)).Return(int64(0), errors.New("db error"))

	_, err := svc.CountByUserID(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
