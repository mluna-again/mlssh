-- +goose Up
-- +goose StatementBegin
CREATE TABLE settings (
  user_pk INTEGER,
  pet_species TEXT,
  pet_color TEXT,
  FOREIGN KEY(user_pk) REFERENCES users(public_key)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE settings;
-- +goose StatementEnd
