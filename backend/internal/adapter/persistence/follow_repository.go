// Package persistence は usecase/repository の port を sqlc(pgx) で実装する adapter 層。
package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// followRepository は [repository.FollowRepository] の sqlc(pgx) 実装。
type followRepository struct {
	q *sqlcgen.Queries
}

// NewFollowRepository は FollowRepository の sqlc(pgx) 実装を返す。
func NewFollowRepository(q *sqlcgen.Queries) repository.FollowRepository {
	return &followRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.FollowRepository = (*followRepository)(nil)

// Follow はフォロー関係を保存する。複合ユニーク索引（follower_id, followee_id）に
// 衝突した場合は domain.ErrConflict を返す。重複フォローの最終防衛は DB の制約に委ね、
// 生の制約違反を 500 として漏らさない。
func (r *followRepository) Follow(ctx context.Context, followerID, followeeID uint) error {
	err := r.q.CreateFollow(ctx, sqlcgen.CreateFollowParams{
		FollowerID: int64(followerID),
		FolloweeID: int64(followeeID),
	})
	if isUniqueViolation(err) {
		return domain.ErrConflict
	}
	return err
}

func (r *followRepository) Unfollow(ctx context.Context, followerID, followeeID uint) error {
	return r.q.DeleteFollow(ctx, sqlcgen.DeleteFollowParams{
		FollowerID: int64(followerID),
		FolloweeID: int64(followeeID),
	})
}

func (r *followRepository) GetFollowers(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	total, err := r.q.CountFollowersByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListFollowers(ctx, sqlcgen.ListFollowersParams{
		FolloweeID: int64(userID),
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	users := make([]model.User, len(rows))
	for i, row := range rows {
		users[i] = toModelUser(row)
	}
	return users, total, nil
}

func (r *followRepository) GetFollowing(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	total, err := r.q.CountFollowingByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListFollowing(ctx, sqlcgen.ListFollowingParams{
		FollowerID: int64(userID),
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	users := make([]model.User, len(rows))
	for i, row := range rows {
		users[i] = toModelUser(row)
	}
	return users, total, nil
}
