-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cw_board (
    id INT PRIMARY KEY,
    name TEXT NOT NULL,
    notify_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    webex_room_id TEXT,
    updated_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS cw_company (
    id INT PRIMARY KEY,
    name TEXT NOT NULL,
    updated_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS cw_contact (
    id INT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT,
    company_id INT REFERENCES cw_company(id),
    updated_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS cw_member (
    id INT PRIMARY KEY,
    identifier TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    primary_email TEXT NOT NULL,
    updated_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS cw_ticket (
    id INT PRIMARY KEY,
    summary TEXT NOT NULL,
    board_id INT REFERENCES cw_board(id) NOT NULL,
    owner_id INT REFERENCES cw_member(id),
    contact_id INT REFERENCES cw_contact(id),
    resources TEXT,
    updated_by TEXT,
    updated_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS cw_ticket_note (
    id INT PRIMARY KEY,
    ticket_id INT REFERENCES cw_ticket(id) NOT NULL,
    member_id INT REFERENCES cw_member(id),
    contact_id INT REFERENCES cw_contact(id),
    notified BOOLEAN NOT NULL DEFAULT FALSE,
    updated_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cw_ticket_note;
DROP TABLE IF EXISTS cw_ticket;
DROP TABLE IF EXISTS cw_member;
DROP TABLE IF EXISTS cw_contact;
DROP TABLE IF EXISTS cw_company;
DROP TABLE IF EXISTS cw_board;
-- +goose StatementEnd
