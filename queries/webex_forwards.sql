-- name: ListWebexUserForwards :many
SELECT * FROM webex_user_forward
ORDER BY id;

-- name: GetWebexUserForward :one
SELECT * FROM webex_user_forward
WHERE id = $1 LIMIT 1;

-- name: ListWebexUserForwardsBySourceRoomID :many
SELECT * FROM webex_user_forward
WHERE source_room_id = $1
ORDER BY id;

-- name: InsertWebexUserForward :one
INSERT INTO webex_user_forward (
    source_room_id, dest_room_id, start_date, end_date, enabled, user_keeps_copy
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteWebexForward :exec
DELETE FROM webex_user_forward
WHERE id = $1;
