package persistence

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// youTubeVideoRepository は [repository.YouTubeVideoRepository] の sqlc(pgx) 実装。
// UpsertVideos は複数件のupsertを1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type youTubeVideoRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewYouTubeVideoRepository は YouTubeVideoRepository の sqlc(pgx) 実装を返す。
func NewYouTubeVideoRepository(pool *pgxpool.Pool) repository.YouTubeVideoRepository {
	return &youTubeVideoRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.YouTubeVideoRepository = (*youTubeVideoRepository)(nil)

// toModelYouTubeVideo は sqlc の生成行を model.YouTubeVideo へ変換する。
func toModelYouTubeVideo(row sqlcgen.YouTubeVideo) model.YouTubeVideo {
	return model.YouTubeVideo{
		ID:           uint(row.ID),
		VideoID:      row.VideoID,
		Title:        row.Title,
		Description:  fromStringPtr(row.Description),
		ChannelID:    fromStringPtr(row.ChannelID),
		ChannelTitle: fromStringPtr(row.ChannelTitle),
		ThumbnailURL: fromStringPtr(row.ThumbnailUrl),
		PublishedAt:  timeValue(fromTimestamptz(row.PublishedAt)),
		CreatedAt:    timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:    timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// toModelYouTubeSearchCache は sqlc の生成行を model.YouTubeSearchCache へ変換する。
func toModelYouTubeSearchCache(row sqlcgen.YouTubeSearchCach) model.YouTubeSearchCache {
	return model.YouTubeSearchCache{
		ID:           uint(row.ID),
		Query:        row.Query,
		Language:     fromStringPtr(row.Language),
		VideoIDs:     fromStringPtr(row.VideoIds),
		CacheExpires: timeValue(fromTimestamptz(row.CacheExpires)),
		CreatedAt:    timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:    timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// UpsertVideos は動画データを一括で Upsert する。
// 移行前のGORMバッチ作成（単一INSERT文）と同じ原子性を保つため、1トランザクションで行う。
func (r *youTubeVideoRepository) UpsertVideos(ctx context.Context, videos []model.YouTubeVideo) error {
	if len(videos) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for i, video := range videos {
		row, err := q.UpsertYouTubeVideo(ctx, sqlcgen.UpsertYouTubeVideoParams{
			VideoID:      video.VideoID,
			Title:        video.Title,
			Description:  &video.Description,
			ChannelID:    &video.ChannelID,
			ChannelTitle: &video.ChannelTitle,
			ThumbnailUrl: &video.ThumbnailURL,
			PublishedAt:  toTimestamptzNotNull(video.PublishedAt),
		})
		if err != nil {
			return err
		}
		videos[i] = toModelYouTubeVideo(row)
	}
	return tx.Commit(ctx)
}

// FindByVideoIDs はビデオ ID の一覧から動画データを取得する。
func (r *youTubeVideoRepository) FindByVideoIDs(ctx context.Context, videoIDs []string) ([]model.YouTubeVideo, error) {
	if len(videoIDs) == 0 {
		return nil, nil
	}
	rows, err := r.q.ListYouTubeVideosByIDs(ctx, videoIDs)
	if err != nil {
		return nil, err
	}
	videos := make([]model.YouTubeVideo, len(rows))
	for i, row := range rows {
		videos[i] = toModelYouTubeVideo(row)
	}
	return videos, nil
}

// FindCachedSearch は有効期限内の検索キャッシュを取得する。不在の場合は (nil, nil) を返す。
func (r *youTubeVideoRepository) FindCachedSearch(ctx context.Context, query, language string) (*model.YouTubeSearchCache, error) {
	lowerQuery := strings.ToLower(query)
	row, err := r.q.GetYouTubeSearchCache(ctx, sqlcgen.GetYouTubeSearchCacheParams{
		Query:        lowerQuery,
		Language:     &language,
		CacheExpires: toTimestamptzNotNull(time.Now()),
	})
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	cache := toModelYouTubeSearchCache(row)
	return &cache, nil
}

// SaveSearchCache は検索キャッシュを保存する（同一クエリ・言語の既存分は上書きする）。
func (r *youTubeVideoRepository) SaveSearchCache(ctx context.Context, cache *model.YouTubeSearchCache) error {
	existing, err := r.q.GetYouTubeSearchCacheExact(ctx, sqlcgen.GetYouTubeSearchCacheExactParams{
		Query:    cache.Query,
		Language: &cache.Language,
	})
	if err == nil {
		row, err := r.q.UpdateYouTubeSearchCache(ctx, sqlcgen.UpdateYouTubeSearchCacheParams{
			ID:           existing.ID,
			VideoIds:     &cache.VideoIDs,
			CacheExpires: toTimestamptzNotNull(cache.CacheExpires),
		})
		if err != nil {
			return err
		}
		*cache = toModelYouTubeSearchCache(row)
		return nil
	}

	row, err := r.q.CreateYouTubeSearchCache(ctx, sqlcgen.CreateYouTubeSearchCacheParams{
		Query:        cache.Query,
		Language:     &cache.Language,
		VideoIds:     &cache.VideoIDs,
		CacheExpires: toTimestamptzNotNull(cache.CacheExpires),
	})
	if err != nil {
		return err
	}
	*cache = toModelYouTubeSearchCache(row)
	return nil
}
