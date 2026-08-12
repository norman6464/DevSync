package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// YouTubeVideoSearcher は YouTube の動画検索を外部へ委ねるための最小の契約。
// API キーや HTTP の詳細は adapter 側に閉じ込める。
type YouTubeVideoSearcher interface {
	// SearchVideos はキーワードに一致する動画を関連度順に返す。
	SearchVideos(ctx context.Context, query string, maxResults int, language string) ([]model.YouTubeVideo, error)
}

// YouTubeVideoRepository は取得済み動画と検索キャッシュの永続化に対する契約。
type YouTubeVideoRepository interface {
	// UpsertVideos は動画を video_id で重複判定して保存する。
	UpsertVideos(ctx context.Context, videos []model.YouTubeVideo) error
	// FindByVideoIDs は指定 ID の動画をまとめて取得する。
	FindByVideoIDs(ctx context.Context, videoIDs []string) ([]model.YouTubeVideo, error)
	// FindCachedSearch は有効期限内の検索キャッシュを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindCachedSearch(ctx context.Context, query, language string) (*model.YouTubeSearchCache, error)
	// SaveSearchCache は同一クエリのキャッシュを上書き保存する。
	SaveSearchCache(ctx context.Context, cache *model.YouTubeSearchCache) error
}
