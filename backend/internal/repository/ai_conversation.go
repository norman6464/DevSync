package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// AIConversationRepository はAI会話セッションデータのDB操作を実装する。
type AIConversationRepository struct {
	db *gorm.DB
}

// NewAIConversationRepository は新しいAIConversationRepositoryインスタンスを生成する。
func NewAIConversationRepository(db *gorm.DB) *AIConversationRepository {
	return &AIConversationRepository{db: db}
}

// CreateConversation は新しいAI会話セッションをDBに作成する。
func (r *AIConversationRepository) CreateConversation(conv *model.AIConversation) error {
	return r.db.Create(conv).Error
}

// FindConversationsByUserID は指定ユーザーIDの会話一覧を更新日時降順で取得する。
func (r *AIConversationRepository) FindConversationsByUserID(userID uint, limit, offset int) ([]model.AIConversation, error) {
	var conversations []model.AIConversation
	err := r.db.Preload("Messages").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&conversations).Error
	return conversations, err
}

// FindConversationByID は指定IDの会話をメッセージ付きで取得する。
func (r *AIConversationRepository) FindConversationByID(id uint) (*model.AIConversation, error) {
	var conv model.AIConversation
	err := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&conv, id).Error
	return &conv, err
}

// AddMessage は会話にメッセージを追加し、会話のUpdatedAtを更新する。
func (r *AIConversationRepository) AddMessage(msg *model.AIMessage) error {
	if err := r.db.Create(msg).Error; err != nil {
		return err
	}
	return r.db.Model(&model.AIConversation{}).
		Where("id = ?", msg.ConversationID).
		Update("updated_at", time.Now()).Error
}

// GetMessages は指定会話IDのメッセージ一覧を時系列順で取得する。
func (r *AIConversationRepository) GetMessages(conversationID uint) ([]model.AIMessage, error) {
	var messages []model.AIMessage
	err := r.db.Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// CountTodayMessages は本日のユーザーメッセージ数を取得する（レート制限用）。
func (r *AIConversationRepository) CountTodayMessages(userID uint) (int64, error) {
	var count int64
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	err := r.db.Model(&model.AIMessage{}).
		Joins("JOIN ai_conversations ON ai_conversations.id = ai_messages.conversation_id").
		Where("ai_conversations.user_id = ? AND ai_messages.role = ? AND ai_messages.created_at >= ?",
			userID, model.AIMessageRoleUser, startOfDay).
		Count(&count).Error
	return count, err
}

// DeleteConversation は指定会話を所有権確認後に削除する（メッセージもカスケード削除）。
func (r *AIConversationRepository) DeleteConversation(id, userID uint) error {
	var conv model.AIConversation
	if err := r.db.First(&conv, id).Error; err != nil {
		return err
	}
	if conv.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	// メッセージを先に削除
	if err := r.db.Where("conversation_id = ?", id).Delete(&model.AIMessage{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&conv).Error
}
