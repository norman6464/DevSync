package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

func TestNoteTemplateCreate_Success(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("Create", mock.AnythingOfType("*model.NoteTemplate")).Return(nil)

	r := newRouter(1)
	r.POST("/note-templates", h.Create)
	w := doRequest(r, "POST", "/note-templates", map[string]interface{}{
		"name":             "テスト",
		"content_template": "テンプレート内容",
	})

	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestNoteTemplateCreate_ServiceError(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("Create", mock.AnythingOfType("*model.NoteTemplate")).Return(service.ErrBadRequest)

	r := newRouter(1)
	r.POST("/note-templates", h.Create)
	w := doRequest(r, "POST", "/note-templates", map[string]interface{}{
		"name":             "テスト",
		"content_template": "テンプレート内容",
	})

	assertStatus(t, w, http.StatusBadRequest)
	svc.AssertExpectations(t)
}

func TestNoteTemplateGetByID_Success(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	tmpl := &model.NoteTemplate{ID: 1, UserID: 1, Name: "テスト"}
	svc.On("GetByID", uint(1)).Return(tmpl, nil)

	r := newRouter(1)
	r.GET("/note-templates/:id", h.GetByID)
	w := doRequest(r, "GET", "/note-templates/1", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteTemplateGetByID_NotFound(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("GetByID", uint(99)).Return(nil, service.ErrNotFound)

	r := newRouter(1)
	r.GET("/note-templates/:id", h.GetByID)
	w := doRequest(r, "GET", "/note-templates/99", nil)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestNoteTemplateGetByUserID_Success(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	templates := []model.NoteTemplate{{ID: 1, UserID: 1, Name: "テスト"}}
	svc.On("GetByUserID", uint(1)).Return(templates, nil)

	r := newRouter(1)
	r.GET("/note-templates", h.GetByUserID)
	w := doRequest(r, "GET", "/note-templates", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteTemplateGetDefault_Success(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	tmpl := &model.NoteTemplate{ID: 1, UserID: 1, Name: "デフォルト", IsDefault: true}
	svc.On("GetDefaultByUserID", uint(1)).Return(tmpl, nil)

	r := newRouter(1)
	r.GET("/note-templates/default", h.GetDefault)
	w := doRequest(r, "GET", "/note-templates/default", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteTemplateUpdate_Success(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	updated := &model.NoteTemplate{ID: 1, UserID: 1, Name: "新名"}
	svc.On("Update", uint(1), uint(1), "新名", "", "", "", "", (*bool)(nil)).Return(updated, nil)

	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)
	w := doRequest(r, "PUT", "/note-templates/1", map[string]interface{}{
		"name": "新名",
	})

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteTemplateUpdate_Forbidden(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("Update", uint(1), uint(1), "変更", "", "", "", "", (*bool)(nil)).Return(nil, service.ErrForbidden)

	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)
	w := doRequest(r, "PUT", "/note-templates/1", map[string]interface{}{
		"name": "変更",
	})

	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestNoteTemplateDelete_Success(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("Delete", uint(1), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/note-templates/:id", h.Delete)
	w := doRequest(r, "DELETE", "/note-templates/1", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteTemplateDelete_Forbidden(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("Delete", uint(1), uint(1)).Return(service.ErrForbidden)

	r := newRouter(1)
	r.DELETE("/note-templates/:id", h.Delete)
	w := doRequest(r, "DELETE", "/note-templates/1", nil)

	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestNoteTemplateUseTemplate_Success(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	note := &model.Note{UserID: 1, Title: "テンプレタイトル", Content: "テンプレ内容", Tags: "tag1"}
	svc.On("UseTemplate", uint(1), uint(1)).Return(note, nil)

	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)
	w := doRequest(r, "POST", "/note-templates/1/use", nil)

	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestNoteTemplateUseTemplate_Forbidden(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("UseTemplate", uint(1), uint(1)).Return(nil, service.ErrForbidden)

	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)
	w := doRequest(r, "POST", "/note-templates/1/use", nil)

	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestNoteTemplateUseTemplate_NotFound(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("UseTemplate", uint(99), uint(1)).Return(nil, service.ErrNotFound)

	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)
	w := doRequest(r, "POST", "/note-templates/99/use", nil)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestNoteTemplateGetByUserID_ServiceError(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("GetByUserID", uint(1)).Return([]model.NoteTemplate(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/note-templates", h.GetByUserID)
	w := doRequest(r, "GET", "/note-templates", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteTemplateGetDefault_ServiceError(t *testing.T) {
	h, svc := setupNoteTemplateHandler()
	svc.On("GetDefaultByUserID", uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/note-templates/default", h.GetDefault)
	w := doRequest(r, "GET", "/note-templates/default", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteTemplateCreate_BadRequest(t *testing.T) {
	h, _ := setupNoteTemplateHandler()

	r := newRouter(1)
	r.POST("/note-templates", h.Create)
	w := doRequest(r, "POST", "/note-templates", map[string]interface{}{})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteTemplateGetByID_InvalidID(t *testing.T) {
	h, _ := setupNoteTemplateHandler()

	r := newRouter(1)
	r.GET("/note-templates/:id", h.GetByID)
	w := doRequest(r, "GET", "/note-templates/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteTemplateUpdate_InvalidID(t *testing.T) {
	h, _ := setupNoteTemplateHandler()

	r := newRouter(1)
	r.PUT("/note-templates/:id", h.Update)
	w := doRequest(r, "PUT", "/note-templates/abc", map[string]interface{}{"name": "テスト"})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteTemplateDelete_InvalidID(t *testing.T) {
	h, _ := setupNoteTemplateHandler()

	r := newRouter(1)
	r.DELETE("/note-templates/:id", h.Delete)
	w := doRequest(r, "DELETE", "/note-templates/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteTemplateUseTemplate_InvalidID(t *testing.T) {
	h, _ := setupNoteTemplateHandler()

	r := newRouter(1)
	r.POST("/note-templates/:id/use", h.UseTemplate)
	w := doRequest(r, "POST", "/note-templates/abc/use", nil)

	assertStatus(t, w, http.StatusBadRequest)
}
