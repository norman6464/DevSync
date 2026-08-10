package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// findNoteVersionOfNote は指定ノートに属するバージョンを取得する。
// 不在の場合も、他ノートのバージョンを指した場合も、同じ 404 を返す（移行前の挙動を維持している）。
func findNoteVersionOfNote(
	ctx context.Context,
	versions repository.NoteVersionRepository,
	noteID, versionID uint,
) (*model.NoteVersion, error) {
	version, err := versions.FindByID(ctx, versionID)
	if err != nil || version == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "バージョンが見つかりません", err)
	}
	if version.NoteID != noteID {
		return nil, domain.NewError(domain.ErrCodeNotFound, "バージョンが見つかりません", nil)
	}
	return version, nil
}

// ListNoteVersionsUseCase はノートのバージョン履歴一覧を取得する。
type ListNoteVersionsUseCase struct {
	versions repository.NoteVersionRepository
	notes    repository.NoteReader
}

// NewListNoteVersionsUseCase は ListNoteVersionsUseCase を生成する。
func NewListNoteVersionsUseCase(versions repository.NoteVersionRepository, notes repository.NoteReader) *ListNoteVersionsUseCase {
	return &ListNoteVersionsUseCase{versions: versions, notes: notes}
}

// Execute はノートの所有者に対してバージョン履歴と総件数を返す。
func (uc *ListNoteVersionsUseCase) Execute(ctx context.Context, noteID, userID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	if _, err := ensureOwner(ctx, uc.notes.FindByID, noteID, userID, noteOwnerOf); err != nil {
		return nil, 0, err
	}
	return uc.versions.FindByNoteID(ctx, noteID, limit, offset)
}

// GetNoteVersionUseCase は指定バージョンの詳細を取得する。
type GetNoteVersionUseCase struct {
	versions repository.NoteVersionRepository
	notes    repository.NoteReader
}

// NewGetNoteVersionUseCase は GetNoteVersionUseCase を生成する。
func NewGetNoteVersionUseCase(versions repository.NoteVersionRepository, notes repository.NoteReader) *GetNoteVersionUseCase {
	return &GetNoteVersionUseCase{versions: versions, notes: notes}
}

// Execute はノートの所有者に対してバージョンの詳細を返す。
func (uc *GetNoteVersionUseCase) Execute(ctx context.Context, noteID, versionID, userID uint) (*model.NoteVersion, error) {
	if _, err := ensureOwner(ctx, uc.notes.FindByID, noteID, userID, noteOwnerOf); err != nil {
		return nil, err
	}
	return findNoteVersionOfNote(ctx, uc.versions, noteID, versionID)
}

// RestoreNoteVersionUseCase は指定バージョンの内容でノートを復元する。
type RestoreNoteVersionUseCase struct {
	versions repository.NoteVersionRepository
	notes    repository.NoteReader
	writer   repository.NoteUpdater
}

// NewRestoreNoteVersionUseCase は RestoreNoteVersionUseCase を生成する。
func NewRestoreNoteVersionUseCase(
	versions repository.NoteVersionRepository,
	notes repository.NoteReader,
	writer repository.NoteUpdater,
) *RestoreNoteVersionUseCase {
	return &RestoreNoteVersionUseCase{versions: versions, notes: notes, writer: writer}
}

// Execute は復元前の状態を新しいバージョンとして残したうえで、ノートを指定バージョンの内容に戻す。
func (uc *RestoreNoteVersionUseCase) Execute(ctx context.Context, noteID, versionID, userID uint) (*model.Note, error) {
	note, err := ensureOwner(ctx, uc.notes.FindByID, noteID, userID, noteOwnerOf)
	if err != nil {
		return nil, err
	}

	version, err := findNoteVersionOfNote(ctx, uc.versions, noteID, versionID)
	if err != nil {
		return nil, err
	}

	// 復元によって失われる現在の内容を、先に新しいバージョンとして残す。
	latestNum, err := uc.versions.GetLatestVersionNumber(ctx, noteID)
	if err != nil {
		return nil, err
	}
	current := &model.NoteVersion{
		NoteID:        noteID,
		VersionNumber: latestNum + 1,
		Title:         note.Title,
		Content:       note.Content,
		Tags:          note.Tags,
	}
	if err := uc.versions.Create(ctx, current); err != nil {
		return nil, err
	}

	note.Title = version.Title
	note.Content = version.Content
	note.Tags = version.Tags
	if err := uc.writer.Update(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}
