package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// messageRepository は [repository.MessageRepository] の sqlc(pgx) 実装。
type messageRepository struct {
	q *sqlcgen.Queries
}

// NewMessageRepository は MessageRepository の sqlc(pgx) 実装を返す。
func NewMessageRepository(q *sqlcgen.Queries) repository.MessageRepository {
	return &messageRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.MessageRepository = (*messageRepository)(nil)

// toModelMessage は sqlc の生成行を model.Message へ変換する（関連の Sender/Receiver は含まない）。
func toModelMessage(row sqlcgen.Message) model.Message {
	return model.Message{
		ID:         uint(row.ID),
		SenderID:   uint(row.SenderID),
		ReceiverID: uint(row.ReceiverID),
		Content:    row.Content,
		Read:       row.Read,
		CreatedAt:  timeValue(fromTimestamptz(row.CreatedAt)),
	}
}

// Create はメッセージを保存する。
func (r *messageRepository) Create(ctx context.Context, msg *model.Message) error {
	row, err := r.q.CreateMessage(ctx, sqlcgen.CreateMessageParams{
		SenderID:   int64(msg.SenderID),
		ReceiverID: int64(msg.ReceiverID),
		Content:    msg.Content,
	})
	if err != nil {
		return err
	}
	*msg = toModelMessage(row)
	return nil
}

// GetConversation は 2 ユーザー間の会話を古い順に取得する。
func (r *messageRepository) GetConversation(ctx context.Context, userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListConversationMessages(ctx, sqlcgen.ListConversationMessagesParams{
		SenderID:   int64(userID),
		ReceiverID: int64(otherUserID),
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
	})
	if err != nil {
		return nil, err
	}
	messages := make([]model.Message, len(rows))
	for i, row := range rows {
		messages[i] = toModelMessage(row.Message)
		messages[i].Sender = toModelUser(row.User)
		messages[i].Receiver = toModelUser(row.User_2)
	}
	return messages, nil
}

// GetConversations は会話相手ごとの最新メッセージと未読件数を取得する。
func (r *messageRepository) GetConversations(ctx context.Context, userID uint) ([]model.ConversationSummary, error) {
	rows, err := r.q.ListConversationSummaries(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	conversations := make([]model.ConversationSummary, len(rows))
	for i, row := range rows {
		otherID, _ := row.UserID.(int64)
		conversations[i] = model.ConversationSummary{
			UserID:      uint(otherID),
			Name:        row.Name,
			AvatarURL:   fromStringPtr(row.AvatarUrl),
			LastMessage: row.LastMessage,
			LastTime:    row.LastTime,
			UnreadCount: int(row.UnreadCount),
		}
	}
	return conversations, nil
}

// MarkAsRead は指定送信者からの未読メッセージをすべて既読にする。
func (r *messageRepository) MarkAsRead(ctx context.Context, senderID, receiverID uint) error {
	return r.q.MarkMessagesAsRead(ctx, sqlcgen.MarkMessagesAsReadParams{
		SenderID:   int64(senderID),
		ReceiverID: int64(receiverID),
	})
}
