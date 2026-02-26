package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// モック定義
// ============================================================

type MockYouTubeClient struct{ mock.Mock }

func (m *MockYouTubeClient) SearchVideos(query string, maxResults int, language string) ([]model.YouTubeVideo, error) {
	args := m.Called(query, maxResults, language)
	return args.Get(0).([]model.YouTubeVideo), args.Error(1)
}

type MockYouTubeVideoRepository struct{ mock.Mock }

func (m *MockYouTubeVideoRepository) UpsertVideos(videos []model.YouTubeVideo) error {
	args := m.Called(videos)
	return args.Error(0)
}

func (m *MockYouTubeVideoRepository) FindByVideoIDs(videoIDs []string) ([]model.YouTubeVideo, error) {
	args := m.Called(videoIDs)
	return args.Get(0).([]model.YouTubeVideo), args.Error(1)
}

func (m *MockYouTubeVideoRepository) FindCachedSearch(query, language string) (*model.YouTubeSearchCache, error) {
	args := m.Called(query, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.YouTubeSearchCache), args.Error(1)
}

func (m *MockYouTubeVideoRepository) SaveSearchCache(cache *model.YouTubeSearchCache) error {
	args := m.Called(cache)
	return args.Error(0)
}

// ============================================================
// ヘルパー
// ============================================================

func newTestYouTubeService(client YouTubeClientInterface) (*YouTubeService, *MockYouTubeVideoRepository, *MockUserRepository) {
	repo := new(MockYouTubeVideoRepository)
	userRepo := new(MockUserRepository)
	svc := NewYouTubeService(repo, userRepo, client)
	return svc, repo, userRepo
}

// ============================================================
// Search テスト
// ============================================================

func TestYouTubeSearch_CacheHit(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, _ := newTestYouTubeService(mockClient)

	cache := &model.YouTubeSearchCache{
		Query:    "go tutorial",
		Language: "ja",
		VideoIDs: "abc123,def456",
	}
	repo.On("FindCachedSearch", "go tutorial", "ja").Return(cache, nil)

	videos := []model.YouTubeVideo{
		{VideoID: "abc123", Title: "Go Tutorial 1"},
		{VideoID: "def456", Title: "Go Tutorial 2"},
	}
	repo.On("FindByVideoIDs", []string{"abc123", "def456"}).Return(videos, nil)

	result, cached, err := svc.Search("Go Tutorial", "ja")
	assert.NoError(t, err)
	assert.True(t, cached)
	assert.Len(t, result, 2)
	mockClient.AssertNotCalled(t, "SearchVideos")
}

func TestYouTubeSearch_CacheMiss(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, _ := newTestYouTubeService(mockClient)

	repo.On("FindCachedSearch", "react hooks", "ja").Return(nil, errors.New("not found"))

	apiVideos := []model.YouTubeVideo{
		{VideoID: "vid1", Title: "React Hooks Tutorial"},
	}
	mockClient.On("SearchVideos", "React Hooks", youtubeMaxResults, "ja").Return(apiVideos, nil)
	repo.On("UpsertVideos", apiVideos).Return(nil)
	repo.On("SaveSearchCache", mock.AnythingOfType("*model.YouTubeSearchCache")).Return(nil)

	result, cached, err := svc.Search("React Hooks", "ja")
	assert.NoError(t, err)
	assert.False(t, cached)
	assert.Len(t, result, 1)
	assert.Equal(t, "React Hooks Tutorial", result[0].Title)
	mockClient.AssertCalled(t, "SearchVideos", "React Hooks", youtubeMaxResults, "ja")
}

func TestYouTubeSearch_EmptyQuery(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, _, _ := newTestYouTubeService(mockClient)

	_, _, err := svc.Search("", "ja")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "検索キーワード")
}

func TestYouTubeSearch_NilClient(t *testing.T) {
	svc, _, _ := newTestYouTubeService(nil)

	_, _, err := svc.Search("Go", "ja")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "YouTube APIが設定されていません")
}

func TestYouTubeSearch_APIError(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, _ := newTestYouTubeService(mockClient)

	repo.On("FindCachedSearch", "python", "ja").Return(nil, errors.New("not found"))
	mockClient.On("SearchVideos", "python", youtubeMaxResults, "ja").Return([]model.YouTubeVideo{}, errors.New("API error"))

	_, _, err := svc.Search("python", "ja")
	assert.Error(t, err)
}

func TestYouTubeSearch_InvalidLanguage(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, _, _ := newTestYouTubeService(mockClient)

	tests := []struct {
		name string
		lang string
	}{
		{"不正な言語コード", "xx"},
		{"インジェクション試行", "ja&key=val"},
		{"長い文字列", "japanese"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.Search("Go tutorial", tt.lang)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "サポートされていない言語コード")
			mockClient.AssertNotCalled(t, "SearchVideos")
		})
	}
}

// ============================================================
// GetRecommendations テスト
// ============================================================

func TestYouTubeGetRecommendations_Success(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, userRepo := newTestYouTubeService(mockClient)

	user := &model.User{SkillsLanguages: "Go", SkillsFrameworks: ""}
	userRepo.On("FindByID", uint(1)).Return(user, nil)

	// "Go プログラミング チュートリアル" の検索
	repo.On("FindCachedSearch", mock.Anything, "ja").Return(nil, errors.New("not found"))
	apiVideos := []model.YouTubeVideo{
		{VideoID: "go1", Title: "Go チュートリアル"},
	}
	mockClient.On("SearchVideos", mock.Anything, youtubeMaxResults, "ja").Return(apiVideos, nil)
	repo.On("UpsertVideos", mock.Anything).Return(nil)
	repo.On("SaveSearchCache", mock.AnythingOfType("*model.YouTubeSearchCache")).Return(nil)

	videos, skills, err := svc.GetRecommendations(1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Go"}, skills)
	assert.NotEmpty(t, videos)
}

func TestYouTubeGetRecommendations_NoSkills(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, _, userRepo := newTestYouTubeService(mockClient)

	user := &model.User{SkillsLanguages: "", SkillsFrameworks: ""}
	userRepo.On("FindByID", uint(1)).Return(user, nil)

	videos, skills, err := svc.GetRecommendations(1)
	assert.NoError(t, err)
	assert.Empty(t, skills)
	assert.Empty(t, videos)
	mockClient.AssertNotCalled(t, "SearchVideos")
}

func TestYouTubeGetRecommendations_UserNotFound(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, _, userRepo := newTestYouTubeService(mockClient)

	userRepo.On("FindByID", uint(1)).Return((*model.User)(nil), errors.New("not found"))

	_, _, err := svc.GetRecommendations(1)
	assert.Error(t, err)
}

func TestYouTubeGetRecommendations_NilClient(t *testing.T) {
	svc, _, _ := newTestYouTubeService(nil)

	_, _, err := svc.GetRecommendations(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "YouTube APIが設定されていません")
}

func TestYouTubeGetRecommendations_SearchError(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, userRepo := newTestYouTubeService(mockClient)

	user := &model.User{SkillsLanguages: "Go,Rust", SkillsFrameworks: ""}
	userRepo.On("FindByID", uint(1)).Return(user, nil)

	// Go検索は失敗
	repo.On("FindCachedSearch", "go プログラミング チュートリアル", "ja").Return(nil, errors.New("not found"))
	mockClient.On("SearchVideos", "Go プログラミング チュートリアル", youtubeMaxResults, "ja").Return([]model.YouTubeVideo{}, errors.New("api error"))

	// Rust検索は成功
	repo.On("FindCachedSearch", "rust プログラミング チュートリアル", "ja").Return(nil, errors.New("not found"))
	rustVideos := []model.YouTubeVideo{
		{VideoID: "rust1", Title: "Rust チュートリアル"},
	}
	mockClient.On("SearchVideos", "Rust プログラミング チュートリアル", youtubeMaxResults, "ja").Return(rustVideos, nil)
	repo.On("UpsertVideos", mock.Anything).Return(nil)
	repo.On("SaveSearchCache", mock.AnythingOfType("*model.YouTubeSearchCache")).Return(nil)

	videos, skills, err := svc.GetRecommendations(1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Go", "Rust"}, skills)
	assert.Len(t, videos, 1)
	assert.Equal(t, "rust1", videos[0].VideoID)
}

// ============================================================
// IsAvailable テスト
// ============================================================

func TestYouTubeIsAvailable_WithClient(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, _, _ := newTestYouTubeService(mockClient)
	assert.True(t, svc.IsAvailable())
}

func TestYouTubeIsAvailable_WithoutClient(t *testing.T) {
	svc, _, _ := newTestYouTubeService(nil)
	assert.False(t, svc.IsAvailable())
}

// ============================================================
// Search 追加テスト（カバレッジ向上）
// ============================================================

// TestYouTubeSearch_CacheHitButVideosFetchFailed はキャッシュエントリが存在するが
// 動画フェッチが失敗した場合にAPIフォールバックが実行されることを確認する。
func TestYouTubeSearch_CacheHitButVideosFetchFailed(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, _ := newTestYouTubeService(mockClient)

	cache := &model.YouTubeSearchCache{
		Query:    "vue js",
		Language: "ja",
		VideoIDs: "vid1,vid2",
	}
	repo.On("FindCachedSearch", "vue js", "ja").Return(cache, nil)
	// キャッシュに対応する動画取得が失敗 → APIフォールバックへ
	repo.On("FindByVideoIDs", []string{"vid1", "vid2"}).Return([]model.YouTubeVideo{}, errors.New("db error"))

	apiVideos := []model.YouTubeVideo{
		{VideoID: "vid1", Title: "Vue.js Tutorial"},
	}
	mockClient.On("SearchVideos", "vue js", youtubeMaxResults, "ja").Return(apiVideos, nil)
	repo.On("UpsertVideos", apiVideos).Return(nil)
	repo.On("SaveSearchCache", mock.AnythingOfType("*model.YouTubeSearchCache")).Return(nil)

	result, cached, err := svc.Search("vue js", "ja")
	assert.NoError(t, err)
	assert.False(t, cached)
	assert.Len(t, result, 1)
	mockClient.AssertCalled(t, "SearchVideos", "vue js", youtubeMaxResults, "ja")
}

// TestYouTubeSearch_EmptyAPIResult はAPIが0件を返した場合、
// キャッシュ保存をスキップして空のスライスを返すことを確認する。
func TestYouTubeSearch_EmptyAPIResult(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, _ := newTestYouTubeService(mockClient)

	repo.On("FindCachedSearch", "unknown query xyz", "ja").Return(nil, errors.New("not found"))
	mockClient.On("SearchVideos", "unknown query xyz", youtubeMaxResults, "ja").Return([]model.YouTubeVideo{}, nil)

	result, cached, err := svc.Search("unknown query xyz", "ja")
	assert.NoError(t, err)
	assert.False(t, cached)
	assert.Empty(t, result)
	// 0件の場合はUpsertVideos/SaveSearchCacheを呼ばない
	repo.AssertNotCalled(t, "UpsertVideos")
	repo.AssertNotCalled(t, "SaveSearchCache")
}

// TestYouTubeSearch_UpsertVideosFails はキャッシュ保存が失敗しても動画を正常に返すことを確認する。
// UpsertVideosのエラーはログ出力のみで処理が続行される。
func TestYouTubeSearch_UpsertVideosFails(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, _ := newTestYouTubeService(mockClient)

	repo.On("FindCachedSearch", "typescript", "ja").Return(nil, errors.New("not found"))
	apiVideos := []model.YouTubeVideo{
		{VideoID: "ts1", Title: "TypeScript Tutorial"},
	}
	mockClient.On("SearchVideos", "typescript", youtubeMaxResults, "ja").Return(apiVideos, nil)
	// UpsertVideosがエラーを返しても処理継続（ログのみ）
	repo.On("UpsertVideos", apiVideos).Return(errors.New("db error"))
	repo.On("SaveSearchCache", mock.AnythingOfType("*model.YouTubeSearchCache")).Return(nil)

	result, cached, err := svc.Search("typescript", "ja")
	assert.NoError(t, err)
	assert.False(t, cached)
	assert.Len(t, result, 1)
}

// TestYouTubeSearch_SaveCacheFails はキャッシュ保存失敗時も動画を正常に返すことを確認する。
// SaveSearchCacheのエラーはログ出力のみで処理が続行される。
func TestYouTubeSearch_SaveCacheFails(t *testing.T) {
	mockClient := new(MockYouTubeClient)
	svc, repo, _ := newTestYouTubeService(mockClient)

	repo.On("FindCachedSearch", "rust lang", "ja").Return(nil, errors.New("not found"))
	apiVideos := []model.YouTubeVideo{
		{VideoID: "rust1", Title: "Rust Programming"},
	}
	mockClient.On("SearchVideos", "rust lang", youtubeMaxResults, "ja").Return(apiVideos, nil)
	repo.On("UpsertVideos", apiVideos).Return(nil)
	// SaveSearchCacheがエラーを返しても処理継続（ログのみ）
	repo.On("SaveSearchCache", mock.AnythingOfType("*model.YouTubeSearchCache")).Return(errors.New("cache error"))

	result, cached, err := svc.Search("rust lang", "ja")
	assert.NoError(t, err)
	assert.False(t, cached)
	assert.Len(t, result, 1)
}
