-- name: CountPostsWithFilter :one
-- タグはAND条件（全タグが付与されている投稿のみ）。空配列なら絞り込みなし
-- （cardinality(タグ配列)=0のときCOUNT(...)側も0になり両辺が一致するため自然に通る）。
SELECT COUNT(*) FROM posts
WHERE (posts.title LIKE $1 OR posts.content LIKE $1)
    AND posts.is_draft = false
    AND (sqlc.narg('date_from')::timestamptz IS NULL OR posts.created_at >= sqlc.narg('date_from'))
    AND (sqlc.narg('date_to')::timestamptz IS NULL OR posts.created_at <= sqlc.narg('date_to'))
    AND (
        SELECT COUNT(DISTINCT pt.tag) FROM post_tags pt
        WHERE pt.post_id = posts.id AND pt.tag = ANY(sqlc.arg('tags')::text[])
    ) = cardinality(sqlc.arg('tags')::text[]);

-- name: SearchPostsWithFilter :many
-- GORMのPreload("User")に相当（CodeSnippetsは別クエリで取得しGo側で結合する）。
-- ソート順はGoの動的Order()呼び出しの代わりに、sort_byごとのCASE式で切り替える
-- （sort_byはクエリ全体で単一の値のため、行ごとにNULLになったり値になったりはしない）。
-- like_count/view_countはpost_metrics側（DEVSYNC-159）。LEFT JOIN + COALESCEで0扱いにする。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
LEFT JOIN post_metrics pm ON pm.post_id = posts.id
WHERE (posts.title LIKE $1 OR posts.content LIKE $1)
    AND posts.is_draft = false
    AND (sqlc.narg('date_from')::timestamptz IS NULL OR posts.created_at >= sqlc.narg('date_from'))
    AND (sqlc.narg('date_to')::timestamptz IS NULL OR posts.created_at <= sqlc.narg('date_to'))
    AND (
        SELECT COUNT(DISTINCT pt.tag) FROM post_tags pt
        WHERE pt.post_id = posts.id AND pt.tag = ANY(sqlc.arg('tags')::text[])
    ) = cardinality(sqlc.arg('tags')::text[])
ORDER BY
    CASE WHEN sqlc.arg('sort_by') = 'popular' THEN COALESCE(pm.like_count, 0) END DESC,
    CASE WHEN sqlc.arg('sort_by') = 'views' THEN COALESCE(pm.view_count, 0) END DESC,
    posts.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
