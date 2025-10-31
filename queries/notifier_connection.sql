-- name: ListNotifierConnections :many
SELECT * FROM notifier_connection
ORDER BY cw_board_id;

-- name: ListRoomsByBoard :many
SELECT w.* FROM webex_room w
    JOIN notifier_connection nc ON nc.webex_room_id = w.id
WHERE nc.cw_board_id = $1 AND nc.notify_enabled = TRUE;

-- name: ListBoardsByRoom :many
SELECT b.* FROM cw_board b
    JOIN notifier_connection nc ON nc.cw_board_id = b.id
WHERE nc.webex_room_id = $1;

-- name: InsertNotifierConnection :one
INSERT INTO notifier_connection(cw_board_id, webex_room_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateNotifierConnection :one
UPDATE notifier_connection
SET notify_enabled = $3
WHERE cw_board_id = $1 AND webex_room_id = $2
RETURNING *;

-- name: SoftDeleteNotifierConnection :exec
UPDATE notifier_connection
SET
    deleted = TRUE,
    updated_on = NOW()
WHERE cw_board_id = $1 AND webex_room_id = $2;

-- name: DeleteNotifierConnection :exec
DELETE FROM notifier_connection
WHERE cw_board_id = $1 AND webex_room_id = $2;
