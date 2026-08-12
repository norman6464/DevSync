package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// spotifyRepository は [repository.SpotifyRepository] の GORM 実装。
type spotifyRepository struct {
	db *gorm.DB
}

// NewSpotifyRepository は SpotifyRepository の GORM 実装を返す。
func NewSpotifyRepository(db *gorm.DB) repository.SpotifyRepository {
	return &spotifyRepository{db: db}
}

var _ repository.SpotifyRepository = (*spotifyRepository)(nil)

// DeleteUserData は指定ユーザーの最近再生した曲を削除する。
func (r *spotifyRepository) DeleteUserData(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.SpotifyRecentTrack{}).Error
}
