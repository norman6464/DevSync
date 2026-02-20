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
// アーカイブ済みノートは除外される。
func (r *NoteRepository) FindByUserID(userID uint, page, limit int) ([]model.Note, error) {
	var notes []model.Note
	offset := (page - 1) * limit
	err := r.db.Preload("Folder").
		Where("user_id = ? AND is_archived = ?", userID, false).
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
// アーカイブ済みノートは除外される。
func (r *NoteRepository) Search(userID uint, query string, limit, offset int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	searchPattern := EscapeLikePattern(query)
	db := r.db.Model(&model.Note{}).Preload("Folder").
		Where("user_id = ? AND is_archived = ? AND (title LIKE ? OR content LIKE ?)", userID, false, searchPattern, searchPattern).
		Order("updated_at DESC")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&notes).Error
	return notes, total, err
}

// CountByUserID は指定ユーザーのノート総数を取得する（アーカイブ済みを除く）。
func (r *NoteRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Note{}).Where("user_id = ? AND is_archived = ?", userID, false).Count(&count).Error
	return count, err
}

// ToggleFavorite はノートのお気に入り状態を切り替える。
func (r *NoteRepository) ToggleFavorite(id uint) error {
	return r.db.Model(&model.Note{}).Where("id = ?", id).
		Update("is_favorite", gorm.Expr("NOT is_favorite")).Error
}

// FindFavorites は指定ユーザーのお気に入りノート一覧をページネーション付きで取得する。
func (r *NoteRepository) FindFavorites(userID uint, page, limit int) ([]model.Note, error) {
	var notes []model.Note
	offset := (page - 1) * limit
	err := r.db.Preload("Folder").
		Where("user_id = ? AND is_favorite = ?", userID, true).
		Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&notes).Error
	return notes, err
}

// CountFavoritesByUserID は指定ユーザーのお気に入りノート総数を取得する。
func (r *NoteRepository) CountFavoritesByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Note{}).Where("user_id = ? AND is_favorite = ?", userID, true).Count(&count).Error
	return count, err
}

// Archive は指定IDのノートをアーカイブする。
func (r *NoteRepository) Archive(id uint) error {
	return r.db.Model(&model.Note{}).Where("id = ?", id).
		Update("is_archived", true).Error
}

// Unarchive は指定IDのノートのアーカイブを解除する。
func (r *NoteRepository) Unarchive(id uint) error {
	return r.db.Model(&model.Note{}).Where("id = ?", id).
		Update("is_archived", false).Error
}

// FindArchived は指定ユーザーのアーカイブ済みノート一覧をページネーション付きで取得する。
func (r *NoteRepository) FindArchived(userID uint, page, limit int) ([]model.Note, error) {
	var notes []model.Note
	offset := (page - 1) * limit
	err := r.db.Preload("Folder").
		Where("user_id = ? AND is_archived = ?", userID, true).
		Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&notes).Error
	return notes, err
}

// CountArchivedByUserID は指定ユーザーのアーカイブ済みノート総数を取得する。
func (r *NoteRepository) CountArchivedByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Note{}).Where("user_id = ? AND is_archived = ?", userID, true).Count(&count).Error
	return count, err
}
