-- +goose Up
CREATE TABLE
  IF NOT EXISTS recipes (
    id TEXT PRIMARY KEY NOT NULL UNIQUE,
    name TEXT NOT NULL,
    minutes INTEGER NOT NULL,
    description TEXT NOT NULL,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    image TEXT NOT NULL,
    recipe_cuisine TEXT NOT NULL,
    user_id TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
  );

-- +goose Down
DROP TABLE IF EXISTS recipes;