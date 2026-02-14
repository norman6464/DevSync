package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NoteFolderRepository はノートフォルダデータへのアクセスを提供するリポジトリ実装。
type NoteFolderRepository struct {
	db *gorm.DB
}

// NewNoteFolderRepository は新しいNoteFolderRepositoryインスタンスを生成する。
func NewNoteFolderRepository(db *gorm.DB) *NoteFolderRepository {
	return &NoteFolderRepository{db: db}
}

// Create は新しいフォルダを作成する。
func (r *NoteFolderRepository) Create(folder *model.NoteFolder) error {
	return r.db.Create(folder).Error
}

// FindByID は指定IDのフォルダを取得する。
func (r *NoteFolderRepository) FindByID(id uint) (*model.NoteFolder, error) {
	var folder model.NoteFolder
	err := r.db.First(&folder, id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// FindByUserID は指定ユーザーの全フォルダを取得する。
func (r *NoteFolderRepository) FindByUserID(userID uint) ([]model.NoteFolder, error) {
	var folders []model.NoteFolder
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&folders).Error
	return folders, err
}

// FindByParentID は指定親フォルダ配下のサブフォルダを取得する。
func (r *NoteFolderRepository) FindByParentID(parentID uint) ([]model.NoteFolder, error) {
	var folders []model.NoteFolder
	err := r.db.Where("parent_id = ?", parentID).Order("created_at DESC").Find(&folders).Error
	return folders, err
}

// GetRootFolders は指定ユーザーのルートフォルダ（親なし）を取得する。
func (r *NoteFolderRepository) GetRootFolders(userID uint) ([]model.NoteFolder, error) {
	var folders []model.NoteFolder
	err := r.db.Where("user_id = ? AND parent_id IS NULL", userID).Order("created_at DESC").Find(&folders).Error
	return folders, err
}

// Update はフォルダ情報を更新する。
func (r *NoteFolderRepository) Update(folder *model.NoteFolder) error {
	return r.db.Save(folder).Error
}

// Delete はフォルダを削除する。
func (r *NoteFolderRepository) Delete(id uint) error {
	return r.db.Delete(&model.NoteFolder{}, id).Error
}
