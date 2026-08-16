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

// mockNoteRepo は usecase/repository.NoteRepository のモック（ctx 付き）。
type mockNoteRepo struct{ mock.Mock }

func (m *mockNoteRepo) Create(ctx context.Context, note *model.Note) error {
	return m.Called(ctx, note).Error(0)
}
func (m *mockNoteRepo) Update(ctx context.Context, note *model.Note) error {
	return m.Called(ctx, note).Error(0)
}
func (m *mockNoteRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteRepo) FindByID(ctx context.Context, id uint) (*model.Note, error) {
	args := m.Called(ctx, id)
	n, _ := args.Get(0).(*model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(ctx, userID, page, limit)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) FindByFolderID(ctx context.Context, folderID, userID uint) ([]model.Note, error) {
	args := m.Called(ctx, folderID, userID)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) Search(ctx context.Context, userID uint, query string, limit, offset int) ([]model.Note, int64, error) {
	args := m.Called(ctx, userID, query, limit, offset)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Get(1).(int64), args.Error(2)
}
func (m *mockNoteRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockNoteRepo) ToggleFavorite(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteRepo) FindFavorites(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(ctx, userID, page, limit)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) CountFavoritesByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockNoteRepo) Archive(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteRepo) Unarchive(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteRepo) FindArchived(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(ctx, userID, page, limit)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) CountArchivedByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// newTestNoteHandler は本物の usecase に port モックを注入した NoteHandler を生成する。
func newTestNoteHandler() (*NoteHandler, *mockNoteRepo) {
	repo := new(mockNoteRepo)
	h := NewNoteHandler(
		usecase.NewCreateNoteUseCase(repo),
		usecase.NewGetNoteUseCase(repo),
		usecase.NewListNotesUseCase(repo),
		usecase.NewListNotesByFolderUseCase(repo),
		usecase.NewUpdateNoteUseCase(repo),
		usecase.NewDeleteNoteUseCase(repo),
		usecase.NewSearchNotesUseCase(repo),
		usecase.NewCountNotesUseCase(repo),
		usecase.NewToggleNoteFavoriteUseCase(repo),
		usecase.NewListFavoriteNotesUseCase(repo),
		usecase.NewCountFavoriteNotesUseCase(repo),
		usecase.NewArchiveNoteUseCase(repo),
		usecase.NewUnarchiveNoteUseCase(repo),
		usecase.NewListArchivedNotesUseCase(repo),
		usecase.NewCountArchivedNotesUseCase(repo),
		usecase.NewListNoteTagsUseCase(repo),
		usecase.NewExportNoteMarkdownUseCase(repo),
		usecase.NewDuplicateNoteUseCase(repo),
	)
	return h, repo
}

// noteOwnedBy は指定ユーザーが所有するノートを返すテスト用ヘルパー。
func noteOwnedBy(id, userID uint) *model.Note {
	return &model.Note{ID: id, UserID: userID, Title: "既存ノート", Content: "本文", Tags: "go,test"}
}

// ============================================================
// Create
// ============================================================

func TestNoteHandler_Create(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes", h.Create)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)

	w := doRequest(r, http.MethodPost, "/notes", map[string]interface{}{
		"title": "新しいノート", "content": "本文です",
	})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

func TestNoteHandler_Create_InvalidJSON(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/notes", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

// タイトルが空なら usecase の検証で 400 になり、リポジトリは呼ばれない。
func TestNoteHandler_Create_ValidationError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes", h.Create)

	w := doRequest(r, http.MethodPost, "/notes", map[string]interface{}{"title": "", "content": "本文"})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestNoteHandler_Create_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes", h.Create)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Note")).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/notes", map[string]interface{}{"title": "題", "content": "本文"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// GetByID
// ============================================================

func TestNoteHandler_GetByID(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodGet, "/notes/1", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, "既存ノート", body["title"])
}

func TestNoteHandler_GetByID_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/notes/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// 他人のノートは 403。
func TestNoteHandler_GetByID_Forbidden(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 999), nil)

	w := doRequest(r, http.MethodGet, "/notes/1", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// 不在は 404 にならず 500 になる（移行前からの挙動）。
func TestNoteHandler_GetByID_MissingReturnsInternalError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id", h.GetByID)

	// port は不在を (nil, nil) で表す。
	repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/notes/1", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ============================================================
// GetByUserID / GetMyCount
// ============================================================

func TestNoteHandler_GetByUserID(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 20).
		Return([]model.Note{{ID: 1}, {ID: 2}}, nil)
	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(2), nil)

	w := doRequest(r, http.MethodGet, "/notes", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_GetByUserID_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 20).
		Return([]model.Note(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notes", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteHandler_GetByUserID_CountError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 20).Return([]model.Note{}, nil)
	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notes", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteHandler_GetMyCount_Success(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/my/count", h.GetMyCount)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(7), nil)

	w := doRequest(r, http.MethodGet, "/notes/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, float64(7), parseJSON(t, w)["count"])
}

func TestNoteHandler_GetMyCount_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/my/count", h.GetMyCount)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notes/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// GetByFolderID
// ============================================================

func TestNoteHandler_GetByFolderID(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/folder/:folderId", h.GetByFolderID)

	repo.On("FindByFolderID", mock.Anything, uint(5), uint(1)).Return([]model.Note{{ID: 1}}, nil)

	w := doRequest(r, http.MethodGet, "/notes/folder/5", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_GetByFolderID_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/folder/:folderId", h.GetByFolderID)

	w := doRequest(r, http.MethodGet, "/notes/folder/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// Update
// ============================================================

func TestNoteHandler_Update(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)

	w := doRequest(r, http.MethodPut, "/notes/1", map[string]interface{}{"title": "更新後"})
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "更新後", parseJSON(t, w)["title"])
}

func TestNoteHandler_Update_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/notes/abc", map[string]interface{}{"title": "X"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_Update_InvalidJSON(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/notes/1", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

// 他人のノートの更新は 403。
func TestNoteHandler_Update_Forbidden(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 999), nil)

	w := doRequest(r, http.MethodPut, "/notes/1", map[string]interface{}{"title": "X"})
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 空白のみのタイトルは専用メッセージで 400。
func TestNoteHandler_Update_BlankTitle(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodPut, "/notes/1", map[string]interface{}{"title": "   "})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "タイトルは空白のみにできません")
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestNoteHandler_Update_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Note")).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/notes/1", map[string]interface{}{"title": "X"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// Delete
// ============================================================

func TestNoteHandler_Delete(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.DELETE("/notes/:id", h.Delete)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)
	repo.On("Delete", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/notes/1", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_Delete_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.DELETE("/notes/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/notes/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_Delete_Forbidden(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.DELETE("/notes/:id", h.Delete)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 999), nil)

	w := doRequest(r, http.MethodDelete, "/notes/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// ============================================================
// Search
// ============================================================

func TestNoteHandler_Search(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/search", h.Search)

	repo.On("Search", mock.Anything, uint(1), "go", 20, 0).
		Return([]model.Note{{ID: 1}}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/notes/search?q=go", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_Search_MissingQuery(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/search", h.Search)

	w := doRequest(r, http.MethodGet, "/notes/search", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestNoteHandler_Search_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/search", h.Search)

	repo.On("Search", mock.Anything, uint(1), "go", 20, 0).
		Return([]model.Note(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notes/search?q=go", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// お気に入り
// ============================================================

func TestNoteHandler_ToggleFavorite(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/favorite", h.ToggleFavorite)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)
	repo.On("ToggleFavorite", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/notes/1/favorite", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_ToggleFavorite_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/favorite", h.ToggleFavorite)

	w := doRequest(r, http.MethodPut, "/notes/abc/favorite", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// 他人のノートのお気に入り切替は 403。
func TestNoteHandler_ToggleFavorite_Forbidden(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/favorite", h.ToggleFavorite)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 999), nil)

	w := doRequest(r, http.MethodPut, "/notes/1/favorite", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "ToggleFavorite", mock.Anything, mock.Anything)
}

func TestNoteHandler_ToggleFavorite_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/favorite", h.ToggleFavorite)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)
	repo.On("ToggleFavorite", mock.Anything, uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/notes/1/favorite", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteHandler_GetFavorites_Success(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/favorites", h.GetFavorites)

	repo.On("FindFavorites", mock.Anything, uint(1), 1, 20).Return([]model.Note{{ID: 1}}, nil)
	repo.On("CountFavoritesByUserID", mock.Anything, uint(1)).Return(int64(1), nil)

	w := doRequest(r, http.MethodGet, "/notes/favorites", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_GetFavorites_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/favorites", h.GetFavorites)

	repo.On("FindFavorites", mock.Anything, uint(1), 1, 20).
		Return([]model.Note(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notes/favorites", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteHandler_GetFavorites_CountError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/favorites", h.GetFavorites)

	repo.On("FindFavorites", mock.Anything, uint(1), 1, 20).Return([]model.Note{}, nil)
	repo.On("CountFavoritesByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notes/favorites", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// アーカイブ
// ============================================================

func TestNoteHandler_Archive(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/archive", h.Archive)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)
	repo.On("Archive", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/notes/1/archive", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_Archive_Forbidden(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/archive", h.Archive)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 999), nil)

	w := doRequest(r, http.MethodPut, "/notes/1/archive", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Archive", mock.Anything, mock.Anything)
}

func TestNoteHandler_Unarchive(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/unarchive", h.Unarchive)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)
	repo.On("Unarchive", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/notes/1/unarchive", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_Unarchive_Forbidden(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/unarchive", h.Unarchive)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 999), nil)

	w := doRequest(r, http.MethodPut, "/notes/1/unarchive", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Unarchive", mock.Anything, mock.Anything)
}

func TestNoteHandler_GetArchived(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/archived", h.GetArchived)

	repo.On("FindArchived", mock.Anything, uint(1), 1, 20).Return([]model.Note{{ID: 1}}, nil)
	repo.On("CountArchivedByUserID", mock.Anything, uint(1)).Return(int64(1), nil)

	w := doRequest(r, http.MethodGet, "/notes/archived", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNoteHandler_GetArchived_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/archived", h.GetArchived)

	repo.On("FindArchived", mock.Anything, uint(1), 1, 20).
		Return([]model.Note(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notes/archived", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// 複製・エクスポート・タグ
// ============================================================

func TestNoteHandler_Duplicate(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes/:id/duplicate", h.Duplicate)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)

	w := doRequest(r, http.MethodPost, "/notes/1/duplicate", nil)
	assertStatus(t, w, http.StatusCreated)
	assert.Equal(t, "既存ノート (コピー)", parseJSON(t, w)["title"])
}

func TestNoteHandler_Duplicate_Forbidden(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes/:id/duplicate", h.Duplicate)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 999), nil)

	w := doRequest(r, http.MethodPost, "/notes/1/duplicate", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestNoteHandler_Export_Success(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id/export", h.Export)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodGet, "/notes/1/export", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "# 既存ノート")
	assert.Contains(t, w.Body.String(), "**Tags:** go,test")
}

func TestNoteHandler_Export_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id/export", h.Export)

	w := doRequest(r, http.MethodGet, "/notes/abc/export", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_Export_Forbidden(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id/export", h.Export)

	repo.On("FindByID", mock.Anything, uint(1)).Return(noteOwnedBy(1, 999), nil)

	w := doRequest(r, http.MethodGet, "/notes/1/export", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestNoteHandler_GetTags(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/tags", h.GetTags)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 1000).
		Return([]model.Note{{Tags: "go, test"}, {Tags: "go,web"}}, nil)

	w := doRequest(r, http.MethodGet, "/notes/tags", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "go")
	assert.Contains(t, w.Body.String(), "web")
}

func TestNoteHandler_GetTags_RepositoryError(t *testing.T) {
	h, repo := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/tags", h.GetTags)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 1000).
		Return([]model.Note(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notes/tags", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}
