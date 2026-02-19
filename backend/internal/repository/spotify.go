package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// SpotifyRepository はSpotify連携データへのアクセスを提供するリポジトリ実装。
type SpotifyRepository struct {
	db *gorm.DB
}

// NewSpotifyRepository は新しいSpotifyRepositoryインスタンスを生成する。
func NewSpotifyRepository(db *gorm.DB) *SpotifyRepository {
	return &SpotifyRepository{db: db}
}

// DeleteUserData は指定ユーザーのSpotify関連データを削除する。
func (r *SpotifyRepository) DeleteUserData(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.SpotifyRecentTrack{}).Error
}
