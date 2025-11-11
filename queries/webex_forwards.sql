-- name: ListWebexUserForwards :many
SELECT sqlc.embed(wf), sqlc.embed(wr)
FROM webex_user_forward wf
JOIN webex_room AS wr ON wr.id = wf.dest_room_id
WHERE (sqlc.narg(email)::text IS NULL OR wf.email = sqlc.narg(email));

-- name: GetWebexUserForward :one
SELECT sqlc.embed(wf), sqlc.embed(wr)
FROM webex_user_forward wf
JOIN webex_room AS wr ON wr.id = wf.source_room_id
WHERE wf.id = $1;

-- name: InsertWebexUserForward :one
INSERT INTO webex_user_forward (
    user_email, dest_room_id, start_date, end_date, enabled, user_keeps_copy
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteWebexForward :exec
DELETE FROM webex_user_forward
WHERE id = $1;