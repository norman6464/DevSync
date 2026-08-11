package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// noteRepository は [repository.NoteRepository] の GORM 実装。
type noteRepository struct {
	db *gorm.DB
}

// NewNoteRepository は NoteRepository の GORM 実装を返す。
func NewNoteRepository(db *gorm.DB) repository.NoteRepository {
	return &noteRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteRepository = (*noteRepository)(nil)

// Create は新しいノートを作成する。
func (r *noteRepository) Create(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

// Update は既存のノートを更新する。
func (r *noteRepository) Update(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

// Delete はノートを削除する。
func (r *noteRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Note{}, id).Error
}

// FindByID は指定 ID のノートをユーザー・フォルダ付きで取得する。不在の場合は (nil, nil) を返す。
func (r *noteRepository) FindByID(ctx context.Context, id uint) (*model.Note, error) {
	var note model.Note
	err := r.db.WithContext(ctx).Preload("User").Preload("Folder").First(&note, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// FindByUserID はアーカイブ済みを除いたノートを更新日の新しい順で取得する。
func (r *noteRepository) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	return r.pagedNotes(ctx, "user_id = ? AND is_archived = ?", []interface{}{userID, false}, page, limit)
}

// FindFavorites はお気に入りのノートを更新日の新しい順で取得する。
func (r *noteRepository) FindFavorites(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	return r.pagedNotes(ctx, "user_id = ? AND is_favorite = ?", []interface{}{userID, true}, page, limit)
}

// FindArchived はアーカイブ済みのノートを更新日の新しい順で取得する。
func (r *noteRepository) FindArchived(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	return r.pagedNotes(ctx, "user_id = ? AND is_archived = ?", []interface{}{userID, true}, page, limit)
}

// pagedNotes は絞り込み条件を受け取り、フォルダ付きでページを取得する共通処理。
func (r *noteRepository) pagedNotes(ctx context.Context, where string, args []interface{}, page, limit int) ([]model.Note, error) {
	var notes []model.Note
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Preload("Folder").
		Where(where, args...).
		Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&notes).Error
	return notes, err
}

// FindByFolderID は指定フォルダ内のノートを更新日の新しい順で取得する。
func (r *noteRepository) FindByFolderID(ctx context.Context, folderID, userID uint) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.WithContext(ctx).
		Where("folder_id = ? AND user_id = ?", folderID, userID).
		Order("updated_at DESC").
		Find(&notes).Error
	return notes, err
}

// Search はアーカイブ済みを除き、タイトルまたは本文への部分一致で検索する。
func (r *noteRepository) Search(ctx context.Context, userID uint, query string, limit, offset int) ([]model.Note, int64, error) {
	pattern := escapeLikePattern(query)
	scope := r.db.WithContext(ctx).Model(&model.Note{}).
		Where("user_id = ? AND is_archived = ? AND (title LIKE ? OR content LIKE ?)",
			userID, false, pattern, pattern)

	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var notes []model.Note
	err := scope.Session(&gorm.Session{}).
		Preload("Folder").
		Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&notes).Error
	return notes, total, err
}

// CountByUserID はアーカイブ済みを除いたノート総数を返す。
func (r *noteRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.countNotes(ctx, "user_id = ? AND is_archived = ?", userID, false)
}

// CountFavoritesByUserID はお気に入りのノート総数を返す。
func (r *noteRepository) CountFavoritesByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.countNotes(ctx, "user_id = ? AND is_favorite = ?", userID, true)
}

// CountArchivedByUserID はアーカイブ済みのノート総数を返す。
func (r *noteRepository) CountArchivedByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.countNotes(ctx, "user_id = ? AND is_archived = ?", userID, true)
}

// countNotes は絞り込み条件に一致するノート数を返す共通処理。
func (r *noteRepository) countNotes(ctx context.Context, where string, args ...interface{}) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Note{}).Where(where, args...).Count(&count).Error
	return count, err
}

// ToggleFavorite はお気に入り状態を反転させる。
func (r *noteRepository) ToggleFavorite(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Note{}).Where("id = ?", id).
		Update("is_favorite", gorm.Expr("NOT is_favorite")).Error
}

// Archive はノートをアーカイブする。
func (r *noteRepository) Archive(ctx context.Context, id uint) error {
	return r.setArchived(ctx, id, true)
}

// Unarchive はノートのアーカイブを解除する。
func (r *noteRepository) Unarchive(ctx context.Context, id uint) error {
	return r.setArchived(ctx, id, false)
}

// setArchived はアーカイブ状態を更新する共通処理。
func (r *noteRepository) setArchived(ctx context.Context, id uint, archived bool) error {
	return r.db.WithContext(ctx).Model(&model.Note{}).Where("id = ?", id).
		Update("is_archived", archived).Error
}
