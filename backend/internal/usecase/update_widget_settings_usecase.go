package usecase

import (
	"context"
	"encoding/json"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// UpdateWidgetSettingsUseCase はダッシュボードウィジェット設定を更新する。
type UpdateWidgetSettingsUseCase struct {
	settings repository.WidgetSettingsRepository
}

// NewUpdateWidgetSettingsUseCase は UpdateWidgetSettingsUseCase を生成する。
func NewUpdateWidgetSettingsUseCase(settings repository.WidgetSettingsRepository) *UpdateWidgetSettingsUseCase {
	return &UpdateWidgetSettingsUseCase{settings: settings}
}

// Execute は設定が JSON 配列であることを検証したうえで作成または更新する。
func (uc *UpdateWidgetSettingsUseCase) Execute(ctx context.Context, userID uint, settings string) error {
	if err := domain.ValidateStringLength(settings, 1, 10000, "設定"); err != nil {
		return err
	}

	if !json.Valid([]byte(settings)) {
		return domain.NewError(domain.ErrCodeBadRequest, "設定は有効なJSONである必要があります", nil)
	}

	// JSON配列であることを検証
	var arr []json.RawMessage
	if err := json.Unmarshal([]byte(settings), &arr); err != nil {
		return domain.NewError(domain.ErrCodeBadRequest, "設定はJSON配列である必要があります", nil)
	}

	return uc.settings.Upsert(ctx, &model.WidgetSettings{
		UserID:   userID,
		Settings: settings,
	})
}
