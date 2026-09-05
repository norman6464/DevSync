-- name: UpsertYouTubeVideo :one
INSERT INTO you_tube_videos (video_id, title, description, channel_id, channel_title, thumbnail_url, published_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
ON CONFLICT (video_id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    channel_title = EXCLUDED.channel_title,
    thumbnail_url = EXCLUDED.thumbnail_url,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: ListYouTubeVideosByIDs :many
SELECT * FROM you_tube_videos WHERE video_id = ANY($1::text[]);

-- name: GetYouTubeSearchCache :one
-- FindCachedSearchはクエリを小文字化して検索する（移行前のGORM実装と同じ）。
SELECT * FROM you_tube_search_caches
WHERE query = $1 AND language = $2 AND cache_expires > $3;

-- name: GetYouTubeSearchCacheExact :one
-- SaveSearchCacheの既存確認用。移行前のGORM実装と同じく、こちらはクエリを小文字化しない。
SELECT * FROM you_tube_search_caches WHERE query = $1 AND language = $2;

-- name: CreateYouTubeSearchCache :one
INSERT INTO you_tube_search_caches (query, language, video_ids, cache_expires, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: UpdateYouTubeSearchCache :one
UPDATE you_tube_search_caches SET video_ids = $2, cache_expires = $3, updated_at = now()
WHERE id = $1
RETURNING *;
