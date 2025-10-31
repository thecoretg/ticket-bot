-- name: GetWebexRoomIDByInternalID :one
SELECT webex_id FROM webex_room
WHERE id = $1;

-- name: GetWebexRoom :one
SELECT * FROM webex_room
WHERE id = $1;

-- name: GetWebexRoomByWebexID :one
SELECT * FROM webex_room
WHERE webex_id = $1;

-- name: ListWebexRooms :many
SELECT * FROM webex_room
ORDER BY id;

-- name: InsertWebexRoom :one
INSERT INTO webex_room
(webex_id, name, type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateWebexRoom :one
UPDATE webex_room
SET
    name = $2,
    type = $3,
    updated_on = NOW()
WHERE id = $1
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

