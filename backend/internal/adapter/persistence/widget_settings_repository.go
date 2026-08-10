package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// widgetSettingsRepository は [repository.WidgetSettingsRepository] の GORM 実装。
type widgetSettingsRepository struct {
	db *gorm.DB
}

// NewWidgetSettingsRepository は WidgetSettingsRepository の GORM 実装を返す。
func NewWidgetSettingsRepository(db *gorm.DB) repository.WidgetSettingsRepository {
	return &widgetSettingsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.WidgetSettingsRepository = (*widgetSettingsRepository)(nil)

// FindByUserID は指定ユーザーのウィジェット設定を取得する。
func (r *widgetSettingsRepository) FindByUserID(ctx context.Context, userID uint) (*model.WidgetSettings, error) {
	var settings model.WidgetSettings
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

// Upsert はウィジェット設定を作成または更新する。
func (r *widgetSettingsRepository) Upsert(ctx context.Context, settings *model.WidgetSettings) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"settings", "updated_at"}),
	}).Create(settings).Error
}
