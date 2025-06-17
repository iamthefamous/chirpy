-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM chirps;

-- name: GetUser :one
SELECT * FROM chirps
WHERE id = $1;