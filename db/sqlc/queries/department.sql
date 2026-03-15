-- name: GetDepartmentByID :one
SELECT id,name FROM departments WHERE id = $1;