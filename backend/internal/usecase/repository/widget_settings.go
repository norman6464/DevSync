package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// WidgetSettingsRepository はダッシュボードウィジェット設定の永続化に対する、usecase 側が要求する契約。
type WidgetSettingsRepository interface {
	FindByUserID(ctx context.Context, userID uint) (*model.WidgetSettings, error)
	Upsert(ctx context.Context, settings *model.WidgetSettings) error
}
