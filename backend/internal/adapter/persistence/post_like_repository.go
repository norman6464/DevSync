package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postLikeRepository は [repository.PostLikeRepository] の GORM 実装。
type postLikeRepository struct {
	db *gorm.DB
}

// NewPostLikeRepository は PostLikeRepository の GORM 実装を返す。
func NewPostLikeRepository(db *gorm.DB) repository.PostLikeRepository {
	return &postLikeRepository{db: db}
}

var _ repository.PostLikeRepository = (*postLikeRepository)(nil)

// Like はいいねを追加し、投稿の like_count を加算する。
func (r *postLikeRepository) Like(ctx context.Context, userID, postID uint) error {
	if err := r.db.WithContext(ctx).Create(&model.Like{UserID: userID, PostID: postID}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", postID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// Unlike はいいねを取り消し、実際に削除できたときだけ like_count をデクリメントする。
func (r *postLikeRepository) Unlike(ctx context.Context, userID, postID uint) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Like{})
	if result.RowsAffected > 0 {
		r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
	}
	return result.Error
}

// HasLiked は指定ユーザーが投稿にいいね済みかを返す。
func (r *postLikeRepository) HasLiked(ctx context.Context, userID, postID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Like{}).
		Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
