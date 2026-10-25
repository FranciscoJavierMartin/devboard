-- name: CreateTask :one
INSERT INTO tasks (
  project_id,
  title,
  description,
  status,
  created_by
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1 LIMIT 1;

-- name: ListTasksByProject :many
SELECT * FROM tasks where project_id = $1;

-- name: UpdateTask :one
UPDATE tasks
SET
  title = COALESCE($2, title),
  description = COALESCE($3, description),
  status = COALESCE($4, status),
  assignee_id = $5,
  updated_at = NOW()
WHERE id = $1
RETURNING *;