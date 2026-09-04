package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// codeSnippetStatsRepository は [repository.CodeSnippetStatsRepository] の sqlc(pgx) 実装。
type codeSnippetStatsRepository struct {
	q *sqlcgen.Queries
}

// NewCodeSnippetStatsRepository は CodeSnippetStatsRepository の sqlc(pgx) 実装を返す。
func NewCodeSnippetStatsRepository(q *sqlcgen.Queries) repository.CodeSnippetStatsRepository {
	return &codeSnippetStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.CodeSnippetStatsRepository = (*codeSnippetStatsRepository)(nil)

// GetCodeSnippetStats は指定ユーザーのコードスニペット活動集計統計を返す。
func (r *codeSnippetStatsRepository) GetCodeSnippetStats(ctx context.Context, userID uint) (*model.CodeSnippetStats, error) {
	total, err := r.q.CountCodeSnippetsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	comments, err := r.q.SumCodeSnippetCommentCountByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	languages, err := r.q.CountCodeSnippetLanguagesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return &model.CodeSnippetStats{
		TotalSnippets: total,
		TotalComments: comments,
		LanguageCount: languages,
	}, nil
}
