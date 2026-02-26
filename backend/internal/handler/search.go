package handler

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostSearchService は投稿の高度な検索サービスのインターフェース。
type PostSearchService interface {
	SearchPosts(params model.PostSearchParams) (*model.PostSearchResult, error)
}

// CircleSearchService はスタディサークル検索のサービスインターフェース。
// Service 層を経由することでクリーンアーキテクチャを維持する。
type CircleSearchService interface {
	SearchCircles(query string, limit, offset int) ([]model.StudyCircle, int64, error)
}

// SearchHandler is the handler for search operations
type SearchHandler struct {
	searchService PostSearchService
	circleService CircleSearchService
}

// NewSearchHandler creates a new search handler instance
func NewSearchHandler(searchService PostSearchService, circleService CircleSearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		circleService: circleService,
	}
}

// SearchPosts handles post search requests with advanced filtering
func (h *SearchHandler) SearchPosts(c *gin.Context) {
	query, ok := parseSearchQuery(c)
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	// タグフィルター（カンマ区切り対応）
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
	sortBy := model.SearchSortBy(c.DefaultQuery("sort_by", string(model.SearchSortByLatest)))

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

	params := model.PostSearchParams{
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
	query, ok := parseSearchQuery(c)
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	results, _, err := h.circleService.SearchCircles(query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(results))
}
