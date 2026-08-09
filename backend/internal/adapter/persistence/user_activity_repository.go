package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// userActivityRepository は [repository.UserActivityRepository] の GORM 実装。
type userActivityRepository struct {
	db *gorm.DB
}

// NewUserActivityRepository は UserActivityRepository の GORM 実装を返す。
func NewUserActivityRepository(db *gorm.DB) repository.UserActivityRepository {
	return &userActivityRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.UserActivityRepository = (*userActivityRepository)(nil)

// FindByUserID は指定ユーザーのアクティビティを時系列（新しい順）で取得する。
// activityType が空でなければ種別で絞り込む。
// created_at が同値の行でもページングが安定するよう、id を第 2 ソートキーにして順序を決定的にする。
func (r *userActivityRepository) FindByUserID(ctx context.Context, userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error) {
	var activities []model.UserActivity
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}

	if err := query.Model(&model.UserActivity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Order("id DESC").Limit(limit).Offset(offset).Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}
