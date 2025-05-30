-- +goose Up
CREATE TABLE
  IF NOT EXISTS users (id TEXT, name TEXT, createdAt TEXT);

-- +goose Down
DROP TABLE IF EXISTS users;