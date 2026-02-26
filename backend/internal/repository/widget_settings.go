package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WidgetSettingsRepository はウィジェット設定データへのアクセスを提供するリポジトリ実装。
type WidgetSettingsRepository struct {
	db *gorm.DB
}

// NewWidgetSettingsRepository は新しいWidgetSettingsRepositoryインスタンスを生成する。
func NewWidgetSettingsRepository(db *gorm.DB) *WidgetSettingsRepository {
	return &WidgetSettingsRepository{db: db}
}

// FindByUserID は指定ユーザーのウィジェット設定を取得する。
func (r *WidgetSettingsRepository) FindByUserID(userID uint) (*model.WidgetSettings, error) {
	var settings model.WidgetSettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// Upsert はウィジェット設定を作成または更新する。
func (r *WidgetSettingsRepository) Upsert(settings *model.WidgetSettings) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"settings", "updated_at"}),
	}).Create(settings).Error
}
