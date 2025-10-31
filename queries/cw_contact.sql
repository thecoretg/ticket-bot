-- name: GetContact :one
SELECT * FROM cw_contact
WHERE id = $1 LIMIT 1;

-- name: ListContacts :many
SELECT * FROM cw_contact
ORDER BY id;

-- name: InsertContact :one
INSERT INTO cw_contact
(id, first_name, last_name, company_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateContact :one
UPDATE cw_contact
SET
    first_name = $2,
    last_name = $3,
    company_id = $4,
    updated_on = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteContact :exec
UPDATE cw_contact
SET deleted = TRUE
WHERE id = $1;

-- name: DeleteContact :exec
DELETE FROM cw_contact
WHERE id = $1;
