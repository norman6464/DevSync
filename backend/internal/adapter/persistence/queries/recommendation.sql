-- name: ListFolloweeIDs :many
-- 指定ユーザーがフォローしているユーザーのIDを返す。
SELECT followee_id FROM follows WHERE follower_id = $1;

-- name: GetRecommendedUserCandidates :many
-- スキルの部分一致で候補を絞り込む。パターンは呼び出し側で "%" + escapeLikeChars(skill) + "%" として
-- 組み立て済みのものを配列で渡す（元のGo実装のエスケープ・ワイルドカード付与をそのまま踏襲するため）。
-- 自分自身とフォロー済みのユーザーは候補から除く。
SELECT * FROM users
WHERE NOT (id = ANY(sqlc.arg('exclude_ids')::bigint[]))
    AND EXISTS (
        SELECT 1 FROM unnest(sqlc.arg('skill_patterns')::text[]) AS pattern
        WHERE skills_languages LIKE pattern OR skills_frameworks LIKE pattern
    );

-- name: ListTrendingPosts :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
-- CodeSnippetsは別途 ListCodeSnippetsByPostIDs（post_bookmark.sql）で取得しGo側で付与する。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
WHERE posts.created_at > NOW() - INTERVAL '1 day' * sqlc.arg('days')::int
ORDER BY (posts.like_count + posts.comment_count) DESC
LIMIT sqlc.arg('limit');

-- name: ListTrendingResources :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
-- learning_resourcesは論理削除があるため、削除済みは除外する。
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.is_public = true
    AND learning_resources.created_at > NOW() - INTERVAL '1 day' * sqlc.arg('days')::int
    AND learning_resources.deleted_at IS NULL
ORDER BY (learning_resources.like_count + learning_resources.save_count) DESC
LIMIT sqlc.arg('limit');
