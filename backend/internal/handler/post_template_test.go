package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockPostTemplateRepo は usecase/repository.PostTemplateRepository のモック（ctx 付き）。
type mockPostTemplateRepo struct{ mock.Mock }

func (m *mockPostTemplateRepo) Create(ctx context.Context, template *model.PostTemplate) error {
	return m.Called(ctx, template).Error(0)
}

func (m *mockPostTemplateRepo) FindByID(ctx context.Context, id uint) (*model.PostTemplate, error) {
	args := m.Called(ctx, id)
	t, _ := args.Get(0).(*model.PostTemplate)
	return t, args.Error(1)
}

func (m *mockPostTemplateRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostTemplate, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	t, _ := args.Get(0).([]model.PostTemplate)
	return t, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostTemplateRepo) Update(ctx context.Context, template *model.PostTemplate) error {
	return m.Called(ctx, template).Error(0)
}

func (m *mockPostTemplateRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// setupPostTemplateHandler は本物の usecase と port モックで PostTemplateHandler を組む。
func setupPostTemplateHandler() (*PostTemplateHandler, *mockPostTemplateRepo) {
	repo := new(mockPostTemplateRepo)
	h := NewPostTemplateHandler(
		usecase.NewCreatePostTemplateUseCase(repo),
		usecase.NewGetPostTemplateUseCase(repo),
		usecase.NewListPostTemplatesUseCase(repo),
		usecase.NewUpdatePostTemplateUseCase(repo),
		usecase.NewDeletePostTemplateUseCase(repo),
	)
	return h, repo
}

// ownedPostTemplate は認証ユーザー（userID=1）が所有するテンプレートを返す。
func ownedPostTemplate() *model.PostTemplate {
	return &model.PostTemplate{UserID: 1, Name: "テンプレ", ContentTemplate: "本文"}
}

// ============================================================
// Create テスト
// ============================================================

func TestPostTemplate_Create_Success(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.PostTemplate")).Return(nil)

	r := newRouter(1)
	r.POST("/post-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/post-templates", map[string]interface{}{
		"name":             "日報テンプレート",
		"title_template":   "日報: {{date}}",
		"content_template": "## 今日の学び\n\n## 明日の予定",
	})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

func TestPostTemplate_Create_ValidationError(t *testing.T) {
	h, _ := setupPostTemplateHandler()

	r := newRouter(1)
	r.POST("/post-templates", h.Create)

	// name と content_template は required
	w := doRequest(r, http.MethodPost, "/post-templates", map[string]interface{}{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTemplate_Create_InvalidJSON(t *testing.T) {
	h, _ := setupPostTemplateHandler()

	r := newRouter(1)
	r.POST("/post-templates", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/post-templates", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTemplate_Create_ServiceError(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.PostTemplate")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/post-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/post-templates", map[string]interface{}{
		"name":             "テスト",
		"content_template": "内容",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// GetMyTemplates テスト
// ============================================================

func TestPostTemplate_GetMyTemplates_Success(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	templates := []model.PostTemplate{
		{Name: "日報"},
		{Name: "週報"},
	}
	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).Return(templates, int64(2), nil)

	r := newRouter(1)
	r.GET("/post-templates", h.GetMyTemplates)

	w := doRequest(r, http.MethodGet, "/post-templates", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostTemplate_GetMyTemplates_Empty(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).Return([]model.PostTemplate{}, int64(0), nil)

	r := newRouter(1)
	r.GET("/post-templates", h.GetMyTemplates)

	w := doRequest(r, http.MethodGet, "/post-templates", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostTemplate_GetMyTemplates_ServiceError(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).Return([]model.PostTemplate(nil), int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/post-templates", h.GetMyTemplates)

	w := doRequest(r, http.MethodGet, "/post-templates", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestPostTemplate_GetByID_Success(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	// 所有権チェックを通すため認証ユーザー（userID=1）所有のテンプレートを返す
	tmpl := &model.PostTemplate{UserID: 1, Name: "テスト", ContentTemplate: "内容"}
	repo.On("FindByID", mock.Anything, uint(1)).Return(tmpl, nil)

	r := newRouter(1)
	r.GET("/post-templates/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/post-templates/1", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostTemplate_GetByID_NotFound(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByID", mock.Anything, uint(99)).Return((*model.PostTemplate)(nil), errors.New("not found"))

	r := newRouter(1)
	r.GET("/post-templates/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/post-templates/99", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestPostTemplate_GetByID_InvalidID(t *testing.T) {
	h, _ := setupPostTemplateHandler()

	r := newRouter(1)
	r.GET("/post-templates/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/post-templates/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// Update テスト
// ============================================================

func TestPostTemplate_Update_Success(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedPostTemplate(), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.PostTemplate")).Return(nil)

	r := newRouter(1)
	r.PUT("/post-templates/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/post-templates/1", map[string]interface{}{
		"name": "更新済み",
	})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostTemplate_Update_InvalidID(t *testing.T) {
	h, _ := setupPostTemplateHandler()

	r := newRouter(1)
	r.PUT("/post-templates/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/post-templates/abc", map[string]interface{}{
		"name": "テスト",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTemplate_Update_InvalidJSON(t *testing.T) {
	h, _ := setupPostTemplateHandler()

	r := newRouter(1)
	r.PUT("/post-templates/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/post-templates/1", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTemplate_Update_ServiceError(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedPostTemplate(), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.PostTemplate")).Return(errors.New("db error"))

	r := newRouter(1)
	r.PUT("/post-templates/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/post-templates/1", map[string]interface{}{
		"name": "更新",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestPostTemplate_Delete_Success(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedPostTemplate(), nil)
	repo.On("Delete", mock.Anything, uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/post-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/post-templates/1", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostTemplate_Delete_InvalidID(t *testing.T) {
	h, _ := setupPostTemplateHandler()

	r := newRouter(1)
	r.DELETE("/post-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/post-templates/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTemplate_Delete_ServiceError(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByID", mock.Anything, uint(99)).Return((*model.PostTemplate)(nil), errors.New("not found"))

	r := newRouter(1)
	r.DELETE("/post-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/post-templates/99", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// 他人のテンプレートは取得できない（403）。旧テストは service をモックしていたため
// この分岐が handler 経由で検証されていなかった。
func TestPostTemplate_GetByID_Forbidden(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(&model.PostTemplate{UserID: 99, Name: "他人の"}, nil)

	r := newRouter(1)
	r.GET("/post-templates/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/post-templates/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertExpectations(t)
}

// 他人のテンプレートは削除できない（403）。
func TestPostTemplate_Delete_Forbidden(t *testing.T) {
	h, repo := setupPostTemplateHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(&model.PostTemplate{UserID: 99, Name: "他人の"}, nil)

	r := newRouter(1)
	r.DELETE("/post-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/post-templates/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Delete")
	repo.AssertExpectations(t)
}
