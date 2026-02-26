package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestPostTemplateService はPostTemplateServiceのテスト用インスタンスを生成するヘルパー。
func newTestPostTemplateService() (*PostTemplateService, *MockPostTemplateRepository) {
	repo := new(MockPostTemplateRepository)
	svc := NewPostTemplateService(repo)
	return svc, repo
}

// ============================================================
// テンプレート作成テスト
// ============================================================

func TestPostTemplateCreate_Success(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	tmpl := &model.PostTemplate{
		UserID:          1,
		Name:            "日報テンプレート",
		TitleTemplate:   "日報: {{date}}",
		ContentTemplate: "## 今日やったこと\n\n## 明日やること\n",
	}
	repo.On("Create", tmpl).Return(nil)

	err := svc.Create(tmpl)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostTemplateCreate_EmptyName(t *testing.T) {
	svc, _ := newTestPostTemplateService()

	tmpl := &model.PostTemplate{
		UserID:          1,
		Name:            "",
		ContentTemplate: "内容",
	}

	err := svc.Create(tmpl)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "テンプレート名")
}

func TestPostTemplateCreate_EmptyContent(t *testing.T) {
	svc, _ := newTestPostTemplateService()

	tmpl := &model.PostTemplate{
		UserID: 1,
		Name:   "テスト",
	}

	err := svc.Create(tmpl)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "テンプレート内容")
}

func TestPostTemplateCreate_RepoError(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	tmpl := &model.PostTemplate{
		UserID:          1,
		Name:            "テスト",
		ContentTemplate: "内容",
	}
	repo.On("Create", tmpl).Return(errors.New("db error"))

	err := svc.Create(tmpl)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// テンプレート取得テスト
// ============================================================

func TestPostTemplateGetByID_Success(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	tmpl := &model.PostTemplate{Name: "日報", UserID: 1}
	tmpl.ID = 1
	repo.On("FindByID", uint(1)).Return(tmpl, nil)

	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "日報", result.Name)
	repo.AssertExpectations(t)
}

func TestPostTemplateGetByID_NotFound(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	repo.On("FindByID", uint(999)).Return((*model.PostTemplate)(nil), errors.New("not found"))

	result, err := svc.GetByID(999, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPostTemplateGetByID_Forbidden(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	tmpl := &model.PostTemplate{UserID: 1}
	tmpl.ID = 1
	repo.On("FindByID", uint(1)).Return(tmpl, nil)

	result, err := svc.GetByID(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
}

// ============================================================
// テンプレート一覧テスト
// ============================================================

func TestPostTemplateGetByUserID_Success(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	templates := []model.PostTemplate{
		{Name: "日報"},
		{Name: "週報"},
	}
	repo.On("FindByUserID", uint(1), 20, 0).Return(templates, int64(2), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestPostTemplateGetByUserID_Empty(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	repo.On("FindByUserID", uint(1), 20, 0).Return([]model.PostTemplate{}, int64(0), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
}

// ============================================================
// テンプレート更新テスト
// ============================================================

func TestPostTemplateUpdate_Success(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	existing := &model.PostTemplate{Name: "旧名前", ContentTemplate: "旧内容", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.PostTemplate{Name: "新名前"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "新名前", result.Name)
	assert.Equal(t, "旧内容", result.ContentTemplate)
	repo.AssertExpectations(t)
}

func TestPostTemplateUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	existing := &model.PostTemplate{UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.PostTemplate{Name: "新名前"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
}

func TestPostTemplateUpdate_NotFound(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	repo.On("FindByID", uint(999)).Return((*model.PostTemplate)(nil), errors.New("not found"))

	result, err := svc.Update(999, 1, &model.PostTemplate{Name: "新名前"})
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// テンプレート削除テスト
// ============================================================

func TestPostTemplateDelete_Success(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	existing := &model.PostTemplate{UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostTemplateDelete_Forbidden(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	existing := &model.PostTemplate{UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestPostTemplateDelete_NotFound(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	repo.On("FindByID", uint(999)).Return((*model.PostTemplate)(nil), errors.New("not found"))

	err := svc.Delete(999, 1)
	assert.Error(t, err)
}

func TestPostTemplateUpdate_RepoError(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	existing := &model.PostTemplate{UserID: 1, Name: "テスト", ContentTemplate: "内容"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.Anything).Return(errors.New("db error"))

	_, err := svc.Update(1, 1, &model.PostTemplate{Name: "更新名"})
	assert.Error(t, err)
}

func TestPostTemplateUpdate_PartialFields(t *testing.T) {
	svc, repo := newTestPostTemplateService()

	existing := &model.PostTemplate{UserID: 1, Name: "元の名前", TitleTemplate: "元のタイトル", ContentTemplate: "元の内容"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.MatchedBy(func(t *model.PostTemplate) bool {
		return t.Name == "新しい名前" && t.TitleTemplate == "元のタイトル" && t.ContentTemplate == "元の内容"
	})).Return(nil)

	result, err := svc.Update(1, 1, &model.PostTemplate{Name: "新しい名前"})
	assert.NoError(t, err)
	assert.Equal(t, "新しい名前", result.Name)
	assert.Equal(t, "元のタイトル", result.TitleTemplate)
}
