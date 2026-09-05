package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteVersionRepository は [repository.NoteVersionRepository] の sqlc(pgx) 実装。
type noteVersionRepository struct {
	q *sqlcgen.Queries
}

// NewNoteVersionRepository は NoteVersionRepository の sqlc(pgx) 実装を返す。
func NewNoteVersionRepository(q *sqlcgen.Queries) repository.NoteVersionRepository {
	return &noteVersionRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteVersionRepository = (*noteVersionRepository)(nil)

// toModelNoteVersion は sqlc の生成行を model.NoteVersion へ変換する。
func toModelNoteVersion(row sqlcgen.NoteVersion) model.NoteVersion {
	return model.NoteVersion{
		ID:            uint(row.ID),
		NoteID:        uint(row.NoteID),
		VersionNumber: int(row.VersionNumber),
		Title:         row.Title,
		Content:       fromStringPtr(row.Content),
		Tags:          fromStringPtr(row.Tags),
		CreatedAt:     row.CreatedAt.Time,
	}
}

// Create は新しいノートバージョンを保存する。
func (r *noteVersionRepository) Create(ctx context.Context, version *model.NoteVersion) error {
	row, err := r.q.CreateNoteVersion(ctx, sqlcgen.CreateNoteVersionParams{
		NoteID:        int64(version.NoteID),
		VersionNumber: int64(version.VersionNumber),
		Title:         version.Title,
		Content:       &version.Content,
		Tags:          &version.Tags,
	})
	if err != nil {
		return err
	}
	*version = toModelNoteVersion(row)
	return nil
}

// FindByNoteID は指定ノートのバージョン履歴を新しい順に取得する。
func (r *noteVersionRepository) FindByNoteID(ctx context.Context, noteID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	total, err := r.q.CountNoteVersionsByNote(ctx, int64(noteID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListNoteVersionsByNote(ctx, sqlcgen.ListNoteVersionsByNoteParams{
		NoteID: int64(noteID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	versions := make([]model.NoteVersion, len(rows))
	for i, row := range rows {
		versions[i] = toModelNoteVersion(row)
	}
	return versions, total, nil
}

// FindByID は指定 ID のバージョンを取得する。不在の場合は (nil, nil) を返す。
func (r *noteVersionRepository) FindByID(ctx context.Context, id uint) (*model.NoteVersion, error) {
	row, err := r.q.GetNoteVersionByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	version := toModelNoteVersion(row)
	return &version, nil
}

// GetLatestVersionNumber は指定ノートの最新バージョン番号を返す。バージョンがない場合は 0 を返す。
func (r *noteVersionRepository) GetLatestVersionNumber(ctx context.Context, noteID uint) (int, error) {
	max, err := r.q.GetLatestNoteVersionNumber(ctx, int64(noteID))
	if err != nil {
		return 0, err
	}
	return int(max), nil
}
