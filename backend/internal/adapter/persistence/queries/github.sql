-- name: UpsertGitHubContribution :one
INSERT INTO git_hub_contributions (user_id, date, count, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
ON CONFLICT (user_id, date) DO UPDATE SET
    count = EXCLUDED.count,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: ListAllGitHubContributionsByUser :many
SELECT * FROM git_hub_contributions WHERE user_id = $1 ORDER BY date ASC;

-- name: UpsertGitHubLanguageStat :one
INSERT INTO git_hub_language_stats (user_id, language, bytes, repo_count, updated_at)
VALUES ($1, $2, $3, $4, now())
ON CONFLICT (user_id, language) DO UPDATE SET
    bytes = EXCLUDED.bytes,
    repo_count = EXCLUDED.repo_count,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: ListGitHubLanguageStatsByUser :many
SELECT * FROM git_hub_language_stats WHERE user_id = $1 ORDER BY bytes DESC;

-- name: UpsertGitHubRepo :one
INSERT INTO git_hub_repositories (user_id, git_hub_repo_id, name, full_name, description, language, stars, forks, is_private, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())
ON CONFLICT (user_id, git_hub_repo_id) DO UPDATE SET
    name = EXCLUDED.name,
    full_name = EXCLUDED.full_name,
    description = EXCLUDED.description,
    language = EXCLUDED.language,
    stars = EXCLUDED.stars,
    forks = EXCLUDED.forks,
    is_private = EXCLUDED.is_private,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: ListGitHubReposByUser :many
SELECT * FROM git_hub_repositories WHERE user_id = $1 ORDER BY stars DESC;

-- GitHub連携の「解除」（アカウント削除ではない）で使う。ユーザー本体は残るため
-- FKのON DELETE CASCADEでは代替できない。GitHubRepository.DeleteUserDataから呼ばれる。
-- name: DeleteGitHubContributionsByUser :exec
DELETE FROM git_hub_contributions WHERE user_id = $1;

-- name: DeleteGitHubLanguageStatsByUser :exec
DELETE FROM git_hub_language_stats WHERE user_id = $1;

-- name: DeleteGitHubReposByUser :exec
DELETE FROM git_hub_repositories WHERE user_id = $1;
