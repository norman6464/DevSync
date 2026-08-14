// Package notify は通知の作成に、受信者へのリアルタイム配信を上乗せする。
// 保存そのものは既存の port（実装は adapter/persistence）に委ね、
// 保存に成功したものだけを配信する。
package notify

import (
	"context"
	"encoding/json"
	"log"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// notificationMessageType は WebSocket のメッセージ種別。
// クライアントはこの値で通知を判別する。
const notificationMessageType = "notification"

// UserSender は特定ユーザーへメッセージを届けるために、この層が必要とする最小の契約。
type UserSender interface {
	// SendToUser は接続中のユーザーへメッセージを送る。接続が無ければ何もしない。
	SendToUser(userID uint, message []byte)
}

// message は通知 1 件の配信ペイロード。
// クライアントは通知の一覧・未読数を取り直すため、描画に要る最小限だけを載せる。
type message struct {
	Type         string `json:"type"`
	ID           uint   `json:"id"`
	Notification string `json:"notification_type"`
	ActorID      uint   `json:"actor_id"`
	PostID       *uint  `json:"post_id,omitempty"`
	QuestionID   *uint  `json:"question_id,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// BroadcastingCreator は通知を保存したうえで受信者へ配信する。
type BroadcastingCreator struct {
	inner  repository.NotificationCreator
	sender UserSender
}

var _ repository.NotificationCreator = (*BroadcastingCreator)(nil)

// NewBroadcastingCreator は保存後に配信する NotificationCreator を返す。
func NewBroadcastingCreator(inner repository.NotificationCreator, sender UserSender) *BroadcastingCreator {
	return &BroadcastingCreator{inner: inner, sender: sender}
}

// Create は通知を保存し、成功したら受信者へ配信する。
func (c *BroadcastingCreator) Create(ctx context.Context, notification *model.Notification) error {
	if err := c.inner.Create(ctx, notification); err != nil {
		return err
	}
	send(c.sender, notification)
	return nil
}

// BroadcastingFollowerNotifier はフォロワーへの一括通知を保存したうえで各受信者へ配信する。
type BroadcastingFollowerNotifier struct {
	inner  repository.FollowerNotifier
	sender UserSender
}

var _ repository.FollowerNotifier = (*BroadcastingFollowerNotifier)(nil)

// NewBroadcastingFollowerNotifier は保存後に配信する FollowerNotifier を返す。
func NewBroadcastingFollowerNotifier(inner repository.FollowerNotifier, sender UserSender) *BroadcastingFollowerNotifier {
	return &BroadcastingFollowerNotifier{inner: inner, sender: sender}
}

// FindFollowerIDs は内側の実装にそのまま委ねる。
func (n *BroadcastingFollowerNotifier) FindFollowerIDs(ctx context.Context, userID uint) ([]uint, error) {
	return n.inner.FindFollowerIDs(ctx, userID)
}

// CreateBatch は通知をまとめて保存し、成功したら受信者ごとに配信する。
func (n *BroadcastingFollowerNotifier) CreateBatch(ctx context.Context, notifications []*model.Notification) error {
	if err := n.inner.CreateBatch(ctx, notifications); err != nil {
		return err
	}
	for _, notification := range notifications {
		send(n.sender, notification)
	}
	return nil
}

// send は通知 1 件を受信者へ配信する。配信の失敗は通知の保存を妨げない。
func send(sender UserSender, notification *model.Notification) {
	if sender == nil || notification == nil {
		return
	}
	data, err := json.Marshal(message{
		Type:         notificationMessageType,
		ID:           notification.ID,
		Notification: string(notification.Type),
		ActorID:      notification.ActorID,
		PostID:       notification.PostID,
		QuestionID:   notification.QuestionID,
		CreatedAt:    notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
	if err != nil {
		log.Printf("通知の配信に失敗しました（ユーザー %d）: %v", notification.UserID, err)
		return
	}
	sender.SendToUser(notification.UserID, data)
}
