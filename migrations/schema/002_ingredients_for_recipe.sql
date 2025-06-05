-- +goose Up
CREATE TABLE IF NOT EXISTS ingredients_for_recipe (
  recipe_id TEXT NOT NULL,
  ingredient TEXT NOT NULL,
  PRIMARY KEY (recipe_id, ingredient)
  FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS ingredients_for_recipe;
