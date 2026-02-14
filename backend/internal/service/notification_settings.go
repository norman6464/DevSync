package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
)

// NotificationSettingsRepository は通知設定リポジトリのインターフェース。
type NotificationSettingsRepository interface {
	CreateOrUpdate(settings *model.NotificationSettings) error
	GetByUserID(userID uint) (*model.NotificationSettings, error)
	GetOrCreateDefault(userID uint) (*model.NotificationSettings, error)
}

// NotificationSettingsService は通知設定のビジネスロジックを提供する。
type NotificationSettingsService struct {
	repo NotificationSettingsRepository
}

// NewNotificationSettingsService は新しいNotificationSettingsServiceインスタンスを生成する。
func NewNotificationSettingsService(repo NotificationSettingsRepository) *NotificationSettingsService {
	return &NotificationSettingsService{repo: repo}
}

// GetSettings は指定ユーザーの通知設定を取得する。
// 設定が存在しない場合はデフォルト設定を作成して返す。
func (s *NotificationSettingsService) GetSettings(userID uint) (*model.NotificationSettings, error) {
	return s.repo.GetOrCreateDefault(userID)
}

// UpdateSettings は指定ユーザーの通知設定を更新する。
func (s *NotificationSettingsService) UpdateSettings(userID uint, updates *model.NotificationSettings) (*model.NotificationSettings, error) {
	// 既存設定を取得
	settings, err := s.repo.GetOrCreateDefault(userID)
	if err != nil {
		return nil, err
	}

	// 更新
	settings.EnableLikes = updates.EnableLikes
	settings.EnableComments = updates.EnableComments
	settings.EnableFollows = updates.EnableFollows
	settings.EnableMessages = updates.EnableMessages
	settings.EnableMentions = updates.EnableMentions
	settings.EnableWebPush = updates.EnableWebPush
	settings.EnableEmail = updates.EnableEmail
	settings.EnableSound = updates.EnableSound

	if err := s.repo.CreateOrUpdate(settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// ShouldNotify は指定ユーザーが指定タイプの通知を受け取るかを判定する。
func (s *NotificationSettingsService) ShouldNotify(userID uint, notificationType model.NotificationType) (bool, error) {
	settings, err := s.repo.GetOrCreateDefault(userID)
	if err != nil {
		return false, err
	}

	switch notificationType {
	case model.NotificationTypeLike:
		return settings.EnableLikes, nil
	case model.NotificationTypeComment:
		return settings.EnableComments, nil
	case model.NotificationTypeFollow:
		return settings.EnableFollows, nil
	case model.NotificationTypeMessage:
		return settings.EnableMessages, nil
	default:
		// その他の通知タイプ（投稿・回答・バッジ・レベルアップ）はデフォルトで有効
		return true, nil
	}
}
