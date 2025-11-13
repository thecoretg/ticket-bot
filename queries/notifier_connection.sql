-- name: ListNotifierConnections :many
SELECT * FROM notifier_connection
ORDER BY id;

-- name: GetNotifierConnection :one
SELECT * FROM notifier_connection
WHERE id = $1 LIMIT 1;

-- name: CheckNotifierExists :one
SELECT EXISTS (
    SELECT 1
    FROM notifier_connection
    WHERE cw_board_id = $1 AND webex_room_id = $2
) AS exists;

-- name: ListNotifierConnectionsByBoard :many
SELECT * FROM notifier_connection
WHERE cw_board_id = $1
ORDER BY id;

-- name: ListNotifierConnectionsByRoom :many
SELECT * FROM notifier_connection
WHERE webex_room_id = $1
ORDER BY id;

-- name: InsertNotifierConnection :one
INSERT INTO notifier_connection(cw_board_id, webex_room_id, notify_enabled)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateNotifierConnection :one
UPDATE notifier_connection
SET
    cw_board_id = $2,
    webex_room_id = $3,
    notify_enabled = $4
WHERE id = $1
RETURNING *;

-- name: SoftDeleteNotifierConnection :exec
UPDATE notifier_connection
SET
    deleted = TRUE,
    updated_on = NOW()
WHERE id = $1;

-- name: DeleteNotifierConnection :exec
DELETE FROM notifier_connection
WHERE id = $1;
