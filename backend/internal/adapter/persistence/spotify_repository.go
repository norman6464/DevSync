package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// spotifyRepository は [repository.SpotifyRepository] の sqlc(pgx) 実装。
type spotifyRepository struct {
	q *sqlcgen.Queries
}

// NewSpotifyRepository は SpotifyRepository の sqlc(pgx) 実装を返す。
func NewSpotifyRepository(q *sqlcgen.Queries) repository.SpotifyRepository {
	return &spotifyRepository{q: q}
}

var _ repository.SpotifyRepository = (*spotifyRepository)(nil)

// DeleteUserData は指定ユーザーの最近再生した曲を削除する。
func (r *spotifyRepository) DeleteUserData(ctx context.Context, userID uint) error {
	return r.q.DeleteSpotifyRecentTracksByUser(ctx, int64(userID))
}
