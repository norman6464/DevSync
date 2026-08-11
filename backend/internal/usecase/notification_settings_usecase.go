package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetNotificationSettingsUseCase は通知設定を取得する。
type GetNotificationSettingsUseCase struct {
	settings repository.NotificationSettingsRepository
}

// NewGetNotificationSettingsUseCase は GetNotificationSettingsUseCase を生成する。
func NewGetNotificationSettingsUseCase(settings repository.NotificationSettingsRepository) *GetNotificationSettingsUseCase {
	return &GetNotificationSettingsUseCase{settings: settings}
}

// Execute は設定を返す。未登録の場合はデフォルト設定が作成される。
func (uc *GetNotificationSettingsUseCase) Execute(ctx context.Context, userID uint) (*model.NotificationSettings, error) {
	if err := domain.ValidateRequiredID(userID, "ユーザーID"); err != nil {
		return nil, err
	}
	return uc.settings.GetOrCreateDefault(ctx, userID)
}

// UpdateNotificationSettingsInput は通知設定更新の入力。
// 8 つの真偽値はすべてリクエストの値で上書きされる（部分更新ではない）。
type UpdateNotificationSettingsInput struct {
	UserID         uint
	EnableLikes    bool
	EnableComments bool
	EnableFollows  bool
	EnableMessages bool
	EnableMentions bool
	EnableWebPush  bool
	EnableEmail    bool
	EnableSound    bool
}

// UpdateNotificationSettingsUseCase は通知設定を更新する。
type UpdateNotificationSettingsUseCase struct {
	settings repository.NotificationSettingsRepository
}

// NewUpdateNotificationSettingsUseCase は UpdateNotificationSettingsUseCase を生成する。
func NewUpdateNotificationSettingsUseCase(settings repository.NotificationSettingsRepository) *UpdateNotificationSettingsUseCase {
	return &UpdateNotificationSettingsUseCase{settings: settings}
}

// Execute は既存設定（無ければデフォルト）に入力を反映して保存し、保存後の設定を返す。
func (uc *UpdateNotificationSettingsUseCase) Execute(ctx context.Context, in UpdateNotificationSettingsInput) (*model.NotificationSettings, error) {
	if err := domain.ValidateRequiredID(in.UserID, "ユーザーID"); err != nil {
		return nil, err
	}

	settings, err := uc.settings.GetOrCreateDefault(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	settings.EnableLikes = in.EnableLikes
	settings.EnableComments = in.EnableComments
	settings.EnableFollows = in.EnableFollows
	settings.EnableMessages = in.EnableMessages
	settings.EnableMentions = in.EnableMentions
	settings.EnableWebPush = in.EnableWebPush
	settings.EnableEmail = in.EnableEmail
	settings.EnableSound = in.EnableSound

	if err := uc.settings.Save(ctx, settings); err != nil {
		return nil, err
	}
	return settings, nil
}
