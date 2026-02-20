package handler

import (
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

func TestNoteLinkCreateLink_Success(t *testing.T) {
	h, svc := setupNoteLinkHandler()
	svc.On("CreateLink", uint(1), uint(2), uint(1)).Return(nil)

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestNoteLinkCreateLink_ServiceError(t *testing.T) {
	h, svc := setupNoteLinkHandler()
	svc.On("CreateLink", uint(1), uint(2), uint(1)).Return(service.ErrConflict)

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusConflict)
	svc.AssertExpectations(t)
}

func TestNoteLinkCreateLink_InvalidID(t *testing.T) {
	h, _ := setupNoteLinkHandler()

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/abc/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteLinkGetLinks_Success(t *testing.T) {
	h, svc := setupNoteLinkHandler()
	links := []model.NoteLink{{ID: 1, SourceNoteID: 1, TargetNoteID: 2}}
	svc.On("GetLinks", uint(1)).Return(links, nil)

	r := newRouter(1)
	r.GET("/notes/:id/links", h.GetLinks)
	w := doRequest(r, "GET", "/notes/1/links", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteLinkGetLinks_ServiceError(t *testing.T) {
	h, svc := setupNoteLinkHandler()
	svc.On("GetLinks", uint(1)).Return([]model.NoteLink(nil), service.ErrNotFound)

	r := newRouter(1)
	r.GET("/notes/:id/links", h.GetLinks)
	w := doRequest(r, "GET", "/notes/1/links", nil)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestNoteLinkGetBacklinks_Success(t *testing.T) {
	h, svc := setupNoteLinkHandler()
	backlinks := []model.NoteLink{{ID: 1, SourceNoteID: 2, TargetNoteID: 1}}
	svc.On("GetBacklinks", uint(1)).Return(backlinks, nil)

	r := newRouter(1)
	r.GET("/notes/:id/backlinks", h.GetBacklinks)
	w := doRequest(r, "GET", "/notes/1/backlinks", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteLinkGetBacklinks_ServiceError(t *testing.T) {
	h, svc := setupNoteLinkHandler()
	svc.On("GetBacklinks", uint(1)).Return([]model.NoteLink(nil), service.ErrNotFound)

	r := newRouter(1)
	r.GET("/notes/:id/backlinks", h.GetBacklinks)
	w := doRequest(r, "GET", "/notes/1/backlinks", nil)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestNoteLinkDeleteLink_Success(t *testing.T) {
	h, svc := setupNoteLinkHandler()
	svc.On("DeleteLink", uint(1), uint(2), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/notes/:id/links/:targetId", h.DeleteLink)
	w := doRequest(r, "DELETE", "/notes/1/links/2", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteLinkDeleteLink_ServiceError(t *testing.T) {
	h, svc := setupNoteLinkHandler()
	svc.On("DeleteLink", uint(1), uint(2), uint(1)).Return(service.ErrNotFound)

	r := newRouter(1)
	r.DELETE("/notes/:id/links/:targetId", h.DeleteLink)
	w := doRequest(r, "DELETE", "/notes/1/links/2", nil)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestNoteLinkCreateLink_BadRequest(t *testing.T) {
	h, _ := setupNoteLinkHandler()

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteLinkGetLinks_InvalidID(t *testing.T) {
	h, _ := setupNoteLinkHandler()

	r := newRouter(1)
	r.GET("/notes/:id/links", h.GetLinks)
	w := doRequest(r, "GET", "/notes/abc/links", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteLinkGetBacklinks_InvalidID(t *testing.T) {
	h, _ := setupNoteLinkHandler()

	r := newRouter(1)
	r.GET("/notes/:id/backlinks", h.GetBacklinks)
	w := doRequest(r, "GET", "/notes/abc/backlinks", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteLinkDeleteLink_InvalidID(t *testing.T) {
	h, _ := setupNoteLinkHandler()

	r := newRouter(1)
	r.DELETE("/notes/:id/links/:targetId", h.DeleteLink)
	w := doRequest(r, "DELETE", "/notes/abc/links/2", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteLinkDeleteLink_InvalidTargetID(t *testing.T) {
	h, _ := setupNoteLinkHandler()

	r := newRouter(1)
	r.DELETE("/notes/:id/links/:targetId", h.DeleteLink)
	w := doRequest(r, "DELETE", "/notes/1/links/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}
