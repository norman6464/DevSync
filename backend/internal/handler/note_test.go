package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockNoteService は NoteService のモック実装。
type MockNoteService struct {
	mock.Mock
}

func (m *MockNoteService) Create(note *model.Note) error {
	return m.Called(note).Error(0)
}

func (m *MockNoteService) GetByID(id uint) (*model.Note, error) {
	args := m.Called(id)
	if note := args.Get(0); note != nil {
		return note.(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteService) GetByUserID(userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Note), args.Error(1)
}

func (m *MockNoteService) GetByFolderID(folderID uint) ([]model.Note, error) {
	args := m.Called(folderID)
	return args.Get(0).([]model.Note), args.Error(1)
}

func (m *MockNoteService) Update(id, userID uint, title, content, tags string, folderID *uint) (*model.Note, error) {
	args := m.Called(id, userID, title, content, tags, folderID)
	if n := args.Get(0); n != nil {
		return n.(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

func (m *MockNoteService) Search(userID uint, query string, page, limit int) ([]model.Note, int64, error) {
	args := m.Called(userID, query, page, limit)
	return args.Get(0).([]model.Note), args.Get(1).(int64), args.Error(2)
}

func (m *MockNoteService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNoteService) ToggleFavorite(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

func (m *MockNoteService) GetFavorites(userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Note), args.Error(1)
}

func (m *MockNoteService) CountFavoritesByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNoteService) Archive(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

func (m *MockNoteService) Unarchive(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

func (m *MockNoteService) GetArchived(userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Note), args.Error(1)
}

func (m *MockNoteService) CountArchivedByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNoteService) Duplicate(id uint, userID uint) (*model.Note, error) {
	args := m.Called(id, userID)
	if note := args.Get(0); note != nil {
		return note.(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

// newTestNoteHandler はテスト用のNoteHandlerを生成する。
func newTestNoteHandler() (*NoteHandler, *MockNoteService) {
	mockService := new(MockNoteService)
	handler := NewNoteHandler(mockService)
	return handler, mockService
}

// ============================================================
// Create テスト
// ============================================================

func TestNoteHandler_Create(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes", h.Create)

	svc.On("Create", mock.AnythingOfType("*model.Note")).Return(nil)

	w := doRequest(r, "POST", "/notes", map[string]interface{}{
		"title":   "テストノート",
		"content": "これはテスト内容です",
		"tags":    "Go,TDD",
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestNoteHandler_GetByID(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id", h.GetByID)

	note := &model.Note{
		ID:      1,
		UserID:  1,
		Title:   "テストノート",
		Content: "内容",
	}

	svc.On("GetByID", uint(1)).Return(note, nil)

	w := doRequest(r, "GET", "/notes/1", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestNoteHandler_GetByUserID(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes", h.GetByUserID)

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "ノート1"},
		{ID: 2, UserID: 1, Title: "ノート2"},
	}

	svc.On("GetByUserID", uint(1), 1, 20).Return(notes, nil)
	svc.On("CountByUserID", uint(1)).Return(int64(2), nil)

	w := doRequest(r, "GET", "/notes?page=1&limit=20", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestNoteHandler_Update(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	updatedNote := &model.Note{
		ID:      1,
		UserID:  1,
		Title:   "更新後タイトル",
		Content: "更新後内容",
	}

	svc.On("Update", uint(1), uint(1), "更新後タイトル", "更新後内容", "", (*uint)(nil)).Return(updatedNote, nil)

	w := doRequest(r, "PUT", "/notes/1", map[string]interface{}{
		"title":   "更新後タイトル",
		"content": "更新後内容",
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteHandler_Delete(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.DELETE("/notes/:id", h.Delete)

	svc.On("Delete", uint(1), uint(1)).Return(nil)

	w := doRequest(r, "DELETE", "/notes/1", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// Search テスト
// ============================================================

func TestNoteHandler_Search(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/search", h.Search)

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "Go学習", Content: "Goを学習中"},
	}

	svc.On("Search", uint(1), "Go", 1, 20).Return(notes, int64(1), nil)

	w := doRequest(r, "GET", "/notes/search?q=Go&page=1&limit=20", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// ToggleFavorite テスト
// ============================================================

func TestNoteHandler_ToggleFavorite(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/favorite", h.ToggleFavorite)

	svc.On("ToggleFavorite", uint(1), uint(1)).Return(nil)

	w := doRequest(r, "PUT", "/notes/1/favorite", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByFolderID テスト
// ============================================================

func TestNoteHandler_GetByFolderID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.GET("/folders/:folderId/notes", h.GetByFolderID)

		notes := []model.Note{
			{ID: 1, UserID: 1, Title: "フォルダ内ノート1"},
			{ID: 2, UserID: 1, Title: "フォルダ内ノート2"},
		}

		svc.On("GetByFolderID", uint(5)).Return(notes, nil)

		w := doRequest(r, "GET", "/folders/5/notes", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("InvalidFolderID", func(t *testing.T) {
		h, _ := newTestNoteHandler()
		r := newRouter(1)
		r.GET("/folders/:folderId/notes", h.GetByFolderID)

		w := doRequest(r, "GET", "/folders/abc/notes", nil)
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("ServiceError", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.GET("/folders/:folderId/notes", h.GetByFolderID)

		svc.On("GetByFolderID", uint(5)).Return([]model.Note{}, errors.New("error"))

		w := doRequest(r, "GET", "/folders/5/notes", nil)
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})
}

// ============================================================
// Archive テスト
// ============================================================

func TestNoteHandler_Archive(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.PUT("/notes/:id/archive", h.Archive)

		svc.On("Archive", uint(1), uint(1)).Return(nil)

		w := doRequest(r, "PUT", "/notes/1/archive", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		h, _ := newTestNoteHandler()
		r := newRouter(1)
		r.PUT("/notes/:id/archive", h.Archive)

		w := doRequest(r, "PUT", "/notes/abc/archive", nil)
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("ServiceError", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.PUT("/notes/:id/archive", h.Archive)

		svc.On("Archive", uint(1), uint(1)).Return(errors.New("error"))

		w := doRequest(r, "PUT", "/notes/1/archive", nil)
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})
}

// ============================================================
// Unarchive テスト
// ============================================================

func TestNoteHandler_Unarchive(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.PUT("/notes/:id/unarchive", h.Unarchive)

		svc.On("Unarchive", uint(1), uint(1)).Return(nil)

		w := doRequest(r, "PUT", "/notes/1/unarchive", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		h, _ := newTestNoteHandler()
		r := newRouter(1)
		r.PUT("/notes/:id/unarchive", h.Unarchive)

		w := doRequest(r, "PUT", "/notes/abc/unarchive", nil)
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("ServiceError", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.PUT("/notes/:id/unarchive", h.Unarchive)

		svc.On("Unarchive", uint(1), uint(1)).Return(errors.New("error"))

		w := doRequest(r, "PUT", "/notes/1/unarchive", nil)
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})
}

// ============================================================
// GetArchived テスト
// ============================================================

func TestNoteHandler_GetArchived(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.GET("/notes/archived", h.GetArchived)

		notes := []model.Note{
			{ID: 3, UserID: 1, Title: "アーカイブノート1"},
		}

		svc.On("GetArchived", uint(1), 1, 20).Return(notes, nil)
		svc.On("CountArchivedByUserID", uint(1)).Return(int64(1), nil)

		w := doRequest(r, "GET", "/notes/archived?page=1&limit=20", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("GetArchivedError", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.GET("/notes/archived", h.GetArchived)

		svc.On("GetArchived", uint(1), 1, 20).Return([]model.Note{}, errors.New("error"))

		w := doRequest(r, "GET", "/notes/archived?page=1&limit=20", nil)
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})

	t.Run("CountError", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.GET("/notes/archived", h.GetArchived)

		notes := []model.Note{
			{ID: 3, UserID: 1, Title: "アーカイブノート1"},
		}

		svc.On("GetArchived", uint(1), 1, 20).Return(notes, nil)
		svc.On("CountArchivedByUserID", uint(1)).Return(int64(0), errors.New("error"))

		w := doRequest(r, "GET", "/notes/archived?page=1&limit=20", nil)
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})
}

// ============================================================
// Duplicate テスト
// ============================================================

func TestNoteHandler_Duplicate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.POST("/notes/:id/duplicate", h.Duplicate)

		duplicate := &model.Note{
			ID:      10,
			UserID:  1,
			Title:   "テストノート (コピー)",
			Content: "元の内容",
		}

		svc.On("Duplicate", uint(1), uint(1)).Return(duplicate, nil)

		w := doRequest(r, "POST", "/notes/1/duplicate", nil)
		assertStatus(t, w, http.StatusCreated)
		svc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		h, _ := newTestNoteHandler()
		r := newRouter(1)
		r.POST("/notes/:id/duplicate", h.Duplicate)

		w := doRequest(r, "POST", "/notes/abc/duplicate", nil)
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("ServiceError", func(t *testing.T) {
		h, svc := newTestNoteHandler()
		r := newRouter(1)
		r.POST("/notes/:id/duplicate", h.Duplicate)

		svc.On("Duplicate", uint(1), uint(1)).Return(nil, errors.New("error"))

		w := doRequest(r, "POST", "/notes/1/duplicate", nil)
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})
}

// ============================================================
// Create エラーパステスト
// ============================================================

func TestNoteHandler_Create_InvalidJSON(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes", h.Create)

	w := doRequestRaw(r, "POST", "/notes", "invalid json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_Create_ServiceError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.POST("/notes", h.Create)

	svc.On("Create", mock.AnythingOfType("*model.Note")).Return(errors.New("error"))

	w := doRequest(r, "POST", "/notes", map[string]interface{}{"title": "テスト", "content": "内容"})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByID エラーパステスト
// ============================================================

func TestNoteHandler_GetByID_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id", h.GetByID)

	w := doRequest(r, "GET", "/notes/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_GetByID_ServiceError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/:id", h.GetByID)

	svc.On("GetByID", uint(1)).Return(nil, errors.New("error"))

	w := doRequest(r, "GET", "/notes/1", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByUserID エラーパステスト
// ============================================================

func TestNoteHandler_GetByUserID_ServiceError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes", h.GetByUserID)

	svc.On("GetByUserID", uint(1), 1, 20).Return([]model.Note{}, errors.New("error"))

	w := doRequest(r, "GET", "/notes?page=1&limit=20", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestNoteHandler_GetByUserID_CountError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes", h.GetByUserID)

	svc.On("GetByUserID", uint(1), 1, 20).Return([]model.Note{}, nil)
	svc.On("CountByUserID", uint(1)).Return(int64(0), errors.New("error"))

	w := doRequest(r, "GET", "/notes?page=1&limit=20", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Update エラーパステスト
// ============================================================

func TestNoteHandler_Update_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	w := doRequest(r, "PUT", "/notes/abc", map[string]interface{}{"title": "t"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_Update_InvalidJSON(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	w := doRequestRaw(r, "PUT", "/notes/1", "invalid")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_Update_ServiceError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id", h.Update)

	svc.On("Update", uint(1), uint(1), "更新", "", "", (*uint)(nil)).Return(nil, errors.New("error"))

	w := doRequest(r, "PUT", "/notes/1", map[string]interface{}{"title": "更新"})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Delete エラーパステスト
// ============================================================

func TestNoteHandler_Delete_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.DELETE("/notes/:id", h.Delete)

	w := doRequest(r, "DELETE", "/notes/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_Delete_ServiceError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.DELETE("/notes/:id", h.Delete)

	svc.On("Delete", uint(1), uint(1)).Return(errors.New("error"))

	w := doRequest(r, "DELETE", "/notes/1", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Search エラーパステスト
// ============================================================

func TestNoteHandler_Search_MissingQuery(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/search", h.Search)

	w := doRequest(r, "GET", "/notes/search", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteHandler_Search_ServiceError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/search", h.Search)

	svc.On("Search", uint(1), "Go", 1, 20).Return([]model.Note{}, int64(0), errors.New("error"))

	w := doRequest(r, "GET", "/notes/search?q=Go&page=1&limit=20", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// ToggleFavorite エラーパステスト
// ============================================================

func TestNoteHandler_ToggleFavorite_InvalidID(t *testing.T) {
	h, _ := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/favorite", h.ToggleFavorite)

	w := doRequest(r, "PUT", "/notes/abc/favorite", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetFavorites テスト
// ============================================================

func TestNoteHandler_GetFavorites_Success(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/favorites", h.GetFavorites)

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "お気に入りノート"},
	}

	svc.On("GetFavorites", uint(1), 1, 20).Return(notes, nil)
	svc.On("CountFavoritesByUserID", uint(1)).Return(int64(1), nil)

	w := doRequest(r, "GET", "/notes/favorites?page=1&limit=20", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNoteHandler_GetFavorites_ServiceError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/favorites", h.GetFavorites)

	svc.On("GetFavorites", uint(1), 1, 20).Return([]model.Note{}, errors.New("error"))

	w := doRequest(r, "GET", "/notes/favorites?page=1&limit=20", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestNoteHandler_GetFavorites_CountError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.GET("/notes/favorites", h.GetFavorites)

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "お気に入りノート"},
	}

	svc.On("GetFavorites", uint(1), 1, 20).Return(notes, nil)
	svc.On("CountFavoritesByUserID", uint(1)).Return(int64(0), errors.New("error"))

	w := doRequest(r, "GET", "/notes/favorites?page=1&limit=20", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestNoteHandler_ToggleFavorite_ServiceError(t *testing.T) {
	h, svc := newTestNoteHandler()
	r := newRouter(1)
	r.PUT("/notes/:id/favorite", h.ToggleFavorite)

	svc.On("ToggleFavorite", uint(1), uint(1)).Return(errors.New("error"))

	w := doRequest(r, "PUT", "/notes/1/favorite", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
