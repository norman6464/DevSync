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
