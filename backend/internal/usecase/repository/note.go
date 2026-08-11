package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteRepository は学習ノートの永続化に対する、usecase 側が要求する契約。
//
// note_link 用の [NoteReader] と note_version 用の [NoteUpdater] は
// それぞれ別スライスが所有する最小 port なので、こちらへは統合していない。
type NoteRepository interface {
	Create(ctx context.Context, note *model.Note) error
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id uint) error
	// FindByID は指定 ID のノートを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.Note, error)

	// FindByUserID はアーカイブ済みを除いたノートを更新日の新しい順で返す。
	FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Note, error)
	FindByFolderID(ctx context.Context, folderID, userID uint) ([]model.Note, error)
	// Search はアーカイブ済みを除き、タイトルまたは本文への部分一致で検索する。
	Search(ctx context.Context, userID uint, query string, limit, offset int) ([]model.Note, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)

	// ToggleFavorite はお気に入り状態を反転させる。
	ToggleFavorite(ctx context.Context, id uint) error
	FindFavorites(ctx context.Context, userID uint, page, limit int) ([]model.Note, error)
	CountFavoritesByUserID(ctx context.Context, userID uint) (int64, error)

	Archive(ctx context.Context, id uint) error
	Unarchive(ctx context.Context, id uint) error
	FindArchived(ctx context.Context, userID uint, page, limit int) ([]model.Note, error)
	CountArchivedByUserID(ctx context.Context, userID uint) (int64, error)
}
