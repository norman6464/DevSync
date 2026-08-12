package persistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// youTubeVideoRepository は [repository.YouTubeVideoRepository] の GORM 実装。
type youTubeVideoRepository struct {
	db *gorm.DB
}

// NewYouTubeVideoRepository は YouTubeVideoRepository の GORM 実装を返す。
func NewYouTubeVideoRepository(db *gorm.DB) repository.YouTubeVideoRepository {
	return &youTubeVideoRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.YouTubeVideoRepository = (*youTubeVideoRepository)(nil)

// UpsertVideos は動画データを一括で Upsert する。
func (r *youTubeVideoRepository) UpsertVideos(ctx context.Context, videos []model.YouTubeVideo) error {
	if len(videos) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "video_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"title", "description", "channel_title", "thumbnail_url", "updated_at",
		}),
	}).Create(&videos).Error
}

// FindByVideoIDs はビデオ ID の一覧から動画データを取得する。
func (r *youTubeVideoRepository) FindByVideoIDs(ctx context.Context, videoIDs []string) ([]model.YouTubeVideo, error) {
	var videos []model.YouTubeVideo
	err := r.db.WithContext(ctx).Where("video_id IN ?", videoIDs).Find(&videos).Error
	return videos, err
}

// FindCachedSearch は有効期限内の検索キャッシュを取得する。不在の場合は (nil, nil) を返す。
func (r *youTubeVideoRepository) FindCachedSearch(ctx context.Context, query, language string) (*model.YouTubeSearchCache, error) {
	var cache model.YouTubeSearchCache
	err := r.db.WithContext(ctx).
		Where("query = ? AND language = ? AND cache_expires > ?", strings.ToLower(query), language, time.Now()).
		First(&cache).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// SaveSearchCache は検索キャッシュを保存する（同一クエリ・言語の既存分は上書きする）。
func (r *youTubeVideoRepository) SaveSearchCache(ctx context.Context, cache *model.YouTubeSearchCache) error {
	db := r.db.WithContext(ctx)

	var existing model.YouTubeSearchCache
	err := db.Where("query = ? AND language = ?", cache.Query, cache.Language).First(&existing).Error
	if err == nil {
		existing.VideoIDs = cache.VideoIDs
		existing.CacheExpires = cache.CacheExpires
		return db.Save(&existing).Error
	}
	return db.Create(cache).Error
}
