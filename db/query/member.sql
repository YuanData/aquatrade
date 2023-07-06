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