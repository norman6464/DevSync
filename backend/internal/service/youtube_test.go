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
