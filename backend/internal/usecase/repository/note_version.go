package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteVersionRepository はノートのバージョン履歴の永続化に対する、usecase 側が要求する契約。
type NoteVersionRepository interface {
	Create(ctx context.Context, version *model.NoteVersion) error
	FindByNoteID(ctx context.Context, noteID uint, limit, offset int) ([]model.NoteVersion, int64, error)
	// FindByID は指定 ID のバージョンを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.NoteVersion, error)
	// GetLatestVersionNumber は最新のバージョン番号を返す。バージョンが無い場合は 0。
	GetLatestVersionNumber(ctx context.Context, noteID uint) (int, error)
}

// NoteUpdater はノート本体の書き戻しだけを切り出した最小 port（-er）。
type NoteUpdater interface {
	Update(ctx context.Context, note *model.Note) error
}
