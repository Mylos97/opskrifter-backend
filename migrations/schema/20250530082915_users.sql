-- +goose Up
CREATE TABLE
  IF NOT EXISTS users (id TEXT, name TEXT, created_at TEXT);

-- +goose Down
DROP TABLE IF EXISTS users;