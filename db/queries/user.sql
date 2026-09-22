-- name: CreateUser :one
INSERT INTO "user" (name, email, password)
VALUES ($1, $2, $3)
RETURNING id, name, email;

-- name: FindUserByEmail :one
SELECT id FROM public."user" WHERE email = $1;

-- name: VerifyPassword :one
SELECT id, password FROM public."user" WHERE email = $1;

-- name: FetchUserByID :one
SELECT id, name, email
FROM public."user"
WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, email FROM public."user";

-- name: AddRefreshToken :execrows
UPDATE "user"
SET refresh_token = $1
WHERE email = $2;

-- name: VerifyRefreshToken :one
SELECT id, name, email FROM public."user" WHERE refresh_token = $1;

-- name: RemoveRefreshToken :execrows
UPDATE "user" SET refresh_token = NULL WHERE refresh_token = $1;
