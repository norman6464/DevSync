package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// UserActivityRepository はユーザーアクティビティのデータアクセスを提供する。
type UserActivityRepository struct {
	db *gorm.DB
}

// NewUserActivityRepository は新しいUserActivityRepositoryインスタンスを生成する。
func NewUserActivityRepository(db *gorm.DB) *UserActivityRepository {
	return &UserActivityRepository{db: db}
}

// Create はアクティビティを記録する。
func (r *UserActivityRepository) Create(activity *model.UserActivity) error {
	return r.db.Create(activity).Error
}

// FindByUserID は指定ユーザーのアクティビティを時系列で取得する。
func (r *UserActivityRepository) FindByUserID(userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error) {
	var activities []model.UserActivity
	var total int64

	query := r.db.Where("user_id = ?", userID)
	if activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}

	if err := query.Model(&model.UserActivity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}
