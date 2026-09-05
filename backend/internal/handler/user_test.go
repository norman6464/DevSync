package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockUserPort は usecase/repository.UserRepository のモック（ctx 付き）。
// user スライスとメール配信設定のテストで共用する。
type mockUserPort struct{ mock.Mock }

func (m *mockUserPort) FindAll(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	u, _ := args.Get(0).([]model.User)
	return u, args.Error(1)
}

func (m *mockUserPort) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockUserPort) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockUserPort) Search(ctx context.Context, query string) ([]model.User, error) {
	args := m.Called(ctx, query)
	u, _ := args.Get(0).([]model.User)
	return u, args.Error(1)
}

func (m *mockUserPort) Update(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

// newTestUserHandler は本物の usecase に port モックを注入した UserHandler を生成する。
func newTestUserHandler() (*UserHandler, *mockUserPort) {
	repo := new(mockUserPort)
	h := NewUserHandler(
		usecase.NewListUsersUseCase(repo),
		usecase.NewGetUserUseCase(repo),
		usecase.NewGetUserByUsernameUseCase(repo),
		usecase.NewUpdateUserProfileUseCase(repo),
		usecase.NewGetProfileCompletenessUseCase(repo),
	)
	return h, repo
}

// ============================================================
// GetAll
// ============================================================

func TestUserHandler_GetAll(t *testing.T) {
	t.Run("検索クエリなしで全ユーザーを取得", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		repo.On("FindAll", mock.Anything).Return([]model.User{
			{ID: 1, Name: "テストユーザー1"}, {ID: 2, Name: "テストユーザー2"},
		}, nil)

		w := doRequest(r, http.MethodGet, "/users", nil)
		assertStatus(t, w, http.StatusOK)
		repo.AssertExpectations(t)
		repo.AssertNotCalled(t, "Search", mock.Anything, mock.Anything)
	})

	t.Run("検索クエリ付きなら検索を使う", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		repo.On("Search", mock.Anything, "テスト").Return([]model.User{{ID: 1}}, nil)

		w := doRequest(r, http.MethodGet, "/users?q=テスト", nil)
		assertStatus(t, w, http.StatusOK)
		repo.AssertExpectations(t)
		repo.AssertNotCalled(t, "FindAll", mock.Anything)
	})

	// 0 件でも null ではなく空配列を返す。
	t.Run("0 件でも空配列を返す", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		repo.On("FindAll", mock.Anything).Return([]model.User(nil), nil)

		w := doRequest(r, http.MethodGet, "/users", nil)
		assertStatus(t, w, http.StatusOK)
		assert.Equal(t, "[]", w.Body.String())
	})

	t.Run("DB 障害で 500 を返す", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		repo.On("FindAll", mock.Anything).Return([]model.User(nil), errors.New("db error"))

		w := doRequest(r, http.MethodGet, "/users", nil)
		assertStatus(t, w, http.StatusInternalServerError)
	})

	t.Run("100 文字超のクエリは 400", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		w := doRequest(r, http.MethodGet, "/users?q="+strings.Repeat("あ", 101), nil)
		assertStatus(t, w, http.StatusBadRequest)
		repo.AssertNotCalled(t, "Search", mock.Anything, mock.Anything)
	})

	t.Run("ちょうど 100 文字のクエリは通る", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		query := strings.Repeat("あ", 100)
		repo.On("Search", mock.Anything, query).Return([]model.User{}, nil)

		w := doRequest(r, http.MethodGet, "/users?q="+query, nil)
		assertStatus(t, w, http.StatusOK)
		repo.AssertExpectations(t)
	})

	t.Run("クエリの前後の空白は落とされる", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users", h.GetAll)

		repo.On("Search", mock.Anything, "テスト").Return([]model.User{}, nil)

		w := doRequest(r, http.MethodGet, "/users?q=%20%20テスト%20%20", nil)
		assertStatus(t, w, http.StatusOK)
		repo.AssertExpectations(t)
	})
}

// ============================================================
// GetByID / GetByUsername
// ============================================================

func TestUserHandler_GetByID(t *testing.T) {
	t.Run("ユーザーを取得できる", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/:id", h.GetByID)

		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.User{ID: 1, Name: "テストユーザー"}, nil)

		w := doRequest(r, http.MethodGet, "/users/1", nil)
		assertStatus(t, w, http.StatusOK)
		assert.Contains(t, w.Body.String(), "テストユーザー")
	})

	t.Run("不在のユーザーは 404", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/:id", h.GetByID)

		repo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

		w := doRequest(r, http.MethodGet, "/users/999", nil)
		assertStatus(t, w, http.StatusNotFound)
		assert.Contains(t, w.Body.String(), "ユーザーが見つかりません")
	})

	// DB 障害も 404 に潰す（移行前から変わらない挙動）。
	t.Run("DB 障害も 404 になる", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/:id", h.GetByID)

		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		w := doRequest(r, http.MethodGet, "/users/1", nil)
		assertStatus(t, w, http.StatusNotFound)
	})

	t.Run("ID が不正なら 400", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/:id", h.GetByID)

		w := doRequest(r, http.MethodGet, "/users/abc", nil)
		assertStatus(t, w, http.StatusBadRequest)
		repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})
}

func TestUserHandler_GetByUsername(t *testing.T) {
	t.Run("ユーザーを取得できる", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/username/:username", h.GetByUsername)

		repo.On("FindByUsername", mock.Anything, "testuser").
			Return(&model.User{ID: 1, Username: "testuser"}, nil)

		w := doRequest(r, http.MethodGet, "/users/username/testuser", nil)
		assertStatus(t, w, http.StatusOK)
		assert.Contains(t, w.Body.String(), "testuser")
	})

	t.Run("不在のユーザー名は 404", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/username/:username", h.GetByUsername)

		repo.On("FindByUsername", mock.Anything, "nonexistent").Return(nil, nil)

		w := doRequest(r, http.MethodGet, "/users/username/nonexistent", nil)
		assertStatus(t, w, http.StatusNotFound)
	})

	t.Run("50 文字超のユーザー名は 400", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/username/:username", h.GetByUsername)

		w := doRequest(r, http.MethodGet, "/users/username/"+strings.Repeat("a", 51), nil)
		assertStatus(t, w, http.StatusBadRequest)
		repo.AssertNotCalled(t, "FindByUsername", mock.Anything, mock.Anything)
	})
}

// ============================================================
// Update
// ============================================================

func TestUserHandler_Update(t *testing.T) {
	t.Run("プロフィールを更新できる", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.User{ID: 1, Name: "旧名前", Bio: "旧自己紹介"}, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Name == "新名前" && u.Bio == "新しい自己紹介"
		})).Return(nil)

		w := doRequest(r, http.MethodPut, "/users/1", map[string]interface{}{
			"name": "新名前", "bio": "新しい自己紹介",
		})
		assertStatus(t, w, http.StatusOK)
		repo.AssertExpectations(t)
	})

	// 自己紹介とアバターは渡された値で必ず上書きする（空文字なら空になる）。
	t.Run("自己紹介とアバターは空文字でも上書きする", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.User{ID: 1, Name: "名前", Bio: "旧自己紹介", AvatarURL: "https://old.example.com/a.png"}, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Bio == "" && u.AvatarURL == "" && u.Name == "名前"
		})).Return(nil)

		w := doRequest(r, http.MethodPut, "/users/1", map[string]interface{}{"name": "名前"})
		assertStatus(t, w, http.StatusOK)
		repo.AssertExpectations(t)
	})

	t.Run("他人のプロフィールは更新できない", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(2)
		r.PUT("/users/:id", h.Update)

		w := doRequest(r, http.MethodPut, "/users/1", map[string]interface{}{"name": "不正更新"})
		assertStatus(t, w, http.StatusForbidden)
		repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("不在のユーザーは 404", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(999)
		r.PUT("/users/:id", h.Update)

		repo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

		w := doRequest(r, http.MethodPut, "/users/999", map[string]interface{}{"name": "テスト"})
		assertStatus(t, w, http.StatusNotFound)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	// 名前が上限超過なら usecase の検証で 400 になり、書き込まれない。
	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Name: "名前"}, nil)

		w := doRequest(r, http.MethodPut, "/users/1", map[string]interface{}{
			"name": strings.Repeat("あ", 101),
		})
		assertStatus(t, w, http.StatusBadRequest)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("ID が不正なら 400", func(t *testing.T) {
		h, _ := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		w := doRequest(r, http.MethodPut, "/users/abc", map[string]interface{}{"name": "テスト"})
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("不正な JSON は 400", func(t *testing.T) {
		h, _ := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		w := doRequestRaw(r, http.MethodPut, "/users/1", "invalid json")
		assertStatus(t, w, http.StatusBadRequest)
	})

	t.Run("書き込みの DB 障害は 500", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.PUT("/users/:id", h.Update)

		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Name: "名前"}, nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

		w := doRequest(r, http.MethodPut, "/users/1", map[string]interface{}{"name": "新名前"})
		assertStatus(t, w, http.StatusInternalServerError)
	})
}

// ============================================================
// GetProfileCompleteness
// ============================================================

func TestUserHandler_GetProfileCompleteness(t *testing.T) {
	t.Run("完成度を返す", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/me/completeness", h.GetProfileCompleteness)

		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{
			ID: 1, AvatarURL: "https://example.com/a.png", Bio: "自己紹介",
			GitHubConnected: true, SkillsLanguages: "Go",
		}, nil)

		w := doRequest(r, http.MethodGet, "/users/me/completeness", nil)
		assertStatus(t, w, http.StatusOK)
		assert.Contains(t, w.Body.String(), `"percentage":100`)
	})

	t.Run("未設定の項目名を返す", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/me/completeness", h.GetProfileCompleteness)

		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Bio: "自己紹介"}, nil)

		w := doRequest(r, http.MethodGet, "/users/me/completeness", nil)
		assertStatus(t, w, http.StatusOK)
		body := w.Body.String()
		assert.Contains(t, body, `"percentage":25`)
		assert.Contains(t, body, "avatar")
		assert.Contains(t, body, "github")
		assert.Contains(t, body, "skills")
	})

	// 完成度だけは取得エラーを 404 に変換せず 500 のままにする（移行前と同じ）。
	t.Run("不在のユーザーは 500", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/me/completeness", h.GetProfileCompleteness)

		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		w := doRequest(r, http.MethodGet, "/users/me/completeness", nil)
		assertStatus(t, w, http.StatusInternalServerError)
	})

	t.Run("DB 障害は 500", func(t *testing.T) {
		h, repo := newTestUserHandler()
		r := newRouter(1)
		r.GET("/users/me/completeness", h.GetProfileCompleteness)

		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		w := doRequest(r, http.MethodGet, "/users/me/completeness", nil)
		assertStatus(t, w, http.StatusInternalServerError)
	})
}
