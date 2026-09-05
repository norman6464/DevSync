package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// userActivityRepository は [repository.UserActivityRepository] の sqlc(pgx) 実装。
type userActivityRepository struct {
	q *sqlcgen.Queries
}

// NewUserActivityRepository は UserActivityRepository の sqlc(pgx) 実装を返す。
func NewUserActivityRepository(q *sqlcgen.Queries) repository.UserActivityRepository {
	return &userActivityRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.UserActivityRepository = (*userActivityRepository)(nil)

// toModelUserActivity は sqlc の生成行を model.UserActivity へ変換する。
func toModelUserActivity(row sqlcgen.UserActivity) model.UserActivity {
	return model.UserActivity{
		ID:           uint(row.ID),
		UserID:       uint(row.UserID),
		ActivityType: model.ActivityType(row.ActivityType),
		TargetType:   row.TargetType,
		TargetID:     uint(row.TargetID),
		Metadata:     fromStringPtr(row.Metadata),
		CreatedAt:    row.CreatedAt.Time,
	}
}

// FindByUserID は指定ユーザーのアクティビティを時系列（新しい順）で取得する。
// activityType が空でなければ種別で絞り込む。
// created_at が同値の行でもページングが安定するよう、id を第 2 ソートキーにして順序を決定的にする。
func (r *userActivityRepository) FindByUserID(ctx context.Context, userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error) {
	var rows []sqlcgen.UserActivity
	var total int64
	var err error

	if activityType == "" {
		total, err = r.q.CountUserActivitiesByUser(ctx, int64(userID))
		if err != nil {
			return nil, 0, err
		}
		rows, err = r.q.ListUserActivitiesByUser(ctx, sqlcgen.ListUserActivitiesByUserParams{
			UserID: int64(userID),
			Limit:  int32Param(limit),
			Offset: int32Param(offset),
		})
	} else {
		total, err = r.q.CountUserActivitiesByUserAndType(ctx, sqlcgen.CountUserActivitiesByUserAndTypeParams{
			UserID:       int64(userID),
			ActivityType: activityType,
		})
		if err != nil {
			return nil, 0, err
		}
		rows, err = r.q.ListUserActivitiesByUserAndType(ctx, sqlcgen.ListUserActivitiesByUserAndTypeParams{
			UserID:       int64(userID),
			ActivityType: activityType,
			Limit:        int32Param(limit),
			Offset:       int32Param(offset),
		})
	}
	if err != nil {
		return nil, 0, err
	}

	activities := make([]model.UserActivity, len(rows))
	for i, row := range rows {
		activities[i] = toModelUserActivity(row)
	}
	return activities, total, nil
}
