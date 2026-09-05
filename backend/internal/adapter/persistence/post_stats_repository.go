package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postStatsRepository は [repository.PostStatsRepository] の sqlc(pgx) 実装。
type postStatsRepository struct {
	q *sqlcgen.Queries
}

// NewPostStatsRepository は PostStatsRepository の sqlc(pgx) 実装を返す。
func NewPostStatsRepository(q *sqlcgen.Queries) repository.PostStatsRepository {
	return &postStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostStatsRepository = (*postStatsRepository)(nil)

// GetPostStats は指定ユーザーの投稿集計統計を返す。
func (r *postStatsRepository) GetPostStats(ctx context.Context, userID uint) (*model.PostStats, error) {
	uid := int64(userID)
	var stats model.PostStats

	totalPosts, err := r.q.CountPostsByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.TotalPosts = totalPosts

	publishedPosts, err := r.q.CountPublishedPostsByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.PublishedPosts = publishedPosts

	draftPosts, err := r.q.CountDraftPostsByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.DraftPosts = draftPosts

	totalLikes, err := r.q.SumPostLikesReceivedByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.TotalLikesReceived = totalLikes

	totalComments, err := r.q.SumPostCommentsReceivedByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.TotalComments = totalComments

	monthStart := domain.StartOfMonth(time.Now())
	postsThisMonth, err := r.q.CountPostsByUserSince(ctx, sqlcgen.CountPostsByUserSinceParams{
		UserID:    uid,
		CreatedAt: toTimestamptzNotNull(monthStart),
	})
	if err != nil {
		return nil, err
	}
	stats.PostsThisMonth = postsThisMonth

	return &stats, nil
}
