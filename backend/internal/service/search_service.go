package service

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
)

const (
	defaultSearchLimit = 20
	maxSearchLimit     = 100
)

// SearchSortBy は検索結果のソート順を表す。
type SearchSortBy string

const (
	SearchSortByLatest  SearchSortBy = "latest"  // 最新順
	SearchSortByPopular SearchSortBy = "popular" // 人気順（いいね数）
	SearchSortByViews   SearchSortBy = "views"   // 閲覧数順
)

// PostSearchParams は投稿検索のパラメータ。
type PostSearchParams struct {
	Query    string
	Tags     []string
	SortBy   SearchSortBy
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}

// PostSearchResult は投稿検索結果のレスポンス。
type PostSearchResult struct {
	Posts  []model.Post `json:"posts"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

// PostAdvancedSearchRepo は投稿の高度な検索リポジトリのインターフェース。
type PostAdvancedSearchRepo interface {
	SearchWithFilter(query string, tags []string, sortBy string, dateFrom, dateTo *time.Time, limit, offset int) ([]model.Post, int64, error)
}

// SearchService は検索に関するビジネスロジックを提供する。
type SearchService struct {
	postSearchRepo PostAdvancedSearchRepo
}

// NewSearchService は新しいSearchServiceインスタンスを生成する。
func NewSearchService(postSearchRepo PostAdvancedSearchRepo) *SearchService {
	return &SearchService{postSearchRepo: postSearchRepo}
}

// SearchPosts は投稿の高度な検索を実行する。
// タグフィルター・日付範囲・ソート順に対応する。
func (s *SearchService) SearchPosts(params PostSearchParams) (*PostSearchResult, error) {
	if params.Query == "" {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "検索クエリは必須です", nil)
	}

	// limitの正規化
	limit := params.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
	} else if limit > maxSearchLimit {
		limit = maxSearchLimit
	}

	// ソート順の正規化
	sortBy := string(params.SortBy)
	if sortBy == "" {
		sortBy = string(SearchSortByLatest)
	}

	posts, total, err := s.postSearchRepo.SearchWithFilter(
		params.Query,
		params.Tags,
		sortBy,
		params.DateFrom,
		params.DateTo,
		limit,
		params.Offset,
	)
	if err != nil {
		return nil, err
	}

	return &PostSearchResult{
		Posts:  posts,
		Total:  total,
		Limit:  limit,
		Offset: params.Offset,
	}, nil
}
