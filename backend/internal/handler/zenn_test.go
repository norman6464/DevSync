package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockZennArticlePort は usecase/repository.ZennArticleRepository のモック。
type mockZennArticlePort struct{ mock.Mock }

func (m *mockZennArticlePort) UpsertArticles(ctx context.Context, userID uint, articles []model.ZennArticle) error {
	return m.Called(ctx, userID, articles).Error(0)
}

func (m *mockZennArticlePort) GetArticles(ctx context.Context, userID uint) ([]model.ZennArticle, error) {
	args := m.Called(ctx, userID)
	a, _ := args.Get(0).([]model.ZennArticle)
	return a, args.Error(1)
}

func (m *mockZennArticlePort) GetStats(ctx context.Context, userID uint) (*model.ZennStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ZennStats)
	return s, args.Error(1)
}

func (m *mockZennArticlePort) DeleteUserArticles(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

// mockZennFetcher は usecase/repository.ZennArticleFetcher のモック。
type mockZennFetcher struct{ mock.Mock }

func (m *mockZennFetcher) FetchArticles(ctx context.Context, username string) ([]model.ZennArticle, error) {
	args := m.Called(ctx, username)
	a, _ := args.Get(0).([]model.ZennArticle)
	return a, args.Error(1)
}

func (m *mockZennFetcher) UserExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

// newTestZennHandler は本物の usecase に port モックを注入したハンドラーを生成する。
// ユーザーの読み書きは user スライスと同じ port モックを共用する。
func newTestZennHandler() (*ArticlePlatformHandler[model.ZennArticle, model.ZennStats], *mockUserPort, *mockZennArticlePort, *mockZennFetcher) {
	users := new(mockUserPort)
	articles := new(mockZennArticlePort)
	fetcher := new(mockZennFetcher)
	h := NewArticlePlatformHandler("Zenn", ArticlePlatformOps[model.ZennArticle, model.ZennStats]{
		Connect:     usecase.NewConnectZennUseCase(users, articles, fetcher).Execute,
		Disconnect:  usecase.NewDisconnectZennUseCase(users, articles).Execute,
		Sync:        usecase.NewSyncZennUseCase(users, articles, fetcher).Execute,
		GetArticles: usecase.NewListZennArticlesUseCase(articles).Execute,
		GetStats:    usecase.NewGetZennStatsUseCase(articles).Execute,
	})
	return h, users, articles, fetcher
}

// ---------- Connect ----------

func TestZenn_Connect_Success(t *testing.T) {
	h, users, articles, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	fetched := []model.ZennArticle{{ZennID: 1, Title: "A"}, {ZennID: 2, Title: "B"}, {ZennID: 3, Title: "C"}}
	fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ZennUsername == "testuser"
	})).Return(nil)
	fetcher.On("FetchArticles", mock.Anything, "testuser").Return(fetched, nil)
	articles.On("UpsertArticles", mock.Anything, uint(1), mock.MatchedBy(func(a []model.ZennArticle) bool {
		// 取り込み時刻は同期時に揃えて入る。
		return len(a) == 3 && !a[0].UpdatedAt.IsZero()
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{"username": "testuser"})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"message":"Zenn connected successfully"`)
	assert.Contains(t, w.Body.String(), `"articles_count":3`)
	users.AssertExpectations(t)
	articles.AssertExpectations(t)
	fetcher.AssertExpectations(t)
}

func TestZenn_Connect_MissingUsername(t *testing.T) {
	h, users, _, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
	fetcher.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// 形式が不正なユーザー名は外部 API を叩かずに 400。
func TestZenn_Connect_InvalidUsernameFormat(t *testing.T) {
	h, users, _, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{"username": "bad!user"})
	assertStatus(t, w, http.StatusBadRequest)
	fetcher.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// Zenn 上に存在しないユーザー名は 400。
func TestZenn_Connect_UnknownZennUser(t *testing.T) {
	h, users, _, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	fetcher.On("UserExists", mock.Anything, "ghost").Return(false, nil)

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{"username": "ghost"})
	assertStatus(t, w, http.StatusBadRequest)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	fetcher.AssertExpectations(t)
}

// 存在確認自体に失敗した場合も 400（移行前と同じ）。
func TestZenn_Connect_UserExistsError(t *testing.T) {
	h, users, _, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	fetcher.On("UserExists", mock.Anything, "testuser").Return(false, errors.New("api error"))

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{"username": "testuser"})
	assertStatus(t, w, http.StatusBadRequest)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	fetcher.AssertExpectations(t)
}

func TestZenn_Connect_UserNotFound(t *testing.T) {
	h, users, articles, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{"username": "testuser"})
	assertStatus(t, w, http.StatusNotFound)
	users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	articles.AssertNotCalled(t, "UpsertArticles", mock.Anything, mock.Anything, mock.Anything)
	users.AssertExpectations(t)
}

// 記事取得に失敗した場合はユーザー名の保存後にエラーを返す（移行前と同じ）。
func TestZenn_Connect_FetchError(t *testing.T) {
	h, users, articles, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/connect", h.Connect)

	fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	users.On("Update", mock.Anything, mock.Anything).Return(nil)
	fetcher.On("FetchArticles", mock.Anything, "testuser").
		Return(nil, domain.NewError(domain.ErrCodeServiceUnavailable, "Zenn記事の取得に失敗", nil))

	w := doRequest(r, http.MethodPost, "/zenn/connect", map[string]string{"username": "testuser"})
	assertStatus(t, w, http.StatusServiceUnavailable)
	articles.AssertNotCalled(t, "UpsertArticles", mock.Anything, mock.Anything, mock.Anything)
	users.AssertExpectations(t)
	fetcher.AssertExpectations(t)
}

// ---------- Disconnect ----------

func TestZenn_Disconnect_Success(t *testing.T) {
	h, users, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.DELETE("/zenn/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, ZennUsername: "testuser"}, nil)
	users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ZennUsername == ""
	})).Return(nil)
	articles.On("DeleteUserArticles", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/zenn/disconnect", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "Zenn disconnected successfully")
	users.AssertExpectations(t)
	articles.AssertExpectations(t)
}

func TestZenn_Disconnect_UserNotFound(t *testing.T) {
	h, users, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.DELETE("/zenn/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/zenn/disconnect", nil)
	assertStatus(t, w, http.StatusNotFound)
	articles.AssertNotCalled(t, "DeleteUserArticles", mock.Anything, mock.Anything)
	users.AssertExpectations(t)
}

func TestZenn_Disconnect_DeleteError(t *testing.T) {
	h, users, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.DELETE("/zenn/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	users.On("Update", mock.Anything, mock.Anything).Return(nil)
	articles.On("DeleteUserArticles", mock.Anything, uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodDelete, "/zenn/disconnect", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	articles.AssertExpectations(t)
}

// ---------- Sync ----------

func TestZenn_Sync_Success(t *testing.T) {
	h, users, articles, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/sync", h.Sync)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, ZennUsername: "testuser"}, nil)
	fetcher.On("FetchArticles", mock.Anything, "testuser").Return([]model.ZennArticle{{ZennID: 1}}, nil)
	articles.On("UpsertArticles", mock.Anything, uint(1), mock.Anything).Return(nil)

	w := doRequest(r, http.MethodPost, "/zenn/sync", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"message":"Zenn synced successfully"`)
	assert.Contains(t, w.Body.String(), `"articles_count":1`)
	users.AssertExpectations(t)
	articles.AssertExpectations(t)
	fetcher.AssertExpectations(t)
}

// 未連携のユーザーは 400。
func TestZenn_Sync_NotConnected(t *testing.T) {
	h, users, _, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/sync", h.Sync)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, ZennUsername: ""}, nil)

	w := doRequest(r, http.MethodPost, "/zenn/sync", nil)
	assertStatus(t, w, http.StatusBadRequest)
	fetcher.AssertNotCalled(t, "FetchArticles", mock.Anything, mock.Anything)
	users.AssertExpectations(t)
}

func TestZenn_Sync_UserNotFound(t *testing.T) {
	h, users, _, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/sync", h.Sync)

	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/zenn/sync", nil)
	assertStatus(t, w, http.StatusNotFound)
	fetcher.AssertNotCalled(t, "FetchArticles", mock.Anything, mock.Anything)
}

func TestZenn_Sync_UpsertError(t *testing.T) {
	h, users, articles, fetcher := newTestZennHandler()
	r := newRouter(1)
	r.POST("/zenn/sync", h.Sync)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, ZennUsername: "testuser"}, nil)
	fetcher.On("FetchArticles", mock.Anything, "testuser").Return([]model.ZennArticle{{ZennID: 1}}, nil)
	articles.On("UpsertArticles", mock.Anything, uint(1), mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/zenn/sync", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	articles.AssertExpectations(t)
}

// ---------- GetArticles ----------

func TestZenn_GetArticles_Success(t *testing.T) {
	h, _, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.GET("/zenn/:userId/articles", h.GetArticles)

	articles.On("GetArticles", mock.Anything, uint(5)).Return([]model.ZennArticle{{Title: "Article 1"}}, nil)

	w := doRequest(r, http.MethodGet, "/zenn/5/articles", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"title":"Article 1"`)
	articles.AssertExpectations(t)
}

// 記事が無ければ空配列を返す。
func TestZenn_GetArticles_Empty(t *testing.T) {
	h, _, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.GET("/zenn/:userId/articles", h.GetArticles)

	articles.On("GetArticles", mock.Anything, uint(5)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/zenn/5/articles", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
	articles.AssertExpectations(t)
}

func TestZenn_GetArticles_InvalidID(t *testing.T) {
	h, _, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.GET("/zenn/:userId/articles", h.GetArticles)

	w := doRequest(r, http.MethodGet, "/zenn/abc/articles", nil)
	assertStatus(t, w, http.StatusBadRequest)
	articles.AssertNotCalled(t, "GetArticles", mock.Anything, mock.Anything)
}

func TestZenn_GetArticles_RepositoryError(t *testing.T) {
	h, _, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.GET("/zenn/:userId/articles", h.GetArticles)

	articles.On("GetArticles", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/zenn/5/articles", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	articles.AssertExpectations(t)
}

// ---------- GetStats ----------

func TestZenn_GetStats_Success(t *testing.T) {
	h, _, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.GET("/zenn/:userId/stats", h.GetStats)

	articles.On("GetStats", mock.Anything, uint(5)).
		Return(&model.ZennStats{TotalArticles: 8, TotalLikes: 20, TotalComments: 3}, nil)

	w := doRequest(r, http.MethodGet, "/zenn/5/stats", nil)
	assertStatus(t, w, http.StatusOK)
	assert.JSONEq(t, `{"total_articles":8,"total_likes":20,"total_comments":3}`, w.Body.String())
	articles.AssertExpectations(t)
}

func TestZenn_GetStats_InvalidID(t *testing.T) {
	h, _, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.GET("/zenn/:userId/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/zenn/abc/stats", nil)
	assertStatus(t, w, http.StatusBadRequest)
	articles.AssertNotCalled(t, "GetStats", mock.Anything, mock.Anything)
}

func TestZenn_GetStats_RepositoryError(t *testing.T) {
	h, _, articles, _ := newTestZennHandler()
	r := newRouter(1)
	r.GET("/zenn/:userId/stats", h.GetStats)

	articles.On("GetStats", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/zenn/5/stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	articles.AssertExpectations(t)
}
