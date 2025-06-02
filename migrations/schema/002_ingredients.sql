-- +goose Up
CREATE TABLE
  IF NOT EXISTS ingredients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
  );

CREATE TABLE
  IF NOT EXISTS ingredient_amounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipe_id TEXT NOT NULL,
    ingredient_id INTEGER NOT NULL,
    amount TEXT NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, ingredient_id)
  );

-- +goose Down
DROP TABLE IF EXISTS ingredient_amounts;
DROP TABLE IF EXISTS ingredients;
