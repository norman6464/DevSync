package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// noteFolderRepository は [repository.NoteFolderRepository] の GORM 実装。
type noteFolderRepository struct {
	db *gorm.DB
}

// NewNoteFolderRepository は NoteFolderRepository の GORM 実装を返す。
func NewNoteFolderRepository(db *gorm.DB) repository.NoteFolderRepository {
	return &noteFolderRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteFolderRepository = (*noteFolderRepository)(nil)

// Create は新しいフォルダを作成する。
func (r *noteFolderRepository) Create(ctx context.Context, folder *model.NoteFolder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

// FindByID は指定 ID のフォルダを取得する。不在の場合は (nil, nil) を返す。
func (r *noteFolderRepository) FindByID(ctx context.Context, id uint) (*model.NoteFolder, error) {
	var folder model.NoteFolder
	err := r.db.WithContext(ctx).First(&folder, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// FindByUserID は指定ユーザーのフォルダをページネーション付きで取得する。
func (r *noteFolderRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.NoteFolder, int64, error) {
	var folders []model.NoteFolder
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.NoteFolder{}).
		Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&folders).Error
	return folders, total, err
}

// FindByParentID は指定親フォルダ配下のサブフォルダを取得する。
func (r *noteFolderRepository) FindByParentID(ctx context.Context, parentID uint) ([]model.NoteFolder, error) {
	var folders []model.NoteFolder
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).Order("created_at DESC").
		Find(&folders).Error
	return folders, err
}

// FindRootsByUserID は指定ユーザーのルートフォルダ（親なし）を取得する。
func (r *noteFolderRepository) FindRootsByUserID(ctx context.Context, userID uint) ([]model.NoteFolder, error) {
	var folders []model.NoteFolder
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND parent_id IS NULL", userID).Order("created_at DESC").
		Find(&folders).Error
	return folders, err
}

// Update はフォルダ情報を更新する。
func (r *noteFolderRepository) Update(ctx context.Context, folder *model.NoteFolder) error {
	return r.db.WithContext(ctx).Save(folder).Error
}

// Delete はフォルダを削除する。
func (r *noteFolderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.NoteFolder{}, id).Error
}

// CountByUserID は指定ユーザーのフォルダ総数を返す。
func (r *noteFolderRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.NoteFolder{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
