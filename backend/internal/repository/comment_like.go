package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// CommentLikeRepository はコメントへのいいね操作のDB実装。
type CommentLikeRepository struct {
	db *gorm.DB
}

// NewCommentLikeRepository は新しいCommentLikeRepositoryインスタンスを生成する。
func NewCommentLikeRepository(db *gorm.DB) *CommentLikeRepository {
	return &CommentLikeRepository{db: db}
}

// Like はコメントにいいねを追加し、コメントのいいね数をインクリメントする。
func (r *CommentLikeRepository) Like(userID, commentID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		like := &model.CommentLike{UserID: userID, CommentID: commentID}
		if err := tx.Create(like).Error; err != nil {
			return err
		}
		return tx.Model(&model.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
	})
}

// Unlike はコメントのいいねを削除し、コメントのいいね数をデクリメントする。
func (r *CommentLikeRepository) Unlike(userID, commentID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND comment_id = ?", userID, commentID).
			Delete(&model.CommentLike{}).Error; err != nil {
			return err
		}
		return tx.Model(&model.Comment{}).Where("id = ? AND like_count > 0", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
	})
}

// HasLiked は指定ユーザーがコメントをいいね済みかどうかを返す。
func (r *CommentLikeRepository) HasLiked(userID, commentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.CommentLike{}).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Count(&count).Error
	return count > 0, err
}

// CountByCommentID はコメントのいいね数を返す。
func (r *CommentLikeRepository) CountByCommentID(commentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.CommentLike{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error
	return count, err
}
