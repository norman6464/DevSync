package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// commentStatsRepository は [repository.CommentStatsRepository] の GORM 実装。
type commentStatsRepository struct {
	db *gorm.DB
}

// NewCommentStatsRepository は CommentStatsRepository の GORM 実装を返す。
func NewCommentStatsRepository(db *gorm.DB) repository.CommentStatsRepository {
	return &commentStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.CommentStatsRepository = (*commentStatsRepository)(nil)

// GetCommentStats は指定ユーザーのコメント活動集計統計を返す。
func (r *commentStatsRepository) GetCommentStats(ctx context.Context, userID uint) (*model.CommentStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.CommentStats

	// ユーザーが投稿したコメント総数（parent_id IS NULL = トップレベルコメント）
	if err := db.Model(&model.Comment{}).Where("user_id = ? AND parent_id IS NULL", userID).Count(&stats.TotalComments).Error; err != nil {
		return nil, err
	}

	// ユーザーが投稿した返信総数（parent_id IS NOT NULL）
	if err := db.Model(&model.Comment{}).Where("user_id = ? AND parent_id IS NOT NULL", userID).Count(&stats.TotalReplies).Error; err != nil {
		return nil, err
	}

	// ユーザーの投稿に付いたコメント数
	if err := db.Model(&model.Comment{}).
		Joins("JOIN posts ON posts.id = comments.post_id").
		Where("posts.user_id = ? AND comments.user_id != ?", userID, userID).
		Count(&stats.CommentsReceived).Error; err != nil {
		return nil, err
	}

	// 今月のコメント数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := db.Model(&model.Comment{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.CommentsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
