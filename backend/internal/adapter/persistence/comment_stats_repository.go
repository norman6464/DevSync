package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// commentStatsRepository は [repository.CommentStatsRepository] の sqlc(pgx) 実装。
type commentStatsRepository struct {
	q *sqlcgen.Queries
}

// NewCommentStatsRepository は CommentStatsRepository の sqlc(pgx) 実装を返す。
func NewCommentStatsRepository(q *sqlcgen.Queries) repository.CommentStatsRepository {
	return &commentStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.CommentStatsRepository = (*commentStatsRepository)(nil)

// GetCommentStats は指定ユーザーのコメント活動集計統計を返す。
func (r *commentStatsRepository) GetCommentStats(ctx context.Context, userID uint) (*model.CommentStats, error) {
	total, err := r.q.CountTopLevelCommentsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	replies, err := r.q.CountRepliesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	received, err := r.q.CountCommentsReceivedByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	startOfMonth := domain.StartOfMonth(time.Now())
	thisMonth, err := r.q.CountCommentsByUserSince(ctx, sqlcgen.CountCommentsByUserSinceParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(startOfMonth),
	})
	if err != nil {
		return nil, err
	}

	return &model.CommentStats{
		TotalComments:     total,
		TotalReplies:      replies,
		CommentsReceived:  received,
		CommentsThisMonth: thisMonth,
	}, nil
}
