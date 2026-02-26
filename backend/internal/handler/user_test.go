package handler

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
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

func (m *MockUserService) UpdateProfile(id, userID uint, input *service.UpdateProfileInput) (*model.User, error) {
	args := m.Called(id, userID, input)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) GetProfileCompleteness(userID uint) (*service.ProfileCompleteness, error) {
	args := m.Called(userID)
	if p := args.Get(0); p != nil {
		return p.(*service.ProfileCompleteness), args.Error(1)
	}
	return nil, args.Error(1)
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
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		users := []model.User{
			{ID: 1, Name: "テストユーザー1", Email: "user1@test.com"},
			{ID: 2, Name: "テストユーザー2", Email: "user2@test.com"},
		}
		svc.On("GetAll", "").Return(users, nil)

		w := doRequest(r, "GET", "/users", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("検索クエリ付きでユーザーを検索", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		users := []model.User{
			{ID: 1, Name: "テストユーザー1", Email: "user1@test.com"},
		}
		svc.On("GetAll", "テスト").Return(users, nil)

		w := doRequest(r, "GET", "/users?q=テスト", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("サービスエラー時に500を返す", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		svc.On("GetAll", "").Return([]model.User{}, errors.New("db error"))

		w := doRequest(r, "GET", "/users", nil)
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})

	t.Run("100文字超のクエリで400エラー", func(t *testing.T) {
		h, _ := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		longQuery := strings.Repeat("あ", 101)
		w := doRequest(r, "GET", "/users?q="+longQuery, nil)
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("ちょうど100文字のクエリで成功", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		query100 := strings.Repeat("あ", 100)
		svc.On("GetAll", query100).Return([]model.User{}, nil)

		w := doRequest(r, "GET", "/users?q="+query100, nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("前後の空白がTrimSpaceされる", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		svc.On("GetAll", "テスト").Return([]model.User{}, nil)

		// 空白はURLエンコードする必要がある
		w := doRequest(r, "GET", "/users?q=%20%20テスト%20%20", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})
}

// ============================================================
// GetByID テスト
// ============================================================

func TestUserHandler_GetByID(t *testing.T) {
	t.Run("正常にユーザーを取得", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/:id", h.GetByID)

		user := &model.User{ID: 1, Name: "テストユーザー", Email: "test@test.com"}
		svc.On("GetByID", uint(1)).Return(user, nil)

		w := doRequest(r, "GET", "/users/1", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("ユーザーが見つからない場合404を返す", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/:id", h.GetByID)

		svc.On("GetByID", uint(999)).Return(nil, domain.NewError(domain.ErrCodeNotFound, "not found", nil))

		w := doRequest(r, "GET", "/users/999", nil)
		assertStatus(t, w, http.StatusNotFound)
		svc.AssertExpectations(t)
	})

	t.Run("無効なIDの場合400を返す", func(t *testing.T) {
		h, _ := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/:id", h.GetByID)

		w := doRequest(r, "GET", "/users/abc", nil)
		assertStatus(t, w, http.StatusBadRequest)
	})
}

// ============================================================
// GetByUsername テスト
// ============================================================

func TestUserHandler_GetByUsername(t *testing.T) {
	t.Run("正常にユーザーを取得", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/username/:username", h.GetByUsername)

		user := &model.User{ID: 1, Name: "テストユーザー", Username: "testuser"}
		svc.On("GetByUsername", "testuser").Return(user, nil)

		w := doRequest(r, "GET", "/users/username/testuser", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("ユーザーが見つからない場合404を返す", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/username/:username", h.GetByUsername)

		svc.On("GetByUsername", "nonexistent").Return(nil, domain.NewError(domain.ErrCodeNotFound, "not found", nil))

		w := doRequest(r, "GET", "/users/username/nonexistent", nil)
		assertStatus(t, w, http.StatusNotFound)
		svc.AssertExpectations(t)
	})
}

// ============================================================
// Update テスト
// ============================================================

func TestUserHandler_Update(t *testing.T) {
	t.Run("正常にプロフィールを更新", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		updated := &model.User{ID: 1, Name: "新名前", Bio: "新しい自己紹介", Email: "test@test.com"}
		svc.On("UpdateProfile", uint(1), uint(1), mock.AnythingOfType("*service.UpdateProfileInput")).Return(updated, nil)

		w := doRequest(r, "PUT", "/users/1", map[string]interface{}{
			"name": "新名前",
			"bio":  "新しい自己紹介",
		})
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("他人のプロフィールは更新できない", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(2)
		r.PUT("/users/:id", h.Update)

		svc.On("UpdateProfile", uint(1), uint(2), mock.AnythingOfType("*service.UpdateProfileInput")).Return(nil, domain.ErrForbidden)

		w := doRequest(r, "PUT", "/users/1", map[string]interface{}{"name": "不正更新"})
		assertStatus(t, w, http.StatusForbidden)
		svc.AssertExpectations(t)
	})

	t.Run("存在しないユーザーの更新で404", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(999)
		r.PUT("/users/:id", h.Update)

		svc.On("UpdateProfile", uint(999), uint(999), mock.AnythingOfType("*service.UpdateProfileInput")).Return(nil, domain.ErrNotFound)

		w := doRequest(r, "PUT", "/users/999", map[string]interface{}{"name": "テスト"})
		assertStatus(t, w, http.StatusNotFound)
		svc.AssertExpectations(t)
	})

	t.Run("無効なIDの場合400を返す", func(t *testing.T) {
		h, _ := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		w := doRequest(r, "PUT", "/users/abc", map[string]interface{}{"name": "テスト"})
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("不正なJSONで400を返す", func(t *testing.T) {
		h, _ := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		w := doRequestRaw(r, "PUT", "/users/1", "invalid json")
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("サービスエラー時に500を返す", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		svc.On("UpdateProfile", uint(1), uint(1), mock.AnythingOfType("*service.UpdateProfileInput")).Return(nil, errors.New("db error"))

		w := doRequest(r, "PUT", "/users/1", map[string]interface{}{"name": "新名前"})
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})
}

// ============================================================
// GetProfileCompleteness テスト
// ============================================================

func TestUserHandler_GetProfileCompleteness(t *testing.T) {
	t.Run("正常にプロフィール完成度を取得", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/me/completeness", h.GetProfileCompleteness)

		result := &service.ProfileCompleteness{Percentage: 80}
		svc.On("GetProfileCompleteness", uint(1)).Return(result, nil)

		w := doRequest(r, "GET", "/users/me/completeness", nil)
		assertStatus(t, w, http.StatusOK)
		svc.AssertExpectations(t)
	})

	t.Run("サービスエラー時に500を返す", func(t *testing.T) {
		h, svc := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/me/completeness", h.GetProfileCompleteness)

		svc.On("GetProfileCompleteness", uint(1)).Return(nil, errors.New("db error"))

		w := doRequest(r, "GET", "/users/me/completeness", nil)
		assertStatus(t, w, http.StatusInternalServerError)
		svc.AssertExpectations(t)
	})
}
