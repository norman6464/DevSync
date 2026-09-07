-- name: CountBookmarksMadeByUser :one
SELECT count(*) FROM post_reactions
WHERE user_id = $1 AND kind = 'bookmark';

-- name: CountBookmarksReceivedByUser :one
SELECT count(*) FROM post_reactions
JOIN posts ON posts.id = post_reactions.post_id
WHERE posts.user_id = $1 AND post_reactions.kind = 'bookmark';

-- name: CountBookmarksMadeByUserSince :one
SELECT count(*) FROM post_reactions
WHERE user_id = $1 AND kind = 'bookmark' AND created_at >= $2;
