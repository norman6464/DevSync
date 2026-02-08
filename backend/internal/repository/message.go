package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// MessageRepository はDMメッセージデータへのアクセスを提供するリポジトリ実装。
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository は新しいMessageRepositoryインスタンスを生成する。
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create は新しいメッセージをデータベースに作成する。
func (r *MessageRepository) Create(msg *model.Message) error {
	return r.db.Create(msg).Error
}

// GetConversation は2ユーザー間の会話をページネーション付きで取得する（古い順）。
func (r *MessageRepository) GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	var messages []model.Message
	offset := (page - 1) * limit
	err := r.db.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, otherUserID, otherUserID, userID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&messages).Error
	return messages, err
}

// ConversationSummary は会話一覧表示用のサマリー情報を表す。
type ConversationSummary struct {
	UserID      uint   `json:"user_id"`     // 会話相手のユーザーID
	Name        string `json:"name"`        // 会話相手の名前
	AvatarURL   string `json:"avatar_url"`  // 会話相手のアバターURL
	LastMessage string `json:"last_message"` // 最新メッセージの内容
	LastTime    string `json:"last_time"`   // 最新メッセージの日時
	UnreadCount int    `json:"unread_count"` // 未読メッセージ数
}

// GetConversations はユーザーの全会話一覧をサマリー形式で取得する。
// DISTINCT ON を使用して各会話相手との最新メッセージを1件ずつ取得する。
func (r *MessageRepository) GetConversations(userID uint) ([]ConversationSummary, error) {
	var conversations []ConversationSummary
	err := r.db.Raw(`
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

// MarkAsRead は指定送信者からのメッセージを全て既読にする。
func (r *MessageRepository) MarkAsRead(senderID, receiverID uint) error {
	return r.db.Model(&model.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND read = false", senderID, receiverID).
		Update("read", true).Error
}
