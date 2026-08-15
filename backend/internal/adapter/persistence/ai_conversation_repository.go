package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// aiConversationRepository は [repository.AIConversationRepository] の GORM 実装。
type aiConversationRepository struct {
	db *gorm.DB
}

// NewAIConversationRepository は AIConversationRepository の GORM 実装を返す。
func NewAIConversationRepository(db *gorm.DB) repository.AIConversationRepository {
	return &aiConversationRepository{db: db}
}

var _ repository.AIConversationRepository = (*aiConversationRepository)(nil)

// CreateConversation は会話を作成する。
func (r *aiConversationRepository) CreateConversation(ctx context.Context, conv *model.AIConversation) error {
	return r.db.WithContext(ctx).Create(conv).Error
}

// FindConversationsByUserID は会話を更新の新しい順にメッセージ付きで返す。
func (r *aiConversationRepository) FindConversationsByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.AIConversation, error) {
	var conversations []model.AIConversation
	err := r.db.WithContext(ctx).Preload("Messages").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&conversations).Error
	return conversations, err
}

// FindConversationByID は会話をメッセージ付き（古い順）で返す。存在しなければ (nil, nil) を返す。
func (r *aiConversationRepository) FindConversationByID(ctx context.Context, id uint) (*model.AIConversation, error) {
	var conv model.AIConversation
	err := r.db.WithContext(ctx).Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&conv, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conv, nil
}

// AddMessage はメッセージを追加し、会話の更新日時を進める。
func (r *aiConversationRepository) AddMessage(ctx context.Context, msg *model.AIMessage) error {
	if err := r.db.WithContext(ctx).Create(msg).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.AIConversation{}).
		Where("id = ?", msg.ConversationID).
		Update("updated_at", time.Now()).Error
}

// CountTodayMessages は当日のユーザー発言数を返す。
func (r *aiConversationRepository) CountTodayMessages(ctx context.Context, userID uint) (int64, error) {
	var count int64
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	err := r.db.WithContext(ctx).Model(&model.AIMessage{}).
		Joins("JOIN ai_conversations ON ai_conversations.id = ai_messages.conversation_id").
		Where("ai_conversations.user_id = ? AND ai_messages.role = ? AND ai_messages.created_at >= ?",
			userID, model.AIMessageRoleUser, startOfDay).
		Count(&count).Error
	return count, err
}

// DeleteConversation は本人の会話をメッセージごと削除する。
// 所有権の判定は usecase 側で済んでいる前提で、本人の会話だけを対象にする。
// 既に無ければ何もしない（冪等）。エラーは DB 障害だけを表す。
func (r *aiConversationRepository) DeleteConversation(ctx context.Context, id, userID uint) error {
	db := r.db.WithContext(ctx)
	var conv model.AIConversation
	err := db.Where("id = ? AND user_id = ?", id, userID).First(&conv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := db.Where("conversation_id = ?", id).Delete(&model.AIMessage{}).Error; err != nil {
		return err
	}
	return db.Delete(&conv).Error
}
