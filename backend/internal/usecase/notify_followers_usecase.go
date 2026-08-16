package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// NotifyFollowersUseCase は指定ユーザーのフォロワー全員へ通知を一括作成する。
type NotifyFollowersUseCase struct {
	notifier repository.FollowerNotifier
}

// NewNotifyFollowersUseCase は NotifyFollowersUseCase を生成する。
func NewNotifyFollowersUseCase(notifier repository.FollowerNotifier) *NotifyFollowersUseCase {
	return &NotifyFollowersUseCase{notifier: notifier}
}

// Execute はフォロワー全員分の通知をまとめて作成する。フォロワーがいなければ何もしない。
func (uc *NotifyFollowersUseCase) Execute(ctx context.Context, actorID, postID uint, notificationType model.NotificationType) error {
	followerIDs, err := uc.notifier.FindFollowerIDs(ctx, actorID)
	if err != nil || len(followerIDs) == 0 {
		return err
	}

	notifications := make([]*model.Notification, 0, len(followerIDs))
	for _, followerID := range followerIDs {
		notifications = append(notifications, &model.Notification{
			UserID:  followerID,
			Type:    notificationType,
			ActorID: actorID,
			PostID:  &postID,
		})
	}
	return uc.notifier.CreateBatch(ctx, notifications)
}

// Notify は通知作成をバックグラウンドで実行する。レスポンスを待たせないための入り口で、
// 失敗しても呼び出し元の処理は継続する（失敗はランナーがログへ残す）。
// 上限付きワーカーで実行するため、通知量が処理能力を超えても goroutine と
// DB 接続は積み上がらない。ジョブにはリクエストと独立した期限付き ctx が渡る。
func (uc *NotifyFollowersUseCase) Notify(ctx context.Context, actorID, postID uint, notificationType model.NotificationType) {
	defaultBackgroundRunner().Submit("notify-followers", func(jobCtx context.Context) error {
		return uc.Execute(jobCtx, actorID, postID, notificationType)
	})
}
