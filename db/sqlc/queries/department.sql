-- name: GetDepartmentByID :one
SELECT department_id,department_name FROM departments WHERE department_id = $1;