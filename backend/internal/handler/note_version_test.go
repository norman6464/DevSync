package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockNoteVersionRepo は usecase/repository.NoteVersionRepository のモック（ctx 付き）。
type mockNoteVersionRepo struct{ mock.Mock }

func (m *mockNoteVersionRepo) Create(ctx context.Context, version *model.NoteVersion) error {
	return m.Called(ctx, version).Error(0)
}

func (m *mockNoteVersionRepo) FindByNoteID(ctx context.Context, noteID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	args := m.Called(ctx, noteID, limit, offset)
	v, _ := args.Get(0).([]model.NoteVersion)
	return v, args.Get(1).(int64), args.Error(2)
}

func (m *mockNoteVersionRepo) FindByID(ctx context.Context, id uint) (*model.NoteVersion, error) {
	args := m.Called(ctx, id)
	v, _ := args.Get(0).(*model.NoteVersion)
	return v, args.Error(1)
}

func (m *mockNoteVersionRepo) GetLatestVersionNumber(ctx context.Context, noteID uint) (int, error) {
	args := m.Called(ctx, noteID)
	return args.Int(0), args.Error(1)
}

// mockNoteUpdater は usecase/repository.NoteUpdater のモック。
type mockNoteUpdater struct{ mock.Mock }

func (m *mockNoteUpdater) Update(ctx context.Context, note *model.Note) error {
	return m.Called(ctx, note).Error(0)
}

// setupNoteVersionHandler は本物の usecase と port モックで NoteVersionHandler を組む。
func setupNoteVersionHandler() (*NoteVersionHandler, *mockNoteVersionRepo, *mockNoteReader, *mockNoteUpdater) {
	versions := new(mockNoteVersionRepo)
	notes := new(mockNoteReader)
	writer := new(mockNoteUpdater)
	h := NewNoteVersionHandler(
		usecase.NewListNoteVersionsUseCase(versions, notes),
		usecase.NewGetNoteVersionUseCase(versions, notes),
		usecase.NewRestoreNoteVersionUseCase(versions, notes, writer),
	)
	return h, versions, notes, writer
}

func TestNoteVersion_GetVersions_Success(t *testing.T) {
	h, versions, notes, _ := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	versions.On("FindByNoteID", mock.Anything, uint(1), 20, 0).
		Return([]model.NoteVersion{{ID: 5, NoteID: 1, VersionNumber: 2}}, int64(2), nil)

	r := newRouter(1)
	r.GET("/notes/:id/versions", h.GetVersions)
	w := doRequest(r, "GET", "/notes/1/versions", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(2), body["total"])
	versions.AssertExpectations(t)
}

// 所有者以外は 403 を返し、履歴を読まない。
func TestNoteVersion_GetVersions_Forbidden(t *testing.T) {
	h, versions, notes, _ := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)

	r := newRouter(1)
	r.GET("/notes/:id/versions", h.GetVersions)
	w := doRequest(r, "GET", "/notes/1/versions", nil)

	assertStatus(t, w, http.StatusForbidden)
	versions.AssertNotCalled(t, "FindByNoteID")
}

// ノートが存在しない場合は 500（移行前の挙動を維持している）。
func TestNoteVersion_GetVersions_NoteNotFound(t *testing.T) {
	h, versions, notes, _ := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.GET("/notes/:id/versions", h.GetVersions)
	w := doRequest(r, "GET", "/notes/1/versions", nil)

	assertStatus(t, w, http.StatusNotFound)
	versions.AssertNotCalled(t, "FindByNoteID")
}

func TestNoteVersion_GetVersion_Success(t *testing.T) {
	h, versions, notes, _ := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	versions.On("FindByID", mock.Anything, uint(5)).
		Return(&model.NoteVersion{ID: 5, NoteID: 1, VersionNumber: 2, Title: "v2"}, nil)

	r := newRouter(1)
	r.GET("/notes/:id/versions/:versionId", h.GetVersion)
	w := doRequest(r, "GET", "/notes/1/versions/5", nil)

	assertStatus(t, w, http.StatusOK)
	versions.AssertExpectations(t)
}

// バージョンが存在しなければ 404。
func TestNoteVersion_GetVersion_NotFound(t *testing.T) {
	h, versions, notes, _ := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	versions.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)

	r := newRouter(1)
	r.GET("/notes/:id/versions/:versionId", h.GetVersion)
	w := doRequest(r, "GET", "/notes/1/versions/5", nil)

	assertStatus(t, w, http.StatusNotFound)
}

// 他ノートのバージョン ID を指定した場合も 404。
func TestNoteVersion_GetVersion_OtherNotesVersion(t *testing.T) {
	h, versions, notes, _ := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	versions.On("FindByID", mock.Anything, uint(5)).
		Return(&model.NoteVersion{ID: 5, NoteID: 999}, nil)

	r := newRouter(1)
	r.GET("/notes/:id/versions/:versionId", h.GetVersion)
	w := doRequest(r, "GET", "/notes/1/versions/5", nil)

	assertStatus(t, w, http.StatusNotFound)
}

func TestNoteVersion_GetVersion_InvalidVersionID(t *testing.T) {
	h, _, _, _ := setupNoteVersionHandler()

	r := newRouter(1)
	r.GET("/notes/:id/versions/:versionId", h.GetVersion)
	w := doRequest(r, "GET", "/notes/1/versions/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteVersion_RestoreVersion_Success(t *testing.T) {
	h, versions, notes, writer := setupNoteVersionHandler()
	current := &model.Note{ID: 1, UserID: 1, Title: "いまの題", Content: "いまの本文"}
	notes.On("FindByID", mock.Anything, uint(1)).Return(current, nil)
	versions.On("FindByID", mock.Anything, uint(5)).
		Return(&model.NoteVersion{ID: 5, NoteID: 1, VersionNumber: 1, Title: "むかしの題", Content: "むかしの本文"}, nil)
	versions.On("GetLatestVersionNumber", mock.Anything, uint(1)).Return(3, nil)
	// 復元前の内容がバージョン 4 として残る
	versions.On("Create", mock.Anything, mock.MatchedBy(func(v *model.NoteVersion) bool {
		return v.NoteID == 1 && v.VersionNumber == 4 && v.Title == "いまの題" && v.Content == "いまの本文"
	})).Return(nil)
	writer.On("Update", mock.Anything, mock.MatchedBy(func(n *model.Note) bool {
		return n.Title == "むかしの題" && n.Content == "むかしの本文"
	})).Return(nil)

	r := newRouter(1)
	r.POST("/notes/:id/versions/:versionId/restore", h.RestoreVersion)
	w := doRequest(r, "POST", "/notes/1/versions/5/restore", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, "むかしの題", body["title"])
	versions.AssertExpectations(t)
	writer.AssertExpectations(t)
}

// 所有者以外の復元は 403 を返し、書き戻さない。
func TestNoteVersion_RestoreVersion_Forbidden(t *testing.T) {
	h, versions, notes, writer := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)

	r := newRouter(1)
	r.POST("/notes/:id/versions/:versionId/restore", h.RestoreVersion)
	w := doRequest(r, "POST", "/notes/1/versions/5/restore", nil)

	assertStatus(t, w, http.StatusForbidden)
	versions.AssertNotCalled(t, "Create")
	writer.AssertNotCalled(t, "Update")
}

// 他ノートのバージョンでの復元は 404 を返し、書き戻さない。
func TestNoteVersion_RestoreVersion_OtherNotesVersion(t *testing.T) {
	h, versions, notes, writer := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	versions.On("FindByID", mock.Anything, uint(5)).
		Return(&model.NoteVersion{ID: 5, NoteID: 999}, nil)

	r := newRouter(1)
	r.POST("/notes/:id/versions/:versionId/restore", h.RestoreVersion)
	w := doRequest(r, "POST", "/notes/1/versions/5/restore", nil)

	assertStatus(t, w, http.StatusNotFound)
	versions.AssertNotCalled(t, "Create")
	writer.AssertNotCalled(t, "Update")
}

// 復元前バージョンの保存に失敗したら書き戻さない。
func TestNoteVersion_RestoreVersion_SnapshotError(t *testing.T) {
	h, versions, notes, writer := setupNoteVersionHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	versions.On("FindByID", mock.Anything, uint(5)).
		Return(&model.NoteVersion{ID: 5, NoteID: 1}, nil)
	versions.On("GetLatestVersionNumber", mock.Anything, uint(1)).Return(1, nil)
	versions.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/notes/:id/versions/:versionId/restore", h.RestoreVersion)
	w := doRequest(r, "POST", "/notes/1/versions/5/restore", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	writer.AssertNotCalled(t, "Update")
}
