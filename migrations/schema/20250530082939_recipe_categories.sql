-- +goose Up
CREATE TABLE
  IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
  );

  CREATE TABLE
  IF NOT EXISTS recipe_categories (
    recipe_id TEXT,
    category_id INTEGER,
    PRIMARY KEY (category_id, recipe_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
  );

-- +goose Down
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS recipe_categories;