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
func (s *SearchService) SearchPosts(params model.PostSearchParams) (*model.PostSearchResult, error) {
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

	// ソート順の正規化・バリデーション
	sortBy := string(params.SortBy)
	if sortBy == "" {
		sortBy = string(model.SearchSortByLatest)
	} else if !model.ValidSearchSortBy[params.SortBy] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なソート順です", nil)
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

	return &model.PostSearchResult{
		Posts:  posts,
		Total:  total,
		Limit:  limit,
		Offset: params.Offset,
	}, nil
}
