package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PostStatsRepository はユーザー投稿集計統計の取得を担当するリポジトリ実装。
type PostStatsRepository struct {
	db *gorm.DB
}

// NewPostStatsRepository は新しいPostStatsRepositoryインスタンスを生成する。
func NewPostStatsRepository(db *gorm.DB) *PostStatsRepository {
	return &PostStatsRepository{db: db}
}

// GetPostStats は指定ユーザーの投稿集計統計を返す。
func (r *PostStatsRepository) GetPostStats(userID uint) (*model.PostStats, error) {
	var stats model.PostStats

	// 総投稿数
	if err := r.db.Model(&model.Post{}).Where("user_id = ?", userID).Count(&stats.TotalPosts).Error; err != nil {
		return nil, err
	}

	// 公開済み投稿数
	if err := r.db.Model(&model.Post{}).Where("user_id = ? AND is_draft = ?", userID, false).Count(&stats.PublishedPosts).Error; err != nil {
		return nil, err
	}

	// 下書き投稿数
	if err := r.db.Model(&model.Post{}).Where("user_id = ? AND is_draft = ?", userID, true).Count(&stats.DraftPosts).Error; err != nil {
		return nil, err
	}

	// 受け取ったいいね総数
	var totalLikes *int64
	if err := r.db.Model(&model.Post{}).Where("user_id = ?", userID).Select("SUM(like_count)").Scan(&totalLikes).Error; err != nil {
		return nil, err
	}
	if totalLikes != nil {
		stats.TotalLikesReceived = *totalLikes
	}

	// 受け取ったコメント総数
	var totalComments *int64
	if err := r.db.Model(&model.Post{}).Where("user_id = ?", userID).Select("SUM(comment_count)").Scan(&totalComments).Error; err != nil {
		return nil, err
	}
	if totalComments != nil {
		stats.TotalComments = *totalComments
	}

	// 今月の投稿数
	monthStart := domain.StartOfMonth(time.Now())
	if err := r.db.Model(&model.Post{}).Where("user_id = ? AND created_at >= ?", userID, monthStart).Count(&stats.PostsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
