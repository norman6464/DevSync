package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// commentLikeRepository は [repository.CommentLikeRepository] の sqlc(pgx) 実装。
// Like/Unlike はいいねの作成・削除とコメントのカウンタ更新を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type commentLikeRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewCommentLikeRepository は CommentLikeRepository の sqlc(pgx) 実装を返す。
func NewCommentLikeRepository(pool *pgxpool.Pool) repository.CommentLikeRepository {
	return &commentLikeRepository{pool: pool, q: sqlcgen.New(pool)}
}

var _ repository.CommentLikeRepository = (*commentLikeRepository)(nil)

// Like はいいねを追加し、コメントのいいね数をインクリメントする（1 トランザクション）。
func (r *commentLikeRepository) Like(ctx context.Context, userID, commentID uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if err := q.CreateCommentLike(ctx, sqlcgen.CreateCommentLikeParams{
		UserID:    int64(userID),
		CommentID: int64(commentID),
	}); err != nil {
		return err
	}
	if err := q.IncrementCommentLikeCount(ctx, int64(commentID)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Unlike はいいねを削除し、コメントのいいね数をデクリメントする（1 トランザクション）。
func (r *commentLikeRepository) Unlike(ctx context.Context, userID, commentID uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if err := q.DeleteCommentLike(ctx, sqlcgen.DeleteCommentLikeParams{
		UserID:    int64(userID),
		CommentID: int64(commentID),
	}); err != nil {
		return err
	}
	if err := q.DecrementCommentLikeCount(ctx, int64(commentID)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// HasLiked は指定ユーザーがコメントをいいね済みかを返す。
func (r *commentLikeRepository) HasLiked(ctx context.Context, userID, commentID uint) (bool, error) {
	count, err := r.q.CountCommentLikesByUserAndComment(ctx, sqlcgen.CountCommentLikesByUserAndCommentParams{
		UserID:    int64(userID),
		CommentID: int64(commentID),
	})
	return count > 0, err
}

// CountByCommentID はコメントのいいね数を返す。
func (r *commentLikeRepository) CountByCommentID(ctx context.Context, commentID uint) (int64, error) {
	return r.q.CountCommentLikesByComment(ctx, int64(commentID))
}
