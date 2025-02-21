-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW() AT TIME ZONE 'UTC',
    NOW() AT TIME ZONE 'UTC',
    $2,
    $3
)
RETURNING *;

-- name: GetUserByRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW() AT TIME ZONE 'UTC';

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW() AT TIME ZONE 'UTC', 
updated_at = NOW() AT TIME ZONE 'UTC'
WHERE token = $1
RETURNING *;