-- name: CreatePostLike :exec
-- post_reactions(kind='like')へのINSERTとpost_metrics.like_countの加算を同一SQL文で
-- 行う（DEVSYNC-159）。CTEの出力列をpost_metrics.post_idと別名にし、sqlcの列名衝突による
-- 誤検出を避ける。DEVSYNC-160でlikes単独テーブルからpost_reactionsへ統合。
WITH inserted AS (
    INSERT INTO post_reactions (user_id, post_id, kind, created_at) VALUES ($1, $2, 'like', now())
    RETURNING post_reactions.post_id AS liked_post_id
)
INSERT INTO post_metrics (post_id, like_count)
SELECT liked_post_id, 1 FROM inserted
ON CONFLICT (post_id) DO UPDATE SET like_count = post_metrics.like_count + 1;

-- name: DeletePostLike :execrows
-- post_reactions(kind='like')の削除とpost_metrics.like_countの減算を同一SQL文で行う。
-- 実際に削除できた（＝insertedに1行ある）ときだけpost_metricsを更新するため、
-- rowsAffectedは従来どおり「実際にいいねを取り消せたか」を表す。
WITH deleted AS (
    DELETE FROM post_reactions
    WHERE post_reactions.user_id = $1 AND post_reactions.post_id = $2 AND post_reactions.kind = 'like'
    RETURNING post_reactions.post_id AS liked_post_id
)
INSERT INTO post_metrics (post_id, like_count)
SELECT liked_post_id, 0 FROM deleted
ON CONFLICT (post_id) DO UPDATE SET like_count = GREATEST(post_metrics.like_count - 1, 0);

-- name: CountPostLikeByUserAndPost :one
SELECT COUNT(*) FROM post_reactions
WHERE post_reactions.user_id = $1 AND post_reactions.post_id = $2 AND post_reactions.kind = 'like';
