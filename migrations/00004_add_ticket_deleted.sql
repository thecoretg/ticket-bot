-- +goose Up
-- +goose StatementBegin
ALTER TABLE tickets
    ADD COLUMN DELETED BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tickets
    DROP COLUMN DELETED;
-- +goose StatementEnd
