package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// ReminderSettingsRepository は学習リマインダー設定データへのアクセスを提供する。
type ReminderSettingsRepository struct {
	db *gorm.DB
}

// NewReminderSettingsRepository は新しいReminderSettingsRepositoryインスタンスを生成する。
func NewReminderSettingsRepository(db *gorm.DB) *ReminderSettingsRepository {
	return &ReminderSettingsRepository{db: db}
}

// CreateOrUpdate はリマインダー設定を作成または更新する。
// IDがセットされている場合は更新、されていない場合は新規作成する。
func (r *ReminderSettingsRepository) CreateOrUpdate(settings *model.ReminderSettings) error {
	return r.db.Save(settings).Error
}

// GetByUserID は指定ユーザーのリマインダー設定を取得する。
func (r *ReminderSettingsRepository) GetByUserID(userID uint) (*model.ReminderSettings, error) {
	var settings model.ReminderSettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// GetOrCreateDefault は指定ユーザーのリマインダー設定を取得し、存在しない場合はデフォルト設定を作成する。
func (r *ReminderSettingsRepository) GetOrCreateDefault(userID uint) (*model.ReminderSettings, error) {
	settings, err := r.GetByUserID(userID)
	if err == nil {
		return settings, nil
	}

	// デフォルト設定を作成
	defaultSettings := &model.ReminderSettings{
		UserID:           userID,
		Enabled:          true,
		Frequency:        model.ReminderFrequencyDaily,
		NotificationTime: "09:00",
		InactiveDays:     3,
		EnableWeb:        true,
		EnableEmail:      false,
	}

	if err := r.db.Create(defaultSettings).Error; err != nil {
		return nil, err
	}

	return defaultSettings, nil
}

// GetEnabledSettings は有効なリマインダー設定を全て取得する（Cron job用）。
func (r *ReminderSettingsRepository) GetEnabledSettings() ([]model.ReminderSettings, error) {
	var settings []model.ReminderSettings
	err := r.db.Where("enabled = ?", true).Find(&settings).Error
	return settings, err
}

// UpdateLastRemindedAt は最後にリマインドした日時を更新する。
func (r *ReminderSettingsRepository) UpdateLastRemindedAt(userID uint) error {
	return r.db.Model(&model.ReminderSettings{}).
		Where("user_id = ?", userID).
		Update("last_reminded_at", gorm.Expr("NOW()")).Error
}
