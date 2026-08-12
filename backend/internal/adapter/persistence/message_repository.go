package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// messageRepository は [repository.MessageRepository] の GORM 実装。
type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository は MessageRepository の GORM 実装を返す。
func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.MessageRepository = (*messageRepository)(nil)

// Create はメッセージを保存する。
func (r *messageRepository) Create(ctx context.Context, msg *model.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

// GetConversation は 2 ユーザー間の会話を古い順に取得する。
func (r *messageRepository) GetConversation(ctx context.Context, userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	var messages []model.Message
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, otherUserID, otherUserID, userID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&messages).Error
	return messages, err
}

// GetConversations は会話相手ごとの最新メッセージと未読件数を取得する。
func (r *messageRepository) GetConversations(ctx context.Context, userID uint) ([]model.ConversationSummary, error) {
	var conversations []model.ConversationSummary
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ON (other_id) other_id as user_id, u.name, u.avatar_url, m.content as last_message, m.created_at as last_time,
			(SELECT COUNT(*) FROM messages WHERE sender_id = other_id AND receiver_id = ? AND read = false) as unread_count
		FROM (
			SELECT CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as other_id, id
			FROM messages
			WHERE sender_id = ? OR receiver_id = ?
		) sub
		JOIN messages m ON m.id = sub.id
		JOIN users u ON u.id = sub.other_id
		ORDER BY other_id, m.created_at DESC
	`, userID, userID, userID, userID).Scan(&conversations).Error
	return conversations, err
}

// MarkAsRead は指定送信者からの未読メッセージをすべて既読にする。
func (r *messageRepository) MarkAsRead(ctx context.Context, senderID, receiverID uint) error {
	return r.db.WithContext(ctx).Model(&model.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND read = false", senderID, receiverID).
		Update("read", true).Error
}
