package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// bookReviewStatsRepository は [repository.BookReviewStatsRepository] の sqlc(pgx) 実装。
type bookReviewStatsRepository struct {
	q *sqlcgen.Queries
}

// NewBookReviewStatsRepository は BookReviewStatsRepository の sqlc(pgx) 実装を返す。
func NewBookReviewStatsRepository(q *sqlcgen.Queries) repository.BookReviewStatsRepository {
	return &bookReviewStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BookReviewStatsRepository = (*bookReviewStatsRepository)(nil)

// GetBookReviewStats は指定ユーザーの書籍レビュー集計統計を返す。
// レビューが0件の場合はCOALESCEにより全項目0が返る。
func (r *bookReviewStatsRepository) GetBookReviewStats(ctx context.Context, userID uint) (*model.BookReviewStats, error) {
	row, err := r.q.GetBookReviewStats(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return &model.BookReviewStats{
		TotalReviews:  row.TotalReviews,
		AverageRating: row.AverageRating,
		MaxRating:     int(row.MaxRating),
		MinRating:     int(row.MinRating),
		FiveStarCount: row.FiveStarCount,
	}, nil
}
