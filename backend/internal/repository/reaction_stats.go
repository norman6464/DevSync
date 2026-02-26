package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// ReactionStatsRepository はユーザーリアクション集計統計の取得を担当するリポジトリ実装。
type ReactionStatsRepository struct {
	db *gorm.DB
}

// NewReactionStatsRepository は新しいReactionStatsRepositoryインスタンスを生成する。
func NewReactionStatsRepository(db *gorm.DB) *ReactionStatsRepository {
	return &ReactionStatsRepository{db: db}
}

// GetReactionStats は指定ユーザーが受け取ったリアクション集計統計を返す。
func (r *ReactionStatsRepository) GetReactionStats(userID uint) (*model.ReactionStats, error) {
	var stats model.ReactionStats

	// 受け取ったリアクション総数（自分の投稿へのリアクション）
	if err := r.db.Model(&model.Reaction{}).
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ?", userID).
		Count(&stats.TotalReactionsReceived).Error; err != nil {
		return nil, err
	}

	// リアクションしたユニークユーザー数
	if err := r.db.Model(&model.Reaction{}).
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ?", userID).
		Distinct("reactions.user_id").
		Count(&stats.UniqueReactors).Error; err != nil {
		return nil, err
	}

	// 今月受け取ったリアクション数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := r.db.Model(&model.Reaction{}).
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ? AND reactions.created_at >= ?", userID, startOfMonth).
		Count(&stats.ReactionsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetEmojiBreakdown は指定ユーザーの全投稿に対する絵文字別リアクション集計を返す。
func (r *ReactionStatsRepository) GetEmojiBreakdown(userID uint) ([]model.ReactionCount, error) {
	var counts []model.ReactionCount
	err := r.db.Model(&model.Reaction{}).
		Select("reactions.emoji, COUNT(*) as count").
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ?", userID).
		Group("reactions.emoji").
		Order("count DESC").
		Find(&counts).Error
	return counts, err
}

// GetTopReactedPosts は指定ユーザーの投稿のうちリアクション数が多い順にlimit件返す。
func (r *ReactionStatsRepository) GetTopReactedPosts(userID uint, limit int) ([]model.TopReactedPost, error) {
	var posts []model.TopReactedPost
	err := r.db.Model(&model.Reaction{}).
		Select("posts.id, posts.title, COUNT(*) as reaction_count").
		Joins("JOIN posts ON posts.id = reactions.post_id").
		Where("posts.user_id = ?", userID).
		Group("posts.id, posts.title").
		Order("reaction_count DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}
