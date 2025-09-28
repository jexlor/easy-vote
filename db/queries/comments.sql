-- name: GetAllComments :many
SELECT * FROM comments ORDER BY created_at DESC;

-- name: CreateComment :one
INSERT INTO comments (user_id, comment) VALUES ($1, $2)
RETURNING *;

-- name: GetCommentByID :one
SELECT * FROM comments WHERE id = $1;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1 AND user_id = $2;
