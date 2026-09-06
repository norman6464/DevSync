-- name: CreateProject :one
INSERT INTO projects (
    user_id, title, description, tech_stack, demo_url, github_url, image_url, role,
    start_date, end_date, featured, is_archived, github_repo_id, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, now(), now()
) RETURNING *;

-- name: UpdateProject :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE projects SET
    title = $2, description = $3, tech_stack = $4, demo_url = $5, github_url = $6,
    image_url = $7, role = $8, start_date = $9, end_date = $10, featured = $11,
    is_archived = $12, github_repo_id = $13, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
-- 依存するマイルストーン等はFKのON DELETE CASCADEでDBが自動的に削除する。
DELETE FROM projects WHERE id = $1;

-- name: CountProjectsByUser :one
-- 件数は FindByUserID の総数取得と CountByUserID(単体メソッド) の両方から利用する。
SELECT COUNT(*) FROM projects WHERE user_id = $1;

-- name: CountArchivedProjectsByUser :one
SELECT COUNT(*) FROM projects WHERE user_id = $1 AND is_archived = true;

-- name: CountAllProjects :one
SELECT COUNT(*) FROM projects;

-- name: CountSearchProjects :one
SELECT COUNT(*) FROM projects
WHERE (title ILIKE $1 OR description ILIKE $1 OR tech_stack ILIKE $1);

-- name: ListProjectsByUserWithRepo :many
-- GithubRepoのみPreloadする（Userは含めない。移行前のFindByUserIDと同じ挙動）。
-- github_repo_idはNULL許容のためLEFT JOIN。git_hub_repositoriesはsqlc.embedを使わず
-- 個別カラム選択にすることでJOINコンテキストのNULL許容性を正しく推論させる。
SELECT sqlc.embed(projects),
    ghr.id AS repo_id,
    ghr.user_id AS repo_user_id,
    ghr.git_hub_repo_id AS repo_git_hub_repo_id,
    ghr.name AS repo_name,
    ghr.full_name AS repo_full_name,
    ghr.description AS repo_description,
    ghr.language AS repo_language,
    ghr.stars AS repo_stars,
    ghr.forks AS repo_forks,
    ghr.is_private AS repo_is_private,
    ghr.updated_at AS repo_updated_at
FROM projects
LEFT JOIN git_hub_repositories ghr ON ghr.id = projects.github_repo_id
WHERE projects.user_id = $1
ORDER BY projects.featured DESC, projects.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListArchivedProjectsByUserWithRepo :many
-- GithubRepoのみPreloadする（Userは含めない。移行前のFindArchivedByUserIDと同じ挙動）。
SELECT sqlc.embed(projects),
    ghr.id AS repo_id,
    ghr.user_id AS repo_user_id,
    ghr.git_hub_repo_id AS repo_git_hub_repo_id,
    ghr.name AS repo_name,
    ghr.full_name AS repo_full_name,
    ghr.description AS repo_description,
    ghr.language AS repo_language,
    ghr.stars AS repo_stars,
    ghr.forks AS repo_forks,
    ghr.is_private AS repo_is_private,
    ghr.updated_at AS repo_updated_at
FROM projects
LEFT JOIN git_hub_repositories ghr ON ghr.id = projects.github_repo_id
WHERE projects.user_id = $1 AND projects.is_archived = true
ORDER BY projects.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAllProjectsWithUserAndRepo :many
-- GORMのPreload("User").Preload("GithubRepo")に相当（移行前のFindAllと同じ挙動）。
-- user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(projects), sqlc.embed(users),
    ghr.id AS repo_id,
    ghr.user_id AS repo_user_id,
    ghr.git_hub_repo_id AS repo_git_hub_repo_id,
    ghr.name AS repo_name,
    ghr.full_name AS repo_full_name,
    ghr.description AS repo_description,
    ghr.language AS repo_language,
    ghr.stars AS repo_stars,
    ghr.forks AS repo_forks,
    ghr.is_private AS repo_is_private,
    ghr.updated_at AS repo_updated_at
FROM projects
JOIN users ON users.id = projects.user_id
LEFT JOIN git_hub_repositories ghr ON ghr.id = projects.github_repo_id
ORDER BY projects.created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchProjectsWithUserAndRepo :many
-- GORMのPreload("User").Preload("GithubRepo")に相当（移行前のSearchと同じ挙動）。
SELECT sqlc.embed(projects), sqlc.embed(users),
    ghr.id AS repo_id,
    ghr.user_id AS repo_user_id,
    ghr.git_hub_repo_id AS repo_git_hub_repo_id,
    ghr.name AS repo_name,
    ghr.full_name AS repo_full_name,
    ghr.description AS repo_description,
    ghr.language AS repo_language,
    ghr.stars AS repo_stars,
    ghr.forks AS repo_forks,
    ghr.is_private AS repo_is_private,
    ghr.updated_at AS repo_updated_at
FROM projects
JOIN users ON users.id = projects.user_id
LEFT JOIN git_hub_repositories ghr ON ghr.id = projects.github_repo_id
WHERE (projects.title ILIKE $1 OR projects.description ILIKE $1 OR projects.tech_stack ILIKE $1)
ORDER BY projects.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListFeaturedProjectsByUserWithRepo :many
-- GithubRepoのみPreloadする（移行前のFindFeaturedByUserIDと同じ挙動）。
SELECT sqlc.embed(projects),
    ghr.id AS repo_id,
    ghr.user_id AS repo_user_id,
    ghr.git_hub_repo_id AS repo_git_hub_repo_id,
    ghr.name AS repo_name,
    ghr.full_name AS repo_full_name,
    ghr.description AS repo_description,
    ghr.language AS repo_language,
    ghr.stars AS repo_stars,
    ghr.forks AS repo_forks,
    ghr.is_private AS repo_is_private,
    ghr.updated_at AS repo_updated_at
FROM projects
LEFT JOIN git_hub_repositories ghr ON ghr.id = projects.github_repo_id
WHERE projects.user_id = $1 AND projects.featured = true
ORDER BY projects.created_at DESC;

-- name: SetProjectArchived :exec
UPDATE projects SET is_archived = $2 WHERE id = $1;
