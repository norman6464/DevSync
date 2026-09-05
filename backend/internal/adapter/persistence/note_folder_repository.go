package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteFolderRepository は [repository.NoteFolderRepository] の sqlc(pgx) 実装。
type noteFolderRepository struct {
	q *sqlcgen.Queries
}

// NewNoteFolderRepository は NoteFolderRepository の sqlc(pgx) 実装を返す。
func NewNoteFolderRepository(q *sqlcgen.Queries) repository.NoteFolderRepository {
	return &noteFolderRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteFolderRepository = (*noteFolderRepository)(nil)

// toModelNoteFolder は sqlc の生成行を model.NoteFolder へ変換する。
func toModelNoteFolder(row sqlcgen.NoteFolder) model.NoteFolder {
	return model.NoteFolder{
		ID:        uint(row.ID),
		UserID:    uint(row.UserID),
		ParentID:  fromInt64PtrToUintPtr(row.ParentID),
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

// Create は新しいフォルダを作成する。
func (r *noteFolderRepository) Create(ctx context.Context, folder *model.NoteFolder) error {
	row, err := r.q.CreateNoteFolder(ctx, sqlcgen.CreateNoteFolderParams{
		UserID:   int64(folder.UserID),
		ParentID: toInt64PtrFromUintPtr(folder.ParentID),
		Name:     folder.Name,
	})
	if err != nil {
		return err
	}
	*folder = toModelNoteFolder(row)
	return nil
}

// FindByID は指定 ID のフォルダを取得する。不在の場合は (nil, nil) を返す。
func (r *noteFolderRepository) FindByID(ctx context.Context, id uint) (*model.NoteFolder, error) {
	row, err := r.q.GetNoteFolderByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	folder := toModelNoteFolder(row)
	return &folder, nil
}

// FindByUserID は指定ユーザーのフォルダをページネーション付きで取得する。
func (r *noteFolderRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.NoteFolder, int64, error) {
	total, err := r.q.CountNoteFoldersByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListNoteFoldersByUser(ctx, sqlcgen.ListNoteFoldersByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	folders := make([]model.NoteFolder, len(rows))
	for i, row := range rows {
		folders[i] = toModelNoteFolder(row)
	}
	return folders, total, nil
}

// FindByParentID は指定親フォルダ配下のサブフォルダを取得する。
func (r *noteFolderRepository) FindByParentID(ctx context.Context, parentID uint) ([]model.NoteFolder, error) {
	pid := int64(parentID)
	rows, err := r.q.ListNoteFoldersByParent(ctx, &pid)
	if err != nil {
		return nil, err
	}
	folders := make([]model.NoteFolder, len(rows))
	for i, row := range rows {
		folders[i] = toModelNoteFolder(row)
	}
	return folders, nil
}

// FindRootsByUserID は指定ユーザーのルートフォルダ（親なし）を取得する。
func (r *noteFolderRepository) FindRootsByUserID(ctx context.Context, userID uint) ([]model.NoteFolder, error) {
	rows, err := r.q.ListRootNoteFoldersByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	folders := make([]model.NoteFolder, len(rows))
	for i, row := range rows {
		folders[i] = toModelNoteFolder(row)
	}
	return folders, nil
}

// Update はフォルダ情報を更新する。
func (r *noteFolderRepository) Update(ctx context.Context, folder *model.NoteFolder) error {
	row, err := r.q.UpdateNoteFolder(ctx, sqlcgen.UpdateNoteFolderParams{
		ID:       int64(folder.ID),
		ParentID: toInt64PtrFromUintPtr(folder.ParentID),
		Name:     folder.Name,
	})
	if err != nil {
		return err
	}
	*folder = toModelNoteFolder(row)
	return nil
}

// Delete はフォルダを削除する。
func (r *noteFolderRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteNoteFolder(ctx, int64(id))
}

// CountByUserID は指定ユーザーのフォルダ総数を返す。
func (r *noteFolderRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountNoteFoldersByUser(ctx, int64(userID))
}
