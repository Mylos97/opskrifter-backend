-- +goose Up
CREATE TABLE
  IF NOT EXISTS recipe_categories (
    recipe_id TEXT,
    category TEXT,
    PRIMARY KEY (recipe_id, category),
    FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE
  );

-- +goose Down
DROP TABLE IF EXISTS recipe_categories;