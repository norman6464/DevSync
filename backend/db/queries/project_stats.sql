-- name: GetProjectStats :one
-- projects は GORM の論理削除（deleted_at）付きモデルのため、GORMの既定スコープに合わせて
-- deleted_at IS NULL を明示する（Unscoped() されていない全クエリと同じ挙動）。
SELECT
  count(*) AS total_projects,
  count(*) FILTER (WHERE featured = true) AS featured_projects,
  count(*) FILTER (WHERE end_date IS NULL) AS ongoing_projects,
  count(*) FILTER (WHERE end_date IS NOT NULL) AS completed_projects
FROM projects
WHERE user_id = $1 AND deleted_at IS NULL;
