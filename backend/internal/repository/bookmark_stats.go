package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// BookmarkStatsRepository はユーザーブックマーク集計統計の取得を担当するリポジトリ実装。
type BookmarkStatsRepository struct {
	db *gorm.DB
}

// NewBookmarkStatsRepository は新しいBookmarkStatsRepositoryインスタンスを生成する。
func NewBookmarkStatsRepository(db *gorm.DB) *BookmarkStatsRepository {
	return &BookmarkStatsRepository{db: db}
}

// GetBookmarkStats は指定ユーザーのブックマーク集計統計を返す。
func (r *BookmarkStatsRepository) GetBookmarkStats(userID uint) (*model.BookmarkStats, error) {
	var stats model.BookmarkStats

	// ブックマークした数
	if err := r.db.Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&stats.TotalBookmarksMade).Error; err != nil {
		return nil, err
	}

	// 投稿がブックマークされた回数
	if err := r.db.Model(&model.Bookmark{}).
		Joins("JOIN posts ON posts.id = bookmarks.post_id").
		Where("posts.user_id = ?", userID).
		Count(&stats.TotalBookmarksReceived).Error; err != nil {
		return nil, err
	}

	// 今月ブックマークした数
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if err := r.db.Model(&model.Bookmark{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.BookmarksThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
