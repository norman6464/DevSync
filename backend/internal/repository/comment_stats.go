package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// CommentStatsRepository はユーザーコメント活動集計統計の取得を担当するリポジトリ実装。
type CommentStatsRepository struct {
	db *gorm.DB
}

// NewCommentStatsRepository は新しいCommentStatsRepositoryインスタンスを生成する。
func NewCommentStatsRepository(db *gorm.DB) *CommentStatsRepository {
	return &CommentStatsRepository{db: db}
}

// GetCommentStats は指定ユーザーのコメント活動集計統計を返す。
func (r *CommentStatsRepository) GetCommentStats(userID uint) (*model.CommentStats, error) {
	var stats model.CommentStats

	// ユーザーが投稿したコメント総数（parent_id IS NULL = トップレベルコメント）
	if err := r.db.Model(&model.Comment{}).Where("user_id = ? AND parent_id IS NULL", userID).Count(&stats.TotalComments).Error; err != nil {
		return nil, err
	}

	// ユーザーが投稿した返信総数（parent_id IS NOT NULL）
	if err := r.db.Model(&model.Comment{}).Where("user_id = ? AND parent_id IS NOT NULL", userID).Count(&stats.TotalReplies).Error; err != nil {
		return nil, err
	}

	// ユーザーの投稿に付いたコメント数
	if err := r.db.Model(&model.Comment{}).
		Joins("JOIN posts ON posts.id = comments.post_id").
		Where("posts.user_id = ? AND comments.user_id != ?", userID, userID).
		Count(&stats.CommentsReceived).Error; err != nil {
		return nil, err
	}

	// 今月のコメント数
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if err := r.db.Model(&model.Comment{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.CommentsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
