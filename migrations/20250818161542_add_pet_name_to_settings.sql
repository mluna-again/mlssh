-- +goose Up
-- +goose StatementBegin
ALTER TABLE settings ADD COLUMN pet_name TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE settings DROP COLUMN pet_name;
-- +goose StatementEnd
