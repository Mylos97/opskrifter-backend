-- +goose Up
CREATE TABLE
  IF NOT EXISTS recipes (
    id TEXT PRIMARY KEY,
    name TEXT,
    minutes INTEGER,
    description TEXT,
    likes INTEGER,
    comments INTEGER,
    image TEXT,
    recipe_cuisine TEXT,
    user TEXT,
    FOREIGN KEY (recipe_cuisine) REFERENCES recipe_cuisines (id) ON DELETE SET NULL 
    FOREIGN KEY (user) REFERENCES users (id) ON DELETE SET NULL
  );

-- +goose Down
DROP TABLE IF EXISTS recipes;