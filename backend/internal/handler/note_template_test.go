package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockNoteTemplateRepo は usecase/repository.NoteTemplateRepository のモック（ctx 付き）。
type mockNoteTemplateRepo struct{ mock.Mock }

func (m *mockNoteTemplateRepo) Create(ctx context.Context, template *model.NoteTemplate) error {
	return m.Called(ctx, template).Error(0)
}
func (m *mockNoteTemplateRepo) Update(ctx context.Context, template *model.NoteTemplate) error {
	return m.Called(ctx, template).Error(0)
}
func (m *mockNoteTemplateRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteTemplateRepo) FindByID(ctx context.Context, id uint) (*model.NoteTemplate, error) {
	args := m.Called(ctx, id)
	t, _ := args.Get(0).(*model.NoteTemplate)
	return t, args.Error(1)
}
func (m *mockNoteTemplateRepo) FindByUserID(ctx context.Context, userID uint) ([]model.NoteTemplate, error) {
	args := m.Called(ctx, userID)
	t, _ := args.Get(0).([]model.NoteTemplate)
	return t, args.Error(1)
}
func (m *mockNoteTemplateRepo) FindDefaultByUserID(ctx context.Context, userID uint) (*model.NoteTemplate, error) {
	args := m.Called(ctx, userID)
	t, _ := args.Get(0).(*model.NoteTemplate)
	return t, args.Error(1)
}
func (m *mockNoteTemplateRepo) ClearDefaultFlag(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *mockNoteTemplateRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// noteTemplatePorts は NoteTemplateHandler が使う port モックの束。
type noteTemplatePorts struct {
	Templates *mockNoteTemplateRepo
	Notes     *mockNoteRepo
}

// newTestNoteTemplateHandler は本物の usecase に port モックを注入した NoteTemplateHandler を生成する。
func newTestNoteTemplateHandler() (*NoteTemplateHandler, noteTemplatePorts) {
	templates := new(mockNoteTemplateRepo)
	notes := new(mockNoteRepo)
	h := NewNoteTemplateHandler(
		usecase.NewCreateNoteTemplateUseCase(templates),
		usecase.NewGetNoteTemplateUseCase(templates),
		usecase.NewListNoteTemplatesUseCase(templates),
		usecase.NewGetDefaultNoteTemplateUseCase(templates),
		usecase.NewUpdateNoteTemplateUseCase(templates),
		usecase.NewDeleteNoteTemplateUseCase(templates),
		usecase.NewCreateNoteFromTemplateUseCase(templates, usecase.NewCreateNoteUseCase(notes)),
		usecase.NewCountNoteTemplatesUseCase(templates),
	)
	return h, noteTemplatePorts{Templates: templates, Notes: notes}
}

// templateOwnedBy は指定ユーザーが所有するテンプレートを返すテスト用ヘルパー。
func templateOwnedBy(id, userID uint) *model.NoteTemplate {
	return &model.NoteTemplate{
		ID: id, UserID: userID,
		Name: "既存テンプレート", ContentTemplate: "## 本文", Description: "説明",
		DefaultTitle: "既定タイトル", DefaultTags: "go,test",
	}
}

// ============================================================
// Create
// ============================================================

func TestNoteTemplateHandler_Create(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates", h.Create)

	p.Templates.On("Create", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).Return(nil)

	w := doRequest(r, http.MethodPost, "/note-templates", map[string]interface{}{
		"name": "テスト", "content_template": "テンプレート内容",
	})
	assertStatus(t, w, http.StatusCreated)
	p.Templates.AssertExpectations(t)
}

// 前後の空白は落としてから保存する。
func TestNoteTemplateHandler_Create_TrimsWhitespace(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates", h.Create)

	p.Templates.On("Create", mock.Anything, mock.MatchedBy(func(tmpl *model.NoteTemplate) bool {
		return tmpl.Name == "テスト" && tmpl.ContentTemplate == "本文" &&
			tmpl.Description == "説明" && tmpl.DefaultTitle == "既定"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/note-templates", map[string]interface{}{
		"name": "  テスト  ", "content_template": " 本文 ",
		"description": "  説明 ", "default_title": " 既定 ",
	})
	assertStatus(t, w, http.StatusCreated)
	p.Templates.AssertExpectations(t)
}

// デフォルト指定つきの作成は、先に既存のデフォルト指定を外してから書き込む。
func TestNoteTemplateHandler_Create_WithDefaultFlag(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates", h.Create)

	p.Templates.On("ClearDefaultFlag", mock.Anything, uint(1)).Return(nil)
	p.Templates.On("Create", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).Return(nil)

	w := doRequest(r, http.MethodPost, "/note-templates", map[string]interface{}{
		"name": "既定", "content_template": "本文", "is_default": true,
	})
	assertStatus(t, w, http.StatusCreated)
	p.Templates.AssertExpectations(t)
}

// デフォルト指定の解除に失敗したら作成しない。
func TestNoteTemplateHandler_Create_ClearDefaultFlagError(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates", h.Create)

	p.Templates.On("ClearDefaultFlag", mock.Anything, uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/note-templates", map[string]interface{}{
		"name": "既定", "content_template": "本文", "is_default": true,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// name / content_template は DTO の binding で必須。
func TestNoteTemplateHandler_Create_BadRequest(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/note-templates", map[string]interface{}{})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 空白のみの名前は binding を通るため、usecase の検証で 400 になる。
func TestNoteTemplateHandler_Create_ValidationError(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/note-templates", map[string]interface{}{
		"name": "   ", "content_template": "本文",
	})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 説明が上限超過なら 400。
func TestNoteTemplateHandler_Create_DescriptionTooLong(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/note-templates", map[string]interface{}{
		"name": "テスト", "content_template": "本文", "description": strings.Repeat("あ", 501),
	})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestNoteTemplateHandler_Create_RepositoryError(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates", h.Create)

	p.Templates.On("Create", mock.Anything, mock.AnythingOfType("*model.NoteTemplate")).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/note-templates", map[string]interface{}{
		"name": "テスト", "content_template": "本文",
	})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// GetByID
// ============================================================

func TestNoteTemplateHandler_GetByID(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/:id", h.GetByID)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodGet, "/note-templates/1", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "既存テンプレート")
	p.Templates.AssertExpectations(t)
}

func TestNoteTemplateHandler_GetByID_Forbidden(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/:id", h.GetByID)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodGet, "/note-templates/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Templates.AssertExpectations(t)
}

// 不在のテンプレートは 500 になる（移行前から変わらない挙動）。
func TestNoteTemplateHandler_GetByID_NotFoundIs500(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/:id", h.GetByID)

	p.Templates.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/note-templates/99", nil)
	assertStatus(t, w, http.StatusNotFound)
	p.Templates.AssertExpectations(t)
}

func TestNoteTemplateHandler_GetByID_InvalidID(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/note-templates/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// GetByUserID
// ============================================================

func TestNoteTemplateHandler_GetByUserID(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates", h.GetByUserID)

	p.Templates.On("FindByUserID", mock.Anything, uint(1)).
		Return([]model.NoteTemplate{*templateOwnedBy(1, 1)}, nil)

	w := doRequest(r, http.MethodGet, "/note-templates", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "既存テンプレート")
	p.Templates.AssertExpectations(t)
}

// 0 件でも null ではなく空配列を返す。
func TestNoteTemplateHandler_GetByUserID_Empty(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates", h.GetByUserID)

	p.Templates.On("FindByUserID", mock.Anything, uint(1)).Return([]model.NoteTemplate(nil), nil)

	w := doRequest(r, http.MethodGet, "/note-templates", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestNoteTemplateHandler_GetByUserID_RepositoryError(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates", h.GetByUserID)

	p.Templates.On("FindByUserID", mock.Anything, uint(1)).
		Return([]model.NoteTemplate(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/note-templates", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// GetDefault
// ============================================================

func TestNoteTemplateHandler_GetDefault(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/default", h.GetDefault)

	tmpl := templateOwnedBy(1, 1)
	tmpl.IsDefault = true
	p.Templates.On("FindDefaultByUserID", mock.Anything, uint(1)).Return(tmpl, nil)

	w := doRequest(r, http.MethodGet, "/note-templates/default", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"is_default":true`)
	p.Templates.AssertExpectations(t)
}

// デフォルト未設定は 500 になる（移行前から変わらない挙動）。
func TestNoteTemplateHandler_GetDefault_NotSetIs500(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/default", h.GetDefault)

	p.Templates.On("FindDefaultByUserID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/note-templates/default", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	p.Templates.AssertExpectations(t)
}

func TestNoteTemplateHandler_GetDefault_RepositoryError(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/default", h.GetDefault)

	p.Templates.On("FindDefaultByUserID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/note-templates/default", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// Update
// ============================================================

func TestNoteTemplateHandler_Update(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 1), nil)
	p.Templates.On("Update", mock.Anything, mock.MatchedBy(func(tmpl *model.NoteTemplate) bool {
		// 指定しなかったフィールドは据え置かれる。
		return tmpl.Name == "新名" && tmpl.ContentTemplate == "## 本文" && tmpl.DefaultTitle == "既定タイトル"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/note-templates/1", map[string]interface{}{"name": "新名"})
	assertStatus(t, w, http.StatusOK)
	p.Templates.AssertExpectations(t)
}

// デフォルト指定つきの更新は、書き込み前に既存のデフォルト指定を外す。
func TestNoteTemplateHandler_Update_WithDefaultFlag(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 1), nil)
	p.Templates.On("ClearDefaultFlag", mock.Anything, uint(1)).Return(nil)
	p.Templates.On("Update", mock.Anything, mock.MatchedBy(func(tmpl *model.NoteTemplate) bool {
		return tmpl.IsDefault
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/note-templates/1", map[string]interface{}{"is_default": true})
	assertStatus(t, w, http.StatusOK)
	p.Templates.AssertExpectations(t)
}

func TestNoteTemplateHandler_Update_Forbidden(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodPut, "/note-templates/1", map[string]interface{}{"name": "変更"})
	assertStatus(t, w, http.StatusForbidden)
	p.Templates.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 不在のテンプレートの更新は 500 になる（移行前から変わらない挙動）。
func TestNoteTemplateHandler_Update_NotFoundIs500(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/note-templates/99", map[string]interface{}{"name": "変更"})
	assertStatus(t, w, http.StatusNotFound)
	p.Templates.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestNoteTemplateHandler_Update_ValidationError(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodPut, "/note-templates/1", map[string]interface{}{
		"name": strings.Repeat("あ", 101),
	})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestNoteTemplateHandler_Update_InvalidID(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/note-templates/abc", map[string]interface{}{"name": "テスト"})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// Delete
// ============================================================

func TestNoteTemplateHandler_Delete(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.DELETE("/note-templates/:id", h.Delete)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 1), nil)
	p.Templates.On("Delete", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/note-templates/1", nil)
	assertStatus(t, w, http.StatusOK)
	p.Templates.AssertExpectations(t)
}

func TestNoteTemplateHandler_Delete_Forbidden(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.DELETE("/note-templates/:id", h.Delete)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodDelete, "/note-templates/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Templates.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// 不在のテンプレートの削除は 500 になる（移行前から変わらない挙動）。
func TestNoteTemplateHandler_Delete_NotFoundIs500(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.DELETE("/note-templates/:id", h.Delete)

	p.Templates.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/note-templates/99", nil)
	assertStatus(t, w, http.StatusNotFound)
	p.Templates.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestNoteTemplateHandler_Delete_InvalidID(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.DELETE("/note-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/note-templates/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// UseTemplate
// ============================================================

func TestNoteTemplateHandler_UseTemplate(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 1), nil)
	p.Notes.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Note) bool {
		return n.UserID == 1 && n.Title == "既定タイトル" && n.Content == "## 本文" && n.Tags == "go,test"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/note-templates/1/use", nil)
	assertStatus(t, w, http.StatusCreated)
	p.Templates.AssertExpectations(t)
	p.Notes.AssertExpectations(t)
}

// デフォルトタイトルが空なら既定のノート名を使う。
func TestNoteTemplateHandler_UseTemplate_FallbackTitle(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)

	tmpl := templateOwnedBy(1, 1)
	tmpl.DefaultTitle = ""
	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(tmpl, nil)
	p.Notes.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Note) bool {
		return n.Title == "新しいノート"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/note-templates/1/use", nil)
	assertStatus(t, w, http.StatusCreated)
	p.Notes.AssertExpectations(t)
}

func TestNoteTemplateHandler_UseTemplate_Forbidden(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodPost, "/note-templates/1/use", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Notes.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 不在のテンプレートの利用は 500 になる（移行前から変わらない挙動）。
func TestNoteTemplateHandler_UseTemplate_NotFoundIs500(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)

	p.Templates.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/note-templates/99/use", nil)
	assertStatus(t, w, http.StatusNotFound)
	p.Notes.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestNoteTemplateHandler_UseTemplate_NoteCreateError(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(templateOwnedBy(1, 1), nil)
	p.Notes.On("Create", mock.Anything, mock.AnythingOfType("*model.Note")).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/note-templates/1/use", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteTemplateHandler_UseTemplate_InvalidID(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)

	w := doRequest(r, http.MethodPost, "/note-templates/abc/use", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// GetMyCount
// ============================================================

func TestNoteTemplateHandler_GetMyCount(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/my/count", h.GetMyCount)

	p.Templates.On("CountByUserID", mock.Anything, uint(1)).Return(int64(5), nil)

	w := doRequest(r, http.MethodGet, "/note-templates/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":5`)
	p.Templates.AssertExpectations(t)
}

func TestNoteTemplateHandler_GetMyCount_RepositoryError(t *testing.T) {
	h, p := newTestNoteTemplateHandler()
	r := newRouter(1)
	r.GET("/note-templates/my/count", h.GetMyCount)

	p.Templates.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/note-templates/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	p.Templates.AssertExpectations(t)
}
