package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// MessageStatsRepository はユーザーメッセージ集計統計の取得を担当するリポジトリ実装。
type MessageStatsRepository struct {
	db *gorm.DB
}

// NewMessageStatsRepository は新しいMessageStatsRepositoryインスタンスを生成する。
func NewMessageStatsRepository(db *gorm.DB) *MessageStatsRepository {
	return &MessageStatsRepository{db: db}
}

// GetMessageStats は指定ユーザーのメッセージ集計統計を返す。
func (r *MessageStatsRepository) GetMessageStats(userID uint) (*model.MessageStats, error) {
	var stats model.MessageStats

	// 送信メッセージ総数
	if err := r.db.Model(&model.Message{}).Where("sender_id = ?", userID).Count(&stats.TotalSent).Error; err != nil {
		return nil, err
	}

	// 受信メッセージ総数
	if err := r.db.Model(&model.Message{}).Where("receiver_id = ?", userID).Count(&stats.TotalReceived).Error; err != nil {
		return nil, err
	}

	// 会話相手数（送信先+受信元のユニーク数）
	if err := r.db.Model(&model.Message{}).
		Select("COUNT(DISTINCT CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END)", userID).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Scan(&stats.ConversationsCount).Error; err != nil {
		return nil, err
	}

	// 今月の送信メッセージ数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := r.db.Model(&model.Message{}).Where("sender_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.MessagesThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
