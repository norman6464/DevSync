package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// MentionStatsRepository はユーザーメンション集計統計の取得を担当するリポジトリ実装。
type MentionStatsRepository struct {
	db *gorm.DB
}

// NewMentionStatsRepository は新しいMentionStatsRepositoryインスタンスを生成する。
func NewMentionStatsRepository(db *gorm.DB) *MentionStatsRepository {
	return &MentionStatsRepository{db: db}
}

// GetMentionStats は指定ユーザーのメンション集計統計を返す。
func (r *MentionStatsRepository) GetMentionStats(userID uint) (*model.MentionStats, error) {
	var stats model.MentionStats

	// メンションされた回数
	if err := r.db.Model(&model.Mention{}).Where("user_id = ?", userID).Count(&stats.MentionsReceived).Error; err != nil {
		return nil, err
	}

	// メンションした回数
	if err := r.db.Model(&model.Mention{}).Where("actor_id = ?", userID).Count(&stats.MentionsMade).Error; err != nil {
		return nil, err
	}

	// 今月メンションされた回数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := r.db.Model(&model.Mention{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.MentionsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
