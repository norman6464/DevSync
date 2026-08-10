package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// messageStatsRepository は [repository.MessageStatsRepository] の GORM 実装。
type messageStatsRepository struct {
	db *gorm.DB
}

// NewMessageStatsRepository は MessageStatsRepository の GORM 実装を返す。
func NewMessageStatsRepository(db *gorm.DB) repository.MessageStatsRepository {
	return &messageStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.MessageStatsRepository = (*messageStatsRepository)(nil)

// GetMessageStats は指定ユーザーのメッセージ集計統計を返す。
func (r *messageStatsRepository) GetMessageStats(ctx context.Context, userID uint) (*model.MessageStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.MessageStats

	// 送信メッセージ総数
	if err := db.Model(&model.Message{}).Where("sender_id = ?", userID).Count(&stats.TotalSent).Error; err != nil {
		return nil, err
	}

	// 受信メッセージ総数
	if err := db.Model(&model.Message{}).Where("receiver_id = ?", userID).Count(&stats.TotalReceived).Error; err != nil {
		return nil, err
	}

	// 会話相手数（送信先+受信元のユニーク数）
	if err := db.Model(&model.Message{}).
		Select("COUNT(DISTINCT CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END)", userID).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Scan(&stats.ConversationsCount).Error; err != nil {
		return nil, err
	}

	// 今月の送信メッセージ数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := db.Model(&model.Message{}).Where("sender_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.MessagesThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
