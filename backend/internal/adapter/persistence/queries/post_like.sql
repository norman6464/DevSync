-- name: CreatePostLike :exec
-- likesへのINSERTとpost_metrics.like_countの加算を同一SQL文で行う（DEVSYNC-159）。
-- CTEの出力列をpost_metrics.post_idと別名にし、sqlcの列名衝突による誤検出を避ける。
WITH inserted AS (
    INSERT INTO likes (user_id, post_id, created_at) VALUES ($1, $2, now())
    RETURNING likes.post_id AS liked_post_id
)
INSERT INTO post_metrics (post_id, like_count)
SELECT liked_post_id, 1 FROM inserted
ON CONFLICT (post_id) DO UPDATE SET like_count = post_metrics.like_count + 1;

-- name: DeletePostLike :execrows
-- likesの削除とpost_metrics.like_countの減算を同一SQL文で行う。実際に削除できた
-- （＝insertedに1行ある）ときだけpost_metricsを更新するため、rowsAffectedは
-- 従来どおり「実際にいいねを取り消せたか」を表す。
WITH deleted AS (
    DELETE FROM likes WHERE likes.user_id = $1 AND likes.post_id = $2
    RETURNING likes.post_id AS liked_post_id
)
INSERT INTO post_metrics (post_id, like_count)
SELECT liked_post_id, 0 FROM deleted
ON CONFLICT (post_id) DO UPDATE SET like_count = GREATEST(post_metrics.like_count - 1, 0);

-- name: CountPostLikeByUserAndPost :one
SELECT COUNT(*) FROM likes WHERE likes.user_id = $1 AND likes.post_id = $2;
