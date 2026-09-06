package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// notificationVerbRepository は [repository.NotificationVerbRepository] の sqlc(pgx) 実装。
type notificationVerbRepository struct {
	q *sqlcgen.Queries
}

// NewNotificationVerbRepository は NotificationVerbRepository の sqlc(pgx) 実装を返す。
func NewNotificationVerbRepository(q *sqlcgen.Queries) repository.NotificationVerbRepository {
	return &notificationVerbRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NotificationVerbRepository = (*notificationVerbRepository)(nil)

// SeedKnownVerbs は既知の通知種別コードをまとめて登録する（既存分はスキップ）。
func (r *notificationVerbRepository) SeedKnownVerbs(ctx context.Context, codes []string) error {
	return r.q.SeedNotificationVerbs(ctx, codes)
}
