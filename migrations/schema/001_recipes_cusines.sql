-- +goose Up
CREATE TABLE
  IF NOT EXISTS recipe_cuisines (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    name TEXT NOT NULL UNIQUE);

-- +goose Down
DROP TABLE IF EXISTS recipe_cuisines;