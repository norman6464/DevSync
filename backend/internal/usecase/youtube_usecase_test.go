package usecase_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockYouTubeVideoRepo は usecase/repository.YouTubeVideoRepository のモック。
type mockYouTubeVideoRepo struct{ mock.Mock }

func (m *mockYouTubeVideoRepo) UpsertVideos(ctx context.Context, videos []model.YouTubeVideo) error {
	return m.Called(ctx, videos).Error(0)
}

func (m *mockYouTubeVideoRepo) FindByVideoIDs(ctx context.Context, videoIDs []string) ([]model.YouTubeVideo, error) {
	args := m.Called(ctx, videoIDs)
	v, _ := args.Get(0).([]model.YouTubeVideo)
	return v, args.Error(1)
}

func (m *mockYouTubeVideoRepo) FindCachedSearch(ctx context.Context, query, language string) (*model.YouTubeSearchCache, error) {
	args := m.Called(ctx, query, language)
	c, _ := args.Get(0).(*model.YouTubeSearchCache)
	return c, args.Error(1)
}

func (m *mockYouTubeVideoRepo) SaveSearchCache(ctx context.Context, cache *model.YouTubeSearchCache) error {
	return m.Called(ctx, cache).Error(0)
}

// mockYouTubeVideoSearcher は usecase/repository.YouTubeVideoSearcher のモック。
type mockYouTubeVideoSearcher struct{ mock.Mock }

func (m *mockYouTubeVideoSearcher) SearchVideos(ctx context.Context, query string, maxResults int, language string) ([]model.YouTubeVideo, error) {
	args := m.Called(ctx, query, maxResults, language)
	v, _ := args.Get(0).([]model.YouTubeVideo)
	return v, args.Error(1)
}

// youtubeVideos は連番の動画一覧を作るテスト用ヘルパー。
func youtubeVideos(prefix string, count int) []model.YouTubeVideo {
	videos := make([]model.YouTubeVideo, 0, count)
	for i := 0; i < count; i++ {
		videos = append(videos, model.YouTubeVideo{VideoID: fmt.Sprintf("%s-%d", prefix, i)})
	}
	return videos
}

func TestRecommendYouTubeVideosUseCase(t *testing.T) {
	t.Run("スキルは先頭 3 件までしか使わない", func(t *testing.T) {
		users := new(mockAccountLinker)
		videos := new(mockYouTubeVideoRepo)
		searcher := new(mockYouTubeVideoSearcher)
		uc := usecase.NewRecommendYouTubeVideosUseCase(users, videos, searcher)

		users.On("FindByID", mock.Anything, uint(1)).
			Return(&model.User{ID: 1, SkillsLanguages: "Go,Rust,TypeScript,Python", SkillsFrameworks: "React"}, nil)
		videos.On("FindCachedSearch", mock.Anything, mock.Anything, "ja").Return(nil, nil)
		for _, skill := range []string{"Go", "Rust", "TypeScript"} {
			searcher.On("SearchVideos", mock.Anything, skill+" プログラミング チュートリアル", 10, "ja").
				Return([]model.YouTubeVideo{}, nil)
		}

		_, skills, err := uc.Execute(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, []string{"Go", "Rust", "TypeScript"}, skills)
		searcher.AssertNumberOfCalls(t, "SearchVideos", 3)
		searcher.AssertExpectations(t)
	})

	t.Run("スキル間で重複した動画は 1 件にまとめる", func(t *testing.T) {
		users := new(mockAccountLinker)
		videos := new(mockYouTubeVideoRepo)
		searcher := new(mockYouTubeVideoSearcher)
		uc := usecase.NewRecommendYouTubeVideosUseCase(users, videos, searcher)

		shared := []model.YouTubeVideo{{VideoID: "same"}, {VideoID: "only-a"}}
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, SkillsLanguages: "Go,Rust"}, nil)
		videos.On("FindCachedSearch", mock.Anything, mock.Anything, "ja").Return(nil, nil)
		videos.On("UpsertVideos", mock.Anything, mock.Anything).Return(nil)
		videos.On("SaveSearchCache", mock.Anything, mock.Anything).Return(nil)
		searcher.On("SearchVideos", mock.Anything, "Go プログラミング チュートリアル", 10, "ja").Return(shared, nil)
		searcher.On("SearchVideos", mock.Anything, "Rust プログラミング チュートリアル", 10, "ja").
			Return([]model.YouTubeVideo{{VideoID: "same"}, {VideoID: "only-b"}}, nil)

		got, _, err := uc.Execute(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, []string{"same", "only-a", "only-b"}, videoIDsOf(got))
	})

	t.Run("動画は 15 件までに切り詰める", func(t *testing.T) {
		users := new(mockAccountLinker)
		videos := new(mockYouTubeVideoRepo)
		searcher := new(mockYouTubeVideoSearcher)
		uc := usecase.NewRecommendYouTubeVideosUseCase(users, videos, searcher)

		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, SkillsLanguages: "Go,Rust"}, nil)
		videos.On("FindCachedSearch", mock.Anything, mock.Anything, "ja").Return(nil, nil)
		videos.On("UpsertVideos", mock.Anything, mock.Anything).Return(nil)
		videos.On("SaveSearchCache", mock.Anything, mock.Anything).Return(nil)
		searcher.On("SearchVideos", mock.Anything, "Go プログラミング チュートリアル", 10, "ja").
			Return(youtubeVideos("go", 10), nil)
		searcher.On("SearchVideos", mock.Anything, "Rust プログラミング チュートリアル", 10, "ja").
			Return(youtubeVideos("rust", 10), nil)

		got, _, err := uc.Execute(context.Background(), 1)
		require.NoError(t, err)
		assert.Len(t, got, 15)
	})

	t.Run("検索クライアント未設定なら 503 でユーザーも読まない", func(t *testing.T) {
		users := new(mockAccountLinker)
		videos := new(mockYouTubeVideoRepo)
		var searcher repository.YouTubeVideoSearcher
		uc := usecase.NewRecommendYouTubeVideosUseCase(users, videos, searcher)

		_, _, err := uc.Execute(context.Background(), 1)
		require.Error(t, err)
		require.NotNil(t, domain.GetDomainError(err))
		assert.Equal(t, domain.ErrCodeServiceUnavailable, domain.GetDomainError(err).Code)
		users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})
}

func TestCheckYouTubeAvailabilityUseCase(t *testing.T) {
	t.Run("検索クライアントがあれば利用可能", func(t *testing.T) {
		assert.True(t, usecase.NewCheckYouTubeAvailabilityUseCase(new(mockYouTubeVideoSearcher)).Execute())
	})

	t.Run("検索クライアントが無ければ利用不可", func(t *testing.T) {
		var searcher repository.YouTubeVideoSearcher
		assert.False(t, usecase.NewCheckYouTubeAvailabilityUseCase(searcher).Execute())
	})
}

func TestSearchYouTubeVideosUseCase_CacheKey(t *testing.T) {
	videos := new(mockYouTubeVideoRepo)
	searcher := new(mockYouTubeVideoSearcher)
	uc := usecase.NewSearchYouTubeVideosUseCase(videos, searcher)

	// キャッシュのキーは小文字化・前後の空白を落としたもの、外部 API へは入力そのまま。
	videos.On("FindCachedSearch", mock.Anything, "golang tutorial", "ja").Return(nil, nil)
	searcher.On("SearchVideos", mock.Anything, " GoLang Tutorial ", 10, "ja").
		Return([]model.YouTubeVideo{{VideoID: "a"}}, nil)
	videos.On("UpsertVideos", mock.Anything, mock.Anything).Return(nil)
	videos.On("SaveSearchCache", mock.Anything, mock.MatchedBy(func(c *model.YouTubeSearchCache) bool {
		return c.Query == "golang tutorial" && strings.Contains(c.VideoIDs, "a")
	})).Return(nil)

	_, cached, err := uc.Execute(context.Background(), " GoLang Tutorial ", "")
	require.NoError(t, err)
	assert.False(t, cached)
	videos.AssertExpectations(t)
	searcher.AssertExpectations(t)
}

// videoIDsOf は動画一覧から VideoID だけを取り出す。
func videoIDsOf(videos []model.YouTubeVideo) []string {
	ids := make([]string, 0, len(videos))
	for _, v := range videos {
		ids = append(ids, v.VideoID)
	}
	return ids
}
