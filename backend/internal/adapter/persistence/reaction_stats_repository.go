package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// reactionStatsRepository は [repository.ReactionStatsRepository] の sqlc(pgx) 実装。
type reactionStatsRepository struct {
	q *sqlcgen.Queries
}

// NewReactionStatsRepository は ReactionStatsRepository の sqlc(pgx) 実装を返す。
func NewReactionStatsRepository(q *sqlcgen.Queries) repository.ReactionStatsRepository {
	return &reactionStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ReactionStatsRepository = (*reactionStatsRepository)(nil)

// GetReactionStats は指定ユーザーが受け取ったリアクション集計統計を返す。
func (r *reactionStatsRepository) GetReactionStats(ctx context.Context, userID uint) (*model.ReactionStats, error) {
	received, err := r.q.CountReactionsReceivedByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	uniqueReactors, err := r.q.CountUniqueReactorsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	startOfMonth := domain.StartOfMonth(time.Now())
	thisMonth, err := r.q.CountReactionsReceivedByUserSince(ctx, sqlcgen.CountReactionsReceivedByUserSinceParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(startOfMonth),
	})
	if err != nil {
		return nil, err
	}

	return &model.ReactionStats{
		TotalReactionsReceived: received,
		UniqueReactors:         uniqueReactors,
		ReactionsThisMonth:     thisMonth,
	}, nil
}

// GetEmojiBreakdown は指定ユーザーの全投稿に対する絵文字別リアクション集計を返す。
func (r *reactionStatsRepository) GetEmojiBreakdown(ctx context.Context, userID uint) ([]model.ReactionCount, error) {
	rows, err := r.q.GetEmojiBreakdownByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	counts := make([]model.ReactionCount, len(rows))
	for i, row := range rows {
		counts[i] = model.ReactionCount{
			Emoji: row.Emoji,
			Count: int(row.Count),
		}
	}
	return counts, nil
}

// GetTopReactedPosts は指定ユーザーの投稿のうちリアクション数が多い順に limit 件返す。
func (r *reactionStatsRepository) GetTopReactedPosts(ctx context.Context, userID uint, limit int) ([]model.TopReactedPost, error) {
	rows, err := r.q.GetTopReactedPostsByUser(ctx, sqlcgen.GetTopReactedPostsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
	})
	if err != nil {
		return nil, err
	}

	posts := make([]model.TopReactedPost, len(rows))
	for i, row := range rows {
		posts[i] = model.TopReactedPost{
			ID:            uint(row.ID),
			Title:         row.Title,
			ReactionCount: int(row.ReactionCount),
		}
	}
	return posts, nil
}
