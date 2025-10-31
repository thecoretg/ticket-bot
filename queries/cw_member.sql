-- name: GetMember :one
SELECT * FROM cw_member
WHERE id = $1 LIMIT 1;

-- name: GetMemberByIdentifier :one
SELECT * FROM cw_member
WHERE identifier = $1 LIMIT 1;

-- name: ListMembers :many
SELECT * FROM cw_member
ORDER BY id;

-- name: InsertMember :one
INSERT INTO cw_member
(id, identifier, first_name, last_name, primary_email)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateMember :one
UPDATE cw_member
SET
    identifier = $2,
    first_name = $3,
    last_name = $4,
    primary_email = $5,
    updated_on = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteMember :exec
UPDATE cw_member
SET deleted = TRUE
WHERE id = $1;

-- name: DeleteMember :exec
DELETE FROM cw_member
WHERE id = $1;
