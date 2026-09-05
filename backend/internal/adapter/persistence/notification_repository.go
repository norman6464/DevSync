package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// notificationRepository は [repository.NotificationReader] / [repository.NotificationCreator] /
// [repository.FollowerNotifier] の sqlc(pgx) 実装。
// CreateBatch は複数件の作成を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type notificationRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewNotificationRepository は NotificationReader の sqlc(pgx) 実装を返す。
func NewNotificationRepository(pool *pgxpool.Pool) repository.NotificationReader {
	return &notificationRepository{pool: pool, q: sqlcgen.New(pool)}
}

// NewNotificationCreator は NotificationCreator の sqlc(pgx) 実装を返す。
// 通知の作成だけを必要とする利用者（バッジ獲得など）はこちらを受け取る。
func NewNotificationCreator(pool *pgxpool.Pool) repository.NotificationCreator {
	return &notificationRepository{pool: pool, q: sqlcgen.New(pool)}
}

// NewFollowerNotifier は FollowerNotifier の sqlc(pgx) 実装を返す。
// フォロワー全員への一括通知を必要とする利用者（投稿の公開など）はこちらを受け取る。
func NewFollowerNotifier(pool *pgxpool.Pool) repository.FollowerNotifier {
	return &notificationRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var (
	_ repository.NotificationReader  = (*notificationRepository)(nil)
	_ repository.NotificationCreator = (*notificationRepository)(nil)
	_ repository.FollowerNotifier    = (*notificationRepository)(nil)
)

// nilIfEmpty は空文字列を nil に変換する。notification_type のような
// 「空文字列なら絞り込みなし」の任意フィルタを sqlc.narg に渡す際に使う。
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toModelNotification(row sqlcgen.Notification) model.Notification {
	return model.Notification{
		ID:         uint(row.ID),
		UserID:     uint(row.UserID),
		Type:       model.NotificationType(row.Type),
		ActorID:    uint(row.ActorID),
		PostID:     fromInt64PtrToUintPtr(row.PostID),
		QuestionID: fromInt64PtrToUintPtr(row.QuestionID),
		BadgeID:    row.BadgeID,
		Read:       fromBoolPtr(row.Read),
		CreatedAt:  timeValue(fromTimestamptz(row.CreatedAt)),
	}
}

// Create は通知を 1 件保存する。
func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	row, err := r.q.CreateNotification(ctx, sqlcgen.CreateNotificationParams{
		UserID:     int64(notification.UserID),
		Type:       string(notification.Type),
		ActorID:    int64(notification.ActorID),
		PostID:     toInt64PtrFromUintPtr(notification.PostID),
		QuestionID: toInt64PtrFromUintPtr(notification.QuestionID),
		BadgeID:    notification.BadgeID,
		Read:       &notification.Read,
	})
	if err != nil {
		return err
	}
	*notification = toModelNotification(row)
	return nil
}

// CreateBatch は通知をまとめて保存する。GORMのCreate(&slice)（単一トランザクション相当）に合わせ、
// 1トランザクション内でループ挿入することでアトミック性を維持する。
func (r *notificationRepository) CreateBatch(ctx context.Context, notifications []*model.Notification) error {
	if len(notifications) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for _, notification := range notifications {
		row, err := q.CreateNotification(ctx, sqlcgen.CreateNotificationParams{
			UserID:     int64(notification.UserID),
			Type:       string(notification.Type),
			ActorID:    int64(notification.ActorID),
			PostID:     toInt64PtrFromUintPtr(notification.PostID),
			QuestionID: toInt64PtrFromUintPtr(notification.QuestionID),
			BadgeID:    notification.BadgeID,
			Read:       &notification.Read,
		})
		if err != nil {
			return err
		}
		*notification = toModelNotification(row)
	}
	return tx.Commit(ctx)
}

// FindFollowerIDs は指定ユーザーをフォローしているユーザーの ID を返す。
func (r *notificationRepository) FindFollowerIDs(ctx context.Context, userID uint) ([]uint, error) {
	ids, err := r.q.FindFollowerIDsByFollowee(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	followerIDs := make([]uint, len(ids))
	for i, id := range ids {
		followerIDs[i] = uint(id)
	}
	return followerIDs, nil
}

// attachNotificationPostAndQuestion は LEFT JOIN で取得した posts / questions の個別カラムを
// model.Notification へ Post / Question として付与する。
func attachNotificationPostAndQuestion(notification *model.Notification, row sqlcgen.FindNotificationsByUserIDRow) {
	if row.PostID2 != nil {
		notification.Post = &model.Post{
			ID:                uint(*row.PostID2),
			UserID:            uint(fromInt64PtrValue(row.PostUserID)),
			Title:             fromStringPtr(row.PostTitle),
			Content:           fromStringPtr(row.PostContent),
			ImageURLs:         fromStringPtr(row.PostImageUrls),
			IsDraft:           fromBoolPtr(row.PostIsDraft),
			LikeCount:         int(fromInt64PtrValue(row.PostLikeCount)),
			CommentCount:      int(fromInt64PtrValue(row.PostCommentCount)),
			BookmarkCount:     int(fromInt64PtrValue(row.PostBookmarkCount)),
			ViewCount:         int(fromInt64PtrValue(row.PostViewCount)),
			EstimatedReadTime: int(fromInt64PtrValue(row.PostEstimatedReadTime)),
			ScheduledAt:       fromTimestamptz(row.PostScheduledAt),
			CreatedAt:         timeValue(fromTimestamptz(row.PostCreatedAt)),
			UpdatedAt:         timeValue(fromTimestamptz(row.PostUpdatedAt)),
		}
	}
	if row.QuestionID2 != nil {
		notification.Question = &model.Question{
			ID:          uint(*row.QuestionID2),
			UserID:      uint(fromInt64PtrValue(row.QuestionUserID)),
			Title:       fromStringPtr(row.QuestionTitle),
			Body:        fromStringPtr(row.QuestionBody),
			Tags:        fromStringPtr(row.QuestionTags),
			VoteCount:   int(fromInt64PtrValue(row.QuestionVoteCount)),
			AnswerCount: int(fromInt64PtrValue(row.QuestionAnswerCount)),
			IsSolved:    fromBoolPtr(row.QuestionIsSolved),
			CreatedAt:   timeValue(fromTimestamptz(row.QuestionCreatedAt)),
			UpdatedAt:   timeValue(fromTimestamptz(row.QuestionUpdatedAt)),
		}
	}
}

// FindByUserID は指定ユーザーの通知を作成日時の降順で取得する。
// 通知の表示に必要な関連（実行者・投稿・質問）を Preload する。
func (r *notificationRepository) FindByUserID(ctx context.Context, userID uint, page, limit int, notificationType string) ([]model.Notification, error) {
	offset := (page - 1) * limit

	rows, err := r.q.FindNotificationsByUserID(ctx, sqlcgen.FindNotificationsByUserIDParams{
		UserID:           int64(userID),
		Limit:            int32Param(limit),
		Offset:           int32Param(offset),
		NotificationType: nilIfEmpty(notificationType),
	})
	if err != nil {
		return nil, err
	}

	notifications := make([]model.Notification, len(rows))
	for i, row := range rows {
		notifications[i] = toModelNotification(row.Notification)
		notifications[i].Actor = toModelUser(row.User)
		attachNotificationPostAndQuestion(&notifications[i], row)
	}
	return notifications, nil
}

// CountByUserID は指定ユーザーの通知総数を取得する。
func (r *notificationRepository) CountByUserID(ctx context.Context, userID uint, notificationType string) (int64, error) {
	return r.q.CountNotificationsByUserID(ctx, sqlcgen.CountNotificationsByUserIDParams{
		UserID:           int64(userID),
		NotificationType: nilIfEmpty(notificationType),
	})
}

// CountUnread は指定ユーザーの未読通知数を取得する。
func (r *notificationRepository) CountUnread(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountUnreadNotifications(ctx, int64(userID))
}

// MarkAsRead は指定 ID の通知を既読にする。本人の通知だけを対象にする。
func (r *notificationRepository) MarkAsRead(ctx context.Context, id, userID uint) error {
	return r.q.MarkNotificationAsRead(ctx, sqlcgen.MarkNotificationAsReadParams{
		ID:     int64(id),
		UserID: int64(userID),
	})
}

// MarkAllAsRead は指定ユーザーの未読通知をすべて既読にする。
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.q.MarkAllNotificationsAsRead(ctx, int64(userID))
}

// Delete は指定 ID の通知を削除する。本人の通知だけを対象にする。
func (r *notificationRepository) Delete(ctx context.Context, id, userID uint) error {
	return r.q.DeleteNotification(ctx, sqlcgen.DeleteNotificationParams{
		ID:     int64(id),
		UserID: int64(userID),
	})
}
