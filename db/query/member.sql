-- name: CreateMember :one
INSERT INTO members (
  membername,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetMember :one
SELECT * FROM members
WHERE membername = $1 LIMIT 1;

-- name: UpdateMember :one
UPDATE members
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  email = COALESCE(sqlc.narg(email), email)
WHERE
  membername = sqlc.arg(membername)
RETURNING *;