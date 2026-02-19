package repository

import (
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// YouTubeVideoRepository はYouTube動画キャッシュのデータアクセスを提供する。
type YouTubeVideoRepository struct {
	db *gorm.DB
}

// NewYouTubeVideoRepository は新しいYouTubeVideoRepositoryインスタンスを生成する。
func NewYouTubeVideoRepository(db *gorm.DB) *YouTubeVideoRepository {
	return &YouTubeVideoRepository{db: db}
}

// UpsertVideos は動画データを一括でUpsertする。
func (r *YouTubeVideoRepository) UpsertVideos(videos []model.YouTubeVideo) error {
	if len(videos) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "video_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "description", "channel_title", "thumbnail_url", "updated_at"}),
	}).Create(&videos).Error
}

// FindByVideoIDs はビデオIDの一覧から動画データを取得する。
func (r *YouTubeVideoRepository) FindByVideoIDs(videoIDs []string) ([]model.YouTubeVideo, error) {
	var videos []model.YouTubeVideo
	err := r.db.Where("video_id IN ?", videoIDs).Find(&videos).Error
	return videos, err
}

// FindCachedSearch は有効なキャッシュを検索する。
func (r *YouTubeVideoRepository) FindCachedSearch(query, language string) (*model.YouTubeSearchCache, error) {
	var cache model.YouTubeSearchCache
	err := r.db.Where("query = ? AND language = ? AND cache_expires > ?",
		strings.ToLower(query), language, time.Now()).
		First(&cache).Error
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// SaveSearchCache は検索キャッシュを保存する（既存の同一クエリは上書き）。
func (r *YouTubeVideoRepository) SaveSearchCache(cache *model.YouTubeSearchCache) error {
	var existing model.YouTubeSearchCache
	err := r.db.Where("query = ? AND language = ?", cache.Query, cache.Language).First(&existing).Error
	if err == nil {
		existing.VideoIDs = cache.VideoIDs
		existing.CacheExpires = cache.CacheExpires
		return r.db.Save(&existing).Error
	}
	return r.db.Create(cache).Error
}
