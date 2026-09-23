-- name: CreateProperty :one
INSERT INTO properties (user_id, name, photo_id)
VALUES ($1, $2, $3)
RETURNING id, user_id, name, photo_id;

-- name: DeleteProperty :exec
DELETE FROM properties WHERE user_id = $1;