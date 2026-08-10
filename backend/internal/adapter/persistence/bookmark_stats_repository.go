package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// bookmarkStatsRepository は [repository.BookmarkStatsRepository] の GORM 実装。
type bookmarkStatsRepository struct {
	db *gorm.DB
}

// NewBookmarkStatsRepository は BookmarkStatsRepository の GORM 実装を返す。
func NewBookmarkStatsRepository(db *gorm.DB) repository.BookmarkStatsRepository {
	return &bookmarkStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BookmarkStatsRepository = (*bookmarkStatsRepository)(nil)

// GetBookmarkStats は指定ユーザーのブックマーク集計統計を返す。
func (r *bookmarkStatsRepository) GetBookmarkStats(ctx context.Context, userID uint) (*model.BookmarkStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.BookmarkStats

	// ブックマークした数
	if err := db.Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&stats.TotalBookmarksMade).Error; err != nil {
		return nil, err
	}

	// 投稿がブックマークされた回数
	if err := db.Model(&model.Bookmark{}).
		Joins("JOIN posts ON posts.id = bookmarks.post_id").
		Where("posts.user_id = ?", userID).
		Count(&stats.TotalBookmarksReceived).Error; err != nil {
		return nil, err
	}

	// 今月ブックマークした数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := db.Model(&model.Bookmark{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.BookmarksThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
