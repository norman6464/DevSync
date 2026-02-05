package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(msg *model.Message) error {
	return r.db.Create(msg).Error
}

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

type ConversationSummary struct {
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	LastMessage string `json:"last_message"`
	LastTime    string `json:"last_time"`
	UnreadCount int    `json:"unread_count"`
}

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

func (r *MessageRepository) MarkAsRead(senderID, receiverID uint) error {
	return r.db.Model(&model.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND read = false", senderID, receiverID).
		Update("read", true).Error
}
