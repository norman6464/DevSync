package usecase

import (
	"context"
	"errors"
	"regexp"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// maxReminderInactiveDays はリマインダーを送るまでの非活動日数の上限。
const maxReminderInactiveDays = 30

// validReminderFrequencies は有効なリマインダー頻度の集合。
var validReminderFrequencies = map[model.ReminderFrequency]bool{
	model.ReminderFrequencyDaily:  true,
	model.ReminderFrequencyWeekly: true,
}

// reminderTimeFormatRegex は通知時間の HH:MM 形式を検証する。
var reminderTimeFormatRegex = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

// GetReminderSettingsUseCase は学習リマインダー設定を取得する。
type GetReminderSettingsUseCase struct {
	settings repository.ReminderSettingsRepository
}

// NewGetReminderSettingsUseCase は GetReminderSettingsUseCase を生成する。
func NewGetReminderSettingsUseCase(settings repository.ReminderSettingsRepository) *GetReminderSettingsUseCase {
	return &GetReminderSettingsUseCase{settings: settings}
}

// Execute は設定を返す。未登録の場合はデフォルト設定が作成される。
func (uc *GetReminderSettingsUseCase) Execute(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	if err := domain.ValidateRequiredID(userID, "ユーザーID"); err != nil {
		return nil, err
	}
	return uc.settings.GetOrCreateDefault(ctx, userID)
}

// UpdateReminderSettingsInput は設定更新の入力。
// 真偽値の 3 項目は常に上書きされ、文字列・数値は「空文字 / 0 なら据え置き」として扱う。
type UpdateReminderSettingsInput struct {
	UserID           uint
	Enabled          bool
	Frequency        model.ReminderFrequency
	NotificationTime string
	InactiveDays     int
	EnableWeb        bool
	EnableEmail      bool
}

// UpdateReminderSettingsUseCase は学習リマインダー設定を更新する。
type UpdateReminderSettingsUseCase struct {
	settings repository.ReminderSettingsRepository
}

// NewUpdateReminderSettingsUseCase は UpdateReminderSettingsUseCase を生成する。
func NewUpdateReminderSettingsUseCase(settings repository.ReminderSettingsRepository) *UpdateReminderSettingsUseCase {
	return &UpdateReminderSettingsUseCase{settings: settings}
}

// Execute は既存設定に入力を反映して保存し、保存後の設定を返す。
func (uc *UpdateReminderSettingsUseCase) Execute(ctx context.Context, in UpdateReminderSettingsInput) (*model.ReminderSettings, error) {
	if err := domain.ValidateRequiredID(in.UserID, "ユーザーID"); err != nil {
		return nil, err
	}

	settings, err := uc.settings.FindByUserID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		// 未登録のまま更新された場合は 500 を返す（移行前の挙動を維持している）。
		return nil, errors.New("リマインダー設定が見つかりません")
	}

	if in.Frequency != "" {
		if !validReminderFrequencies[in.Frequency] {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "頻度はdailyまたはweeklyのみ有効です", nil)
		}
		settings.Frequency = in.Frequency
	}
	if in.NotificationTime != "" {
		if !reminderTimeFormatRegex.MatchString(in.NotificationTime) {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "通知時間はHH:MM形式で指定してください", nil)
		}
		settings.NotificationTime = in.NotificationTime
	}
	if in.InactiveDays > 0 {
		if in.InactiveDays > maxReminderInactiveDays {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "非活動日数は1〜30の範囲で指定してください", nil)
		}
		settings.InactiveDays = in.InactiveDays
	}
	settings.Enabled = in.Enabled
	settings.EnableWeb = in.EnableWeb
	settings.EnableEmail = in.EnableEmail

	if err := uc.settings.Save(ctx, settings); err != nil {
		return nil, err
	}
	return settings, nil
}
