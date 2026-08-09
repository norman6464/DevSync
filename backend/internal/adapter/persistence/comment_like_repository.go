package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// commentLikeRepository は [repository.CommentLikeRepository] の GORM 実装。
type commentLikeRepository struct {
	db *gorm.DB
}

// NewCommentLikeRepository は CommentLikeRepository の GORM 実装を返す。
func NewCommentLikeRepository(db *gorm.DB) repository.CommentLikeRepository {
	return &commentLikeRepository{db: db}
}

var _ repository.CommentLikeRepository = (*commentLikeRepository)(nil)

// Like はいいねを追加し、コメントのいいね数をインクリメントする（1 トランザクション）。
func (r *commentLikeRepository) Like(ctx context.Context, userID, commentID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		like := &model.CommentLike{UserID: userID, CommentID: commentID}
		if err := tx.Create(like).Error; err != nil {
			return err
		}
		return tx.Model(&model.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
	})
}

// Unlike はいいねを削除し、コメントのいいね数をデクリメントする（1 トランザクション）。
func (r *commentLikeRepository) Unlike(ctx context.Context, userID, commentID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND comment_id = ?", userID, commentID).
			Delete(&model.CommentLike{}).Error; err != nil {
			return err
		}
		return tx.Model(&model.Comment{}).Where("id = ? AND like_count > 0", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
	})
}

// HasLiked は指定ユーザーがコメントをいいね済みかを返す。
func (r *commentLikeRepository) HasLiked(ctx context.Context, userID, commentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CommentLike{}).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Count(&count).Error
	return count > 0, err
}

// CountByCommentID はコメントのいいね数を返す。
func (r *commentLikeRepository) CountByCommentID(ctx context.Context, commentID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CommentLike{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error
	return count, err
}
