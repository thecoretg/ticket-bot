-- name: GetTicket :one
SELECT * FROM cw_ticket
WHERE id = $1 LIMIT 1;

-- name: ListTickets :many
SELECT * FROM cw_ticket
ORDER BY id;


-- name: InsertTicket :one
INSERT INTO cw_ticket
(id, summary, board_id, owner_id, company_id, contact_id, resources, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateTicket :one
UPDATE cw_ticket
SET
    summary = $2,
    board_id = $3,
    owner_id = $4,
    company_id = $5,
    contact_id = $6,
    resources = $7,
    updated_by = $8,
    updated_on = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteTicket :exec
UPDATE cw_ticket
SET
    deleted = TRUE,
    updated_on = NOW()
WHERE id = $1;

-- name: DeleteTicket :exec
DELETE FROM cw_ticket
WHERE id = $1;

-- name: GetTicketNote :one
SELECT * FROM cw_ticket_note
WHERE id = $1 LIMIT 1;

-- name: ListAllTicketNotes :many
SELECT * FROM cw_ticket_note
ORDER BY id;

-- name: ListTicketNotesByTicket :many
SELECT * FROM cw_ticket_note
WHERE ticket_id = $1
ORDER BY id;

-- name: InsertTicketNote :one
INSERT INTO cw_ticket_note
(id, ticket_id, member_id, contact_id, notified, skipped_notify)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateTicketNote :one
UPDATE cw_ticket_note
SET
    ticket_id = $2,
    member_id = $3,
    contact_id = $4,
    notified = $5,
    skipped_notify = $6,
    updated_on = NOW()
WHERE id = $1
RETURNING *;

-- name: SetNoteNotified :one
UPDATE cw_ticket_note
SET
    notified = $2
WHERE id = $1
RETURNING *;

-- name: SetNoteSkippedNotify :one
UPDATE cw_ticket_note
SET
    skipped_notify = $2
WHERE id = $1
RETURNING *;

-- name: SoftDeleteTicketNote :exec
UPDATE cw_ticket_note
SET deleted = TRUE
WHERE id = $1;

-- name: DeleteTicketNote :exec
DELETE FROM cw_ticket_note
WHERE id = $1;
