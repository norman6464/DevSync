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

// MockNoteFolderService は NoteFolderService のモック実装。
type MockNoteFolderService struct {
	mock.Mock
}

func (m *MockNoteFolderService) Create(folder *model.NoteFolder) error {
	return m.Called(folder).Error(0)
}

func (m *MockNoteFolderService) GetByID(id uint) (*model.NoteFolder, error) {
	args := m.Called(id)
	if folder := args.Get(0); folder != nil {
		return folder.(*model.NoteFolder), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteFolderService) GetByUserID(userID uint) ([]model.NoteFolder, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderService) GetChildren(parentID uint) ([]model.NoteFolder, error) {
	args := m.Called(parentID)
	return args.Get(0).([]model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderService) GetRootFolders(userID uint) ([]model.NoteFolder, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderService) Update(folder *model.NoteFolder) error {
	return m.Called(folder).Error(0)
}

func (m *MockNoteFolderService) Delete(id uint) error {
	return m.Called(id).Error(0)
}

// newTestNoteFolderHandler はテスト用のNoteFolderHandlerを生成する。
func newTestNoteFolderHandler() (*NoteFolderHandler, *MockNoteFolderService) {
	mockService := new(MockNoteFolderService)
	handler := NewNoteFolderHandler(mockService)
	return handler, mockService
}

// ============================================================
// Create テスト
// ============================================================

func TestNoteFolderHandler_Create(t *testing.T) {
	handler, mockService := newTestNoteFolderHandler()
	router := setupRouter()

	router.POST("/folders", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Create(c)
	})

	input := CreateNoteFolderInput{
		Name:     "新規フォルダ",
		ParentID: nil,
	}
	body, _ := json.Marshal(input)

	mockService.On("Create", mock.AnythingOfType("*model.NoteFolder")).Return(nil)

	req, _ := http.NewRequest("POST", "/folders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestNoteFolderHandler_GetByID(t *testing.T) {
	handler, mockService := newTestNoteFolderHandler()
	router := setupRouter()

	router.GET("/folders/:id", handler.GetByID)

	folder := &model.NoteFolder{
		ID:     1,
		UserID: 1,
		Name:   "テストフォルダ",
	}

	mockService.On("GetByID", uint(1)).Return(folder, nil)

	req, _ := http.NewRequest("GET", "/folders/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestNoteFolderHandler_GetByUserID(t *testing.T) {
	handler, mockService := newTestNoteFolderHandler()
	router := setupRouter()

	router.GET("/folders", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetByUserID(c)
	})

	folders := []model.NoteFolder{
		{ID: 1, UserID: 1, Name: "フォルダ1"},
		{ID: 2, UserID: 1, Name: "フォルダ2"},
	}

	mockService.On("GetByUserID", uint(1)).Return(folders, nil)

	req, _ := http.NewRequest("GET", "/folders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// GetChildren テスト
// ============================================================

func TestNoteFolderHandler_GetChildren(t *testing.T) {
	handler, mockService := newTestNoteFolderHandler()
	router := setupRouter()

	router.GET("/folders/:id/children", handler.GetChildren)

	parentID := uint(1)
	children := []model.NoteFolder{
		{ID: 2, UserID: 1, Name: "子フォルダ1", ParentID: &parentID},
		{ID: 3, UserID: 1, Name: "子フォルダ2", ParentID: &parentID},
	}

	mockService.On("GetChildren", uint(1)).Return(children, nil)

	req, _ := http.NewRequest("GET", "/folders/1/children", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// GetRootFolders テスト
// ============================================================

func TestNoteFolderHandler_GetRootFolders(t *testing.T) {
	handler, mockService := newTestNoteFolderHandler()
	router := setupRouter()

	router.GET("/folders/root", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetRootFolders(c)
	})

	rootFolders := []model.NoteFolder{
		{ID: 1, UserID: 1, Name: "ルートフォルダ1", ParentID: nil},
		{ID: 2, UserID: 1, Name: "ルートフォルダ2", ParentID: nil},
	}

	mockService.On("GetRootFolders", uint(1)).Return(rootFolders, nil)

	req, _ := http.NewRequest("GET", "/folders/root", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestNoteFolderHandler_Update(t *testing.T) {
	handler, mockService := newTestNoteFolderHandler()
	router := setupRouter()

	router.PUT("/folders/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Update(c)
	})

	input := UpdateNoteFolderInput{
		Name: "更新後フォルダ名",
	}
	body, _ := json.Marshal(input)

	updatedFolder := &model.NoteFolder{
		ID:     1,
		UserID: 1,
		Name:   "更新後フォルダ名",
	}

	mockService.On("GetByID", uint(1)).Return(updatedFolder, nil)
	mockService.On("Update", mock.AnythingOfType("*model.NoteFolder")).Return(nil)

	req, _ := http.NewRequest("PUT", "/folders/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteFolderHandler_Delete(t *testing.T) {
	handler, mockService := newTestNoteFolderHandler()
	router := setupRouter()

	router.DELETE("/folders/:id", handler.Delete)

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/folders/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
