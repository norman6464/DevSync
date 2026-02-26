package service

import (
	"log"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

const (
	youtubeSearchCacheTTL = 24 * time.Hour
	youtubeMaxResults     = 10
	youtubeRecommendMax   = 15
)

// YouTubeService はYouTube動画検索・おすすめのビジネスロジックを提供する。
type YouTubeService struct {
	repo     repository.YouTubeVideoRepositoryInterface
	userRepo repository.UserRepositoryInterface
	client   YouTubeClientInterface
}

// NewYouTubeService は新しいYouTubeServiceインスタンスを生成する。
func NewYouTubeService(
	repo repository.YouTubeVideoRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	client YouTubeClientInterface,
) *YouTubeService {
	return &YouTubeService{repo: repo, userRepo: userRepo, client: client}
}

// Search はキーワードでYouTube動画を検索する（DBキャッシュ優先）。
// 戻り値: 動画一覧, キャッシュヒットかどうか, エラー
func (s *YouTubeService) Search(query, language string) ([]model.YouTubeVideo, bool, error) {
	if s.client == nil {
		return nil, false, domain.NewError(domain.ErrCodeServiceUnavailable, "YouTube APIが設定されていません", nil)
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

	// 1. DBキャッシュを検索
	cache, err := s.repo.FindCachedSearch(normalizedQuery, language)
	if err == nil && cache != nil {
		videoIDs := strings.Split(cache.VideoIDs, ",")
		videos, err := s.repo.FindByVideoIDs(videoIDs)
		if err == nil && len(videos) > 0 {
			return videos, true, nil
		}
	}

	// 2. キャッシュミス → YouTube API呼び出し
	videos, err := s.client.SearchVideos(query, youtubeMaxResults, language)
	if err != nil {
		return nil, false, err
	}

	// 3. 結果をDBにキャッシュ
	if len(videos) > 0 {
		if err := s.repo.UpsertVideos(videos); err != nil {
			log.Printf("YouTube動画キャッシュ保存失敗: %v", err)
		}

		videoIDs := make([]string, len(videos))
		for i, v := range videos {
			videoIDs[i] = v.VideoID
		}

		searchCache := &model.YouTubeSearchCache{
			Query:        normalizedQuery,
			Language:     language,
			VideoIDs:     strings.Join(videoIDs, ","),
			CacheExpires: time.Now().Add(youtubeSearchCacheTTL),
		}
		if err := s.repo.SaveSearchCache(searchCache); err != nil {
			log.Printf("YouTube検索キャッシュ保存失敗: %v", err)
		}
	}

	return videos, false, nil
}

// GetRecommendations はユーザーのスキルに基づくおすすめ動画を返す。
func (s *YouTubeService) GetRecommendations(userID uint) ([]model.YouTubeVideo, []string, error) {
	if s.client == nil {
		return nil, nil, domain.NewError(domain.ErrCodeServiceUnavailable, "YouTube APIが設定されていません", nil)
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}

	skills := parseSkills(user.SkillsLanguages, user.SkillsFrameworks)
	if len(skills) == 0 {
		return []model.YouTubeVideo{}, []string{}, nil
	}

	// 最大3スキルで検索
	maxSkills := 3
	if len(skills) < maxSkills {
		maxSkills = len(skills)
	}
	selectedSkills := skills[:maxSkills]

	var allVideos []model.YouTubeVideo
	seen := make(map[string]bool)

	for _, skill := range selectedSkills {
		query := skill + " プログラミング チュートリアル"
		videos, _, err := s.Search(query, "ja")
		if err != nil {
			log.Printf("おすすめ検索失敗 (skill=%s): %v", skill, err)
			continue
		}
		for _, v := range videos {
			if !seen[v.VideoID] {
				seen[v.VideoID] = true
				allVideos = append(allVideos, v)
			}
		}
	}

	if len(allVideos) > youtubeRecommendMax {
		allVideos = allVideos[:youtubeRecommendMax]
	}

	return allVideos, selectedSkills, nil
}

// IsAvailable はYouTube API機能が利用可能かどうかを返す。
func (s *YouTubeService) IsAvailable() bool {
	return s.client != nil
}
