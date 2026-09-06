package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteLinkRepository は [repository.NoteLinkRepository] の sqlc(pgx) 実装。
type noteLinkRepository struct {
	q *sqlcgen.Queries
}

// NewNoteLinkRepository は NoteLinkRepository の sqlc(pgx) 実装を返す。
func NewNoteLinkRepository(q *sqlcgen.Queries) repository.NoteLinkRepository {
	return &noteLinkRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteLinkRepository = (*noteLinkRepository)(nil)

// toModelNote は sqlc の生成行を model.Note へ変換する（GORMの浅いPreloadと同じく関連は含まない）。
func toModelNote(row sqlcgen.Note) model.Note {
	return model.Note{
		ID:         uint(row.ID),
		UserID:     uint(row.UserID),
		FolderID:   fromInt64PtrToUintPtr(row.FolderID),
		Title:      row.Title,
		Content:    fromStringPtr(row.Content),
		Tags:       fromStringPtr(row.Tags),
		IsFavorite: row.IsFavorite,
		IsArchived: row.IsArchived,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
}

// Create は新しいリンクを作成する。
func (r *noteLinkRepository) Create(ctx context.Context, link *model.NoteLink) error {
	return r.q.CreateNoteLink(ctx, sqlcgen.CreateNoteLinkParams{
		SourceNoteID: int64(link.SourceNoteID),
		TargetNoteID: int64(link.TargetNoteID),
	})
}

// FindBySourceNoteID は指定ノートからのリンク一覧を取得する（リンク先を Preload する）。
func (r *noteLinkRepository) FindBySourceNoteID(ctx context.Context, sourceNoteID uint) ([]model.NoteLink, error) {
	rows, err := r.q.ListNoteLinksBySource(ctx, int64(sourceNoteID))
	if err != nil {
		return nil, err
	}
	links := make([]model.NoteLink, len(rows))
	for i, row := range rows {
		targetNote := toModelNote(row.Note)
		links[i] = model.NoteLink{
			ID:           uint(row.NoteLink.ID),
			SourceNoteID: uint(row.NoteLink.SourceNoteID),
			TargetNoteID: uint(row.NoteLink.TargetNoteID),
			TargetNote:   &targetNote,
			CreatedAt:    row.NoteLink.CreatedAt.Time,
		}
	}
	return links, nil
}

// FindByTargetNoteID は指定ノートへのリンク一覧（バックリンク）を取得する（リンク元を Preload する）。
func (r *noteLinkRepository) FindByTargetNoteID(ctx context.Context, targetNoteID uint) ([]model.NoteLink, error) {
	rows, err := r.q.ListNoteLinksByTarget(ctx, int64(targetNoteID))
	if err != nil {
		return nil, err
	}
	links := make([]model.NoteLink, len(rows))
	for i, row := range rows {
		sourceNote := toModelNote(row.Note)
		links[i] = model.NoteLink{
			ID:           uint(row.NoteLink.ID),
			SourceNoteID: uint(row.NoteLink.SourceNoteID),
			TargetNoteID: uint(row.NoteLink.TargetNoteID),
			SourceNote:   &sourceNote,
			CreatedAt:    row.NoteLink.CreatedAt.Time,
		}
	}
	return links, nil
}

// Delete は指定のリンクを削除する。
func (r *noteLinkRepository) Delete(ctx context.Context, sourceNoteID, targetNoteID uint) error {
	return r.q.DeleteNoteLink(ctx, sqlcgen.DeleteNoteLinkParams{
		SourceNoteID: int64(sourceNoteID),
		TargetNoteID: int64(targetNoteID),
	})
}

// Exists は指定のリンクが既に存在するかを返す。
func (r *noteLinkRepository) Exists(ctx context.Context, sourceNoteID, targetNoteID uint) (bool, error) {
	count, err := r.q.CountNoteLinksBetween(ctx, sqlcgen.CountNoteLinksBetweenParams{
		SourceNoteID: int64(sourceNoteID),
		TargetNoteID: int64(targetNoteID),
	})
	return count > 0, err
}

// CountBySourceNoteID は指定ノートからのリンク数を返す。
func (r *noteLinkRepository) CountBySourceNoteID(ctx context.Context, noteID uint) (int64, error) {
	return r.q.CountNoteLinksBySource(ctx, int64(noteID))
}

// CountByTargetNoteID は指定ノートへのリンク数（バックリンク数）を返す。
func (r *noteLinkRepository) CountByTargetNoteID(ctx context.Context, noteID uint) (int64, error) {
	return r.q.CountNoteLinksByTarget(ctx, int64(noteID))
}
