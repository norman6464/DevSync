package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NotificationSettingsRepository は通知設定データへのアクセスを提供するリポジトリ実装。
type NotificationSettingsRepository struct {
	db *gorm.DB
}

// NewNotificationSettingsRepository は新しいNotificationSettingsRepositoryインスタンスを生成する。
func NewNotificationSettingsRepository(db *gorm.DB) *NotificationSettingsRepository {
	return &NotificationSettingsRepository{db: db}
}

// CreateOrUpdate は通知設定を作成または更新する。
// ユーザーIDがユニークキーなので、既存の場合は更新される。
func (r *NotificationSettingsRepository) CreateOrUpdate(settings *model.NotificationSettings) error {
	return r.db.Save(settings).Error
}

// GetByUserID は指定ユーザーの通知設定を取得する。
func (r *NotificationSettingsRepository) GetByUserID(userID uint) (*model.NotificationSettings, error) {
	var settings model.NotificationSettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// GetOrCreateDefault は指定ユーザーの通知設定を取得する。
// 設定が存在しない場合はデフォルト設定を作成して返す。
func (r *NotificationSettingsRepository) GetOrCreateDefault(userID uint) (*model.NotificationSettings, error) {
	var settings model.NotificationSettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error

	if err == gorm.ErrRecordNotFound {
		// デフォルト設定を作成
		settings = model.NotificationSettings{
			UserID:         userID,
			EnableLikes:    true,
			EnableComments: true,
			EnableFollows:  true,
			EnableMessages: true,
			EnableMentions: true,
			EnableWebPush:  true,
			EnableEmail:    true,
			EnableSound:    true,
		}
		if createErr := r.db.Create(&settings).Error; createErr != nil {
			return nil, createErr
		}
		return &settings, nil
	}

	if err != nil {
		return nil, err
	}

	return &settings, nil
}
