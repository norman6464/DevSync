package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// mentionRepository は [repository.MentionRepository] の GORM 実装。
type mentionRepository struct {
	db *gorm.DB
}

// NewMentionRepository は MentionRepository の GORM 実装を返す。
func NewMentionRepository(db *gorm.DB) repository.MentionRepository {
	return &mentionRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.MentionRepository = (*mentionRepository)(nil)

// Create はメンションを保存する。
func (r *mentionRepository) Create(ctx context.Context, mention *model.Mention) error {
	return r.db.WithContext(ctx).Create(mention).Error
}

// FindByUserID は指定ユーザー宛のメンションを作成日時の降順で取得する。
func (r *mentionRepository) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Mention, error) {
	var mentions []model.Mention
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Actor").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&mentions).Error
	return mentions, err
}

// FindByPostID は指定投稿に紐づくメンションを取得する。
func (r *mentionRepository) FindByPostID(ctx context.Context, postID uint) ([]model.Mention, error) {
	var mentions []model.Mention
	err := r.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Preload("User").Preload("Actor").
		Find(&mentions).Error
	return mentions, err
}

// FindByCommentID は指定コメントに紐づくメンションを取得する。
func (r *mentionRepository) FindByCommentID(ctx context.Context, commentID uint) ([]model.Mention, error) {
	var mentions []model.Mention
	err := r.db.WithContext(ctx).
		Where("comment_id = ?", commentID).
		Preload("User").Preload("Actor").
		Find(&mentions).Error
	return mentions, err
}

// DeleteByPostID は指定投稿に紐づくメンションをすべて削除する。
func (r *mentionRepository) DeleteByPostID(ctx context.Context, postID uint) error {
	return r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&model.Mention{}).Error
}

// DeleteByCommentID は指定コメントに紐づくメンションをすべて削除する。
func (r *mentionRepository) DeleteByCommentID(ctx context.Context, commentID uint) error {
	return r.db.WithContext(ctx).Where("comment_id = ?", commentID).Delete(&model.Mention{}).Error
}
