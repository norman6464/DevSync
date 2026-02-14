package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NoteRepository はノートデータへのアクセスを提供するリポジトリ実装。
type NoteRepository struct {
	db *gorm.DB
}

// NewNoteRepository は新しいNoteRepositoryインスタンスを生成する。
func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// Create は新しいノートをデータベースに作成する。
func (r *NoteRepository) Create(note *model.Note) error {
	return r.db.Create(note).Error
}

// FindByID は指定IDのノートを取得する。
func (r *NoteRepository) FindByID(id uint) (*model.Note, error) {
	var note model.Note
	err := r.db.Preload("User").Preload("Folder").First(&note, id).Error
	return &note, err
}

// FindByUserID は指定ユーザーのノート一覧をページネーション付きで取得する（新しい順）。
func (r *NoteRepository) FindByUserID(userID uint, page, limit int) ([]model.Note, error) {
	var notes []model.Note
	offset := (page - 1) * limit
	err := r.db.Preload("Folder").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&notes).Error
	return notes, err
}

// FindByFolderID は指定フォルダ内のノート一覧を取得する。
func (r *NoteRepository) FindByFolderID(folderID uint) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.Where("folder_id = ?", folderID).
		Order("updated_at DESC").
		Find(&notes).Error
	return notes, err
}

// Update は既存のノート情報を更新する。
func (r *NoteRepository) Update(note *model.Note) error {
	return r.db.Save(note).Error
}

// Delete は指定IDのノートを削除する。
func (r *NoteRepository) Delete(id uint) error {
	return r.db.Delete(&model.Note{}, id).Error
}

// Search はキーワードでノートを検索する（タイトルまたは本文に部分一致）。
func (r *NoteRepository) Search(userID uint, query string, limit, offset int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	searchPattern := "%" + query + "%"
	db := r.db.Model(&model.Note{}).Preload("Folder").
		Where("user_id = ? AND (title LIKE ? OR content LIKE ?)", userID, searchPattern, searchPattern).
		Order("updated_at DESC")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&notes).Error
	return notes, total, err
}

// CountByUserID は指定ユーザーのノート総数を取得する。
func (r *NoteRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Note{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ToggleFavorite はノートのお気に入り状態を切り替える。
func (r *NoteRepository) ToggleFavorite(id uint) error {
	return r.db.Model(&model.Note{}).Where("id = ?", id).
		Update("is_favorite", gorm.Expr("NOT is_favorite")).Error
}
