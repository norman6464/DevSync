package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PostTagRepository は PostTagRepositoryInterface の GORM 実装。
type PostTagRepository struct {
	db *gorm.DB
}

// NewPostTagRepository は新しい PostTagRepository を生成する。
func NewPostTagRepository(db *gorm.DB) *PostTagRepository {
	return &PostTagRepository{db: db}
}

// SetTags は投稿のタグを全て置き換える。
func (r *PostTagRepository) SetTags(postID uint, tags []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 既存タグ削除
		if err := tx.Where("post_id = ?", postID).Delete(&model.PostTag{}).Error; err != nil {
			return err
		}
		// 新タグ挿入
		for _, tag := range tags {
			postTag := &model.PostTag{PostID: postID, Tag: tag}
			if err := tx.Create(postTag).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByPostID は投稿のタグ一覧を取得する。
func (r *PostTagRepository) GetByPostID(postID uint) ([]string, error) {
	var postTags []model.PostTag
	if err := r.db.Where("post_id = ?", postID).Find(&postTags).Error; err != nil {
		return nil, err
	}
	tags := make([]string, len(postTags))
	for i, pt := range postTags {
		tags[i] = pt.Tag
	}
	return tags, nil
}

// FindPostsByTag はタグで投稿を検索する。
func (r *PostTagRepository) FindPostsByTag(tag string, limit, offset int) ([]model.Post, int64, error) {
	var count int64
	r.db.Model(&model.PostTag{}).Where("tag = ?", tag).
		Joins("JOIN posts ON posts.id = post_tags.post_id AND posts.is_draft = false").
		Count(&count)

	var postTags []model.PostTag
	if err := r.db.Where("tag = ?", tag).
		Joins("JOIN posts ON posts.id = post_tags.post_id AND posts.is_draft = false").
		Order("post_tags.id DESC").
		Limit(limit).Offset(offset).
		Find(&postTags).Error; err != nil {
		return nil, 0, err
	}

	if len(postTags) == 0 {
		return []model.Post{}, 0, nil
	}

	postIDs := make([]uint, len(postTags))
	for i, pt := range postTags {
		postIDs[i] = pt.PostID
	}

	var posts []model.Post
	if err := r.db.Preload("User").Where("id IN ?", postIDs).
		Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

// GetPopularTags は使用回数の多いタグ一覧を取得する。
func (r *PostTagRepository) GetPopularTags(limit int) ([]model.TagCount, error) {
	var results []model.TagCount
	if err := r.db.Model(&model.PostTag{}).
		Select("tag, COUNT(*) as count").
		Group("tag").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
