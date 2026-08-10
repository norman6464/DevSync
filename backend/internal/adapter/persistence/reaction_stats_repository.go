package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// reactionStatsRepository は [repository.ReactionStatsRepository] の GORM 実装。
type reactionStatsRepository struct {
	db *gorm.DB
}

// NewReactionStatsRepository は ReactionStatsRepository の GORM 実装を返す。
func NewReactionStatsRepository(db *gorm.DB) repository.ReactionStatsRepository {
	return &reactionStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ReactionStatsRepository = (*reactionStatsRepository)(nil)

// GetReactionStats は指定ユーザーが受け取ったリアクション集計統計を返す。
func (r *reactionStatsRepository) GetReactionStats(ctx context.Context, userID uint) (*model.ReactionStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.ReactionStats

	// 受け取ったリアクション総数（自分の投稿へのリアクション）
	if err := db.Model(&model.Reaction{}).
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ?", userID).
		Count(&stats.TotalReactionsReceived).Error; err != nil {
		return nil, err
	}

	// リアクションしたユニークユーザー数
	if err := db.Model(&model.Reaction{}).
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ?", userID).
		Distinct("reactions.user_id").
		Count(&stats.UniqueReactors).Error; err != nil {
		return nil, err
	}

	// 今月受け取ったリアクション数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := db.Model(&model.Reaction{}).
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ? AND reactions.created_at >= ?", userID, startOfMonth).
		Count(&stats.ReactionsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetEmojiBreakdown は指定ユーザーの全投稿に対する絵文字別リアクション集計を返す。
func (r *reactionStatsRepository) GetEmojiBreakdown(ctx context.Context, userID uint) ([]model.ReactionCount, error) {
	var counts []model.ReactionCount
	err := r.db.WithContext(ctx).Model(&model.Reaction{}).
		Select("reactions.emoji, COUNT(*) as count").
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ?", userID).
		Group("reactions.emoji").
		Order("count DESC").
		Order("reactions.emoji ASC").
		Find(&counts).Error
	return counts, err
}

// GetTopReactedPosts は指定ユーザーの投稿のうちリアクション数が多い順に limit 件返す。
func (r *reactionStatsRepository) GetTopReactedPosts(ctx context.Context, userID uint, limit int) ([]model.TopReactedPost, error) {
	var posts []model.TopReactedPost
	err := r.db.WithContext(ctx).Model(&model.Reaction{}).
		Select("posts.id, posts.title, COUNT(*) as reaction_count").
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ?", userID).
		Group("posts.id, posts.title").
		Order("reaction_count DESC").
		Order("posts.id DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}
