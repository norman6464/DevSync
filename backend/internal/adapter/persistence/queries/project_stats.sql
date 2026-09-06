-- name: GetProjectStats :one
SELECT
  count(*) AS total_projects,
  count(*) FILTER (WHERE featured = true) AS featured_projects,
  count(*) FILTER (WHERE end_date IS NULL) AS ongoing_projects,
  count(*) FILTER (WHERE end_date IS NOT NULL) AS completed_projects
FROM projects
WHERE user_id = $1;
