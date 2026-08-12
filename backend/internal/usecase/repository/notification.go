package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NotificationReader は通知の参照・既読・削除に対する、usecase 側が要求する契約。
// 通知の作成側（フォロワーへの一括通知など）は利用者のスライスが移行する際に切り出す。
type NotificationReader interface {
	// FindByUserID は指定ユーザーの通知を新しい順にページネーションして返す。
	// notificationType が空でなければ通知種別で絞り込む。
	FindByUserID(ctx context.Context, userID uint, page, limit int, notificationType string) ([]model.Notification, error)
	// CountByUserID は指定ユーザーの通知総数を返す（種別での絞り込みは FindByUserID と同じ条件）。
	CountByUserID(ctx context.Context, userID uint, notificationType string) (int64, error)
	// CountUnread は未読通知の数を返す。
	CountUnread(ctx context.Context, userID uint) (int64, error)
	// MarkAsRead は本人の通知 1 件を既読にする。
	MarkAsRead(ctx context.Context, id, userID uint) error
	// MarkAllAsRead は本人の未読通知をすべて既読にする。
	MarkAllAsRead(ctx context.Context, userID uint) error
	// Delete は本人の通知 1 件を削除する。
	Delete(ctx context.Context, id, userID uint) error
}

// FollowerNotifier はフォロワー全員へ通知を一括作成するための最小の契約。
// 投稿の公開など、1 アクションで複数ユーザーに通知するスライスが利用する。
type FollowerNotifier interface {
	// FindFollowerIDs は指定ユーザーをフォローしているユーザーの ID を返す。
	FindFollowerIDs(ctx context.Context, userID uint) ([]uint, error)
	// CreateBatch は通知をまとめて作成する。空スライスなら何もしない。
	CreateBatch(ctx context.Context, notifications []*model.Notification) error
}

// NotificationCreator は通知を 1 件作成するための最小の契約。
// バッジ獲得など、他スライスからの通知作成はこの契約で受ける。
type NotificationCreator interface {
	Create(ctx context.Context, notification *model.Notification) error
}
