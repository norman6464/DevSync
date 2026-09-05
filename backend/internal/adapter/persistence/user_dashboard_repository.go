package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// userDashboardRepository は [repository.UserDashboardRepository] の sqlc(pgx) 実装。
type userDashboardRepository struct {
	q *sqlcgen.Queries
}

// NewUserDashboardRepository は UserDashboardRepository の sqlc(pgx) 実装を返す。
func NewUserDashboardRepository(q *sqlcgen.Queries) repository.UserDashboardRepository {
	return &userDashboardRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.UserDashboardRepository = (*userDashboardRepository)(nil)

// GetDashboardStats は指定ユーザーのダッシュボード統計情報を集計して返す。
// 投稿数・受信いいね数・受信コメント数・受信閲覧数・フォロワー数・フォロー数を返す。
func (r *userDashboardRepository) GetDashboardStats(ctx context.Context, userID uint) (*model.UserDashboardStats, error) {
	uid := int64(userID)
	stats := &model.UserDashboardStats{}

	postCount, err := r.q.CountPublishedPostsByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.PostCount = postCount

	likesReceived, err := r.q.SumPostLikesReceivedByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.LikesReceived = likesReceived

	commentsReceived, err := r.q.SumPostCommentsReceivedByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.CommentsReceived = commentsReceived

	viewsReceived, err := r.q.CountPostViewsReceivedByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.ViewsReceived = viewsReceived

	followerCount, err := r.q.CountFollowersByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.FollowerCount = followerCount

	followingCount, err := r.q.CountFollowingByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	stats.FollowingCount = followingCount

	return stats, nil
}
