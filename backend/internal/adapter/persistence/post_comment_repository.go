package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postCommentRepository は [repository.PostCommentRepository] の GORM 実装。
// 単体取得は sqlc(pgx) 実装へ移行済みの commentReader を埋め込んで再利用する。
type postCommentRepository struct {
	repository.CommentReader
	db *gorm.DB
}

// NewPostCommentRepository は PostCommentRepository の GORM 実装を返す。
func NewPostCommentRepository(db *gorm.DB, q *sqlcgen.Queries) repository.PostCommentRepository {
	return &postCommentRepository{CommentReader: NewCommentReader(q), db: db}
}

var _ repository.PostCommentRepository = (*postCommentRepository)(nil)

// Create はコメントを作成し、投稿の comment_count を加算する。
func (r *postCommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", comment.PostID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

// Update はコメントを更新する。
func (r *postCommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// Delete はコメントを削除し、投稿の comment_count をデクリメントする。
// 所有権チェックは usecase 層で実施済みであること。
func (r *postCommentRepository) Delete(ctx context.Context, id uint) error {
	var comment model.Comment
	if err := r.db.WithContext(ctx).First(&comment, id).Error; err != nil {
		return err
	}
	r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", comment.PostID).
		UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))
	return r.db.WithContext(ctx).Delete(&comment).Error
}

// ListByPostID は指定投稿のトップレベルコメントをユーザー情報・返信付きで取得する（古い順）。
func (r *postCommentRepository) ListByPostID(ctx context.Context, postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.WithContext(ctx).Preload("User").Preload("Replies").Preload("Replies.User").
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// ListReplies は指定コメントへの返信をユーザー情報付きで取得する（古い順）。
func (r *postCommentRepository) ListReplies(ctx context.Context, parentID uint) ([]model.Comment, error) {
	var replies []model.Comment
	err := r.db.WithContext(ctx).Preload("User").
		Where("parent_id = ?", parentID).Order("created_at ASC").Find(&replies).Error
	return replies, err
}
