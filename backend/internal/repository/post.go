package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").First(&post, id).Error
	return &post, err
}

func (r *PostRepository) FindAll(page, limit int) ([]model.Post, error) {
	var posts []model.Post
	offset := (page - 1) * limit
	err := r.db.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, err
}

func (r *PostRepository) FindByUserID(userID uint) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

func (r *PostRepository) Timeline(userID uint, page, limit int) ([]model.Post, error) {
	var posts []model.Post
	offset := (page - 1) * limit
	err := r.db.Preload("User").
		Where("user_id IN (SELECT followee_id FROM follows WHERE follower_id = ?) OR user_id = ?", userID, userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *PostRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}

// Like
func (r *PostRepository) Like(userID, postID uint) error {
	err := r.db.Create(&model.Like{UserID: userID, PostID: postID}).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *PostRepository) Unlike(userID, postID uint) error {
	result := r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Like{})
	if result.RowsAffected > 0 {
		r.db.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
	}
	return result.Error
}

func (r *PostRepository) HasLiked(userID, postID uint) bool {
	var count int64
	r.db.Model(&model.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

// Comment
func (r *PostRepository) CreateComment(comment *model.Comment) error {
	err := r.db.Create(comment).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

func (r *PostRepository) GetComments(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("User").Where("post_id = ?", postID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

func (r *PostRepository) DeleteComment(id, userID uint) error {
	var comment model.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return err
	}
	if comment.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	r.db.Model(&model.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))
	return r.db.Delete(&comment).Error
}
