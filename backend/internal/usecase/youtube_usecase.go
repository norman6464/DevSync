package usecase

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// YouTube 検索の既定値。
const (
	youtubeSearchCacheTTL  = 24 * time.Hour
	youtubeMaxResults      = 10
	youtubeRecommendMax    = 15
	youtubeRecommendSkills = 3
)

// msgYouTubeUnavailable は API キー未設定のときに返すメッセージ。
const msgYouTubeUnavailable = "YouTube APIが設定されていません"

// SearchYouTubeVideosUseCase はキーワードで動画を検索する（DB キャッシュ優先）。
type SearchYouTubeVideosUseCase struct {
	videos   repository.YouTubeVideoRepository
	searcher repository.YouTubeVideoSearcher
}

// NewSearchYouTubeVideosUseCase は SearchYouTubeVideosUseCase を生成する。
func NewSearchYouTubeVideosUseCase(
	videos repository.YouTubeVideoRepository,
	searcher repository.YouTubeVideoSearcher,
) *SearchYouTubeVideosUseCase {
	return &SearchYouTubeVideosUseCase{videos: videos, searcher: searcher}
}

// Execute は動画一覧とキャッシュヒットしたかどうかを返す。
func (uc *SearchYouTubeVideosUseCase) Execute(ctx context.Context, query, language string) ([]model.YouTubeVideo, bool, error) {
	return searchYouTubeVideos(ctx, uc.videos, uc.searcher, query, language)
}

// searchYouTubeVideos はキャッシュを見てから検索し、結果をキャッシュへ書き戻す。
// キャッシュの保存に失敗しても検索結果は返す（ログのみ）。
func searchYouTubeVideos(
	ctx context.Context,
	videos repository.YouTubeVideoRepository,
	searcher repository.YouTubeVideoSearcher,
	query, language string,
) ([]model.YouTubeVideo, bool, error) {
	if searcher == nil {
		return nil, false, domain.NewError(domain.ErrCodeServiceUnavailable, msgYouTubeUnavailable, nil)
	}

	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	if normalizedQuery == "" {
		return nil, false, domain.NewError(domain.ErrCodeValidation, "検索キーワードを入力してください", nil)
	}
	if language == "" {
		language = "ja"
	}
	if err := domain.ValidateLanguageCode(language); err != nil {
		return nil, false, err
	}

	if cached, ok := cachedYouTubeVideos(ctx, videos, normalizedQuery, language); ok {
		return cached, true, nil
	}

	found, err := searcher.SearchVideos(ctx, query, youtubeMaxResults, language)
	if err != nil {
		return nil, false, err
	}
	if len(found) > 0 {
		cacheYouTubeVideos(ctx, videos, normalizedQuery, language, found)
	}
	return found, false, nil
}

// cachedYouTubeVideos は有効なキャッシュがあれば動画一覧を返す。
// キャッシュの読み取りに失敗した場合は「キャッシュ無し」として扱う（移行前と同じ）。
func cachedYouTubeVideos(
	ctx context.Context,
	videos repository.YouTubeVideoRepository,
	normalizedQuery, language string,
) ([]model.YouTubeVideo, bool) {
	cache, err := videos.FindCachedSearch(ctx, normalizedQuery, language)
	if err != nil || cache == nil {
		return nil, false
	}

	cachedVideos, err := videos.FindByVideoIDs(ctx, strings.Split(cache.VideoIDs, ","))
	if err != nil || len(cachedVideos) == 0 {
		return nil, false
	}
	return cachedVideos, true
}

// cacheYouTubeVideos は検索結果と検索キャッシュを保存する。失敗してもログのみで続行する。
func cacheYouTubeVideos(
	ctx context.Context,
	videos repository.YouTubeVideoRepository,
	normalizedQuery, language string,
	found []model.YouTubeVideo,
) {
	if err := videos.UpsertVideos(ctx, found); err != nil {
		log.Printf("YouTube動画キャッシュ保存失敗: %v", err)
	}

	videoIDs := make([]string, len(found))
	for i, v := range found {
		videoIDs[i] = v.VideoID
	}

	cache := &model.YouTubeSearchCache{
		Query:        normalizedQuery,
		Language:     language,
		VideoIDs:     strings.Join(videoIDs, ","),
		CacheExpires: time.Now().Add(youtubeSearchCacheTTL),
	}
	if err := videos.SaveSearchCache(ctx, cache); err != nil {
		log.Printf("YouTube検索キャッシュ保存失敗: %v", err)
	}
}

// RecommendYouTubeVideosUseCase はユーザーのスキルに基づくおすすめ動画を返す。
type RecommendYouTubeVideosUseCase struct {
	users    repository.UserSkillsReader
	videos   repository.YouTubeVideoRepository
	searcher repository.YouTubeVideoSearcher
}

// NewRecommendYouTubeVideosUseCase は RecommendYouTubeVideosUseCase を生成する。
func NewRecommendYouTubeVideosUseCase(
	users repository.UserSkillsReader,
	videos repository.YouTubeVideoRepository,
	searcher repository.YouTubeVideoSearcher,
) *RecommendYouTubeVideosUseCase {
	return &RecommendYouTubeVideosUseCase{users: users, videos: videos, searcher: searcher}
}

// Execute はプロフィールのスキル（最大 3 件）で検索し、重複を除いた動画と使ったスキルを返す。
// スキルごとの検索失敗は無視して次のスキルへ進む（移行前と同じ）。
func (uc *RecommendYouTubeVideosUseCase) Execute(ctx context.Context, userID uint) ([]model.YouTubeVideo, []string, error) {
	if uc.searcher == nil {
		return nil, nil, domain.NewError(domain.ErrCodeServiceUnavailable, msgYouTubeUnavailable, nil)
	}

	user, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errOwnedEntityNotFound
	}

	// スキル文字列の分解はレコメンドスライスと同じものを使う。
	skills := ParseSkills(user.SkillsLanguages, user.SkillsFrameworks)
	if len(skills) == 0 {
		return []model.YouTubeVideo{}, []string{}, nil
	}
	if len(skills) > youtubeRecommendSkills {
		skills = skills[:youtubeRecommendSkills]
	}

	var all []model.YouTubeVideo
	seen := make(map[string]bool)
	for _, skill := range skills {
		found, _, err := searchYouTubeVideos(ctx, uc.videos, uc.searcher, skill+" プログラミング チュートリアル", "ja")
		if err != nil {
			log.Printf("おすすめ検索失敗 (skill=%s): %v", skill, err)
			continue
		}
		for _, v := range found {
			if !seen[v.VideoID] {
				seen[v.VideoID] = true
				all = append(all, v)
			}
		}
	}

	if len(all) > youtubeRecommendMax {
		all = all[:youtubeRecommendMax]
	}
	return all, skills, nil
}

// CheckYouTubeAvailabilityUseCase は YouTube 連携が利用可能かどうかを返す。
type CheckYouTubeAvailabilityUseCase struct {
	searcher repository.YouTubeVideoSearcher
}

// NewCheckYouTubeAvailabilityUseCase は CheckYouTubeAvailabilityUseCase を生成する。
func NewCheckYouTubeAvailabilityUseCase(searcher repository.YouTubeVideoSearcher) *CheckYouTubeAvailabilityUseCase {
	return &CheckYouTubeAvailabilityUseCase{searcher: searcher}
}

// Execute は検索クライアントが注入されているかどうかを返す。
// API キー未設定のときは DI が nil を渡すため false になる。
func (uc *CheckYouTubeAvailabilityUseCase) Execute() bool {
	return uc.searcher != nil
}
