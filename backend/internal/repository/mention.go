package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// MentionRepository はメンションのGORM実装。
type MentionRepository struct {
	db *gorm.DB
}

// NewMentionRepository は新しいMentionRepositoryインスタンスを生成する。
func NewMentionRepository(db *gorm.DB) *MentionRepository {
	return &MentionRepository{db: db}
}

func (r *MentionRepository) Create(mention *model.Mention) error {
	return r.db.Create(mention).Error
}

func (r *MentionRepository) FindByUserID(userID uint, page, limit int) ([]model.Mention, error) {
	var mentions []model.Mention
	offset := (page - 1) * limit
	err := r.db.Where("user_id = ?", userID).
		Preload("Actor").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&mentions).Error
	return mentions, err
}

func (r *MentionRepository) FindByPostID(postID uint) ([]model.Mention, error) {
	var mentions []model.Mention
	err := r.db.Where("post_id = ?", postID).
		Preload("User").Preload("Actor").
		Find(&mentions).Error
	return mentions, err
}

func (r *MentionRepository) FindByCommentID(commentID uint) ([]model.Mention, error) {
	var mentions []model.Mention
	err := r.db.Where("comment_id = ?", commentID).
		Preload("User").Preload("Actor").
		Find(&mentions).Error
	return mentions, err
}

func (r *MentionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Mention{}, id).Error
}

func (r *MentionRepository) DeleteByPostID(postID uint) error {
	return r.db.Where("post_id = ?", postID).Delete(&model.Mention{}).Error
}

func (r *MentionRepository) DeleteByCommentID(commentID uint) error {
	return r.db.Where("comment_id = ?", commentID).Delete(&model.Mention{}).Error
}
