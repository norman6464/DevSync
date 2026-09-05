package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postSearchRepository は [repository.PostSearchRepository] の sqlc(pgx) 実装。
type postSearchRepository struct {
	q *sqlcgen.Queries
}

// NewPostSearchRepository は PostSearchRepository の sqlc(pgx) 実装を返す。
func NewPostSearchRepository(q *sqlcgen.Queries) repository.PostSearchRepository {
	return &postSearchRepository{q: q}
}

var _ repository.PostSearchRepository = (*postSearchRepository)(nil)

// SearchWithFilter はタグ・日付範囲・ソート順による高度な投稿検索を実行する。
// 下書きは検索対象外。タグは AND 条件で絞り込む。
func (r *postSearchRepository) SearchWithFilter(ctx context.Context, params model.PostSearchParams) ([]model.Post, int64, error) {
	searchPattern := escapeLikePattern(params.Query)
	// nil スライスを渡すと SQL 上は NULL になり、cardinality(NULL) が NULL になって
	// タグ絞り込みなしの意図（cardinality=0との比較）が壊れるため、必ず空スライス以上を渡す。
	tags := params.Tags
	if tags == nil {
		tags = []string{}
	}
	dateFrom := toTimestamptz(params.DateFrom)
	dateTo := toTimestamptz(params.DateTo)

	total, err := r.q.CountPostsWithFilter(ctx, sqlcgen.CountPostsWithFilterParams{
		Title:    searchPattern,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Tags:     tags,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.SearchPostsWithFilter(ctx, sqlcgen.SearchPostsWithFilterParams{
		Title:    searchPattern,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Tags:     tags,
		SortBy:   string(params.SortBy),
		Offset:   int32Param(params.Offset),
		Limit:    int32Param(params.Limit),
	})
	if err != nil {
		return nil, 0, err
	}

	posts := make([]model.Post, len(rows))
	postIDs := make([]int64, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
		postIDs[i] = row.Post.ID
	}

	if len(postIDs) > 0 {
		snippetRows, err := r.q.ListCodeSnippetsByPostIDs(ctx, postIDs)
		if err != nil {
			return nil, 0, err
		}
		snippetsByPostID := make(map[uint][]model.CodeSnippet)
		for _, row := range snippetRows {
			postID := uint(row.PostID)
			snippetsByPostID[postID] = append(snippetsByPostID[postID], toModelCodeSnippet(row))
		}
		for i := range posts {
			posts[i].CodeSnippets = snippetsByPostID[posts[i].ID]
		}
	}

	return posts, total, nil
}
