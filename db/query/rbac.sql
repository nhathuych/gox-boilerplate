-- name: ListPermissionsByRoleID :many
SELECT
  p.name
FROM permissions p
  INNER JOIN role_permissions rp ON rp.permission_id = p.id
WHERE
  rp.role_id = $1
ORDER BY
  p.name;
