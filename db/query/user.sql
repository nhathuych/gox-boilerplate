-- name: GetUserByID :one
SELECT
  u.id,
  u.email,
  u.password_hash,
  u.role_id,
  u.created_at,
  u.updated_at,
  r.name AS role_name
FROM users u
  INNER JOIN roles r ON r.id = u.role_id
WHERE
  u.id = $1;

-- name: GetUserByEmail :one
SELECT
  u.id,
  u.email,
  u.password_hash,
  u.role_id,
  u.created_at,
  u.updated_at,
  r.name AS role_name
FROM users u
  INNER JOIN roles r ON r.id = u.role_id
WHERE
  u.email = $1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, role_id)
VALUES ($1, $2, $3)
RETURNING
  id,
  email,
  password_hash,
  role_id,
  created_at,
  updated_at;
