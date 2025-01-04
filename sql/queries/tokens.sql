-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW() AT TIME ZONE 'UTC',
    NOW() AT TIME ZONE 'UTC',
    $2,
    $3
);

-- name: GetUserByRefreshToken :one
SELECT expires_at, user_id FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET expires_at = NOW() AT TIME ZONE 'UTC', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE token = $1;