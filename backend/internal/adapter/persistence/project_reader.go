package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// projectReader は [repository.ProjectReader] の GORM 実装。
// 所有権チェックに必要なプロジェクト読み取りだけを提供する。
type projectReader struct {
	db *gorm.DB
}

// NewProjectReader は ProjectReader の GORM 実装を返す。
func NewProjectReader(db *gorm.DB) repository.ProjectReader {
	return &projectReader{db: db}
}

var _ repository.ProjectReader = (*projectReader)(nil)

// FindByID は ID でプロジェクトを取得する（既存 ProjectRepository.FindByID と同じ preload）。
// 不在の場合は (nil, nil) を返す。
func (r *projectReader) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	var project model.Project
	err := r.db.WithContext(ctx).Preload("User").Preload("GithubRepo").First(&project, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}
