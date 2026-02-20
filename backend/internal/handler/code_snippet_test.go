package handler

import (
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCodeSnippetCreate_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	created := &model.CodeSnippet{Language: "Go", FileName: "main.go", Code: "package main"}
	svc.On("Create", mock.AnythingOfType("*model.CodeSnippet")).Return(created, nil)

	r := newRouter(1)
	r.POST("/posts/:id/snippets", h.Create)
	w := doRequest(r, "POST", "/posts/1/snippets", map[string]string{
		"language": "Go", "file_name": "main.go", "code": "package main",
	})

	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestCodeSnippetCreate_BadRequest(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.POST("/posts/:id/snippets", h.Create)
	w := doRequest(r, "POST", "/posts/1/snippets", map[string]string{})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetCreate_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("Create", mock.AnythingOfType("*model.CodeSnippet")).Return(nil, service.ErrBadRequest)

	r := newRouter(1)
	r.POST("/posts/:id/snippets", h.Create)
	w := doRequest(r, "POST", "/posts/1/snippets", map[string]string{
		"language": "Go", "code": "package main",
	})

	assertStatus(t, w, http.StatusBadRequest)
	svc.AssertExpectations(t)
}

func TestCodeSnippetCreate_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.POST("/posts/:id/snippets", h.Create)
	w := doRequest(r, "POST", "/posts/abc/snippets", map[string]string{
		"language": "Go", "code": "package main",
	})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetGetByPostID_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	snippets := []model.CodeSnippet{{Language: "Go", Code: "package main"}}
	svc.On("GetByPostID", uint(1)).Return(snippets, nil)

	r := newRouter(1)
	r.GET("/posts/:id/snippets", h.GetByPostID)
	w := doRequest(r, "GET", "/posts/1/snippets", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCodeSnippetGetByPostID_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("GetByPostID", uint(1)).Return([]model.CodeSnippet(nil), service.ErrNotFound)

	r := newRouter(1)
	r.GET("/posts/:id/snippets", h.GetByPostID)
	w := doRequest(r, "GET", "/posts/1/snippets", nil)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestCodeSnippetGetByID_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	snippets := []model.CodeSnippet{{Language: "Go"}}
	svc.On("GetByPostID", uint(5)).Return(snippets, nil)

	r := newRouter(1)
	r.GET("/snippets/:id", h.GetByID)
	w := doRequest(r, "GET", "/snippets/5", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCodeSnippetUpdate_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	updated := &model.CodeSnippet{Language: "Python", FileName: "app.py", Code: "print('hi')"}
	svc.On("Update", uint(1), uint(1), "Python", "app.py", "print('hi')").Return(updated, nil)

	r := newRouter(1)
	r.PUT("/snippets/:id", h.Update)
	w := doRequest(r, "PUT", "/snippets/1", map[string]string{
		"language": "Python", "file_name": "app.py", "code": "print('hi')",
	})

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCodeSnippetUpdate_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("Update", uint(1), uint(1), "Python", "", "").Return(nil, service.ErrForbidden)

	r := newRouter(1)
	r.PUT("/snippets/:id", h.Update)
	w := doRequest(r, "PUT", "/snippets/1", map[string]string{"language": "Python"})

	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestCodeSnippetDelete_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("Delete", uint(1), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/snippets/:id", h.Delete)
	w := doRequest(r, "DELETE", "/snippets/1", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCodeSnippetDelete_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("Delete", uint(1), uint(1)).Return(service.ErrForbidden)

	r := newRouter(1)
	r.DELETE("/snippets/:id", h.Delete)
	w := doRequest(r, "DELETE", "/snippets/1", nil)

	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestCodeSnippetGetComments_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	comments := []model.SnippetComment{{Content: "nice code", LineNumber: 1}}
	svc.On("GetComments", uint(1)).Return(comments, nil)

	r := newRouter(1)
	r.GET("/snippets/:id/comments", h.GetComments)
	w := doRequest(r, "GET", "/snippets/1/comments", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCodeSnippetGetComments_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("GetComments", uint(1)).Return([]model.SnippetComment(nil), service.ErrNotFound)

	r := newRouter(1)
	r.GET("/snippets/:id/comments", h.GetComments)
	w := doRequest(r, "GET", "/snippets/1/comments", nil)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestCodeSnippetCreateComment_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("CreateComment", mock.AnythingOfType("*model.SnippetComment")).Return(nil)

	r := newRouter(1)
	r.POST("/snippets/:id/comments", h.CreateComment)
	w := doRequest(r, "POST", "/snippets/1/comments", map[string]interface{}{
		"line_number": 5, "content": "great code",
	})

	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestCodeSnippetCreateComment_BadRequest(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.POST("/snippets/:id/comments", h.CreateComment)
	w := doRequest(r, "POST", "/snippets/1/comments", map[string]string{})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetCreateComment_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("CreateComment", mock.AnythingOfType("*model.SnippetComment")).Return(service.ErrBadRequest)

	r := newRouter(1)
	r.POST("/snippets/:id/comments", h.CreateComment)
	w := doRequest(r, "POST", "/snippets/1/comments", map[string]interface{}{
		"line_number": 5, "content": "great code",
	})

	assertStatus(t, w, http.StatusBadRequest)
	svc.AssertExpectations(t)
}

func TestCodeSnippetDeleteComment_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("DeleteComment", uint(1), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/snippets/:id/comments/:commentId", h.DeleteComment)
	w := doRequest(r, "DELETE", "/snippets/1/comments/1", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCodeSnippetDeleteComment_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("DeleteComment", uint(1), uint(1)).Return(service.ErrForbidden)

	r := newRouter(1)
	r.DELETE("/snippets/:id/comments/:commentId", h.DeleteComment)
	w := doRequest(r, "DELETE", "/snippets/1/comments/1", nil)

	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByUserLanguage テスト
// ============================================================

func TestCodeSnippet_GetByUserLanguage_Success(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	r := newRouter(1)
	r.GET("/snippets/language/:language", h.GetByUserLanguage)

	snippets := []model.CodeSnippet{{Language: "Go", FileName: "main.go"}}
	svc.On("GetByUserLanguage", uint(1), "Go").Return(snippets, nil)

	w := doRequest(r, http.MethodGet, "/snippets/language/Go", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCodeSnippet_GetByUserLanguage_NilResult(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	r := newRouter(1)
	r.GET("/snippets/language/:language", h.GetByUserLanguage)

	svc.On("GetByUserLanguage", uint(1), "Rust").Return([]model.CodeSnippet(nil), nil)

	w := doRequest(r, http.MethodGet, "/snippets/language/Rust", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
	svc.AssertExpectations(t)
}

func TestCodeSnippet_GetByUserLanguage_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	r := newRouter(1)
	r.GET("/snippets/language/:language", h.GetByUserLanguage)

	svc.On("GetByUserLanguage", uint(1), "invalid").Return([]model.CodeSnippet(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/snippets/language/invalid", nil)
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestCodeSnippetGetByID_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.GET("/snippets/:id", h.GetByID)
	w := doRequest(r, "GET", "/snippets/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetGetByID_ServiceError(t *testing.T) {
	h, svc := setupCodeSnippetHandler()
	svc.On("GetByPostID", uint(1)).Return([]model.CodeSnippet(nil), service.ErrNotFound)

	r := newRouter(1)
	r.GET("/snippets/:id", h.GetByID)
	w := doRequest(r, "GET", "/snippets/1", nil)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestCodeSnippetGetByPostID_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.GET("/posts/:id/snippets", h.GetByPostID)
	w := doRequest(r, "GET", "/posts/abc/snippets", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetUpdate_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.PUT("/snippets/:id", h.Update)
	w := doRequest(r, "PUT", "/snippets/abc", map[string]string{"language": "Go"})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetDelete_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.DELETE("/snippets/:id", h.Delete)
	w := doRequest(r, "DELETE", "/snippets/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetGetComments_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.GET("/snippets/:id/comments", h.GetComments)
	w := doRequest(r, "GET", "/snippets/abc/comments", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetCreateComment_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.POST("/snippets/:id/comments", h.CreateComment)
	w := doRequest(r, "POST", "/snippets/abc/comments", map[string]interface{}{
		"line_number": 5, "content": "test",
	})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetDeleteComment_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetHandler()

	r := newRouter(1)
	r.DELETE("/snippets/:id/comments/:commentId", h.DeleteComment)
	w := doRequest(r, "DELETE", "/snippets/1/comments/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}
