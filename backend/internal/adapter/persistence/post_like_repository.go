package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postLikeRepository は [repository.PostLikeRepository] の sqlc(pgx) 実装。
type postLikeRepository struct {
	q *sqlcgen.Queries
}

// NewPostLikeRepository は PostLikeRepository の sqlc(pgx) 実装を返す。
func NewPostLikeRepository(q *sqlcgen.Queries) repository.PostLikeRepository {
	return &postLikeRepository{q: q}
}

var _ repository.PostLikeRepository = (*postLikeRepository)(nil)

// Like はいいねを追加し、投稿の like_count を加算する。
// 移行前の GORM 実装と同じくトランザクションでは括らない（元実装も2つの独立した操作だったため）。
func (r *postLikeRepository) Like(ctx context.Context, userID, postID uint) error {
	if err := r.q.CreatePostLike(ctx, sqlcgen.CreatePostLikeParams{
		UserID: int64(userID),
		PostID: int64(postID),
	}); err != nil {
		return err
	}
	return r.q.IncrementPostLikeCount(ctx, int64(postID))
}

// Unlike はいいねを取り消し、実際に削除できたときだけ like_count をデクリメントする。
// 移行前の GORM 実装と同じく、デクリメント自体のエラーは呼び出し元へ返さない。
func (r *postLikeRepository) Unlike(ctx context.Context, userID, postID uint) error {
	rowsAffected, err := r.q.DeletePostLike(ctx, sqlcgen.DeletePostLikeParams{
		UserID: int64(userID),
		PostID: int64(postID),
	})
	if rowsAffected > 0 {
		_ = r.q.DecrementPostLikeCount(ctx, int64(postID))
	}
	return err
}

// HasLiked は指定ユーザーが投稿にいいね済みかを返す。
func (r *postLikeRepository) HasLiked(ctx context.Context, userID, postID uint) (bool, error) {
	count, err := r.q.CountPostLikeByUserAndPost(ctx, sqlcgen.CountPostLikeByUserAndPostParams{
		UserID: int64(userID),
		PostID: int64(postID),
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
