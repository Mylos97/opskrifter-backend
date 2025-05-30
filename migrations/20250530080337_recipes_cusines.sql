-- +goose Up
CREATE TABLE
  IF NOT EXISTS recipe_cuisines (id TEXT PRIMARY KEY, name TEXT);

-- +goose Down
DROP TABLE IF EXISTS recipe_cuisines;