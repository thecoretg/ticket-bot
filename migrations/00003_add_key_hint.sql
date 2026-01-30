-- +goose Up
-- +goose StatementBegin
ALTER TABLE api_key ADD COLUMN key_hint TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE api_key DROP COLUMN key_hint;
-- +goose StatementEnd
