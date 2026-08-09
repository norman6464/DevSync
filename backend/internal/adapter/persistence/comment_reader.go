package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// commentReader は [repository.CommentReader] の GORM 実装。
// コメントの所有権・存在チェックに必要な読み取りだけを提供する。
type commentReader struct {
	db *gorm.DB
}

// NewCommentReader は CommentReader の GORM 実装を返す。
func NewCommentReader(db *gorm.DB) repository.CommentReader {
	return &commentReader{db: db}
}

var _ repository.CommentReader = (*commentReader)(nil)

// FindCommentByID は ID でコメントを取得する。
func (r *commentReader) FindCommentByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.WithContext(ctx).First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}
