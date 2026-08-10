package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// mentionStatsRepository は [repository.MentionStatsRepository] の GORM 実装。
type mentionStatsRepository struct {
	db *gorm.DB
}

// NewMentionStatsRepository は MentionStatsRepository の GORM 実装を返す。
func NewMentionStatsRepository(db *gorm.DB) repository.MentionStatsRepository {
	return &mentionStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.MentionStatsRepository = (*mentionStatsRepository)(nil)

// GetMentionStats は指定ユーザーのメンション集計統計を返す。
func (r *mentionStatsRepository) GetMentionStats(ctx context.Context, userID uint) (*model.MentionStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.MentionStats

	// メンションされた回数
	if err := db.Model(&model.Mention{}).Where("user_id = ?", userID).Count(&stats.MentionsReceived).Error; err != nil {
		return nil, err
	}

	// メンションした回数
	if err := db.Model(&model.Mention{}).Where("actor_id = ?", userID).Count(&stats.MentionsMade).Error; err != nil {
		return nil, err
	}

	// 今月メンションされた回数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := db.Model(&model.Mention{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.MentionsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
