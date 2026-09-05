package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// widgetSettingsRepository は [repository.WidgetSettingsRepository] の sqlc(pgx) 実装。
type widgetSettingsRepository struct {
	q *sqlcgen.Queries
}

// NewWidgetSettingsRepository は WidgetSettingsRepository の sqlc(pgx) 実装を返す。
func NewWidgetSettingsRepository(q *sqlcgen.Queries) repository.WidgetSettingsRepository {
	return &widgetSettingsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.WidgetSettingsRepository = (*widgetSettingsRepository)(nil)

// FindByUserID は指定ユーザーのウィジェット設定を取得する。
// 呼び出し元の usecase は「不在も含めあらゆるエラー」でデフォルト設定にフォールバックする
// 契約のため、他ドメインと異なり (nil, nil) には正規化しない（元のGORM実装と同じ挙動）。
func (r *widgetSettingsRepository) FindByUserID(ctx context.Context, userID uint) (*model.WidgetSettings, error) {
	row, err := r.q.GetWidgetSettingsByUserID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return &model.WidgetSettings{
		ID:        uint(row.ID),
		UserID:    uint(row.UserID),
		Settings:  row.Settings,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

// Upsert はウィジェット設定を作成または更新する。
func (r *widgetSettingsRepository) Upsert(ctx context.Context, settings *model.WidgetSettings) error {
	return r.q.UpsertWidgetSettings(ctx, sqlcgen.UpsertWidgetSettingsParams{
		UserID:   int64(settings.UserID),
		Settings: settings.Settings,
	})
}
