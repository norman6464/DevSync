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

// mockCodeSnippetRepo は usecase/repository.CodeSnippetRepository のモック（ctx 付き）。
type mockCodeSnippetRepo struct{ mock.Mock }

func (m *mockCodeSnippetRepo) Create(ctx context.Context, s *model.CodeSnippet) error {
	return m.Called(ctx, s).Error(0)
}

func (m *mockCodeSnippetRepo) FindByID(ctx context.Context, id uint) (*model.CodeSnippet, error) {
	args := m.Called(ctx, id)
	s, _ := args.Get(0).(*model.CodeSnippet)
	return s, args.Error(1)
}

func (m *mockCodeSnippetRepo) FindByPostID(ctx context.Context, postID uint) ([]model.CodeSnippet, error) {
	args := m.Called(ctx, postID)
	s, _ := args.Get(0).([]model.CodeSnippet)
	return s, args.Error(1)
}

func (m *mockCodeSnippetRepo) FindByUserIDAndLanguage(ctx context.Context, userID uint, language string) ([]model.CodeSnippet, error) {
	args := m.Called(ctx, userID, language)
	s, _ := args.Get(0).([]model.CodeSnippet)
	return s, args.Error(1)
}

func (m *mockCodeSnippetRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.CodeSnippet, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	s, _ := args.Get(0).([]model.CodeSnippet)
	return s, args.Get(1).(int64), args.Error(2)
}

func (m *mockCodeSnippetRepo) Update(ctx context.Context, s *model.CodeSnippet) error {
	return m.Called(ctx, s).Error(0)
}

func (m *mockCodeSnippetRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockCodeSnippetRepo) CreateComment(ctx context.Context, c *model.SnippetComment) error {
	return m.Called(ctx, c).Error(0)
}

func (m *mockCodeSnippetRepo) GetComments(ctx context.Context, snippetID uint) ([]model.SnippetComment, error) {
	args := m.Called(ctx, snippetID)
	c, _ := args.Get(0).([]model.SnippetComment)
	return c, args.Error(1)
}

func (m *mockCodeSnippetRepo) DeleteComment(ctx context.Context, id, userID uint) error {
	return m.Called(ctx, id, userID).Error(0)
}

func (m *mockCodeSnippetRepo) IncrementForkCount(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockCodeSnippetRepo) Favorite(ctx context.Context, userID, snippetID uint) error {
	return m.Called(ctx, userID, snippetID).Error(0)
}

func (m *mockCodeSnippetRepo) Unfavorite(ctx context.Context, userID, snippetID uint) error {
	return m.Called(ctx, userID, snippetID).Error(0)
}

func (m *mockCodeSnippetRepo) HasFavorited(ctx context.Context, userID, snippetID uint) (bool, error) {
	args := m.Called(ctx, userID, snippetID)
	return args.Bool(0), args.Error(1)
}

func (m *mockCodeSnippetRepo) FindFavoritedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.CodeSnippet, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	s, _ := args.Get(0).([]model.CodeSnippet)
	return s, args.Get(1).(int64), args.Error(2)
}

func (m *mockCodeSnippetRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupCodeSnippetHandler は本物の usecase と port モックで CodeSnippetHandler を組む。
func setupCodeSnippetHandler() (*CodeSnippetHandler, *mockCodeSnippetRepo, *mockPostReader) {
	snippets := new(mockCodeSnippetRepo)
	posts := new(mockPostReader)
	h := NewCodeSnippetHandler(
		usecase.NewCreateCodeSnippetUseCase(snippets, posts),
		usecase.NewListCodeSnippetsByPostUseCase(snippets),
		usecase.NewListCodeSnippetsByLanguageUseCase(snippets),
		usecase.NewUpdateCodeSnippetUseCase(snippets),
		usecase.NewDeleteCodeSnippetUseCase(snippets),
		usecase.NewListSnippetCommentsUseCase(snippets),
		usecase.NewCreateSnippetCommentUseCase(snippets),
		usecase.NewDeleteSnippetCommentUseCase(snippets),
		usecase.NewSearchCodeSnippetsUseCase(snippets),
		usecase.NewForkCodeSnippetUseCase(snippets, posts),
		usecase.NewFavoriteCodeSnippetUseCase(snippets),
		usecase.NewUnfavoriteCodeSnippetUseCase(snippets),
		usecase.NewListFavoritedCodeSnippetsUseCase(snippets),
		usecase.NewCountCodeSnippetsUseCase(snippets),
	)
	return h, snippets, posts
}

// ownedSnippet は所有者が userID=1 のスニペットを返す。
func ownedSnippet(id uint) *model.CodeSnippet {
	s := &model.CodeSnippet{Language: "go", Code: "package main", FileName: "main.go"}
	s.ID = id
	s.UserID = 1
	return s
}

// --- Create ---

func TestCodeSnippet_Create_Success(t *testing.T) {
	h, snippets, posts := setupCodeSnippetHandler()
	posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPost(5), nil)
	snippets.On("Create", mock.Anything, mock.MatchedBy(func(s *model.CodeSnippet) bool {
		return s.PostID == 5 && s.Language == "go" && s.Code == "package main"
	})).Return(nil)
	snippets.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).Return(ownedSnippet(1), nil)

	r := newRouter(1)
	r.POST("/posts/:id/snippets", h.Create)
	w := doRequest(r, http.MethodPost, "/posts/5/snippets", map[string]interface{}{
		"language": "go", "file_name": "main.go", "code": "package main",
	})

	assertStatus(t, w, http.StatusCreated)
	snippets.AssertExpectations(t)
}

// 投稿が存在しなければ 404 を返し、作成しない。
func TestCodeSnippet_Create_PostNotFound(t *testing.T) {
	h, snippets, posts := setupCodeSnippetHandler()
	posts.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/posts/:id/snippets", h.Create)
	w := doRequest(r, http.MethodPost, "/posts/5/snippets", map[string]interface{}{
		"language": "go", "code": "x",
	})

	assertStatus(t, w, http.StatusNotFound)
	snippets.AssertNotCalled(t, "Create")
}

// 言語やコードが不正なら 400 を返し、投稿も読まない。
func TestCodeSnippet_Create_Invalid(t *testing.T) {
	cases := map[string]map[string]interface{}{
		"言語が空":       {"language": "", "code": "x"},
		"コードが空":      {"language": "go", "code": ""},
		"言語が 101 文字": {"language": strings.Repeat("a", 101), "code": "x"},
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			h, snippets, posts := setupCodeSnippetHandler()

			r := newRouter(1)
			r.POST("/posts/:id/snippets", h.Create)
			w := doRequest(r, http.MethodPost, "/posts/5/snippets", body)

			assertStatus(t, w, http.StatusBadRequest)
			posts.AssertNotCalled(t, "FindByID")
			snippets.AssertNotCalled(t, "Create")
		})
	}
}

// --- List / Search ---

func TestCodeSnippet_GetByPostID_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByPostID", mock.Anything, uint(5)).
		Return([]model.CodeSnippet{*ownedSnippet(1)}, nil)

	r := newRouter(1)
	r.GET("/posts/:id/snippets", h.GetByPostID)
	w := doRequest(r, http.MethodGet, "/posts/5/snippets", nil)

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

func TestCodeSnippet_Search_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("Search", mock.Anything, "go", 20, 0).
		Return([]model.CodeSnippet{*ownedSnippet(1)}, int64(1), nil)

	r := newRouter(1)
	r.GET("/snippets/search", h.Search)
	w := doRequest(r, http.MethodGet, "/snippets/search?q=go", nil)

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

// 言語で絞り込んで取得する。
// なお空文字のときの 400 はこの経路からは到達しない（パスパラメータが空になり得ないため）。
// 空文字の分岐は usecase テストで検証している。
func TestCodeSnippet_GetByUserLanguage_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByUserIDAndLanguage", mock.Anything, uint(1), "go").
		Return([]model.CodeSnippet{*ownedSnippet(1)}, nil)

	r := newRouter(1)
	r.GET("/snippets/language/:language", h.GetByUserLanguage)
	w := doRequest(r, http.MethodGet, "/snippets/language/go", nil)

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

// --- Update / Delete ---

func TestCodeSnippet_Update_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippet(1), nil)
	snippets.On("Update", mock.Anything, mock.MatchedBy(func(s *model.CodeSnippet) bool {
		// 空の項目は据え置かれる
		return s.Language == "rust" && s.Code == "package main" && s.FileName == "main.go"
	})).Return(nil)

	r := newRouter(1)
	r.PUT("/snippets/:id", h.Update)
	w := doRequest(r, http.MethodPut, "/snippets/1", map[string]interface{}{"language": "rust"})

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

// 所有者以外の更新は 403 を返し、保存しない。
func TestCodeSnippet_Update_Forbidden(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	other := ownedSnippet(1)
	other.UserID = 999
	snippets.On("FindByID", mock.Anything, uint(1)).Return(other, nil)

	r := newRouter(1)
	r.PUT("/snippets/:id", h.Update)
	w := doRequest(r, http.MethodPut, "/snippets/1", map[string]interface{}{"language": "rust"})

	assertStatus(t, w, http.StatusForbidden)
	snippets.AssertNotCalled(t, "Update")
}

func TestCodeSnippet_Delete_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippet(1), nil)
	snippets.On("Delete", mock.Anything, uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/snippets/:id", h.Delete)
	w := doRequest(r, http.MethodDelete, "/snippets/1", nil)

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

func TestCodeSnippet_Delete_Forbidden(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	other := ownedSnippet(1)
	other.UserID = 999
	snippets.On("FindByID", mock.Anything, uint(1)).Return(other, nil)

	r := newRouter(1)
	r.DELETE("/snippets/:id", h.Delete)
	w := doRequest(r, http.MethodDelete, "/snippets/1", nil)

	assertStatus(t, w, http.StatusForbidden)
	snippets.AssertNotCalled(t, "Delete")
}

// --- Comments ---

func TestCodeSnippet_CreateComment_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippet(1), nil)
	snippets.On("CreateComment", mock.Anything, mock.MatchedBy(func(c *model.SnippetComment) bool {
		return c.SnippetID == 1 && c.Content == "コメント"
	})).Return(nil)

	r := newRouter(1)
	r.POST("/snippets/:id/comments", h.CreateComment)
	w := doRequest(r, http.MethodPost, "/snippets/1/comments", map[string]interface{}{
		"content": "コメント", "line_number": 3,
	})

	assertStatus(t, w, http.StatusCreated)
	snippets.AssertExpectations(t)
}

// スニペットが存在しなければ 404 を返し、コメントを作らない。
func TestCodeSnippet_CreateComment_SnippetNotFound(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/snippets/:id/comments", h.CreateComment)
	w := doRequest(r, http.MethodPost, "/snippets/1/comments", map[string]interface{}{
		"content": "コメント", "line_number": 1,
	})

	assertStatus(t, w, http.StatusNotFound)
	snippets.AssertNotCalled(t, "CreateComment")
}

func TestCodeSnippet_GetComments_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("GetComments", mock.Anything, uint(1)).
		Return([]model.SnippetComment{{SnippetID: 1, Content: "c"}}, nil)

	r := newRouter(1)
	r.GET("/snippets/:id/comments", h.GetComments)
	w := doRequest(r, http.MethodGet, "/snippets/1/comments", nil)

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

// --- Fork ---

func TestCodeSnippet_Fork_Success(t *testing.T) {
	h, snippets, posts := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippet(1), nil)
	posts.On("FindByID", mock.Anything, uint(9)).Return(ownedPost(9), nil)
	snippets.On("Create", mock.Anything, mock.MatchedBy(func(s *model.CodeSnippet) bool {
		return s.PostID == 9 && s.ForkedFromID != nil && *s.ForkedFromID == 1
	})).Return(nil)
	snippets.On("IncrementForkCount", mock.Anything, uint(1)).Return(nil)

	r := newRouter(1)
	r.POST("/snippets/:id/fork", h.Fork)
	w := doRequest(r, http.MethodPost, "/snippets/1/fork", map[string]interface{}{"target_post_id": 9})

	assertStatus(t, w, http.StatusCreated)
	snippets.AssertExpectations(t)
}

// 他ユーザーの投稿へのフォークは 403 を返し、専用メッセージを使う。
func TestCodeSnippet_Fork_ForbiddenMessage(t *testing.T) {
	h, snippets, posts := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippet(1), nil)
	otherPost := ownedPost(9)
	otherPost.UserID = 999
	posts.On("FindByID", mock.Anything, uint(9)).Return(otherPost, nil)

	r := newRouter(1)
	r.POST("/snippets/:id/fork", h.Fork)
	w := doRequest(r, http.MethodPost, "/snippets/1/fork", map[string]interface{}{"target_post_id": 9})

	assertStatus(t, w, http.StatusForbidden)
	body := parseJSON(t, w)
	assert.Equal(t, "自分の投稿にのみフォークできます。投稿の編集権限がありません", body["error"])
	snippets.AssertNotCalled(t, "Create")
}

// フォーク元が存在しなければ 404。
func TestCodeSnippet_Fork_SnippetNotFound(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/snippets/:id/fork", h.Fork)
	w := doRequest(r, http.MethodPost, "/snippets/1/fork", map[string]interface{}{"target_post_id": 9})

	assertStatus(t, w, http.StatusNotFound)
	snippets.AssertNotCalled(t, "Create")
}

// --- Favorite ---

func TestCodeSnippet_Favorite_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippet(1), nil)
	snippets.On("HasFavorited", mock.Anything, uint(1), uint(1)).Return(false, nil)
	snippets.On("Favorite", mock.Anything, uint(1), uint(1)).Return(nil)

	r := newRouter(1)
	r.POST("/snippets/:id/favorite", h.Favorite)
	w := doRequest(r, http.MethodPost, "/snippets/1/favorite", nil)

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

// 既にお気に入り済みなら 409 を返し、追加しない。
func TestCodeSnippet_Favorite_Conflict(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippet(1), nil)
	snippets.On("HasFavorited", mock.Anything, uint(1), uint(1)).Return(true, nil)

	r := newRouter(1)
	r.POST("/snippets/:id/favorite", h.Favorite)
	w := doRequest(r, http.MethodPost, "/snippets/1/favorite", nil)

	assertStatus(t, w, http.StatusConflict)
	snippets.AssertNotCalled(t, "Favorite")
}

// スニペットが存在しなければ 404。
func TestCodeSnippet_Favorite_NotFound(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/snippets/:id/favorite", h.Favorite)
	w := doRequest(r, http.MethodPost, "/snippets/1/favorite", nil)

	assertStatus(t, w, http.StatusNotFound)
	snippets.AssertNotCalled(t, "Favorite")
}

func TestCodeSnippet_Unfavorite_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("Unfavorite", mock.Anything, uint(1), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/snippets/:id/favorite", h.Unfavorite)
	w := doRequest(r, http.MethodDelete, "/snippets/1/favorite", nil)

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

func TestCodeSnippet_GetFavorites_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("FindFavoritedByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.CodeSnippet{*ownedSnippet(1)}, int64(1), nil)

	r := newRouter(1)
	r.GET("/snippets/favorites", h.GetFavorites)
	w := doRequest(r, http.MethodGet, "/snippets/favorites", nil)

	assertStatus(t, w, http.StatusOK)
	snippets.AssertExpectations(t)
}

// --- Count ---

func TestCodeSnippet_GetMyCount_Success(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("CountByUserID", mock.Anything, uint(1)).Return(int64(7), nil)

	r := newRouter(1)
	r.GET("/snippets/my/count", h.GetMyCount)
	w := doRequest(r, http.MethodGet, "/snippets/my/count", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(7), body["count"])
}

func TestCodeSnippet_GetMyCount_RepoError(t *testing.T) {
	h, snippets, _ := setupCodeSnippetHandler()
	snippets.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/snippets/my/count", h.GetMyCount)
	w := doRequest(r, http.MethodGet, "/snippets/my/count", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}
