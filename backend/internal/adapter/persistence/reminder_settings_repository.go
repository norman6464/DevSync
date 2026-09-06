package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// デフォルト設定の初期値。未登録ユーザーが初めて設定を開いたときに作られる。
const (
	defaultReminderNotificationTime = "09:00"
	defaultReminderInactiveDays     = 3
)

// reminderSettingsRepository は [repository.ReminderSettingsRepository] の sqlc(pgx) 実装。
type reminderSettingsRepository struct {
	q *sqlcgen.Queries
}

// NewReminderSettingsRepository は ReminderSettingsRepository の sqlc(pgx) 実装を返す。
func NewReminderSettingsRepository(q *sqlcgen.Queries) repository.ReminderSettingsRepository {
	return &reminderSettingsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ReminderSettingsRepository = (*reminderSettingsRepository)(nil)

// toModelReminderSettings は sqlc の生成行を model.ReminderSettings へ変換する。
func toModelReminderSettings(row sqlcgen.ReminderSetting) model.ReminderSettings {
	return model.ReminderSettings{
		ID:               uint(row.ID),
		UserID:           uint(row.UserID),
		Enabled:          row.Enabled,
		Frequency:        model.ReminderFrequency(fromStringPtr(row.Frequency)),
		NotificationTime: fromStringPtr(row.NotificationTime),
		InactiveDays:     fromInt64Ptr(row.InactiveDays),
		EnableWeb:        row.EnableWeb,
		EnableEmail:      row.EnableEmail,
		LastRemindedAt:   fromTimestamptz(row.LastRemindedAt),
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}

// GetOrCreateDefault は指定ユーザーの設定を取得し、未登録ならデフォルト設定を作成して返す。
func (r *reminderSettingsRepository) GetOrCreateDefault(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	settings, err := r.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if settings != nil {
		return settings, nil
	}

	enabled := true
	frequency := string(model.ReminderFrequencyDaily)
	notificationTime := defaultReminderNotificationTime
	inactiveDays := int64(defaultReminderInactiveDays)
	enableWeb := true
	enableEmail := false

	// 同一ユーザーの初回取得が同時に走ると複数リクエストが「不在」と判定するため、
	// user_id の一意制約に任せて ON CONFLICT DO NOTHING で挿入し、
	// 競合に負けた側は先に作られた行を読み直して返す（失敗させない）。
	rows, err := r.q.CreateDefaultReminderSettings(ctx, sqlcgen.CreateDefaultReminderSettingsParams{
		UserID:           int64(userID),
		Enabled:          enabled,
		Frequency:        &frequency,
		NotificationTime: &notificationTime,
		InactiveDays:     &inactiveDays,
		EnableWeb:        enableWeb,
		EnableEmail:      enableEmail,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return r.FindByUserID(ctx, userID)
	}
	created := toModelReminderSettings(rows[0])
	return &created, nil
}

// FindByUserID は指定ユーザーの設定を取得する。未登録の場合は (nil, nil) を返す。
func (r *reminderSettingsRepository) FindByUserID(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	row, err := r.q.GetReminderSettingsByUserID(ctx, int64(userID))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	settings := toModelReminderSettings(row)
	return &settings, nil
}

// Save は設定を保存する。呼び出し元は必ず GetOrCreateDefault を経てから渡すため、
// ID は常にセットされている（GORMの Save が担っていた「ID未設定なら新規作成」という
// 分岐は実際には使われていないため、更新のみを実装する）。
func (r *reminderSettingsRepository) Save(ctx context.Context, settings *model.ReminderSettings) error {
	frequency := string(settings.Frequency)
	_, err := r.q.UpdateReminderSettings(ctx, sqlcgen.UpdateReminderSettingsParams{
		ID:               int64(settings.ID),
		Enabled:          settings.Enabled,
		Frequency:        &frequency,
		NotificationTime: &settings.NotificationTime,
		InactiveDays:     toInt64Ptr(settings.InactiveDays),
		EnableWeb:        settings.EnableWeb,
		EnableEmail:      settings.EnableEmail,
		LastRemindedAt:   toTimestamptz(settings.LastRemindedAt),
	})
	return err
}
