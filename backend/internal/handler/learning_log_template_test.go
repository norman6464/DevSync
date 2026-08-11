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

// mockLearningLogTemplateRepo は usecase/repository.LearningLogTemplateRepository のモック（ctx 付き）。
type mockLearningLogTemplateRepo struct{ mock.Mock }

func (m *mockLearningLogTemplateRepo) Create(ctx context.Context, template *model.LearningLogTemplate) error {
	return m.Called(ctx, template).Error(0)
}
func (m *mockLearningLogTemplateRepo) Update(ctx context.Context, template *model.LearningLogTemplate) error {
	return m.Called(ctx, template).Error(0)
}
func (m *mockLearningLogTemplateRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockLearningLogTemplateRepo) FindByID(ctx context.Context, id uint) (*model.LearningLogTemplate, error) {
	args := m.Called(ctx, id)
	t, _ := args.Get(0).(*model.LearningLogTemplate)
	return t, args.Error(1)
}
func (m *mockLearningLogTemplateRepo) FindByUserID(ctx context.Context, userID uint) ([]model.LearningLogTemplate, error) {
	args := m.Called(ctx, userID)
	t, _ := args.Get(0).([]model.LearningLogTemplate)
	return t, args.Error(1)
}
func (m *mockLearningLogTemplateRepo) FindDefaultByUserID(ctx context.Context, userID uint) (*model.LearningLogTemplate, error) {
	args := m.Called(ctx, userID)
	t, _ := args.Get(0).(*model.LearningLogTemplate)
	return t, args.Error(1)
}
func (m *mockLearningLogTemplateRepo) ClearDefaultFlag(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *mockLearningLogTemplateRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// learningLogTemplatePorts は LearningLogTemplateHandler が使う port モックの束。
type learningLogTemplatePorts struct {
	Templates *mockLearningLogTemplateRepo
	Logs      *mockLearningLogRepo
}

// newTestLearningLogTemplateHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestLearningLogTemplateHandler() (*LearningLogTemplateHandler, learningLogTemplatePorts) {
	templates := new(mockLearningLogTemplateRepo)
	logs := new(mockLearningLogRepo)
	goals := new(mockLearningGoalLinker)
	h := NewLearningLogTemplateHandler(
		usecase.NewCreateLearningLogTemplateUseCase(templates),
		usecase.NewGetLearningLogTemplateUseCase(templates),
		usecase.NewListLearningLogTemplatesUseCase(templates),
		usecase.NewGetDefaultLearningLogTemplateUseCase(templates),
		usecase.NewUpdateLearningLogTemplateUseCase(templates),
		usecase.NewDeleteLearningLogTemplateUseCase(templates),
		usecase.NewCreateLearningLogFromTemplateUseCase(templates, usecase.NewCreateLearningLogUseCase(logs, goals)),
		usecase.NewCountLearningLogTemplatesUseCase(templates),
	)
	return h, learningLogTemplatePorts{Templates: templates, Logs: logs}
}

// logTemplateOwnedBy は指定ユーザーが所有するテンプレートを返すテスト用ヘルパー。
func logTemplateOwnedBy(id, userID uint) *model.LearningLogTemplate {
	return &model.LearningLogTemplate{
		ID: id, UserID: userID,
		Name: "既存テンプレート", DefaultTitle: "既定タイトル", DefaultContent: "既定本文",
		DefaultCategory: model.LogCategoryCoding, DefaultDuration: 60,
	}
}

// ============================================================
// Create
// ============================================================

func TestLearningLogTemplateHandler_Create(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	p.Templates.On("Create", mock.Anything, mock.MatchedBy(func(tm *model.LearningLogTemplate) bool {
		return tm.UserID == 1 && tm.Name == "週報"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{"name": "週報"})
	assertStatus(t, w, http.StatusCreated)
	p.Templates.AssertExpectations(t)
}

// 前後の空白は落としてから保存する。
func TestLearningLogTemplateHandler_Create_TrimsWhitespace(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	p.Templates.On("Create", mock.Anything, mock.MatchedBy(func(tm *model.LearningLogTemplate) bool {
		return tm.Name == "週報" && tm.DefaultTitle == "題" && tm.DefaultContent == "本文"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{
		"name": " 週報 ", "default_title": "  題 ", "default_content": " 本文  ",
	})
	assertStatus(t, w, http.StatusCreated)
	p.Templates.AssertExpectations(t)
}

// デフォルト指定つきの作成は、先に既存の指定を外してから書き込む。
func TestLearningLogTemplateHandler_Create_WithDefaultFlag(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	p.Templates.On("ClearDefaultFlag", mock.Anything, uint(1)).Return(nil)
	p.Templates.On("Create", mock.Anything, mock.Anything).Return(nil)

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{
		"name": "既定", "is_default": true,
	})
	assertStatus(t, w, http.StatusCreated)
	p.Templates.AssertExpectations(t)
}

func TestLearningLogTemplateHandler_Create_ClearDefaultFlagError(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	p.Templates.On("ClearDefaultFlag", mock.Anything, uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{
		"name": "既定", "is_default": true,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// name は DTO の binding で必須。
func TestLearningLogTemplateHandler_Create_BadRequest(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_Create_InvalidJSON(t *testing.T) {
	h, _ := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/log-templates", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

// 空白のみの名前は binding を通るため、usecase の検証で 400 になる。
func TestLearningLogTemplateHandler_Create_ValidationError(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{"name": "   "})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_Create_InvalidCategory(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{
		"name": "週報", "default_category": "unknown",
	})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "無効なカテゴリです")
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_Create_InvalidDuration(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{
		"name": "週報", "default_duration": 1441,
	})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_Create_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates", h.Create)

	p.Templates.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/log-templates", map[string]interface{}{"name": "週報"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// GetByID / 一覧 / デフォルト
// ============================================================

func TestLearningLogTemplateHandler_GetByID(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates/:id", h.GetByID)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodGet, "/log-templates/1", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "既存テンプレート")
}

func TestLearningLogTemplateHandler_GetByID_Forbidden(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates/:id", h.GetByID)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodGet, "/log-templates/1", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// 不在のテンプレートは 500 になる（移行前から変わらない挙動）。
func TestLearningLogTemplateHandler_GetByID_NotFoundIs500(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates/:id", h.GetByID)

	p.Templates.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/log-templates/99", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestLearningLogTemplateHandler_GetByID_InvalidID(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/log-templates/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_GetByUserID(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates", h.GetByUserID)

	p.Templates.On("FindByUserID", mock.Anything, uint(1)).
		Return([]model.LearningLogTemplate{*logTemplateOwnedBy(1, 1)}, nil)

	w := doRequest(r, http.MethodGet, "/log-templates", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "既存テンプレート")
}

// 0 件でも null ではなく空配列を返す。
func TestLearningLogTemplateHandler_GetByUserID_Empty(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates", h.GetByUserID)

	p.Templates.On("FindByUserID", mock.Anything, uint(1)).Return([]model.LearningLogTemplate(nil), nil)

	w := doRequest(r, http.MethodGet, "/log-templates", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestLearningLogTemplateHandler_GetByUserID_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates", h.GetByUserID)

	p.Templates.On("FindByUserID", mock.Anything, uint(1)).
		Return([]model.LearningLogTemplate(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/log-templates", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestLearningLogTemplateHandler_GetDefault(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates/default", h.GetDefault)

	tmpl := logTemplateOwnedBy(1, 1)
	tmpl.IsDefault = true
	p.Templates.On("FindDefaultByUserID", mock.Anything, uint(1)).Return(tmpl, nil)

	w := doRequest(r, http.MethodGet, "/log-templates/default", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"is_default":true`)
}

// デフォルト未設定は 500 になる（移行前から変わらない挙動）。
func TestLearningLogTemplateHandler_GetDefault_NotSetIs500(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates/default", h.GetDefault)

	p.Templates.On("FindDefaultByUserID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/log-templates/default", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// Update
// ============================================================

func TestLearningLogTemplateHandler_Update(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.PUT("/log-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 1), nil)
	p.Templates.On("Update", mock.Anything, mock.MatchedBy(func(tm *model.LearningLogTemplate) bool {
		// 指定しなかったフィールドは据え置かれる。
		return tm.Name == "新名" && tm.DefaultTitle == "既定タイトル" && tm.DefaultDuration == 60
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/log-templates/1", map[string]interface{}{"name": "新名"})
	assertStatus(t, w, http.StatusOK)
	p.Templates.AssertExpectations(t)
}

// デフォルト指定に true を渡したときだけ既存の指定を外す。
func TestLearningLogTemplateHandler_Update_WithDefaultFlag(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.PUT("/log-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 1), nil)
	p.Templates.On("ClearDefaultFlag", mock.Anything, uint(1)).Return(nil)
	p.Templates.On("Update", mock.Anything, mock.MatchedBy(func(tm *model.LearningLogTemplate) bool {
		return tm.IsDefault
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/log-templates/1", map[string]interface{}{"is_default": true})
	assertStatus(t, w, http.StatusOK)
	p.Templates.AssertExpectations(t)
}

func TestLearningLogTemplateHandler_Update_UnsetDefaultFlag(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.PUT("/log-templates/:id", h.Update)

	existing := logTemplateOwnedBy(1, 1)
	existing.IsDefault = true
	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
	p.Templates.On("Update", mock.Anything, mock.MatchedBy(func(tm *model.LearningLogTemplate) bool {
		return !tm.IsDefault
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/log-templates/1", map[string]interface{}{"is_default": false})
	assertStatus(t, w, http.StatusOK)
	p.Templates.AssertNotCalled(t, "ClearDefaultFlag", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_Update_Forbidden(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.PUT("/log-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodPut, "/log-templates/1", map[string]interface{}{"name": "新名"})
	assertStatus(t, w, http.StatusForbidden)
	p.Templates.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_Update_ValidationError(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.PUT("/log-templates/:id", h.Update)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodPut, "/log-templates/1", map[string]interface{}{
		"name": strings.Repeat("あ", 101),
	})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_Update_InvalidID(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.PUT("/log-templates/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/log-templates/abc", map[string]interface{}{"name": "新名"})
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// Delete
// ============================================================

func TestLearningLogTemplateHandler_Delete(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.DELETE("/log-templates/:id", h.Delete)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 1), nil)
	p.Templates.On("Delete", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/log-templates/1", nil)
	assertStatus(t, w, http.StatusOK)
	p.Templates.AssertExpectations(t)
}

func TestLearningLogTemplateHandler_Delete_Forbidden(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.DELETE("/log-templates/:id", h.Delete)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodDelete, "/log-templates/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Templates.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_Delete_InvalidID(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.DELETE("/log-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/log-templates/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// UseTemplate
// ============================================================

func TestLearningLogTemplateHandler_UseTemplate(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates/:id/use", h.UseTemplate)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 1), nil)
	p.Logs.On("Create", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
		return l.UserID == 1 && l.Title == "既定タイトル" && l.Content == "既定本文" &&
			l.Category == model.LogCategoryCoding && l.Duration == 60 && l.Source == model.LogSourceManual
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/log-templates/1/use", nil)
	assertStatus(t, w, http.StatusCreated)
	p.Logs.AssertExpectations(t)
}

// デフォルトタイトルが空ならテンプレート名を使う。
func TestLearningLogTemplateHandler_UseTemplate_FallbackTitle(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates/:id/use", h.UseTemplate)

	tmpl := logTemplateOwnedBy(1, 1)
	tmpl.DefaultTitle = ""
	tmpl.DefaultCategory = ""
	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(tmpl, nil)
	p.Logs.On("Create", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
		return l.Title == "既存テンプレート" && l.Category == model.LogCategoryOther
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/log-templates/1/use", nil)
	assertStatus(t, w, http.StatusCreated)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogTemplateHandler_UseTemplate_Forbidden(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates/:id/use", h.UseTemplate)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodPost, "/log-templates/1/use", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 不在のテンプレートの利用は 500 になる（移行前から変わらない挙動）。
func TestLearningLogTemplateHandler_UseTemplate_NotFoundIs500(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates/:id/use", h.UseTemplate)

	p.Templates.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/log-templates/99/use", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestLearningLogTemplateHandler_UseTemplate_LogCreateError(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates/:id/use", h.UseTemplate)

	p.Templates.On("FindByID", mock.Anything, uint(1)).Return(logTemplateOwnedBy(1, 1), nil)
	p.Logs.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/log-templates/1/use", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestLearningLogTemplateHandler_UseTemplate_InvalidID(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.POST("/log-templates/:id/use", h.UseTemplate)

	w := doRequest(r, http.MethodPost, "/log-templates/abc/use", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Templates.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// GetMyCount
// ============================================================

func TestLearningLogTemplateHandler_GetMyCount(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates/my/count", h.GetMyCount)

	p.Templates.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)

	w := doRequest(r, http.MethodGet, "/log-templates/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":3`)
	p.Templates.AssertExpectations(t)
}

func TestLearningLogTemplateHandler_GetMyCount_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogTemplateHandler()
	r := newRouter(1)
	r.GET("/log-templates/my/count", h.GetMyCount)

	p.Templates.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/log-templates/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}
