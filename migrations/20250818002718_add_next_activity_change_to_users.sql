-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN next_activity_change_at INTEGER NOT NULL DEFAULT 1755477161;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN next_activity_change_at;
-- +goose StatementEnd
