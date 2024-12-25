-- name: GetUsers :many
SELECT * FROM users 
WHERE name LIKE $1
LIMIT $2 OFFSET $3;

-- name: CreateUser :one
INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id;

-- name: UpdateUser :exec
UPDATE users SET name = $1 WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUsersWithRoles :many
SELECT * FROM users LEFT JOIN users_roles ON users.id = users_roles.user_id;

-- name: GetUserWithRoles :one
SELECT * FROM users LEFT JOIN users_roles ON users.id = users_roles.user_id WHERE users.id = $1;
