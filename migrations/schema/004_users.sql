-- +goose Up
CREATE TABLE
  IF NOT EXISTS users (
  id TEXT NOT NULL UNIQUE, 
  name TEXT NOT NULL, 
  email TEXT NOT NULL UNIQUE, 
  created_at TEXT);

-- +goose Down
DROP TABLE IF EXISTS users;