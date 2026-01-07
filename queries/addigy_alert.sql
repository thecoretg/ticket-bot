-- name: GetAddigyAlertConfig :one
SELECT * FROM addigy_alert_config
WHERE id = 1;

-- name: UpsertAddigyAlertConfig :one
INSERT INTO addigy_alert_config (
    id,
    cw_board_id,
    unattended_status_id,
    acknowledged_status_id,
    mute_1_day_status_id,
    mute_5_day_status_id,
    mute_10_day_status_id,
    mute_30_day_status_id
) 
VALUES (1, $1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE SET
    cw_board_id = EXCLUDED.cw_board_id,
    unattended_status_id = EXCLUDED.unattended_status_id,
    acknowledged_status_id = EXCLUDED.acknowledged_status_id,
    mute_1_day_status_id = EXCLUDED.mute_1_day_status_id,
    mute_5_day_status_id = EXCLUDED.mute_5_day_status_id,
    mute_10_day_status_id = EXCLUDED.mute_10_day_status_id,
    mute_30_day_status_id = EXCLUDED.mute_30_day_status_id,
    updated_on = NOW()
RETURNING *;

-- name: DeleteAddigyAlertConfig :exec
DELETE FROM addigy_alert_config
WHERE id = 1;

-- name: GetAddigyAlert :one
SELECT * FROM addigy_alert
WHERE id = $1 LIMIT 1;

-- name: ListAddigyAlerts :many
SELECT * FROM addigy_alert
ORDER BY added_on DESC;

-- name: ListAddigyAlertsByStatus :many
SELECT * FROM addigy_alert
WHERE status = $1
ORDER BY added_on DESC;

-- name: ListUnresolvedAddigyAlerts :many
SELECT * FROM addigy_alert
WHERE resolved_on IS NULL
ORDER BY added_on DESC;

-- name: ListAddigyAlertsByTicket :many
SELECT * FROM addigy_alert
WHERE ticket_id = $1
ORDER BY added_on DESC;

-- name: CreateAddigyAlert :one
INSERT INTO addigy_alert (
    id,
    ticket_id,
    level,
    category,
    name,
    fact_name,
    fact_identifier,
    fact_type,
    selector,
    status,
    value,
    muted,
    remediation,
    resolved_by_email,
    resolved_on,
    acknowledged_on,
    added_on
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
RETURNING *;

-- name: UpdateAddigyAlert :one
UPDATE addigy_alert
SET
    ticket_id = $2,
    level = $3,
    category = $4,
    name = $5,
    fact_name = $6,
    fact_identifier = $7,
    fact_type = $8,
    selector = $9,
    status = $10,
    value = $11,
    muted = $12,
    remediation = $13,
    resolved_by_email = $14,
    resolved_on = $15,
    acknowledged_on = $16
WHERE id = $1
RETURNING *;

-- name: UpdateAddigyAlertTicket :exec
UPDATE addigy_alert
SET ticket_id = $2
WHERE id = $1;

-- name: UpdateAddigyAlertStatus :exec
UPDATE addigy_alert
SET status = $2
WHERE id = $1;

-- name: AcknowledgeAddigyAlert :exec
UPDATE addigy_alert
SET acknowledged_on = $2
WHERE id = $1;

-- name: ResolveAddigyAlert :exec
UPDATE addigy_alert
SET
    resolved_on = $2,
    resolved_by_email = $3
WHERE id = $1;

-- name: DeleteAddigyAlert :exec
DELETE FROM addigy_alert
WHERE id = $1;
