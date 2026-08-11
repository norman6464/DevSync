package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// notificationSettingsRepository は [repository.NotificationSettingsRepository] の GORM 実装。
type notificationSettingsRepository struct {
	db *gorm.DB
}

// NewNotificationSettingsRepository は NotificationSettingsRepository の GORM 実装を返す。
func NewNotificationSettingsRepository(db *gorm.DB) repository.NotificationSettingsRepository {
	return &notificationSettingsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NotificationSettingsRepository = (*notificationSettingsRepository)(nil)

// GetOrCreateDefault は設定を取得し、未登録ならすべて有効なデフォルト設定を作成して返す。
func (r *notificationSettingsRepository) GetOrCreateDefault(ctx context.Context, userID uint) (*model.NotificationSettings, error) {
	var settings model.NotificationSettings
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if err == nil {
		return &settings, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

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
	if err := r.db.WithContext(ctx).Create(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

// Save は設定を保存する。ID がセットされていれば更新、無ければ新規作成になる。
func (r *notificationSettingsRepository) Save(ctx context.Context, settings *model.NotificationSettings) error {
	return r.db.WithContext(ctx).Save(settings).Error
}
