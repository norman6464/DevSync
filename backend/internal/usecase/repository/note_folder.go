package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteFolderRepository はノートフォルダの永続化に対する、usecase 側が要求する契約。
type NoteFolderRepository interface {
	Create(ctx context.Context, folder *model.NoteFolder) error
	// FindByID は指定 ID のフォルダを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	// usecase 側が永続化技術のエラー型（gorm.ErrRecordNotFound 等）を知らずに済むようにするための契約。
	FindByID(ctx context.Context, id uint) (*model.NoteFolder, error)
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.NoteFolder, int64, error)
	FindByParentID(ctx context.Context, parentID uint) ([]model.NoteFolder, error)
	FindRootsByUserID(ctx context.Context, userID uint) ([]model.NoteFolder, error)
	Update(ctx context.Context, folder *model.NoteFolder) error
	Delete(ctx context.Context, id uint) error
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}
