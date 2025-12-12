-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: NewUser :exec
INSERT INTO users (first_name, last_name, email, password) VALUES (?,?,?,?);
