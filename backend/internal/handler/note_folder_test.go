package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
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

func (m *MockNoteFolderService) GetByUserID(userID uint, limit, offset int) ([]model.NoteFolder, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.NoteFolder), args.Get(1).(int64), args.Error(2)
}

func (m *MockNoteFolderService) GetChildren(parentID uint) ([]model.NoteFolder, error) {
	args := m.Called(parentID)
	return args.Get(0).([]model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderService) GetRootFolders(userID uint) ([]model.NoteFolder, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.NoteFolder), args.Error(1)
}

func (m *MockNoteFolderService) Update(id, userID uint, name string, parentID *uint) (*model.NoteFolder, error) {
	args := m.Called(id, userID, name, parentID)
	if f := args.Get(0); f != nil {
		return f.(*model.NoteFolder), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteFolderService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

func (m *MockNoteFolderService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
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
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.POST("/folders", h.Create)

	svc.On("Create", mock.AnythingOfType("*model.NoteFolder")).Return(nil)

	w := doRequest(r, "POST", "/folders", map[string]interface{}{"name": "新規フォルダ"})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestNoteFolderHandler_GetByID(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id", h.GetByID)

	folder := &model.NoteFolder{
		ID:     1,
		UserID: 1,
		Name:   "テストフォルダ",
	}

	svc.On("GetByID", uint(1)).Return(folder, nil)

	w := doRequest(r, "GET", "/folders/1", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestNoteFolderHandler_GetByUserID(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders", h.GetByUserID)

	folders := []model.NoteFolder{
		{ID: 1, UserID: 1, Name: "フォルダ1"},
		{ID: 2, UserID: 1, Name: "フォルダ2"},
	}

	svc.On("GetByUserID", uint(1), 20, 0).Return(folders, int64(2), nil)

	w := doRequest(r, "GET", "/folders", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// GetChildren テスト
// ============================================================

func TestNoteFolderHandler_GetChildren(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id/children", h.GetChildren)

	parentID := uint(1)
	children := []model.NoteFolder{
		{ID: 2, UserID: 1, Name: "子フォルダ1", ParentID: &parentID},
		{ID: 3, UserID: 1, Name: "子フォルダ2", ParentID: &parentID},
	}

	svc.On("GetChildren", uint(1)).Return(children, nil)

	w := doRequest(r, "GET", "/folders/1/children", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// GetRootFolders テスト
// ============================================================

func TestNoteFolderHandler_GetRootFolders(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/root", h.GetRootFolders)

	rootFolders := []model.NoteFolder{
		{ID: 1, UserID: 1, Name: "ルートフォルダ1", ParentID: nil},
		{ID: 2, UserID: 1, Name: "ルートフォルダ2", ParentID: nil},
	}

	svc.On("GetRootFolders", uint(1)).Return(rootFolders, nil)

	w := doRequest(r, "GET", "/folders/root", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestNoteFolderHandler_Update(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	updatedFolder := &model.NoteFolder{
		ID:     1,
		UserID: 1,
		Name:   "更新後フォルダ名",
	}

	svc.On("Update", uint(1), uint(1), "更新後フォルダ名", (*uint)(nil)).Return(updatedFolder, nil)

	w := doRequest(r, "PUT", "/folders/1", map[string]interface{}{"name": "更新後フォルダ名"})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteFolderHandler_Delete(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.DELETE("/folders/:id", h.Delete)

	svc.On("Delete", uint(1), uint(1)).Return(nil)

	w := doRequest(r, "DELETE", "/folders/1", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// ServiceError / InvalidID テスト
// ============================================================

func TestNoteFolderHandler_Create_ServiceError(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.POST("/folders", h.Create)

	svc.On("Create", mock.AnythingOfType("*model.NoteFolder")).Return(errors.New("db error"))

	w := doRequest(r, "POST", "/folders", map[string]interface{}{"name": "テスト"})
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteFolderHandler_GetByID_InvalidID(t *testing.T) {
	h, _ := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id", h.GetByID)

	w := doRequest(r, "GET", "/folders/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteFolderHandler_GetByID_ServiceError(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id", h.GetByID)

	svc.On("GetByID", uint(1)).Return(nil, service.ErrNotFound)

	w := doRequest(r, "GET", "/folders/1", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestNoteFolderHandler_GetByUserID_ServiceError(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders", h.GetByUserID)

	svc.On("GetByUserID", uint(1), 20, 0).Return([]model.NoteFolder(nil), int64(0), errors.New("db error"))

	w := doRequest(r, "GET", "/folders", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteFolderHandler_GetChildren_InvalidID(t *testing.T) {
	h, _ := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id/children", h.GetChildren)

	w := doRequest(r, "GET", "/folders/abc/children", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteFolderHandler_GetChildren_ServiceError(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id/children", h.GetChildren)

	svc.On("GetChildren", uint(1)).Return([]model.NoteFolder(nil), errors.New("db error"))

	w := doRequest(r, "GET", "/folders/1/children", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteFolderHandler_GetRootFolders_ServiceError(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/root", h.GetRootFolders)

	svc.On("GetRootFolders", uint(1)).Return([]model.NoteFolder(nil), errors.New("db error"))

	w := doRequest(r, "GET", "/folders/root", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteFolderHandler_Update_InvalidID(t *testing.T) {
	h, _ := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	w := doRequest(r, "PUT", "/folders/abc", map[string]interface{}{"name": "テスト"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteFolderHandler_Update_ServiceError(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	svc.On("Update", uint(1), uint(1), "変更", (*uint)(nil)).Return(nil, service.ErrForbidden)

	w := doRequest(r, "PUT", "/folders/1", map[string]interface{}{"name": "変更"})
	assertStatus(t, w, http.StatusForbidden)
}

func TestNoteFolderHandler_Delete_InvalidID(t *testing.T) {
	h, _ := newTestNoteFolderHandler()
	r := newRouter(1)
	r.DELETE("/folders/:id", h.Delete)

	w := doRequest(r, "DELETE", "/folders/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteFolderHandler_Delete_ServiceError(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.DELETE("/folders/:id", h.Delete)

	svc.On("Delete", uint(1), uint(1)).Return(service.ErrForbidden)

	w := doRequest(r, "DELETE", "/folders/1", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ============================================================
// GetMyCount テスト
// ============================================================

func TestNoteFolderHandler_GetMyCount_Success(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/my/count", h.GetMyCount)

	svc.On("CountByUserID", uint(1)).Return(int64(4), nil)

	w := doRequest(r, "GET", "/folders/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(4), body["count"])
}

func TestNoteFolderHandler_GetMyCount_ServiceError(t *testing.T) {
	h, svc := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/my/count", h.GetMyCount)

	svc.On("CountByUserID", uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, "GET", "/folders/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}
