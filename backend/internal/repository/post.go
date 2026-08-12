package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PostRepository は投稿データへのアクセスを提供するリポジトリ実装。
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository は新しいPostRepositoryインスタンスを生成する。
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// SearchWithFilter はタグ・日付範囲・ソート順による高度な投稿検索を実行する。
// 下書きは検索対象外。タグはAND条件で絞り込む。
func (r *PostRepository) SearchWithFilter(
	query string,
	tags []string,
	sortBy string,
	dateFrom, dateTo *time.Time,
	limit, offset int,
) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	searchPattern := EscapeLikePattern(query)
	db := r.db.Preload("User").Preload("CodeSnippets").
		Where("(title LIKE ? OR content LIKE ?) AND is_draft = ?", searchPattern, searchPattern, false)

	// 日付範囲フィルター
	if dateFrom != nil {
		db = db.Where("created_at >= ?", dateFrom)
	}
	if dateTo != nil {
		db = db.Where("created_at <= ?", dateTo)
	}

	// タグフィルター（AND条件：全タグが付与されている投稿のみ）
	for _, tag := range tags {
		db = db.Where("id IN (SELECT post_id FROM post_tags WHERE tag = ?)", tag)
	}

	// ソート順
	switch sortBy {
	case "popular":
		db = db.Order("like_count DESC")
	case "views":
		db = db.Order("view_count DESC")
	default:
		db = db.Order("created_at DESC")
	}

	if err := db.Model(&model.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}
