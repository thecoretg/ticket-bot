-- name: GetWebexRoom :one
SELECT * FROM webex_room
WHERE id = $1;

-- name: GetWebexRoomByWebexID :one
SELECT * FROM webex_room
WHERE webex_id = $1;

-- name: ListWebexRooms :many
SELECT * FROM webex_room
ORDER BY id;

-- name: ListByEmail :many
SELECT * FROM webex_room
WHERE email = $1;

-- name: UpsertWebexRoom :one
INSERT INTO webex_room
(webex_id, name, type, email, last_activity)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (webex_id) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    email = EXCLUDED.email,
    last_activity = EXCLUDED.last_activity,
    updated_on = NOW()
RETURNING *;

-- name: SoftDeleteWebexRoom :exec
UPDATE cw_board
SET
    deleted = TRUE,
    updated_on = NOW()
WHERE id = $1;

-- name: DeleteWebexRoom :exec
DELETE FROM webex_room
WHERE id = $1;

