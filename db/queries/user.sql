-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id, name, email;

-- name: FindUserByEmail :one
SELECT id FROM users WHERE email = $1;

-- name: VerifyPassword :one
SELECT id, password FROM users WHERE email = $1;

-- name: FetchUserByID :one
SELECT id, name, email
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, email FROM users;

-- name: AddRefreshToken :execrows
UPDATE users
SET refresh_token = $1
WHERE email = $2;

-- name: VerifyRefreshToken :one
SELECT id, name, email FROM users WHERE refresh_token = $1;

-- name: RemoveRefreshToken :execrows
UPDATE users SET refresh_token = NULL WHERE refresh_token = $1;

