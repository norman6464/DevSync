-- name: CreateProjectMilestone :one
INSERT INTO project_milestones (project_id, title, description, status, due_date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: GetProjectMilestoneByID :one
SELECT * FROM project_milestones
WHERE id = $1;

-- name: ListProjectMilestonesByProject :many
SELECT * FROM project_milestones
WHERE project_id = $1
ORDER BY due_date ASC NULLS LAST, created_at ASC;

-- name: UpdateProjectMilestone :one
UPDATE project_milestones
SET title = $2,
    description = $3,
    status = $4,
    due_date = $5,
    completed_at = $6,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProjectMilestone :exec
DELETE FROM project_milestones
WHERE id = $1;
