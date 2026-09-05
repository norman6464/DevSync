package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// notificationSettingsRepository は [repository.NotificationSettingsRepository] の sqlc(pgx) 実装。
type notificationSettingsRepository struct {
	q *sqlcgen.Queries
}

// NewNotificationSettingsRepository は NotificationSettingsRepository の sqlc(pgx) 実装を返す。
func NewNotificationSettingsRepository(q *sqlcgen.Queries) repository.NotificationSettingsRepository {
	return &notificationSettingsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NotificationSettingsRepository = (*notificationSettingsRepository)(nil)

// toModelNotificationSettings は sqlc の生成行を model.NotificationSettings へ変換する。
func toModelNotificationSettings(row sqlcgen.NotificationSetting) model.NotificationSettings {
	return model.NotificationSettings{
		ID:             uint(row.ID),
		UserID:         uint(row.UserID),
		EnableLikes:    fromBoolPtr(row.EnableLikes),
		EnableComments: fromBoolPtr(row.EnableComments),
		EnableFollows:  fromBoolPtr(row.EnableFollows),
		EnableMessages: fromBoolPtr(row.EnableMessages),
		EnableMentions: fromBoolPtr(row.EnableMentions),
		EnableWebPush:  fromBoolPtr(row.EnableWebPush),
		EnableEmail:    fromBoolPtr(row.EnableEmail),
		EnableSound:    fromBoolPtr(row.EnableSound),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
}

// GetOrCreateDefault は設定を取得し、未登録ならすべて有効なデフォルト設定を作成して返す。
func (r *notificationSettingsRepository) GetOrCreateDefault(ctx context.Context, userID uint) (*model.NotificationSettings, error) {
	row, err := r.q.GetNotificationSettingsByUserID(ctx, int64(userID))
	if err == nil {
		settings := toModelNotificationSettings(row)
		return &settings, nil
	}
	if !isNoRows(err) {
		return nil, err
	}

	allEnabled := true
	created, err := r.q.CreateNotificationSettings(ctx, sqlcgen.CreateNotificationSettingsParams{
		UserID:         int64(userID),
		EnableLikes:    &allEnabled,
		EnableComments: &allEnabled,
		EnableFollows:  &allEnabled,
		EnableMessages: &allEnabled,
		EnableMentions: &allEnabled,
		EnableWebPush:  &allEnabled,
		EnableEmail:    &allEnabled,
		EnableSound:    &allEnabled,
	})
	if err != nil {
		return nil, err
	}
	settings := toModelNotificationSettings(created)
	return &settings, nil
}

// Save は設定を保存する。呼び出し元は必ず GetOrCreateDefault を経てから渡すため、
// ID は常にセットされている（GORMの Save が担っていた「ID未設定なら新規作成」という
// 分岐は実際には使われていないため、更新のみを実装する）。
func (r *notificationSettingsRepository) Save(ctx context.Context, settings *model.NotificationSettings) error {
	_, err := r.q.UpdateNotificationSettings(ctx, sqlcgen.UpdateNotificationSettingsParams{
		ID:             int64(settings.ID),
		EnableLikes:    &settings.EnableLikes,
		EnableComments: &settings.EnableComments,
		EnableFollows:  &settings.EnableFollows,
		EnableMessages: &settings.EnableMessages,
		EnableMentions: &settings.EnableMentions,
		EnableWebPush:  &settings.EnableWebPush,
		EnableEmail:    &settings.EnableEmail,
		EnableSound:    &settings.EnableSound,
	})
	return err
}
