package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ReminderSettingsRepository は学習リマインダー設定の永続化に対する、usecase 側が要求する契約。
type ReminderSettingsRepository interface {
	// GetOrCreateDefault は設定を返す。未登録の場合はデフォルト設定を作成して返す。
	GetOrCreateDefault(ctx context.Context, userID uint) (*model.ReminderSettings, error)
	// FindByUserID は設定を返す。
	// 未登録の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	// usecase 側が永続化技術のエラー型（gorm.ErrRecordNotFound 等）を知らずに済むようにするための契約。
	FindByUserID(ctx context.Context, userID uint) (*model.ReminderSettings, error)
	Save(ctx context.Context, settings *model.ReminderSettings) error
}
