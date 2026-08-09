package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// learningResourceReader は [repository.LearningResourceReader] の GORM 実装。
// レビューの存在確認に必要な読み取りだけを担う最小 adapter。
type learningResourceReader struct {
	db *gorm.DB
}

// NewLearningResourceReader は LearningResourceReader の GORM 実装を返す。
func NewLearningResourceReader(db *gorm.DB) repository.LearningResourceReader {
	return &learningResourceReader{db: db}
}

// コンパイル時に port を満たすことを保証する。
var _ repository.LearningResourceReader = (*learningResourceReader)(nil)

// FindByID は指定 ID の学習リソースをユーザー情報付きで取得する。
func (r *learningResourceReader) FindByID(ctx context.Context, id uint) (*model.LearningResource, error) {
	var resource model.LearningResource
	if err := r.db.WithContext(ctx).Preload("User").First(&resource, id).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}
