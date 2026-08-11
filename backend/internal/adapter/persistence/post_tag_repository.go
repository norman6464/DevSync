package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postTagRepository は [repository.PostTagRepository] の GORM 実装。
type postTagRepository struct {
	db *gorm.DB
}

// NewPostTagRepository は PostTagRepository の GORM 実装を返す。
func NewPostTagRepository(db *gorm.DB) repository.PostTagRepository {
	return &postTagRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostTagRepository = (*postTagRepository)(nil)

// SetTags は投稿のタグを全て置き換える（削除と挿入を 1 トランザクションで行う）。
func (r *postTagRepository) SetTags(ctx context.Context, postID uint, tags []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("post_id = ?", postID).Delete(&model.PostTag{}).Error; err != nil {
			return err
		}
		for _, tag := range tags {
			if err := tx.Create(&model.PostTag{PostID: postID, Tag: tag}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByPostID は投稿のタグ一覧を取得する。
func (r *postTagRepository) GetByPostID(ctx context.Context, postID uint) ([]string, error) {
	var postTags []model.PostTag
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&postTags).Error; err != nil {
		return nil, err
	}
	tags := make([]string, len(postTags))
	for i, pt := range postTags {
		tags[i] = pt.Tag
	}
	return tags, nil
}

// FindPostsByTag はタグで投稿を検索する（下書きは除外する）。
func (r *postTagRepository) FindPostsByTag(ctx context.Context, tag string, limit, offset int) ([]model.Post, int64, error) {
	db := r.db.WithContext(ctx)

	var count int64
	if err := db.Model(&model.PostTag{}).Where("tag = ?", tag).
		Joins("JOIN posts ON posts.id = post_tags.post_id AND posts.is_draft = false").
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var postTags []model.PostTag
	if err := db.Where("tag = ?", tag).
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
	if err := db.Preload("User").Where("id IN ?", postIDs).
		Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

// GetPopularTags は使用回数の多いタグ一覧を取得する。
func (r *postTagRepository) GetPopularTags(ctx context.Context, limit int) ([]model.TagCount, error) {
	var results []model.TagCount
	if err := r.db.WithContext(ctx).Model(&model.PostTag{}).
		Select("tag, COUNT(*) as count").
		Group("tag").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
