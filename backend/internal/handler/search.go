package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// PostSearchService は投稿の高度な検索サービスのインターフェース。
type PostSearchService interface {
	SearchPosts(params service.PostSearchParams) (*service.PostSearchResult, error)
}

// StudyCircleSearchRepository はスタディサークル検索のリポジトリインターフェース
type StudyCircleSearchRepository interface {
	Search(query string, limit, offset int) (interface{}, int64, error)
}

// SearchHandler is the handler for search operations
type SearchHandler struct {
	searchService PostSearchService
	circleRepo    StudyCircleSearchRepository
}

// NewSearchHandler creates a new search handler instance
func NewSearchHandler(searchService PostSearchService, circleRepo StudyCircleSearchRepository) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		circleRepo:    circleRepo,
	}
}

// SearchPosts handles post search requests with advanced filtering
func (h *SearchHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		respondBadRequest(c, "query parameter 'q' is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// タグフィルター（カンマ区切り or 複数パラメータ対応）
	var tags []string
	if rawTags := c.Query("tags"); rawTags != "" {
		for _, t := range strings.Split(rawTags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	// ソート順
	sortBy := service.SearchSortBy(c.DefaultQuery("sort_by", string(service.SearchSortByLatest)))

	// 日付範囲フィルター
	var dateFrom, dateTo *time.Time
	if raw := c.Query("date_from"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			dateFrom = &t
		}
	}
	if raw := c.Query("date_to"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			dateTo = &t
		}
	}

	params := service.PostSearchParams{
		Query:    query,
		Tags:     tags,
		SortBy:   sortBy,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Limit:    limit,
		Offset:   offset,
	}

	result, err := h.searchService.SearchPosts(params)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, result)
}

// SearchCircles handles study circle search requests
func (h *SearchHandler) SearchCircles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		respondBadRequest(c, "query parameter 'q' is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	results, _, err := h.circleRepo.Search(query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, results)
}

// PostSearchRepository は handler/search_test.go で使われる後方互換インターフェース。
// 既存のhandlerテストとの互換性のために残す。
type PostSearchRepository interface {
	Search(query string, limit, offset int) (interface{}, int64, error)
}

// legacyPostSearchService は旧インターフェース互換のラッパー。
type legacyPostSearchService struct {
	repo PostSearchRepository
}

func (s *legacyPostSearchService) SearchPosts(params service.PostSearchParams) (*service.PostSearchResult, error) {
	results, total, err := s.repo.Search(params.Query, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	var posts []model.Post
	if p, ok := results.([]model.Post); ok {
		posts = p
	}

	return &service.PostSearchResult{
		Posts:  posts,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}, nil
}

// NewSearchHandlerWithRepo は旧インターフェース（リポジトリ直接）からSearchHandlerを生成する。
// di/container.goからの後方互換のために提供する。
func NewSearchHandlerWithRepo(postRepo PostSearchRepository, circleRepo StudyCircleSearchRepository) *SearchHandler {
	svc := &legacyPostSearchService{repo: postRepo}
	return &SearchHandler{
		searchService: svc,
		circleRepo:    circleRepo,
	}
}
