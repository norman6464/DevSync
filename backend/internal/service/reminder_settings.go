package service

import (
	"regexp"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// validFrequencies は有効なリマインダー頻度の集合。
var validFrequencies = map[string]bool{
	string(model.ReminderFrequencyDaily):  true,
	string(model.ReminderFrequencyWeekly): true,
}

// timeFormatRegex はHH:MM形式の正規表現。
var timeFormatRegex = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

// ReminderSettingsService は学習リマインダー設定のビジネスロジックを提供する。
type ReminderSettingsService struct {
	repo repository.ReminderSettingsRepositoryInterface
}

// NewReminderSettingsService は新しいReminderSettingsServiceインスタンスを生成する。
func NewReminderSettingsService(repo repository.ReminderSettingsRepositoryInterface) *ReminderSettingsService {
	return &ReminderSettingsService{repo: repo}
}

// GetSettings は指定ユーザーのリマインダー設定を取得する。存在しない場合はデフォルト設定を作成する。
func (s *ReminderSettingsService) GetSettings(userID uint) (*model.ReminderSettings, error) {
	return s.repo.GetOrCreateDefault(userID)
}

// UpdateSettings は指定ユーザーのリマインダー設定を更新する。
func (s *ReminderSettingsService) UpdateSettings(userID uint, updates *model.ReminderSettings) (*model.ReminderSettings, error) {
	// 既存設定を取得
	settings, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 更新内容を反映（バリデーション付き）
	if updates.Frequency != "" {
		if !validFrequencies[string(updates.Frequency)] {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "頻度はdailyまたはweeklyのみ有効です", nil)
		}
		settings.Frequency = updates.Frequency
	}
	if updates.NotificationTime != "" {
		if !timeFormatRegex.MatchString(updates.NotificationTime) {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "通知時間はHH:MM形式で指定してください", nil)
		}
		settings.NotificationTime = updates.NotificationTime
	}
	if updates.InactiveDays > 0 {
		if updates.InactiveDays > 30 {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "非活動日数は1〜30の範囲で指定してください", nil)
		}
		settings.InactiveDays = updates.InactiveDays
	}
	settings.Enabled = updates.Enabled
	settings.EnableWeb = updates.EnableWeb
	settings.EnableEmail = updates.EnableEmail

	// 保存
	if err := s.repo.CreateOrUpdate(settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// ShouldRemind は指定ユーザーにリマインダーを送信すべきかどうかを判定する。
// 最後の学習活動からInactiveDays以上経過している場合にtrueを返す。
func (s *ReminderSettingsService) ShouldRemind(userID uint, lastActivity time.Time) bool {
	settings, err := s.repo.GetByUserID(userID)
	if err != nil || !settings.Enabled {
		return false
	}

	// 最後の学習活動からの経過時間を計算
	daysSinceLastActivity := int(time.Since(lastActivity).Hours() / 24)

	return daysSinceLastActivity >= settings.InactiveDays
}

// SendReminder はリマインダー通知を送信し、LastRemindedAtを更新する。
func (s *ReminderSettingsService) SendReminder(userID uint) error {
	// TODO: 通知システムと統合してリマインダー通知を送信
	// 現時点では LastRemindedAt の更新のみ実装
	return s.repo.UpdateLastRemindedAt(userID)
}
