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
    user_id TEXT NOT NULL);

CREATE TABLE IF NOT EXISTS ingredients (
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS ingredients_for_recipe (
  recipe_id TEXT NOT NULL,
  ingredient_id TEXT NOT NULL,
  PRIMARY KEY (recipe_id, ingredient_id),
  FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
  FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS ingredients_for_recipe;
DROP TABLE IF EXISTS ingredients;
