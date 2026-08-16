package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockYouTubeVideoPort は usecase/repository.YouTubeVideoRepository のモック。
type mockYouTubeVideoPort struct{ mock.Mock }

func (m *mockYouTubeVideoPort) UpsertVideos(ctx context.Context, videos []model.YouTubeVideo) error {
	return m.Called(ctx, videos).Error(0)
}
func (m *mockYouTubeVideoPort) FindByVideoIDs(ctx context.Context, videoIDs []string) ([]model.YouTubeVideo, error) {
	args := m.Called(ctx, videoIDs)
	v, _ := args.Get(0).([]model.YouTubeVideo)
	return v, args.Error(1)
}
func (m *mockYouTubeVideoPort) FindCachedSearch(ctx context.Context, query, language string) (*model.YouTubeSearchCache, error) {
	args := m.Called(ctx, query, language)
	c, _ := args.Get(0).(*model.YouTubeSearchCache)
	return c, args.Error(1)
}
func (m *mockYouTubeVideoPort) SaveSearchCache(ctx context.Context, cache *model.YouTubeSearchCache) error {
	return m.Called(ctx, cache).Error(0)
}

// mockYouTubeSearcher は usecase/repository.YouTubeVideoSearcher のモック。
type mockYouTubeSearcher struct{ mock.Mock }

func (m *mockYouTubeSearcher) SearchVideos(ctx context.Context, query string, maxResults int, language string) ([]model.YouTubeVideo, error) {
	args := m.Called(ctx, query, maxResults, language)
	v, _ := args.Get(0).([]model.YouTubeVideo)
	return v, args.Error(1)
}

// newTestYouTubeHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestYouTubeHandler() (*YouTubeHandler, *mockUserPort, *mockYouTubeVideoPort, *mockYouTubeSearcher) {
	users := new(mockUserPort)
	videos := new(mockYouTubeVideoPort)
	searcher := new(mockYouTubeSearcher)
	h := NewYouTubeHandler(
		usecase.NewSearchYouTubeVideosUseCase(videos, searcher),
		usecase.NewRecommendYouTubeVideosUseCase(users, videos, searcher),
		usecase.NewCheckYouTubeAvailabilityUseCase(searcher),
	)
	return h, users, videos, searcher
}

// newUnavailableYouTubeHandler は API キー未設定（検索クライアント無し）のハンドラーを生成する。
func newUnavailableYouTubeHandler() (*YouTubeHandler, *mockUserPort, *mockYouTubeVideoPort) {
	users := new(mockUserPort)
	videos := new(mockYouTubeVideoPort)
	var searcher repository.YouTubeVideoSearcher
	h := NewYouTubeHandler(
		usecase.NewSearchYouTubeVideosUseCase(videos, searcher),
		usecase.NewRecommendYouTubeVideosUseCase(users, videos, searcher),
		usecase.NewCheckYouTubeAvailabilityUseCase(searcher),
	)
	return h, users, videos
}

// ---------- Search ----------

func TestYouTubeSearch_Success(t *testing.T) {
	h, _, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	found := []model.YouTubeVideo{{VideoID: "abc123", Title: "Go Tutorial"}}
	videos.On("FindCachedSearch", mock.Anything, "golang", "en").Return(nil, nil)
	searcher.On("SearchVideos", mock.Anything, "golang", 10, "en").Return(found, nil)
	videos.On("UpsertVideos", mock.Anything, found).Return(nil)
	videos.On("SaveSearchCache", mock.Anything, mock.MatchedBy(func(c *model.YouTubeSearchCache) bool {
		return c.Query == "golang" && c.Language == "en" && c.VideoIDs == "abc123" && c.CacheExpires.After(time.Now())
	})).Return(nil)

	w := doRequest(r, http.MethodGet, "/youtube/search?q=golang&lang=en", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "query", "golang")
	assertJSONEqual(t, body, "total", float64(1))
	assertJSONEqual(t, body, "cached", false)
	videos.AssertExpectations(t)
	searcher.AssertExpectations(t)
}

// キャッシュがあれば外部 API を叩かない。
func TestYouTubeSearch_CacheHit(t *testing.T) {
	h, _, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	videos.On("FindCachedSearch", mock.Anything, "golang", "ja").
		Return(&model.YouTubeSearchCache{VideoIDs: "abc123,def456"}, nil)
	videos.On("FindByVideoIDs", mock.Anything, []string{"abc123", "def456"}).
		Return([]model.YouTubeVideo{{VideoID: "abc123"}, {VideoID: "def456"}}, nil)

	w := doRequest(r, http.MethodGet, "/youtube/search?q=golang", nil)
	assertStatus(t, w, http.StatusOK)
	assertJSONEqual(t, parseJSON(t, w), "cached", true)
	searcher.AssertNotCalled(t, "SearchVideos", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	videos.AssertExpectations(t)
}

// キャッシュに紐づく動画が取れなければ検索し直す。
func TestYouTubeSearch_CacheWithoutVideos(t *testing.T) {
	h, _, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	videos.On("FindCachedSearch", mock.Anything, "golang", "ja").
		Return(&model.YouTubeSearchCache{VideoIDs: "abc123"}, nil)
	videos.On("FindByVideoIDs", mock.Anything, []string{"abc123"}).Return([]model.YouTubeVideo{}, nil)
	searcher.On("SearchVideos", mock.Anything, "golang", 10, "ja").Return([]model.YouTubeVideo{}, nil)

	w := doRequest(r, http.MethodGet, "/youtube/search?q=golang", nil)
	assertStatus(t, w, http.StatusOK)
	assertJSONEqual(t, parseJSON(t, w), "cached", false)
	// 結果が空ならキャッシュは書かない。
	videos.AssertNotCalled(t, "SaveSearchCache", mock.Anything, mock.Anything)
	searcher.AssertExpectations(t)
}

// 検索キーワードは小文字化してキャッシュを引く。
func TestYouTubeSearch_NormalizesQuery(t *testing.T) {
	h, _, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	videos.On("FindCachedSearch", mock.Anything, "golang", "ja").Return(nil, nil)
	// 外部 API へは入力そのままのキーワードを渡す。
	searcher.On("SearchVideos", mock.Anything, "  GoLang  ", 10, "ja").Return([]model.YouTubeVideo{}, nil)

	w := doRequest(r, http.MethodGet, "/youtube/search?q=%20%20GoLang%20%20", nil)
	assertStatus(t, w, http.StatusOK)
	videos.AssertExpectations(t)
	searcher.AssertExpectations(t)
}

func TestYouTubeSearch_MissingQuery(t *testing.T) {
	h, _, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	w := doRequest(r, http.MethodGet, "/youtube/search", nil)
	assertStatus(t, w, http.StatusBadRequest)
	videos.AssertNotCalled(t, "FindCachedSearch", mock.Anything, mock.Anything, mock.Anything)
	searcher.AssertNotCalled(t, "SearchVideos", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// 空白だけのキーワードは usecase 側で 400 になる。
func TestYouTubeSearch_BlankQuery(t *testing.T) {
	h, _, _, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	w := doRequest(r, http.MethodGet, "/youtube/search?q=%20%20", nil)
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "検索キーワードを入力してください")
	searcher.AssertNotCalled(t, "SearchVideos", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// 未対応の言語コードは 400。
func TestYouTubeSearch_InvalidLanguage(t *testing.T) {
	h, _, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	w := doRequest(r, http.MethodGet, "/youtube/search?q=golang&lang=xx", nil)
	assertStatus(t, w, http.StatusBadRequest)
	videos.AssertNotCalled(t, "FindCachedSearch", mock.Anything, mock.Anything, mock.Anything)
	searcher.AssertNotCalled(t, "SearchVideos", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestYouTubeSearch_SearcherError(t *testing.T) {
	h, _, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	videos.On("FindCachedSearch", mock.Anything, "test", "ja").Return(nil, nil)
	searcher.On("SearchVideos", mock.Anything, "test", 10, "ja").Return(nil, errors.New("api error"))

	w := doRequest(r, http.MethodGet, "/youtube/search?q=test", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	searcher.AssertExpectations(t)
}

// キャッシュ保存に失敗しても検索結果は返す。
func TestYouTubeSearch_CacheSaveFailureIsIgnored(t *testing.T) {
	h, _, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	found := []model.YouTubeVideo{{VideoID: "abc123"}}
	videos.On("FindCachedSearch", mock.Anything, "golang", "ja").Return(nil, nil)
	searcher.On("SearchVideos", mock.Anything, "golang", 10, "ja").Return(found, nil)
	videos.On("UpsertVideos", mock.Anything, found).Return(errors.New("db error"))
	videos.On("SaveSearchCache", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/youtube/search?q=golang", nil)
	assertStatus(t, w, http.StatusOK)
	assertJSONEqual(t, parseJSON(t, w), "total", float64(1))
	videos.AssertExpectations(t)
}

// APIキー未設定なら 503。
func TestYouTubeSearch_Unavailable(t *testing.T) {
	h, _, videos := newUnavailableYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/search", h.Search)

	w := doRequest(r, http.MethodGet, "/youtube/search?q=golang", nil)
	assertStatus(t, w, http.StatusServiceUnavailable)
	assert.Contains(t, w.Body.String(), "YouTube APIが設定されていません")
	videos.AssertNotCalled(t, "FindCachedSearch", mock.Anything, mock.Anything, mock.Anything)
}

// ---------- Recommend ----------

func TestYouTubeRecommend_Success(t *testing.T) {
	h, users, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/recommend", h.Recommend)

	users.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, SkillsLanguages: "Go, TypeScript", SkillsFrameworks: "React"}, nil)
	// 最大 3 スキルまで検索する。
	for _, skill := range []string{"Go", "TypeScript", "React"} {
		query := skill + " プログラミング チュートリアル"
		videos.On("FindCachedSearch", mock.Anything, strings.ToLower(query), "ja").
			Return(&model.YouTubeSearchCache{VideoIDs: skill}, nil)
		videos.On("FindByVideoIDs", mock.Anything, []string{skill}).
			Return([]model.YouTubeVideo{{VideoID: skill}}, nil)
	}

	w := doRequest(r, http.MethodGet, "/youtube/recommend", nil)
	assertStatus(t, w, http.StatusOK)
	body := w.Body.String()
	assert.Contains(t, body, `"skills":["Go","TypeScript","React"]`)
	assert.Contains(t, body, `"available":true`)
	users.AssertExpectations(t)
	videos.AssertExpectations(t)
	searcher.AssertNotCalled(t, "SearchVideos", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// スキル未登録なら空配列を返す。
func TestYouTubeRecommend_NoSkills(t *testing.T) {
	h, users, videos, _ := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/recommend", h.Recommend)

	users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)

	w := doRequest(r, http.MethodGet, "/youtube/recommend", nil)
	assertStatus(t, w, http.StatusOK)
	body := w.Body.String()
	assert.Contains(t, body, `"videos":[]`)
	assert.Contains(t, body, `"skills":[]`)
	videos.AssertNotCalled(t, "FindCachedSearch", mock.Anything, mock.Anything, mock.Anything)
	users.AssertExpectations(t)
}

// スキルごとの検索失敗は無視して続行する。
func TestYouTubeRecommend_SkipsFailedSkill(t *testing.T) {
	h, users, videos, searcher := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/recommend", h.Recommend)

	users.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, SkillsLanguages: "Go,Rust"}, nil)
	videos.On("FindCachedSearch", mock.Anything, mock.Anything, "ja").Return(nil, nil)
	searcher.On("SearchVideos", mock.Anything, "Go プログラミング チュートリアル", 10, "ja").
		Return(nil, errors.New("api error"))
	found := []model.YouTubeVideo{{VideoID: "rust1"}}
	searcher.On("SearchVideos", mock.Anything, "Rust プログラミング チュートリアル", 10, "ja").Return(found, nil)
	videos.On("UpsertVideos", mock.Anything, found).Return(nil)
	videos.On("SaveSearchCache", mock.Anything, mock.Anything).Return(nil)

	w := doRequest(r, http.MethodGet, "/youtube/recommend", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"rust1"`)
	searcher.AssertExpectations(t)
}

func TestYouTubeRecommend_UserNotFound(t *testing.T) {
	h, users, videos, _ := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/recommend", h.Recommend)

	users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/youtube/recommend", nil)
	assertStatus(t, w, http.StatusNotFound)
	videos.AssertNotCalled(t, "FindCachedSearch", mock.Anything, mock.Anything, mock.Anything)
	users.AssertExpectations(t)
}

func TestYouTubeRecommend_Unavailable(t *testing.T) {
	h, users, _ := newUnavailableYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/recommend", h.Recommend)

	w := doRequest(r, http.MethodGet, "/youtube/recommend", nil)
	assertStatus(t, w, http.StatusServiceUnavailable)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ---------- Status ----------

func TestYouTubeStatus_Available(t *testing.T) {
	h, _, _, _ := newTestYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/status", h.Status)

	w := doRequest(r, http.MethodGet, "/youtube/status", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"available":true`)
}

func TestYouTubeStatus_Unavailable(t *testing.T) {
	h, _, _ := newUnavailableYouTubeHandler()
	r := newRouter(1)
	r.GET("/youtube/status", h.Status)

	w := doRequest(r, http.MethodGet, "/youtube/status", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"available":false`)
}

// assertJSONEqual はレスポンス JSON の指定キーが期待値と一致することを検証する。
func assertJSONEqual(t *testing.T, body map[string]interface{}, key string, expected interface{}) {
	t.Helper()
	if body[key] != expected {
		t.Errorf("expected %s=%v, got %v", key, expected, body[key])
	}
}
