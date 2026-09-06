package repository

import "context"

// NotificationVerbRepository は notification_verbs（通知種別のlookupテーブル、
// DEVSYNC-159でnotifications.typeのFK参照先とした）に対する、usecase側が要求する契約。
type NotificationVerbRepository interface {
	// SeedKnownVerbs は既知の通知種別コードをまとめて登録する（既存分はスキップ）。
	SeedKnownVerbs(ctx context.Context, codes []string) error
}
