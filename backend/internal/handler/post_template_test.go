package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Create テスト
// ============================================================

func TestPostTemplate_Create_Success(t *testing.T) {
	h, svc := setupPostTemplateHandler()
	svc.On("Create", mock.AnythingOfType("*model.PostTemplate")).Return(nil)

	r := newRouter(1)
	r.POST("/post-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/post-templates", map[string]interface{}{
		"name":             "日報テンプレート",
		"title_template":   "日報: {{date}}",
		"content_template": "## 今日の学び\n\n## 明日の予定",
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
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
	h, svc := setupPostTemplateHandler()
	svc.On("Create", mock.AnythingOfType("*model.PostTemplate")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/post-templates", h.Create)

	w := doRequest(r, http.MethodPost, "/post-templates", map[string]interface{}{
		"name":             "テスト",
		"content_template": "内容",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetMyTemplates テスト
// ============================================================

func TestPostTemplate_GetMyTemplates_Success(t *testing.T) {
	h, svc := setupPostTemplateHandler()
	templates := []model.PostTemplate{
		{Name: "日報"},
		{Name: "週報"},
	}
	svc.On("GetByUserID", uint(1), 20, 0).Return(templates, int64(2), nil)

	r := newRouter(1)
	r.GET("/post-templates", h.GetMyTemplates)

	w := doRequest(r, http.MethodGet, "/post-templates", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostTemplate_GetMyTemplates_Empty(t *testing.T) {
	h, svc := setupPostTemplateHandler()
	svc.On("GetByUserID", uint(1), 20, 0).Return([]model.PostTemplate{}, int64(0), nil)

	r := newRouter(1)
	r.GET("/post-templates", h.GetMyTemplates)

	w := doRequest(r, http.MethodGet, "/post-templates", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostTemplate_GetMyTemplates_ServiceError(t *testing.T) {
	h, svc := setupPostTemplateHandler()
	svc.On("GetByUserID", uint(1), 20, 0).Return([]model.PostTemplate(nil), int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/post-templates", h.GetMyTemplates)

	w := doRequest(r, http.MethodGet, "/post-templates", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestPostTemplate_GetByID_Success(t *testing.T) {
	h, svc := setupPostTemplateHandler()
	tmpl := &model.PostTemplate{Name: "テスト", ContentTemplate: "内容"}
	svc.On("GetByID", uint(1), uint(1)).Return(tmpl, nil)

	r := newRouter(1)
	r.GET("/post-templates/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/post-templates/1", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostTemplate_GetByID_NotFound(t *testing.T) {
	h, svc := setupPostTemplateHandler()
	svc.On("GetByID", uint(99), uint(1)).Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.GET("/post-templates/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/post-templates/99", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
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
	h, svc := setupPostTemplateHandler()
	updated := &model.PostTemplate{Name: "更新済み"}
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.PostTemplate")).Return(updated, nil)

	r := newRouter(1)
	r.PUT("/post-templates/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/post-templates/1", map[string]interface{}{
		"name": "更新済み",
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
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
	h, svc := setupPostTemplateHandler()
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.PostTemplate")).Return(nil, errors.New("forbidden"))

	r := newRouter(1)
	r.PUT("/post-templates/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/post-templates/1", map[string]interface{}{
		"name": "更新",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestPostTemplate_Delete_Success(t *testing.T) {
	h, svc := setupPostTemplateHandler()
	svc.On("Delete", uint(1), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/post-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/post-templates/1", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostTemplate_Delete_InvalidID(t *testing.T) {
	h, _ := setupPostTemplateHandler()

	r := newRouter(1)
	r.DELETE("/post-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/post-templates/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTemplate_Delete_ServiceError(t *testing.T) {
	h, svc := setupPostTemplateHandler()
	svc.On("Delete", uint(99), uint(1)).Return(errors.New("not found"))

	r := newRouter(1)
	r.DELETE("/post-templates/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/post-templates/99", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
