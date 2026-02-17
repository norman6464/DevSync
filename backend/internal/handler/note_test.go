package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
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

func (m *MockNoteService) Update(note *model.Note) error {
	return m.Called(note).Error(0)
}

func (m *MockNoteService) Delete(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockNoteService) Search(userID uint, query string, page, limit int) ([]model.Note, int64, error) {
	args := m.Called(userID, query, page, limit)
	return args.Get(0).([]model.Note), args.Get(1).(int64), args.Error(2)
}

func (m *MockNoteService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNoteService) ToggleFavorite(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockNoteService) Archive(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockNoteService) Unarchive(id uint) error {
	return m.Called(id).Error(0)
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

// setupRouter はテスト用のGinルーターをセットアップする。
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// ============================================================
// Create テスト
// ============================================================

func TestNoteHandler_Create(t *testing.T) {
	handler, mockService := newTestNoteHandler()
	router := setupRouter()

	router.POST("/notes", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Create(c)
	})

	input := CreateNoteInput{
		Title:    "テストノート",
		Content:  "これはテスト内容です",
		Tags:     "Go,TDD",
		FolderID: nil,
	}
	body, _ := json.Marshal(input)

	mockService.On("Create", mock.AnythingOfType("*model.Note")).Return(nil)

	req, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestNoteHandler_GetByID(t *testing.T) {
	handler, mockService := newTestNoteHandler()
	router := setupRouter()

	router.GET("/notes/:id", handler.GetByID)

	note := &model.Note{
		ID:      1,
		UserID:  1,
		Title:   "テストノート",
		Content: "内容",
	}

	mockService.On("GetByID", uint(1)).Return(note, nil)

	req, _ := http.NewRequest("GET", "/notes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestNoteHandler_GetByUserID(t *testing.T) {
	handler, mockService := newTestNoteHandler()
	router := setupRouter()

	router.GET("/notes", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetByUserID(c)
	})

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "ノート1"},
		{ID: 2, UserID: 1, Title: "ノート2"},
	}

	mockService.On("GetByUserID", uint(1), 1, 20).Return(notes, nil)
	mockService.On("CountByUserID", uint(1)).Return(int64(2), nil)

	req, _ := http.NewRequest("GET", "/notes?page=1&limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestNoteHandler_Update(t *testing.T) {
	handler, mockService := newTestNoteHandler()
	router := setupRouter()

	router.PUT("/notes/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Update(c)
	})

	input := UpdateNoteInput{
		Title:   "更新後タイトル",
		Content: "更新後内容",
	}
	body, _ := json.Marshal(input)

	updatedNote := &model.Note{
		ID:      1,
		UserID:  1,
		Title:   "更新後タイトル",
		Content: "更新後内容",
	}

	mockService.On("GetByID", uint(1)).Return(updatedNote, nil)
	mockService.On("Update", mock.AnythingOfType("*model.Note")).Return(nil)

	req, _ := http.NewRequest("PUT", "/notes/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteHandler_Delete(t *testing.T) {
	handler, mockService := newTestNoteHandler()
	router := setupRouter()

	router.DELETE("/notes/:id", handler.Delete)

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/notes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// Search テスト
// ============================================================

func TestNoteHandler_Search(t *testing.T) {
	handler, mockService := newTestNoteHandler()
	router := setupRouter()

	router.GET("/notes/search", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Search(c)
	})

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "Go学習", Content: "Goを学習中"},
	}

	mockService.On("Search", uint(1), "Go", 1, 20).Return(notes, int64(1), nil)

	req, _ := http.NewRequest("GET", "/notes/search?q=Go&page=1&limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// ToggleFavorite テスト
// ============================================================

func TestNoteHandler_ToggleFavorite(t *testing.T) {
	handler, mockService := newTestNoteHandler()
	router := setupRouter()

	router.PUT("/notes/:id/favorite", handler.ToggleFavorite)

	mockService.On("ToggleFavorite", uint(1)).Return(nil)

	req, _ := http.NewRequest("PUT", "/notes/1/favorite", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// GetByFolderID テスト
// ============================================================

func TestNoteHandler_GetByFolderID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.GET("/folders/:folderId/notes", handler.GetByFolderID)

		notes := []model.Note{
			{ID: 1, UserID: 1, Title: "フォルダ内ノート1"},
			{ID: 2, UserID: 1, Title: "フォルダ内ノート2"},
		}

		mockService.On("GetByFolderID", uint(5)).Return(notes, nil)

		req, _ := http.NewRequest("GET", "/folders/5/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidFolderID", func(t *testing.T) {
		handler, _ := newTestNoteHandler()
		router := setupRouter()

		router.GET("/folders/:folderId/notes", handler.GetByFolderID)

		req, _ := http.NewRequest("GET", "/folders/abc/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.GET("/folders/:folderId/notes", handler.GetByFolderID)

		mockService.On("GetByFolderID", uint(5)).Return([]model.Note{}, assert.AnError)

		req, _ := http.NewRequest("GET", "/folders/5/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

// ============================================================
// Archive テスト
// ============================================================

func TestNoteHandler_Archive(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.PUT("/notes/:id/archive", handler.Archive)

		mockService.On("Archive", uint(1)).Return(nil)

		req, _ := http.NewRequest("PUT", "/notes/1/archive", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := newTestNoteHandler()
		router := setupRouter()

		router.PUT("/notes/:id/archive", handler.Archive)

		req, _ := http.NewRequest("PUT", "/notes/abc/archive", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.PUT("/notes/:id/archive", handler.Archive)

		mockService.On("Archive", uint(1)).Return(assert.AnError)

		req, _ := http.NewRequest("PUT", "/notes/1/archive", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

// ============================================================
// Unarchive テスト
// ============================================================

func TestNoteHandler_Unarchive(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.PUT("/notes/:id/unarchive", handler.Unarchive)

		mockService.On("Unarchive", uint(1)).Return(nil)

		req, _ := http.NewRequest("PUT", "/notes/1/unarchive", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := newTestNoteHandler()
		router := setupRouter()

		router.PUT("/notes/:id/unarchive", handler.Unarchive)

		req, _ := http.NewRequest("PUT", "/notes/abc/unarchive", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.PUT("/notes/:id/unarchive", handler.Unarchive)

		mockService.On("Unarchive", uint(1)).Return(assert.AnError)

		req, _ := http.NewRequest("PUT", "/notes/1/unarchive", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

// ============================================================
// GetArchived テスト
// ============================================================

func TestNoteHandler_GetArchived(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.GET("/notes/archived", func(c *gin.Context) {
			c.Set("userID", uint(1))
			handler.GetArchived(c)
		})

		notes := []model.Note{
			{ID: 3, UserID: 1, Title: "アーカイブノート1"},
		}

		mockService.On("GetArchived", uint(1), 1, 20).Return(notes, nil)
		mockService.On("CountArchivedByUserID", uint(1)).Return(int64(1), nil)

		req, _ := http.NewRequest("GET", "/notes/archived?page=1&limit=20", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GetArchivedError", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.GET("/notes/archived", func(c *gin.Context) {
			c.Set("userID", uint(1))
			handler.GetArchived(c)
		})

		mockService.On("GetArchived", uint(1), 1, 20).Return([]model.Note{}, assert.AnError)

		req, _ := http.NewRequest("GET", "/notes/archived?page=1&limit=20", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("CountError", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.GET("/notes/archived", func(c *gin.Context) {
			c.Set("userID", uint(1))
			handler.GetArchived(c)
		})

		notes := []model.Note{
			{ID: 3, UserID: 1, Title: "アーカイブノート1"},
		}

		mockService.On("GetArchived", uint(1), 1, 20).Return(notes, nil)
		mockService.On("CountArchivedByUserID", uint(1)).Return(int64(0), assert.AnError)

		req, _ := http.NewRequest("GET", "/notes/archived?page=1&limit=20", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

// ============================================================
// Duplicate テスト
// ============================================================

func TestNoteHandler_Duplicate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.POST("/notes/:id/duplicate", func(c *gin.Context) {
			c.Set("userID", uint(1))
			handler.Duplicate(c)
		})

		duplicate := &model.Note{
			ID:      10,
			UserID:  1,
			Title:   "テストノート (コピー)",
			Content: "元の内容",
		}

		mockService.On("Duplicate", uint(1), uint(1)).Return(duplicate, nil)

		req, _ := http.NewRequest("POST", "/notes/1/duplicate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler, _ := newTestNoteHandler()
		router := setupRouter()

		router.POST("/notes/:id/duplicate", func(c *gin.Context) {
			c.Set("userID", uint(1))
			handler.Duplicate(c)
		})

		req, _ := http.NewRequest("POST", "/notes/abc/duplicate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		handler, mockService := newTestNoteHandler()
		router := setupRouter()

		router.POST("/notes/:id/duplicate", func(c *gin.Context) {
			c.Set("userID", uint(1))
			handler.Duplicate(c)
		})

		mockService.On("Duplicate", uint(1), uint(1)).Return(nil, assert.AnError)

		req, _ := http.NewRequest("POST", "/notes/1/duplicate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
