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

// mockQiitaArticlePort は usecase/repository.QiitaArticleRepository のモック。
type mockQiitaArticlePort struct{ mock.Mock }

func (m *mockQiitaArticlePort) UpsertArticles(ctx context.Context, userID uint, articles []model.QiitaArticle) error {
	return m.Called(ctx, userID, articles).Error(0)
}

func (m *mockQiitaArticlePort) GetArticles(ctx context.Context, userID uint) ([]model.QiitaArticle, error) {
	args := m.Called(ctx, userID)
	a, _ := args.Get(0).([]model.QiitaArticle)
	return a, args.Error(1)
}

func (m *mockQiitaArticlePort) GetStats(ctx context.Context, userID uint) (*model.QiitaStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.QiitaStats)
	return s, args.Error(1)
}

func (m *mockQiitaArticlePort) DeleteUserArticles(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

// mockQiitaFetcher は usecase/repository.QiitaArticleFetcher のモック。
type mockQiitaFetcher struct{ mock.Mock }

func (m *mockQiitaFetcher) FetchArticles(ctx context.Context, username string) ([]model.QiitaArticle, error) {
	args := m.Called(ctx, username)
	a, _ := args.Get(0).([]model.QiitaArticle)
	return a, args.Error(1)
}

func (m *mockQiitaFetcher) UserExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

// newTestQiitaHandler は本物の usecase に port モックを注入したハンドラーを生成する。
// ユーザーの読み書きは user スライスと同じ port モックを共用する。
func newTestQiitaHandler() (*ArticlePlatformHandler[model.QiitaArticle, model.QiitaStats], *mockUserPort, *mockQiitaArticlePort, *mockQiitaFetcher) {
	users := new(mockUserPort)
	articles := new(mockQiitaArticlePort)
	fetcher := new(mockQiitaFetcher)
	h := NewArticlePlatformHandler("Qiita", ArticlePlatformOps[model.QiitaArticle, model.QiitaStats]{
		Connect:     usecase.NewConnectQiitaUseCase(users, articles, fetcher).Execute,
		Disconnect:  usecase.NewDisconnectQiitaUseCase(users, articles).Execute,
		Sync:        usecase.NewSyncQiitaUseCase(users, articles, fetcher).Execute,
		GetArticles: usecase.NewListQiitaArticlesUseCase(articles).Execute,
		GetStats:    usecase.NewGetQiitaStatsUseCase(articles).Execute,
	})
	return h, users, articles, fetcher
}

// ---------- Connect ----------

func TestQiita_Connect_Success(t *testing.T) {
	h, users, articles, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	fetched := []model.QiitaArticle{{QiitaID: "a1"}, {QiitaID: "a2"}, {QiitaID: "a3"}, {QiitaID: "a4"}, {QiitaID: "a5"}}
	fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.QiitaUsername == "testuser"
	})).Return(nil)
	fetcher.On("FetchArticles", mock.Anything, "testuser").Return(fetched, nil)
	articles.On("UpsertArticles", mock.Anything, uint(1), mock.MatchedBy(func(a []model.QiitaArticle) bool {
		// 取り込み時刻は同期時に揃えて入る。
		return len(a) == 5 && !a[0].UpdatedAt.IsZero()
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{"username": "testuser"})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"message":"Qiita connected successfully"`)
	assert.Contains(t, w.Body.String(), `"articles_count":5`)
	users.AssertExpectations(t)
	articles.AssertExpectations(t)
	fetcher.AssertExpectations(t)
}

func TestQiita_Connect_MissingUsername(t *testing.T) {
	h, users, _, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
	fetcher.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// 形式が不正なユーザー名は外部 API を叩かずに 400。
func TestQiita_Connect_InvalidUsernameFormat(t *testing.T) {
	h, users, _, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{"username": "bad!user"})
	assertStatus(t, w, http.StatusBadRequest)
	fetcher.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// Qiita 上に存在しないユーザー名は 400。
func TestQiita_Connect_UnknownQiitaUser(t *testing.T) {
	h, users, _, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	fetcher.On("UserExists", mock.Anything, "ghost").Return(false, nil)

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{"username": "ghost"})
	assertStatus(t, w, http.StatusBadRequest)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	fetcher.AssertExpectations(t)
}

// 存在確認自体に失敗した場合も 400（移行前と同じ）。
func TestQiita_Connect_UserExistsError(t *testing.T) {
	h, users, _, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	fetcher.On("UserExists", mock.Anything, "testuser").Return(false, errors.New("api error"))

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{"username": "testuser"})
	assertStatus(t, w, http.StatusBadRequest)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	fetcher.AssertExpectations(t)
}

func TestQiita_Connect_UserNotFound(t *testing.T) {
	h, users, articles, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{"username": "testuser"})
	assertStatus(t, w, http.StatusNotFound)
	users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	articles.AssertNotCalled(t, "UpsertArticles", mock.Anything, mock.Anything, mock.Anything)
	users.AssertExpectations(t)
}

// 記事取得に失敗した場合はユーザー名の保存後にエラーを返す（移行前と同じ）。
func TestQiita_Connect_FetchError(t *testing.T) {
	h, users, articles, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/connect", h.Connect)

	fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	users.On("Update", mock.Anything, mock.Anything).Return(nil)
	fetcher.On("FetchArticles", mock.Anything, "testuser").
		Return(nil, domain.NewError(domain.ErrCodeNotFound, "Qiitaユーザーが見つかりません", nil))

	w := doRequest(r, http.MethodPost, "/qiita/connect", map[string]string{"username": "testuser"})
	assertStatus(t, w, http.StatusNotFound)
	articles.AssertNotCalled(t, "UpsertArticles", mock.Anything, mock.Anything, mock.Anything)
	users.AssertExpectations(t)
	fetcher.AssertExpectations(t)
}

// ---------- Disconnect ----------

func TestQiita_Disconnect_Success(t *testing.T) {
	h, users, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.DELETE("/qiita/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, QiitaUsername: "testuser"}, nil)
	users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.QiitaUsername == ""
	})).Return(nil)
	articles.On("DeleteUserArticles", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/qiita/disconnect", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "Qiita disconnected successfully")
	users.AssertExpectations(t)
	articles.AssertExpectations(t)
}

func TestQiita_Disconnect_UserNotFound(t *testing.T) {
	h, users, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.DELETE("/qiita/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/qiita/disconnect", nil)
	assertStatus(t, w, http.StatusNotFound)
	articles.AssertNotCalled(t, "DeleteUserArticles", mock.Anything, mock.Anything)
	users.AssertExpectations(t)
}

func TestQiita_Disconnect_DeleteError(t *testing.T) {
	h, users, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.DELETE("/qiita/disconnect", h.Disconnect)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	users.On("Update", mock.Anything, mock.Anything).Return(nil)
	articles.On("DeleteUserArticles", mock.Anything, uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodDelete, "/qiita/disconnect", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	articles.AssertExpectations(t)
}

// ---------- Sync ----------

func TestQiita_Sync_Success(t *testing.T) {
	h, users, articles, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/sync", h.Sync)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, QiitaUsername: "testuser"}, nil)
	fetcher.On("FetchArticles", mock.Anything, "testuser").Return([]model.QiitaArticle{{QiitaID: "a1"}}, nil)
	articles.On("UpsertArticles", mock.Anything, uint(1), mock.Anything).Return(nil)

	w := doRequest(r, http.MethodPost, "/qiita/sync", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"message":"Qiita synced successfully"`)
	assert.Contains(t, w.Body.String(), `"articles_count":1`)
	users.AssertExpectations(t)
	articles.AssertExpectations(t)
	fetcher.AssertExpectations(t)
}

// 未連携のユーザーは 400。
func TestQiita_Sync_NotConnected(t *testing.T) {
	h, users, _, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/sync", h.Sync)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, QiitaUsername: ""}, nil)

	w := doRequest(r, http.MethodPost, "/qiita/sync", nil)
	assertStatus(t, w, http.StatusBadRequest)
	fetcher.AssertNotCalled(t, "FetchArticles", mock.Anything, mock.Anything)
	users.AssertExpectations(t)
}

func TestQiita_Sync_UserNotFound(t *testing.T) {
	h, users, _, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/sync", h.Sync)

	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/qiita/sync", nil)
	assertStatus(t, w, http.StatusNotFound)
	fetcher.AssertNotCalled(t, "FetchArticles", mock.Anything, mock.Anything)
}

func TestQiita_Sync_UpsertError(t *testing.T) {
	h, users, articles, fetcher := newTestQiitaHandler()
	r := newRouter(1)
	r.POST("/qiita/sync", h.Sync)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, QiitaUsername: "testuser"}, nil)
	fetcher.On("FetchArticles", mock.Anything, "testuser").Return([]model.QiitaArticle{{QiitaID: "a1"}}, nil)
	articles.On("UpsertArticles", mock.Anything, uint(1), mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/qiita/sync", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	articles.AssertExpectations(t)
}

// ---------- GetArticles ----------

func TestQiita_GetArticles_Success(t *testing.T) {
	h, _, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/articles", h.GetArticles)

	articles.On("GetArticles", mock.Anything, uint(5)).
		Return([]model.QiitaArticle{{Title: "Article 1", Tags: "go,api"}}, nil)

	w := doRequest(r, http.MethodGet, "/qiita/5/articles", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"title":"Article 1"`)
	assert.Contains(t, w.Body.String(), `"tags":"go,api"`)
	articles.AssertExpectations(t)
}

// 記事が無ければ空配列を返す。
func TestQiita_GetArticles_Empty(t *testing.T) {
	h, _, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/articles", h.GetArticles)

	articles.On("GetArticles", mock.Anything, uint(5)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/qiita/5/articles", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
	articles.AssertExpectations(t)
}

func TestQiita_GetArticles_InvalidID(t *testing.T) {
	h, _, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/articles", h.GetArticles)

	w := doRequest(r, http.MethodGet, "/qiita/abc/articles", nil)
	assertStatus(t, w, http.StatusBadRequest)
	articles.AssertNotCalled(t, "GetArticles", mock.Anything, mock.Anything)
}

func TestQiita_GetArticles_RepositoryError(t *testing.T) {
	h, _, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/articles", h.GetArticles)

	articles.On("GetArticles", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/qiita/5/articles", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	articles.AssertExpectations(t)
}

// ---------- GetStats ----------

func TestQiita_GetStats_Success(t *testing.T) {
	h, _, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/stats", h.GetStats)

	articles.On("GetStats", mock.Anything, uint(5)).
		Return(&model.QiitaStats{TotalArticles: 10, TotalLikes: 50, TotalComments: 4}, nil)

	w := doRequest(r, http.MethodGet, "/qiita/5/stats", nil)
	assertStatus(t, w, http.StatusOK)
	assert.JSONEq(t, `{"total_articles":10,"total_likes":50,"total_comments":4}`, w.Body.String())
	articles.AssertExpectations(t)
}

func TestQiita_GetStats_InvalidID(t *testing.T) {
	h, _, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/qiita/abc/stats", nil)
	assertStatus(t, w, http.StatusBadRequest)
	articles.AssertNotCalled(t, "GetStats", mock.Anything, mock.Anything)
}

func TestQiita_GetStats_RepositoryError(t *testing.T) {
	h, _, articles, _ := newTestQiitaHandler()
	r := newRouter(1)
	r.GET("/qiita/:userId/stats", h.GetStats)

	articles.On("GetStats", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/qiita/5/stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	articles.AssertExpectations(t)
}
