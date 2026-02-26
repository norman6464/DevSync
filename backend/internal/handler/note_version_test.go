package handler

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MockNoteVersionService
// ============================================================

type MockNoteVersionService struct{ mock.Mock }

func (m *MockNoteVersionService) SaveVersion(noteID, userID uint) error {
	return m.Called(noteID, userID).Error(0)
}

func (m *MockNoteVersionService) GetVersions(noteID, userID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	args := m.Called(noteID, userID, limit, offset)
	return args.Get(0).([]model.NoteVersion), args.Get(1).(int64), args.Error(2)
}

func (m *MockNoteVersionService) GetVersion(noteID, versionID, userID uint) (*model.NoteVersion, error) {
	args := m.Called(noteID, versionID, userID)
	if v := args.Get(0); v != nil {
		return v.(*model.NoteVersion), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteVersionService) RestoreVersion(noteID, versionID, userID uint) (*model.Note, error) {
	args := m.Called(noteID, versionID, userID)
	if v := args.Get(0); v != nil {
		return v.(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupNoteVersionHandler() (*MockNoteVersionService, *NoteVersionHandler) {
	svc := new(MockNoteVersionService)
	h := NewNoteVersionHandler(svc)
	return svc, h
}

// ============================================================
// GetVersions テスト
// ============================================================

func TestNoteVersion_GetVersions_Success(t *testing.T) {
	svc, h := setupNoteVersionHandler()
	versions := []model.NoteVersion{
		{ID: 2, NoteID: 1, VersionNumber: 2, Title: "V2", CreatedAt: time.Now()},
		{ID: 1, NoteID: 1, VersionNumber: 1, Title: "V1", CreatedAt: time.Now()},
	}
	svc.On("GetVersions", uint(1), uint(1), 20, 0).Return(versions, int64(2), nil)

	r := newRouter(1)
	r.GET("/notes/:id/versions", h.GetVersions)
	w := doRequest(r, http.MethodGet, "/notes/1/versions", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestNoteVersion_GetVersions_InvalidID(t *testing.T) {
	_, h := setupNoteVersionHandler()

	r := newRouter(1)
	r.GET("/notes/:id/versions", h.GetVersions)
	w := doRequest(r, http.MethodGet, "/notes/abc/versions", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteVersion_GetVersions_ServiceError(t *testing.T) {
	svc, h := setupNoteVersionHandler()
	svc.On("GetVersions", uint(1), uint(1), 20, 0).Return([]model.NoteVersion{}, int64(0), errors.New("not found"))

	r := newRouter(1)
	r.GET("/notes/:id/versions", h.GetVersions)
	w := doRequest(r, http.MethodGet, "/notes/1/versions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// GetVersion テスト
// ============================================================

func TestNoteVersion_GetVersion_Success(t *testing.T) {
	svc, h := setupNoteVersionHandler()
	version := &model.NoteVersion{ID: 5, NoteID: 1, VersionNumber: 3, Title: "V3"}
	svc.On("GetVersion", uint(1), uint(5), uint(1)).Return(version, nil)

	r := newRouter(1)
	r.GET("/notes/:id/versions/:versionId", h.GetVersion)
	w := doRequest(r, http.MethodGet, "/notes/1/versions/5", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestNoteVersion_GetVersion_ServiceError(t *testing.T) {
	svc, h := setupNoteVersionHandler()
	svc.On("GetVersion", uint(1), uint(5), uint(1)).Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.GET("/notes/:id/versions/:versionId", h.GetVersion)
	w := doRequest(r, http.MethodGet, "/notes/1/versions/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// RestoreVersion テスト
// ============================================================

func TestNoteVersion_Restore_Success(t *testing.T) {
	svc, h := setupNoteVersionHandler()
	note := &model.Note{ID: 1, UserID: 1, Title: "Restored Title", Content: "Restored Content"}
	svc.On("RestoreVersion", uint(1), uint(5), uint(1)).Return(note, nil)

	r := newRouter(1)
	r.POST("/notes/:id/versions/:versionId/restore", h.RestoreVersion)
	w := doRequest(r, http.MethodPost, "/notes/1/versions/5/restore", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestNoteVersion_Restore_ServiceError(t *testing.T) {
	svc, h := setupNoteVersionHandler()
	svc.On("RestoreVersion", uint(1), uint(5), uint(1)).Return(nil, errors.New("forbidden"))

	r := newRouter(1)
	r.POST("/notes/:id/versions/:versionId/restore", h.RestoreVersion)
	w := doRequest(r, http.MethodPost, "/notes/1/versions/5/restore", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}
