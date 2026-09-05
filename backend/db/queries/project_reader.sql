-- name: GetProjectWithUserAndRepoByID :one
-- GORMのPreload("User").Preload("GithubRepo")に相当。
-- projectsはGORMの論理削除（deleted_at）付きモデルのためdeleted_at IS NULLを明示する。
-- github_repo_idはNULL許容のためLEFT JOIN。git_hub_repositoriesはsqlc.embedを使わず
-- 個別カラム選択にすることで、sqlcのJOIN文脈依存のnull推論を効かせる
-- （sqlc.embedは対象テーブル自身のスキーマ上のnull許容性をそのまま使ってしまい、
-- LEFT JOINで欠落しうる行を表現できないため）。
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
WHERE projects.id = $1 AND projects.deleted_at IS NULL;
