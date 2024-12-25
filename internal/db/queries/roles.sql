-- name: GetRole :one
SELECT * FROM roles WHERE id = $1;

-- name: GetUserRoles :many
SELECT roles.id as id, roles.name as name, user_id, permissions.id as permission_id, permissions.name as permission_name FROM roles
INNER JOIN users_roles ON roles.id = users_roles.role_id 
INNER JOIN roles_permissions ON roles.id = roles_permissions.role_id
INNER JOIN permissions ON roles_permissions.permission_id = permissions.id
WHERE users_roles.user_id = $1;

