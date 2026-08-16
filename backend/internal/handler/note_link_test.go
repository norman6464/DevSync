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

// mockNoteLinkRepo は usecase/repository.NoteLinkRepository のモック（ctx 付き）。
type mockNoteLinkRepo struct{ mock.Mock }

func (m *mockNoteLinkRepo) Create(ctx context.Context, link *model.NoteLink) error {
	return m.Called(ctx, link).Error(0)
}

func (m *mockNoteLinkRepo) FindBySourceNoteID(ctx context.Context, sourceNoteID uint) ([]model.NoteLink, error) {
	args := m.Called(ctx, sourceNoteID)
	l, _ := args.Get(0).([]model.NoteLink)
	return l, args.Error(1)
}

func (m *mockNoteLinkRepo) FindByTargetNoteID(ctx context.Context, targetNoteID uint) ([]model.NoteLink, error) {
	args := m.Called(ctx, targetNoteID)
	l, _ := args.Get(0).([]model.NoteLink)
	return l, args.Error(1)
}

func (m *mockNoteLinkRepo) Delete(ctx context.Context, sourceNoteID, targetNoteID uint) error {
	return m.Called(ctx, sourceNoteID, targetNoteID).Error(0)
}

func (m *mockNoteLinkRepo) Exists(ctx context.Context, sourceNoteID, targetNoteID uint) (bool, error) {
	args := m.Called(ctx, sourceNoteID, targetNoteID)
	return args.Bool(0), args.Error(1)
}

func (m *mockNoteLinkRepo) CountBySourceNoteID(ctx context.Context, noteID uint) (int64, error) {
	args := m.Called(ctx, noteID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockNoteLinkRepo) CountByTargetNoteID(ctx context.Context, noteID uint) (int64, error) {
	args := m.Called(ctx, noteID)
	return args.Get(0).(int64), args.Error(1)
}

// mockNoteReader は usecase/repository.NoteReader のモック。
type mockNoteReader struct{ mock.Mock }

func (m *mockNoteReader) FindByID(ctx context.Context, id uint) (*model.Note, error) {
	args := m.Called(ctx, id)
	n, _ := args.Get(0).(*model.Note)
	return n, args.Error(1)
}

// setupNoteLinkHandler は本物の usecase と port モックで NoteLinkHandler を組む。
func setupNoteLinkHandler() (*NoteLinkHandler, *mockNoteLinkRepo, *mockNoteReader) {
	links := new(mockNoteLinkRepo)
	notes := new(mockNoteReader)
	h := NewNoteLinkHandler(
		usecase.NewCreateNoteLinkUseCase(links, notes),
		usecase.NewListNoteLinksUseCase(links),
		usecase.NewListNoteBacklinksUseCase(links),
		usecase.NewGetNoteLinkStatsUseCase(links, notes),
		usecase.NewDeleteNoteLinkUseCase(links, notes),
	)
	return h, links, notes
}

// ownedNote は所有者が userID=1 のノートを返す。
func ownedNote(id uint) *model.Note {
	return &model.Note{ID: id, UserID: 1}
}

func TestNoteLinkCreateLink_Success(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	notes.On("FindByID", mock.Anything, uint(2)).Return(ownedNote(2), nil)
	links.On("Exists", mock.Anything, uint(1), uint(2)).Return(false, nil)
	links.On("Create", mock.Anything, mock.MatchedBy(func(l *model.NoteLink) bool {
		return l.SourceNoteID == 1 && l.TargetNoteID == 2
	})).Return(nil)

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusCreated)
	links.AssertExpectations(t)
	notes.AssertExpectations(t)
}

// 既に同じリンクがあれば 400 を返し、作成しない。
func TestNoteLinkCreateLink_Duplicate(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	notes.On("FindByID", mock.Anything, uint(2)).Return(ownedNote(2), nil)
	links.On("Exists", mock.Anything, uint(1), uint(2)).Return(true, nil)

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusBadRequest)
	links.AssertNotCalled(t, "Create")
}

// 同じノートへのリンクは 400 を返し、ノートも読まない。
func TestNoteLinkCreateLink_SelfLink(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{"target_note_id": 1})

	assertStatus(t, w, http.StatusBadRequest)
	notes.AssertNotCalled(t, "FindByID")
	links.AssertNotCalled(t, "Create")
}

// ソースノートの所有者でなければ 403 を返す。
func TestNoteLinkCreateLink_Forbidden(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusForbidden)
	links.AssertNotCalled(t, "Create")
}

// ソースノートが存在しなければ 404 を返す。
func TestNoteLinkCreateLink_SourceNotFound(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusNotFound)
	links.AssertNotCalled(t, "Create")
}

// リンク先ノートが存在しなければ 404 を返す。
func TestNoteLinkCreateLink_TargetNotFound(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	notes.On("FindByID", mock.Anything, uint(2)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/1/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusNotFound)
	links.AssertNotCalled(t, "Create")
}

func TestNoteLinkCreateLink_InvalidID(t *testing.T) {
	h, _, _ := setupNoteLinkHandler()

	r := newRouter(1)
	r.POST("/notes/:id/links", h.CreateLink)
	w := doRequest(r, "POST", "/notes/abc/links", map[string]interface{}{"target_note_id": 2})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteLinkGetLinks_Success(t *testing.T) {
	h, links, _ := setupNoteLinkHandler()
	links.On("FindBySourceNoteID", mock.Anything, uint(1)).
		Return([]model.NoteLink{{ID: 1, SourceNoteID: 1, TargetNoteID: 2}}, nil)

	r := newRouter(1)
	r.GET("/notes/:id/links", h.GetLinks)
	w := doRequest(r, "GET", "/notes/1/links", nil)

	assertStatus(t, w, http.StatusOK)
	links.AssertExpectations(t)
}

func TestNoteLinkGetLinks_RepoError(t *testing.T) {
	h, links, _ := setupNoteLinkHandler()
	links.On("FindBySourceNoteID", mock.Anything, uint(1)).
		Return([]model.NoteLink(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/notes/:id/links", h.GetLinks)
	w := doRequest(r, "GET", "/notes/1/links", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteLinkGetBacklinks_Success(t *testing.T) {
	h, links, _ := setupNoteLinkHandler()
	links.On("FindByTargetNoteID", mock.Anything, uint(2)).
		Return([]model.NoteLink{{ID: 1, SourceNoteID: 1, TargetNoteID: 2}}, nil)

	r := newRouter(1)
	r.GET("/notes/:id/backlinks", h.GetBacklinks)
	w := doRequest(r, "GET", "/notes/2/backlinks", nil)

	assertStatus(t, w, http.StatusOK)
	links.AssertExpectations(t)
}

func TestNoteLinkGetLinkStats_Success(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	links.On("CountBySourceNoteID", mock.Anything, uint(1)).Return(int64(3), nil)
	links.On("CountByTargetNoteID", mock.Anything, uint(1)).Return(int64(5), nil)

	r := newRouter(1)
	r.GET("/notes/:id/link-stats", h.GetLinkStats)
	w := doRequest(r, "GET", "/notes/1/link-stats", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(3), body["forward_link_count"])
	assert.Equal(t, float64(5), body["backlink_count"])
	links.AssertExpectations(t)
}

// 統計はノートの所有者以外には 403 を返す。
func TestNoteLinkGetLinkStats_Forbidden(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)

	r := newRouter(1)
	r.GET("/notes/:id/link-stats", h.GetLinkStats)
	w := doRequest(r, "GET", "/notes/1/link-stats", nil)

	assertStatus(t, w, http.StatusForbidden)
	links.AssertNotCalled(t, "CountBySourceNoteID")
}

// 統計はノートが存在しない場合、作成・削除と違い 500 を返す（移行前の挙動を維持している）。
func TestNoteLinkGetLinkStats_NotFound(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.GET("/notes/:id/link-stats", h.GetLinkStats)
	w := doRequest(r, "GET", "/notes/1/link-stats", nil)

	assertStatus(t, w, http.StatusNotFound)
	links.AssertNotCalled(t, "CountBySourceNoteID")
}

func TestNoteLinkDeleteLink_Success(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNote(1), nil)
	links.On("Delete", mock.Anything, uint(1), uint(2)).Return(nil)

	r := newRouter(1)
	r.DELETE("/notes/:id/links/:targetId", h.DeleteLink)
	w := doRequest(r, "DELETE", "/notes/1/links/2", nil)

	assertStatus(t, w, http.StatusOK)
	links.AssertExpectations(t)
}

// 所有者以外の削除は 403 を返し、削除しない。
func TestNoteLinkDeleteLink_Forbidden(t *testing.T) {
	h, links, notes := setupNoteLinkHandler()
	notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)

	r := newRouter(1)
	r.DELETE("/notes/:id/links/:targetId", h.DeleteLink)
	w := doRequest(r, "DELETE", "/notes/1/links/2", nil)

	assertStatus(t, w, http.StatusForbidden)
	links.AssertNotCalled(t, "Delete")
}

func TestNoteLinkDeleteLink_InvalidTargetID(t *testing.T) {
	h, _, _ := setupNoteLinkHandler()

	r := newRouter(1)
	r.DELETE("/notes/:id/links/:targetId", h.DeleteLink)
	w := doRequest(r, "DELETE", "/notes/1/links/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}
