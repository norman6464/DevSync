-- name: SumContributionsInRange :one
-- 期間の合計にも日別集計（開始日=当日0時、終了日=翌日0時）にも同じクエリを使う。
SELECT COALESCE(SUM(count), 0)::bigint FROM git_hub_contributions
WHERE user_id = $1 AND contributed_on >= $2 AND contributed_on < $3;

-- name: CountPostsInRange :one
-- 期間の合計にも日別集計にも同じクエリを使う。
SELECT COUNT(*) FROM posts
WHERE user_id = $1 AND created_at >= $2 AND created_at < $3;

-- name: CountCommentsInRange :one
-- 期間の合計にも日別集計にも同じクエリを使う。
SELECT COUNT(*) FROM comments
WHERE user_id = $1 AND created_at >= $2 AND created_at < $3;

-- name: CountLikesReceivedInRange :one
SELECT COUNT(*) FROM post_reactions
JOIN posts ON post_reactions.post_id = posts.id
WHERE posts.user_id = $1 AND post_reactions.kind = 'like'
    AND post_reactions.created_at >= $2 AND post_reactions.created_at < $3;

-- name: CountCompletedGoalsInRange :one
SELECT COUNT(*) FROM learning_goals
WHERE user_id = $1 AND status = 'completed' AND completed_at >= $2 AND completed_at < $3;

-- name: CountNewFollowersInRange :one
SELECT COUNT(*) FROM follows
WHERE followee_id = $1 AND created_at >= $2 AND created_at < $3;

-- name: CountMessagesSentInRange :one
SELECT COUNT(*) FROM messages
WHERE sender_id = $1 AND created_at >= $2 AND created_at < $3;

-- name: CountMessagesReceivedInRange :one
SELECT COUNT(*) FROM messages
WHERE receiver_id = $1 AND created_at >= $2 AND created_at < $3;
