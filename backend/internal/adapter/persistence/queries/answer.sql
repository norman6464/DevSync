-- name: CreateAnswer :one
INSERT INTO answers (user_id, question_id, body, vote_count, is_best, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: IncrementQuestionAnswerCount :exec
UPDATE questions SET answer_count = answer_count + 1 WHERE id = $1;

-- name: DecrementQuestionAnswerCountFloored :exec
-- 0未満にはしない（GORMのGREATEST(answer_count - 1, 0)に相当）。
UPDATE questions SET answer_count = GREATEST(answer_count - 1, 0) WHERE id = $1;

-- name: GetAnswerWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(answers), sqlc.embed(users)
FROM answers
JOIN users ON users.id = answers.user_id
WHERE answers.id = $1;

-- name: UpdateAnswer :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE answers SET body = $2, vote_count = $3, is_best = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: LockQuestionForAnswerChange :one
-- Delete/SetBestAnswerでのロック順序（質問→回答）をGORM実装と揃えるための行ロック
-- （GORMの clause.Locking{Strength: "UPDATE"} に相当）。質問が存在しなければ
-- pgx.ErrNoRows を返し、呼び出し側のトランザクションを失敗させる。
SELECT id FROM questions WHERE id = $1 FOR UPDATE;

-- name: DeleteAnswer :exec
-- 依存する投票等はFKのON DELETE CASCADEでDBが自動的に削除する。
DELETE FROM answers WHERE id = $1;

-- name: ListAnswersByQuestionID :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(answers), sqlc.embed(users)
FROM answers
JOIN users ON users.id = answers.user_id
WHERE answers.question_id = $1
ORDER BY answers.is_best DESC, answers.vote_count DESC, answers.created_at ASC;

-- name: ListAnswersByVoteRange :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(answers), sqlc.embed(users)
FROM answers
JOIN users ON users.id = answers.user_id
WHERE answers.question_id = $1 AND answers.vote_count >= $2 AND answers.vote_count <= $3
ORDER BY answers.vote_count DESC, answers.created_at ASC;

-- name: ClearBestAnswer :exec
UPDATE answers SET is_best = false
WHERE question_id = $1 AND is_best = true;

-- name: SetAnswerBest :exec
UPDATE answers SET is_best = true WHERE id = $1;

-- name: SetQuestionSolved :exec
UPDATE questions SET is_solved = true WHERE id = $1;

-- name: LockAnswerForVoteChange :one
-- Vote/RemoveVoteでの差分計算が並行実行で古い値を読まないようにするための行ロック
-- （GORMの clause.Locking{Strength: "UPDATE"} に相当）。回答が存在しなければ
-- pgx.ErrNoRows を返し、呼び出し側のトランザクションを失敗させる。
SELECT id FROM answers WHERE id = $1 FOR UPDATE;

-- name: GetAnswerVoteByUserAndAnswer :one
SELECT * FROM answer_votes WHERE user_id = $1 AND answer_id = $2;

-- name: CreateAnswerVote :one
INSERT INTO answer_votes (user_id, answer_id, value, created_at)
VALUES ($1, $2, $3, now())
RETURNING *;

-- name: UpdateAnswerVoteValue :exec
UPDATE answer_votes SET value = $3 WHERE user_id = $1 AND answer_id = $2;

-- name: DeleteAnswerVote :exec
DELETE FROM answer_votes WHERE user_id = $1 AND answer_id = $2;

-- name: AdjustAnswerVoteCount :exec
UPDATE answers SET vote_count = vote_count + sqlc.arg('diff')::bigint WHERE id = sqlc.arg('id');
