package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

const (
	// defaultPostSearchLimit は limit 未指定時の取得件数。
	defaultPostSearchLimit = 20
	// maxPostSearchLimit は 1 回の検索で取得できる上限。
	maxPostSearchLimit = 100
)

// SearchPostsUseCase は投稿の高度な検索を実行する。
type SearchPostsUseCase struct {
	posts repository.PostSearchRepository
}

// NewSearchPostsUseCase は SearchPostsUseCase を生成する。
func NewSearchPostsUseCase(posts repository.PostSearchRepository) *SearchPostsUseCase {
	return &SearchPostsUseCase{posts: posts}
}

// Execute は検索条件を正規化したうえで投稿を検索する。
// 検索クエリは必須。ソート順が不正な場合はエラーを返す。
func (uc *SearchPostsUseCase) Execute(ctx context.Context, params model.PostSearchParams) (*model.PostSearchResult, error) {
	if params.Query == "" {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "検索クエリは必須です", nil)
	}

	if params.Limit <= 0 {
		params.Limit = defaultPostSearchLimit
	} else if params.Limit > maxPostSearchLimit {
		params.Limit = maxPostSearchLimit
	}

	if params.SortBy == "" {
		params.SortBy = model.SearchSortByLatest
	} else if !model.ValidSearchSortBy[params.SortBy] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なソート順です", nil)
	}

	posts, total, err := uc.posts.SearchWithFilter(ctx, params)
	if err != nil {
		return nil, err
	}

	return &model.PostSearchResult{
		Posts:  posts,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}, nil
}
