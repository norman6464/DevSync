package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postSearchRepository は [repository.PostSearchRepository] の GORM 実装。
type postSearchRepository struct {
	db *gorm.DB
}

// NewPostSearchRepository は PostSearchRepository の GORM 実装を返す。
func NewPostSearchRepository(db *gorm.DB) repository.PostSearchRepository {
	return &postSearchRepository{db: db}
}

var _ repository.PostSearchRepository = (*postSearchRepository)(nil)

// SearchWithFilter はタグ・日付範囲・ソート順による高度な投稿検索を実行する。
// 下書きは検索対象外。タグは AND 条件で絞り込む。
func (r *postSearchRepository) SearchWithFilter(ctx context.Context, params model.PostSearchParams) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	searchPattern := escapeLikePattern(params.Query)
	db := r.db.WithContext(ctx).Preload("User").Preload("CodeSnippets").
		Where("(title LIKE ? OR content LIKE ?) AND is_draft = ?", searchPattern, searchPattern, false)

	// 日付範囲フィルター
	if params.DateFrom != nil {
		db = db.Where("created_at >= ?", params.DateFrom)
	}
	if params.DateTo != nil {
		db = db.Where("created_at <= ?", params.DateTo)
	}

	// タグフィルター（AND 条件：全タグが付与されている投稿のみ）
	for _, tag := range params.Tags {
		db = db.Where("id IN (SELECT post_id FROM post_tags WHERE tag = ?)", tag)
	}

	// ソート順
	switch params.SortBy {
	case model.SearchSortByPopular:
		db = db.Order("like_count DESC")
	case model.SearchSortByViews:
		db = db.Order("view_count DESC")
	default:
		db = db.Order("created_at DESC")
	}

	if err := db.Model(&model.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(params.Offset).Limit(params.Limit).Find(&posts).Error
	return posts, total, err
}
