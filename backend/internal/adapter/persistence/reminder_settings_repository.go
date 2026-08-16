package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// デフォルト設定の初期値。未登録ユーザーが初めて設定を開いたときに作られる。
const (
	defaultReminderNotificationTime = "09:00"
	defaultReminderInactiveDays     = 3
)

// reminderSettingsRepository は [repository.ReminderSettingsRepository] の GORM 実装。
type reminderSettingsRepository struct {
	db *gorm.DB
}

// NewReminderSettingsRepository は ReminderSettingsRepository の GORM 実装を返す。
func NewReminderSettingsRepository(db *gorm.DB) repository.ReminderSettingsRepository {
	return &reminderSettingsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ReminderSettingsRepository = (*reminderSettingsRepository)(nil)

// GetOrCreateDefault は指定ユーザーの設定を取得し、未登録ならデフォルト設定を作成して返す。
func (r *reminderSettingsRepository) GetOrCreateDefault(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	settings, err := r.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if settings != nil {
		return settings, nil
	}

	defaultSettings := &model.ReminderSettings{
		UserID:           userID,
		Enabled:          true,
		Frequency:        model.ReminderFrequencyDaily,
		NotificationTime: defaultReminderNotificationTime,
		InactiveDays:     defaultReminderInactiveDays,
		EnableWeb:        true,
		EnableEmail:      false,
	}
	// 同一ユーザーの初回取得が同時に走ると複数リクエストが「不在」と判定するため、
	// user_id の一意制約に任せて ON CONFLICT DO NOTHING で挿入し、
	// 競合に負けた側は先に作られた行を読み直して返す（失敗させない）。
	res := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoNothing: true,
	}).Create(defaultSettings)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return r.FindByUserID(ctx, userID)
	}
	return defaultSettings, nil
}

// FindByUserID は指定ユーザーの設定を取得する。未登録の場合は (nil, nil) を返す。
func (r *reminderSettingsRepository) FindByUserID(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	var settings model.ReminderSettings
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// Save は設定を保存する。ID がセットされていれば更新、無ければ新規作成になる。
func (r *reminderSettingsRepository) Save(ctx context.Context, settings *model.ReminderSettings) error {
	return r.db.WithContext(ctx).Save(settings).Error
}
