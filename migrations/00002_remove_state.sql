-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS app_state;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app_state (
    id INT PRIMARY KEY DEFAULT 1,
    syncing_tickets BOOLEAN NOT NULL DEFAULT false,
    syncing_webex_rooms BOOLEAN NOT NULL DEFAULT false
);
-- +goose StatementEnd
