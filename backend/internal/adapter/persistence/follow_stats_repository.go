package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// followStatsRepository は [repository.FollowStatsRepository] の sqlc(pgx) 実装。
type followStatsRepository struct {
	q *sqlcgen.Queries
}

// NewFollowStatsRepository は FollowStatsRepository の sqlc(pgx) 実装を返す。
func NewFollowStatsRepository(q *sqlcgen.Queries) repository.FollowStatsRepository {
	return &followStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.FollowStatsRepository = (*followStatsRepository)(nil)

// GetFollowStats は指定ユーザーのフォロー関係集計統計を返す。
func (r *followStatsRepository) GetFollowStats(ctx context.Context, userID uint) (*model.FollowStats, error) {
	followers, err := r.q.CountFollowersByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	following, err := r.q.CountFollowingByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return &model.FollowStats{
		FollowerCount:  followers,
		FollowingCount: following,
	}, nil
}
