package handler

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// SearchHandler は検索エンドポイントのハンドラ。
type SearchHandler struct {
	postSearch   *usecase.SearchPostsUseCase
	circleSearch *usecase.SearchStudyCirclesUseCase
}

// NewSearchHandler は新しい SearchHandler インスタンスを生成する。
func NewSearchHandler(postSearch *usecase.SearchPostsUseCase, circleSearch *usecase.SearchStudyCirclesUseCase) *SearchHandler {
	return &SearchHandler{
		postSearch:   postSearch,
		circleSearch: circleSearch,
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

	// ソート順（Handler層で無効値をデフォルトに正規化）
	sortBy := model.SearchSortBy(c.DefaultQuery("sort_by", string(model.SearchSortByLatest)))
	if !model.ValidSearchSortBy[sortBy] {
		sortBy = model.SearchSortByLatest
	}

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

	result, err := h.postSearch.Execute(c.Request.Context(), params)
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

	results, _, err := h.circleSearch.Execute(c.Request.Context(), query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(results))
}
