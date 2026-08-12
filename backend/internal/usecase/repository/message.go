package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// MessageRepository はダイレクトメッセージの永続化に対する、usecase 側が要求する契約。
type MessageRepository interface {
	Create(ctx context.Context, msg *model.Message) error
	// GetConversation は 2 ユーザー間のメッセージを古い順にページネーションして返す。
	GetConversation(ctx context.Context, userID, otherUserID uint, page, limit int) ([]model.Message, error)
	// GetConversations は会話相手ごとの最新メッセージと未読件数を返す。
	GetConversations(ctx context.Context, userID uint) ([]model.ConversationSummary, error)
	// MarkAsRead は指定送信者から受信者への未読メッセージをすべて既読にする。
	MarkAsRead(ctx context.Context, senderID, receiverID uint) error
}
