-- name: CountBookmarksMadeByUser :one
SELECT count(*) FROM bookmarks
WHERE user_id = $1;

-- name: CountBookmarksReceivedByUser :one
SELECT count(*) FROM bookmarks
JOIN posts ON posts.id = bookmarks.post_id
WHERE posts.user_id = $1;

-- name: CountBookmarksMadeByUserSince :one
SELECT count(*) FROM bookmarks
WHERE user_id = $1 AND created_at >= $2;
