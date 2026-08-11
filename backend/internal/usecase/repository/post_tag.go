package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostTagRepository は投稿タグの永続化に対する、usecase 側が要求する契約。
type PostTagRepository interface {
	// SetTags は投稿のタグを与えられた内容で全て置き換える。
	SetTags(ctx context.Context, postID uint, tags []string) error
	GetByPostID(ctx context.Context, postID uint) ([]string, error)
	FindPostsByTag(ctx context.Context, tag string, limit, offset int) ([]model.Post, int64, error)
	GetPopularTags(ctx context.Context, limit int) ([]model.TagCount, error)
}
