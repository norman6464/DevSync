-- name: CreateQuestion :one
INSERT INTO questions (
    user_id, title, body, tags, vote_count, answer_count, is_solved, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, now(), now()
) RETURNING *;

-- name: GetQuestionWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(questions), sqlc.embed(users)
FROM questions
JOIN users ON users.id = questions.user_id
WHERE questions.id = $1;

-- name: UpdateQuestion :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE questions SET
    title = $2, body = $3, tags = $4, vote_count = $5, answer_count = $6, is_solved = $7,
    updated_at = now()
WHERE questions.id = $1
RETURNING *;

-- name: DeleteQuestion :exec
-- 依存する回答・投票・ブックマーク等はFKのON DELETE CASCADEでDBが自動的に削除する。
DELETE FROM questions WHERE questions.id = $1;

-- name: ListQuestionsWithUser :many
-- GORMのPreload("User")に相当。tag_patternは呼び出し側で "%\"" + escapeLikeChars(tag) + "\"%"
-- として組み立て済みのものを渡す（元のGo実装のエスケープ・引用符付与をそのまま踏襲するため）。
-- sortが"unanswered"のときだけanswer_count=0で絞り込み、"votes"のときだけ投票数降順にする。
SELECT sqlc.embed(questions), sqlc.embed(users)
FROM questions
JOIN users ON users.id = questions.user_id
WHERE (sqlc.narg('tag_pattern')::text IS NULL OR questions.tags ILIKE sqlc.narg('tag_pattern'))
    AND (sqlc.arg('sort')::text != 'unanswered' OR questions.answer_count = 0)
ORDER BY
    CASE WHEN sqlc.arg('sort')::text = 'votes' THEN questions.vote_count END DESC,
    questions.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountQuestions :one
SELECT COUNT(*) FROM questions
WHERE (sqlc.narg('tag_pattern')::text IS NULL OR questions.tags ILIKE sqlc.narg('tag_pattern'))
    AND (sqlc.arg('sort')::text != 'unanswered' OR questions.answer_count = 0);

-- name: SearchQuestionsWithUser :many
SELECT sqlc.embed(questions), sqlc.embed(users)
FROM questions
JOIN users ON users.id = questions.user_id
WHERE (questions.title ILIKE $1 OR questions.body ILIKE $1 OR questions.tags ILIKE $1)
ORDER BY questions.vote_count DESC, questions.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchQuestions :one
SELECT COUNT(*) FROM questions
WHERE (questions.title ILIKE $1 OR questions.body ILIKE $1 OR questions.tags ILIKE $1);

-- name: ListQuestionsByUser :many
-- Userは含めない（移行前からの挙動。一覧系の中でこれだけPreloadしない）。
SELECT * FROM questions
WHERE questions.user_id = $1
ORDER BY questions.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListSolvedQuestionsWithUser :many
SELECT sqlc.embed(questions), sqlc.embed(users)
FROM questions
JOIN users ON users.id = questions.user_id
WHERE questions.is_solved = true
ORDER BY questions.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSolvedQuestions :one
SELECT COUNT(*) FROM questions WHERE questions.is_solved = true;

-- name: ListUnansweredQuestionsWithUser :many
SELECT sqlc.embed(questions), sqlc.embed(users)
FROM questions
JOIN users ON users.id = questions.user_id
WHERE questions.answer_count = 0
ORDER BY questions.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUnansweredQuestions :one
SELECT COUNT(*) FROM questions WHERE questions.answer_count = 0;

-- name: ListBookmarkedQuestionsWithUser :many
SELECT sqlc.embed(questions), sqlc.embed(users)
FROM questions
JOIN users ON users.id = questions.user_id
WHERE questions.id IN (
    SELECT qb.question_id FROM question_bookmarks qb WHERE qb.user_id = $1
)
ORDER BY questions.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountBookmarkedQuestions :one
SELECT COUNT(*) FROM questions
WHERE questions.id IN (
    SELECT qb.question_id FROM question_bookmarks qb WHERE qb.user_id = $1
);

-- name: GetQuestionVoteByUserAndQuestion :one
SELECT * FROM question_votes
WHERE question_votes.user_id = $1 AND question_votes.question_id = $2;

-- name: CreateQuestionVote :one
INSERT INTO question_votes (user_id, question_id, value, created_at)
VALUES ($1, $2, $3, now())
RETURNING *;

-- name: UpdateQuestionVoteValue :exec
UPDATE question_votes SET value = $3
WHERE question_votes.user_id = $1 AND question_votes.question_id = $2;

-- name: DeleteQuestionVote :exec
DELETE FROM question_votes
WHERE question_votes.user_id = $1 AND question_votes.question_id = $2;

-- name: AdjustQuestionVoteCount :exec
UPDATE questions SET vote_count = vote_count + sqlc.arg('diff')::bigint
WHERE questions.id = sqlc.arg('id');

-- name: CreateQuestionBookmark :exec
INSERT INTO question_bookmarks (user_id, question_id, created_at) VALUES ($1, $2, now());

-- name: DeleteQuestionBookmark :exec
DELETE FROM question_bookmarks
WHERE question_bookmarks.user_id = $1 AND question_bookmarks.question_id = $2;

-- name: CountQuestionBookmark :one
SELECT COUNT(*) FROM question_bookmarks
WHERE question_bookmarks.user_id = $1 AND question_bookmarks.question_id = $2;
