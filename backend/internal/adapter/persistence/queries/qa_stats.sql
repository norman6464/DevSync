-- name: CountQuestionsByUser :one
SELECT count(*) FROM questions
WHERE user_id = $1;

-- name: CountAnswersByUser :one
SELECT count(*) FROM answers
WHERE user_id = $1;

-- name: CountBestAnswersByUser :one
SELECT count(*) FROM answers
WHERE user_id = $1 AND is_best = true;

-- name: SumQuestionVotesByUser :one
SELECT COALESCE(SUM(vote_count), 0)::bigint FROM questions
WHERE user_id = $1;

-- name: SumAnswerVotesByUser :one
SELECT COALESCE(SUM(vote_count), 0)::bigint FROM answers
WHERE user_id = $1;
