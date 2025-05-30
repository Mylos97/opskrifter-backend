-- +goose Up
CREATE TABLE
  IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
  );

  CREATE TABLE
  IF NOT EXISTS recipe_categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipe_id TEXT,
    category INTEGER,
    FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE
    FOREIGN KEY (category) REFERENCES categories (id) ON DELETE CASCADE
  );

-- +goose Down
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS recipe_categories;