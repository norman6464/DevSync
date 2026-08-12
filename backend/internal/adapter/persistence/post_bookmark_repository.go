package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postBookmarkRepository は [repository.PostBookmarkRepository] の GORM 実装。
type postBookmarkRepository struct {
	db *gorm.DB
}

// NewPostBookmarkRepository は PostBookmarkRepository の GORM 実装を返す。
func NewPostBookmarkRepository(db *gorm.DB) repository.PostBookmarkRepository {
	return &postBookmarkRepository{db: db}
}

var _ repository.PostBookmarkRepository = (*postBookmarkRepository)(nil)

// Bookmark は投稿をブックマークし、投稿の bookmark_count を加算する。
func (r *postBookmarkRepository) Bookmark(ctx context.Context, userID, postID uint) error {
	if err := r.db.WithContext(ctx).Create(&model.Bookmark{UserID: userID, PostID: postID}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", postID).
		UpdateColumn("bookmark_count", gorm.Expr("bookmark_count + 1")).Error
}

// Unbookmark はブックマークを解除し、実際に削除できたときだけ bookmark_count をデクリメントする。
func (r *postBookmarkRepository) Unbookmark(ctx context.Context, userID, postID uint) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Bookmark{})
	if result.RowsAffected > 0 {
		r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", postID).
			UpdateColumn("bookmark_count", gorm.Expr("GREATEST(bookmark_count - 1, 0)"))
	}
	return result.Error
}

// HasBookmarked は指定ユーザーが投稿をブックマーク済みかを返す。
func (r *postBookmarkRepository) HasBookmarked(ctx context.Context, userID, postID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Bookmark{}).
		Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindBookmarkedByUserID は指定ユーザーのブックマーク済み投稿をページネーション付きで取得する。
func (r *postBookmarkRepository) FindBookmarkedByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	offset := (page - 1) * limit

	subQuery := r.db.WithContext(ctx).Model(&model.Bookmark{}).Select("post_id").Where("user_id = ?", userID)

	if err := r.db.WithContext(ctx).Model(&model.Post{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).Preload("User").Preload("CodeSnippets").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, total, err
}

// CountBookmarkedByUserID は指定ユーザーのブックマーク件数を返す。
func (r *postBookmarkRepository) CountBookmarkedByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
