package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postReader は [repository.PostReader] の GORM 実装。
// 所有権チェックに必要な投稿読み取りだけを提供する。
type postReader struct {
	db *gorm.DB
}

// NewPostReader は PostReader の GORM 実装を返す。
func NewPostReader(db *gorm.DB) repository.PostReader {
	return &postReader{db: db}
}

var _ repository.PostReader = (*postReader)(nil)

// FindByID は ID で投稿を取得する（既存 PostRepository.FindByID と同じ preload）。
func (r *postReader) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.WithContext(ctx).Preload("User").Preload("CodeSnippets").First(&post, id).Error
	return &post, err
}
