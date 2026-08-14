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

// mockNoteFolderRepo は usecase/repository.NoteFolderRepository のモック（ctx 付き）。
type mockNoteFolderRepo struct{ mock.Mock }

func (m *mockNoteFolderRepo) Create(ctx context.Context, folder *model.NoteFolder) error {
	return m.Called(ctx, folder).Error(0)
}

func (m *mockNoteFolderRepo) FindByID(ctx context.Context, id uint) (*model.NoteFolder, error) {
	args := m.Called(ctx, id)
	f, _ := args.Get(0).(*model.NoteFolder)
	return f, args.Error(1)
}

func (m *mockNoteFolderRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.NoteFolder, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	f, _ := args.Get(0).([]model.NoteFolder)
	return f, args.Get(1).(int64), args.Error(2)
}

func (m *mockNoteFolderRepo) FindByParentID(ctx context.Context, parentID uint) ([]model.NoteFolder, error) {
	args := m.Called(ctx, parentID)
	f, _ := args.Get(0).([]model.NoteFolder)
	return f, args.Error(1)
}

func (m *mockNoteFolderRepo) FindRootsByUserID(ctx context.Context, userID uint) ([]model.NoteFolder, error) {
	args := m.Called(ctx, userID)
	f, _ := args.Get(0).([]model.NoteFolder)
	return f, args.Error(1)
}

func (m *mockNoteFolderRepo) Update(ctx context.Context, folder *model.NoteFolder) error {
	return m.Called(ctx, folder).Error(0)
}

func (m *mockNoteFolderRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockNoteFolderRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// newTestNoteFolderHandler は本物の usecase と port モックで NoteFolderHandler を組む。
func newTestNoteFolderHandler() (*NoteFolderHandler, *mockNoteFolderRepo) {
	repo := new(mockNoteFolderRepo)
	h := NewNoteFolderHandler(
		usecase.NewCreateNoteFolderUseCase(repo),
		usecase.NewGetNoteFolderUseCase(repo),
		usecase.NewListNoteFoldersUseCase(repo),
		usecase.NewListChildNoteFoldersUseCase(repo),
		usecase.NewListRootNoteFoldersUseCase(repo),
		usecase.NewUpdateNoteFolderUseCase(repo),
		usecase.NewCountNoteFoldersUseCase(repo),
		usecase.NewDeleteNoteFolderUseCase(repo),
	)
	return h, repo
}

// ownedFolder は所有者が userID=1 のフォルダを返す。
func ownedFolder() *model.NoteFolder {
	return &model.NoteFolder{ID: 1, UserID: 1, Name: "テストフォルダ"}
}

// ============================================================
// Create テスト
// ============================================================

func TestNoteFolderHandler_Create(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.POST("/folders", h.Create)

	repo.On("Create", mock.Anything, mock.MatchedBy(func(f *model.NoteFolder) bool {
		return f.UserID == 1 && f.Name == "新規フォルダ"
	})).Return(nil)

	w := doRequest(r, "POST", "/folders", map[string]interface{}{"name": "新規フォルダ"})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

// 名前が長すぎる場合は 400 を返し、作成しない。
func TestNoteFolderHandler_Create_NameTooLong(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.POST("/folders", h.Create)

	long := ""
	for i := 0; i < 101; i++ {
		long += "a"
	}

	w := doRequest(r, "POST", "/folders", map[string]interface{}{"name": long})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Create")
}

// 名前が空白のみの場合は 400 を返し、作成しない。
func TestNoteFolderHandler_Create_BlankName(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.POST("/folders", h.Create)

	w := doRequest(r, "POST", "/folders", map[string]interface{}{"name": "   "})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Create")
}

// ============================================================
// GetByID テスト
// ============================================================

func TestNoteFolderHandler_GetByID(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedFolder(), nil)

	w := doRequest(r, "GET", "/folders/1", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 所有権を検証しないため、他ユーザーのフォルダも取得できる（移行前の挙動を維持している）。
// TestNoteFolderHandler_GetByID_OtherUsersFolder は他ユーザーのフォルダを取得できないことを確認する。
func TestNoteFolderHandler_GetByID_OtherUsersFolder(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.NoteFolder{ID: 1, UserID: 999, Name: "他人のフォルダ"}, nil)

	w := doRequest(r, "GET", "/folders/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	assert.NotContains(t, w.Body.String(), "他人のフォルダ", "フォルダ名を漏らさない")
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestNoteFolderHandler_GetByUserID(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders", h.GetByUserID)

	folders := []model.NoteFolder{
		{ID: 1, UserID: 1, Name: "フォルダ1"},
		{ID: 2, UserID: 1, Name: "フォルダ2"},
	}
	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).Return(folders, int64(2), nil)

	w := doRequest(r, "GET", "/folders", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(2), body["total"])
	repo.AssertExpectations(t)
}

// ============================================================
// GetChildren テスト
// ============================================================

func TestNoteFolderHandler_GetChildren(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id/children", h.GetChildren)

	parentID := uint(1)
	children := []model.NoteFolder{
		{ID: 2, UserID: 1, Name: "子フォルダ1", ParentID: &parentID},
		{ID: 3, UserID: 1, Name: "子フォルダ2", ParentID: &parentID},
	}
	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.NoteFolder{ID: 1, UserID: 1, Name: "親フォルダ"}, nil)
	repo.On("FindByParentID", mock.Anything, uint(1)).Return(children, nil)

	w := doRequest(r, "GET", "/folders/1/children", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// TestNoteFolderHandler_GetChildren_OtherUsersFolder は他ユーザーのフォルダの
// 子一覧を辿れないことを確認する。
func TestNoteFolderHandler_GetChildren_OtherUsersFolder(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id/children", h.GetChildren)

	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.NoteFolder{ID: 1, UserID: 999, Name: "他人のフォルダ"}, nil)

	w := doRequest(r, "GET", "/folders/1/children", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "FindByParentID", mock.Anything, mock.Anything)
}

// ============================================================
// GetRootFolders テスト
// ============================================================

func TestNoteFolderHandler_GetRootFolders(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/root", h.GetRootFolders)

	roots := []model.NoteFolder{
		{ID: 1, UserID: 1, Name: "ルートフォルダ1"},
		{ID: 2, UserID: 1, Name: "ルートフォルダ2"},
	}
	repo.On("FindRootsByUserID", mock.Anything, uint(1)).Return(roots, nil)

	w := doRequest(r, "GET", "/folders/root", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestNoteFolderHandler_Update(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedFolder(), nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(f *model.NoteFolder) bool {
		return f.Name == "更新後フォルダ名"
	})).Return(nil)

	w := doRequest(r, "PUT", "/folders/1", map[string]interface{}{"name": "更新後フォルダ名"})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 所有者以外の更新は 403 を返し、保存しない。
func TestNoteFolderHandler_Update_Forbidden(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.NoteFolder{ID: 1, UserID: 999, Name: "他人のフォルダ"}, nil)

	w := doRequest(r, "PUT", "/folders/1", map[string]interface{}{"name": "変更"})
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Update")
}

// 自分自身を親にする更新は 400 を返し、保存しない。
func TestNoteFolderHandler_Update_SelfParent(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedFolder(), nil)

	w := doRequest(r, "PUT", "/folders/1", map[string]interface{}{"parent_id": 1})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Update")
}

// 自分の子孫を親にする更新は 400 を返し、保存しない。
func TestNoteFolderHandler_Update_Cycle(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedFolder(), nil)
	repo.On("FindByParentID", mock.Anything, uint(1)).
		Return([]model.NoteFolder{{ID: 2, UserID: 1, Name: "子"}}, nil)

	w := doRequest(r, "PUT", "/folders/1", map[string]interface{}{"parent_id": 2})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Update")
}

// 名前が空白のみの更新は 400 を返し、保存しない。
func TestNoteFolderHandler_Update_BlankName(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedFolder(), nil)

	w := doRequest(r, "PUT", "/folders/1", map[string]interface{}{"name": "   "})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Update")
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteFolderHandler_Delete(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.DELETE("/folders/:id", h.Delete)

	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedFolder(), nil)
	repo.On("Delete", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, "DELETE", "/folders/1", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 所有者以外の削除は 403 を返し、削除しない。
func TestNoteFolderHandler_Delete_Forbidden(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.DELETE("/folders/:id", h.Delete)

	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.NoteFolder{ID: 1, UserID: 999}, nil)

	w := doRequest(r, "DELETE", "/folders/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Delete")
}

// ============================================================
// 不在 / 不正 ID / DB 障害
// ============================================================

// 存在しないフォルダは 500 を返す（移行前の挙動を維持している）。
func TestNoteFolderHandler_GetByID_NotFound(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, "GET", "/folders/1", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestNoteFolderHandler_Update_NotFound(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.PUT("/folders/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, "PUT", "/folders/1", map[string]interface{}{"name": "変更"})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertNotCalled(t, "Update")
}

func TestNoteFolderHandler_Delete_NotFound(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.DELETE("/folders/:id", h.Delete)

	repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, "DELETE", "/folders/1", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertNotCalled(t, "Delete")
}

func TestNoteFolderHandler_Create_RepoError(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.POST("/folders", h.Create)

	repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

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

func TestNoteFolderHandler_GetByUserID_RepoError(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.NoteFolder(nil), int64(0), errors.New("db error"))

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

func TestNoteFolderHandler_GetChildren_RepoError(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/:id/children", h.GetChildren)

	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.NoteFolder{ID: 1, UserID: 1, Name: "親フォルダ"}, nil)
	repo.On("FindByParentID", mock.Anything, uint(1)).
		Return([]model.NoteFolder(nil), errors.New("db error"))

	w := doRequest(r, "GET", "/folders/1/children", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestNoteFolderHandler_GetRootFolders_RepoError(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/root", h.GetRootFolders)

	repo.On("FindRootsByUserID", mock.Anything, uint(1)).
		Return([]model.NoteFolder(nil), errors.New("db error"))

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

func TestNoteFolderHandler_Delete_InvalidID(t *testing.T) {
	h, _ := newTestNoteFolderHandler()
	r := newRouter(1)
	r.DELETE("/folders/:id", h.Delete)

	w := doRequest(r, "DELETE", "/folders/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetMyCount テスト
// ============================================================

func TestNoteFolderHandler_GetMyCount_Success(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/my/count", h.GetMyCount)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(4), nil)

	w := doRequest(r, "GET", "/folders/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(4), body["count"])
}

func TestNoteFolderHandler_GetMyCount_RepoError(t *testing.T) {
	h, repo := newTestNoteFolderHandler()
	r := newRouter(1)
	r.GET("/folders/my/count", h.GetMyCount)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, "GET", "/folders/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}
