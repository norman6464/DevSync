-- name: CountQuestionsByUser :one
-- questions/answers は GORM の論理削除（deleted_at）付きモデルのため、GORMの既定スコープに
-- 合わせて deleted_at IS NULL を明示する（Unscoped() されていない全クエリと同じ挙動）。
SELECT count(*) FROM questions
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: CountAnswersByUser :one
SELECT count(*) FROM answers
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: CountBestAnswersByUser :one
SELECT count(*) FROM answers
WHERE user_id = $1 AND is_best = true AND deleted_at IS NULL;

-- name: SumQuestionVotesByUser :one
SELECT COALESCE(SUM(vote_count), 0)::bigint FROM questions
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: SumAnswerVotesByUser :one
SELECT COALESCE(SUM(vote_count), 0)::bigint FROM answers
WHERE user_id = $1 AND deleted_at IS NULL;
