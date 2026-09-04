-- name: CreateWeeklyChallenge :one
INSERT INTO weekly_challenges (
  user_id, year, week, challenge_type, description, target_value,
  current_value, is_completed, completed_at, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now())
RETURNING *;

-- name: GetWeeklyChallengeByUserAndWeek :one
SELECT * FROM weekly_challenges
WHERE user_id = $1 AND year = $2 AND week = $3
ORDER BY id ASC
LIMIT 1;

-- name: UpdateWeeklyChallenge :one
UPDATE weekly_challenges
SET user_id = $2,
    year = $3,
    week = $4,
    challenge_type = $5,
    description = $6,
    target_value = $7,
    current_value = $8,
    is_completed = $9,
    completed_at = $10,
    updated_at = now()
WHERE id = $1
RETURNING *;
