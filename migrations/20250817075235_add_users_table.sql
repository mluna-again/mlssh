-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  public_key TEXT PRIMARY KEY NOT NULL,
  name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
