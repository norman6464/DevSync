package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NotificationSettingsRepository は通知設定の永続化に対する、usecase 側が要求する契約。
type NotificationSettingsRepository interface {
	// GetOrCreateDefault は設定を返す。未登録の場合はデフォルト設定を作成して返す。
	GetOrCreateDefault(ctx context.Context, userID uint) (*model.NotificationSettings, error)
	Save(ctx context.Context, settings *model.NotificationSettings) error
}
