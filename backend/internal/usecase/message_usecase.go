package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// messageMaxLen はダイレクトメッセージ本文の最大文字数。
const messageMaxLen = 5000

// ListConversationsUseCase は会話相手ごとのサマリー一覧を返す。
type ListConversationsUseCase struct {
	messages repository.MessageRepository
}

// NewListConversationsUseCase は ListConversationsUseCase を生成する。
func NewListConversationsUseCase(messages repository.MessageRepository) *ListConversationsUseCase {
	return &ListConversationsUseCase{messages: messages}
}

// Execute は会話相手ごとの最新メッセージと未読件数を返す。
func (uc *ListConversationsUseCase) Execute(ctx context.Context, userID uint) ([]model.ConversationSummary, error) {
	return uc.messages.GetConversations(ctx, userID)
}

// GetConversationUseCase は 2 ユーザー間のメッセージを取得する。
type GetConversationUseCase struct {
	messages repository.MessageRepository
}

// NewGetConversationUseCase は GetConversationUseCase を生成する。
func NewGetConversationUseCase(messages repository.MessageRepository) *GetConversationUseCase {
	return &GetConversationUseCase{messages: messages}
}

// Execute は相手からのメッセージを既読にしてから、会話を古い順に返す。
// 既読の失敗は取得結果に影響させない（移行前と同じ）。
func (uc *GetConversationUseCase) Execute(ctx context.Context, userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	_ = uc.messages.MarkAsRead(ctx, otherUserID, userID)
	return uc.messages.GetConversation(ctx, userID, otherUserID, page, limit)
}

// MarkMessagesAsReadUseCase は指定送信者からのメッセージを既読にする。
type MarkMessagesAsReadUseCase struct {
	messages repository.MessageRepository
}

// NewMarkMessagesAsReadUseCase は MarkMessagesAsReadUseCase を生成する。
func NewMarkMessagesAsReadUseCase(messages repository.MessageRepository) *MarkMessagesAsReadUseCase {
	return &MarkMessagesAsReadUseCase{messages: messages}
}

// Execute は指定送信者から受信者への未読メッセージをすべて既読にする。
func (uc *MarkMessagesAsReadUseCase) Execute(ctx context.Context, senderID, receiverID uint) error {
	return uc.messages.MarkAsRead(ctx, senderID, receiverID)
}

// SendMessageUseCase はダイレクトメッセージを送信し、受信者へ通知する。
type SendMessageUseCase struct {
	messages      repository.MessageRepository
	notifications repository.NotificationCreator
}

// NewSendMessageUseCase は SendMessageUseCase を生成する。
func NewSendMessageUseCase(
	messages repository.MessageRepository,
	notifications repository.NotificationCreator,
) *SendMessageUseCase {
	return &SendMessageUseCase{messages: messages, notifications: notifications}
}

// Execute は本文を検証してメッセージを保存し、受信者への通知を非同期に作成する。
// 自分自身へは送信できない。
func (uc *SendMessageUseCase) Execute(ctx context.Context, msg *model.Message) error {
	if err := domain.ValidateStringLength(msg.Content, 1, messageMaxLen, "メッセージ内容"); err != nil {
		return err
	}
	msg.Content = strings.TrimSpace(msg.Content)
	if msg.SenderID == msg.ReceiverID {
		return domain.ErrBadRequest
	}
	if err := uc.messages.Create(ctx, msg); err != nil {
		return err
	}

	// 受信者への通知は送信結果に影響させない（非同期・失敗はランナーがログへ残す）。
	// 上限付きワーカーで実行し、ジョブにはリクエストと独立した期限付き ctx が渡る。
	uc.notifyReceiver(msg.SenderID, msg.ReceiverID)

	return nil
}

// notifyReceiver は受信者へメッセージ受信の通知をバックグラウンドで作成する。
func (uc *SendMessageUseCase) notifyReceiver(senderID, receiverID uint) {
	defaultBackgroundRunner().Submit("message-notification", func(ctx context.Context) error {
		return uc.notifications.Create(ctx, &model.Notification{
			UserID:  receiverID,
			Type:    model.NotificationTypeMessage,
			ActorID: senderID,
		})
	})
}
