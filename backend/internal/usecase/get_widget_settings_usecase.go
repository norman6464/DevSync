package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// defaultWidgetSettings はダッシュボードウィジェットのデフォルト設定。
const defaultWidgetSettings = `[{"key":"userProfile","visible":true,"order":0},{"key":"level","visible":true,"order":1},{"key":"streak","visible":true,"order":2},{"key":"dailyChallenge","visible":true,"order":3},{"key":"weeklyChallenge","visible":true,"order":4},{"key":"studyCircle","visible":true,"order":5},{"key":"quickEntry","visible":true,"order":6},{"key":"quickActions","visible":true,"order":7},{"key":"recommendedUsers","visible":true,"order":8},{"key":"trending","visible":true,"order":9},{"key":"aiAdvice","visible":true,"order":10},{"key":"goalsProgress","visible":true,"order":11},{"key":"recentNotifications","visible":true,"order":12},{"key":"quickStats","visible":true,"order":13}]`

// GetWidgetSettingsUseCase は指定ユーザーのダッシュボードウィジェット設定を取得する。
type GetWidgetSettingsUseCase struct {
	settings repository.WidgetSettingsRepository
}

// NewGetWidgetSettingsUseCase は GetWidgetSettingsUseCase を生成する。
func NewGetWidgetSettingsUseCase(settings repository.WidgetSettingsRepository) *GetWidgetSettingsUseCase {
	return &GetWidgetSettingsUseCase{settings: settings}
}

// Execute はウィジェット設定を返す。未登録の場合はデフォルト設定を返す。
func (uc *GetWidgetSettingsUseCase) Execute(ctx context.Context, userID uint) (*model.WidgetSettings, error) {
	settings, err := uc.settings.FindByUserID(ctx, userID)
	if err != nil {
		// レコード未登録の場合はデフォルト設定を返す
		return &model.WidgetSettings{
			UserID:   userID,
			Settings: defaultWidgetSettings,
		}, nil
	}
	return settings, nil
}
