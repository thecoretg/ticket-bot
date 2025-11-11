-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS webex_user_forward (
    id SERIAL PRIMARY KEY,
    user_email TEXT NOT NULL,
    dest_room_id INT NOT NULL REFERENCES webex_room(id) ON DELETE CASCADE,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    user_keeps_copy BOOLEAN NOT NULL DEFAULT TRUE,
    created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_email, dest_room_id, start_date, end_date),
    CHECK (start_date < end_date)
);

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
