-- +goose Up
-- +goose StatementBegin
ALTER TABLE settings ADD COLUMN inserted_at INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE settings DROP COLUMN inserted_at;
-- +goose StatementEnd
