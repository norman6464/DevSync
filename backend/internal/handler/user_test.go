package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService は UserServiceInterface のモック実装。
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll(query string) ([]model.User, error) {
	args := m.Called(query)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserService) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Update(user *model.User) error {
	return m.Called(user).Error(0)
}

func newTestUserHandler() (*UserHandler, *MockUserService) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	return handler, mockService
}

// ============================================================
// GetAll テスト
// ============================================================

func TestUserHandler_GetAll(t *testing.T) {
	t.Run("検索クエリなしで全ユーザーを取得", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.GET("/users", handler.GetAll)

		users := []model.User{
			{ID: 1, Name: "テストユーザー1", Email: "user1@test.com"},
			{ID: 2, Name: "テストユーザー2", Email: "user2@test.com"},
		}
		mockService.On("GetAll", "").Return(users, nil)

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("検索クエリ付きでユーザーを検索", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.GET("/users", handler.GetAll)

		users := []model.User{
			{ID: 1, Name: "テストユーザー1", Email: "user1@test.com"},
		}
		mockService.On("GetAll", "テスト").Return(users, nil)

		req, _ := http.NewRequest("GET", "/users?q=テスト", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("サービスエラー時に500を返す", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.GET("/users", handler.GetAll)

		mockService.On("GetAll", "").Return([]model.User{}, errors.New("db error"))

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

// ============================================================
// GetByID テスト
// ============================================================

func TestUserHandler_GetByID(t *testing.T) {
	t.Run("正常にユーザーを取得", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.GET("/users/:id", handler.GetByID)

		user := &model.User{ID: 1, Name: "テストユーザー", Email: "test@test.com"}
		mockService.On("GetByID", uint(1)).Return(user, nil)

		req, _ := http.NewRequest("GET", "/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("ユーザーが見つからない場合404を返す", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.GET("/users/:id", handler.GetByID)

		mockService.On("GetByID", uint(999)).Return(nil, errors.New("not found"))

		req, _ := http.NewRequest("GET", "/users/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("無効なIDの場合400を返す", func(t *testing.T) {
		handler, _ := newTestUserHandler()
		router := setupRouter()
		router.GET("/users/:id", handler.GetByID)

		req, _ := http.NewRequest("GET", "/users/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// ============================================================
// GetByUsername テスト
// ============================================================

func TestUserHandler_GetByUsername(t *testing.T) {
	t.Run("正常にユーザーを取得", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.GET("/users/username/:username", handler.GetByUsername)

		user := &model.User{ID: 1, Name: "テストユーザー", Username: "testuser"}
		mockService.On("GetByUsername", "testuser").Return(user, nil)

		req, _ := http.NewRequest("GET", "/users/username/testuser", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("ユーザーが見つからない場合404を返す", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.GET("/users/username/:username", handler.GetByUsername)

		mockService.On("GetByUsername", "nonexistent").Return(nil, errors.New("not found"))

		req, _ := http.NewRequest("GET", "/users/username/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

// ============================================================
// Update テスト
// ============================================================

func TestUserHandler_Update(t *testing.T) {
	t.Run("正常にプロフィールを更新", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.PUT("/users/:id", func(c *gin.Context) {
			c.Set("userID", uint(1))
			handler.Update(c)
		})

		existing := &model.User{ID: 1, Name: "旧名前", Email: "test@test.com"}
		mockService.On("GetByID", uint(1)).Return(existing, nil)
		mockService.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

		input := map[string]interface{}{
			"name": "新名前",
			"bio":  "新しい自己紹介",
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("他人のプロフィールは更新できない", func(t *testing.T) {
		handler, _ := newTestUserHandler()
		router := setupRouter()
		router.PUT("/users/:id", func(c *gin.Context) {
			c.Set("userID", uint(2))
			handler.Update(c)
		})

		input := map[string]interface{}{"name": "不正更新"}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("存在しないユーザーの更新で404", func(t *testing.T) {
		handler, mockService := newTestUserHandler()
		router := setupRouter()
		router.PUT("/users/:id", func(c *gin.Context) {
			c.Set("userID", uint(999))
			handler.Update(c)
		})

		mockService.On("GetByID", uint(999)).Return(nil, errors.New("not found"))

		input := map[string]interface{}{"name": "テスト"}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest("PUT", "/users/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
