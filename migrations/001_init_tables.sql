-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY NOT NULL UNIQUE,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  created_at TEXT NOT NULL, 
  status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS recipes (
  id TEXT PRIMARY KEY NOT NULL UNIQUE,
  name TEXT NOT NULL,
  minutes INTEGER NOT NULL,
  description TEXT NOT NULL,
  likes INTEGER DEFAULT 0,
  comments INTEGER DEFAULT 0,
  views INTEGER DEFAULT 0,
  image TEXT NOT NULL,
  recipe_cuisine TEXT NOT NULL,
  created_at TEXT NOT NULL,
  user_id TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS user_liked_recipes (
  user_id TEXT NOT NULL,
  recipe_id TEXT NOT NULL, 
  PRIMARY KEY (user_id, recipe_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ingredients (
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS ingredients_for_recipe (
  recipe_id TEXT NOT NULL,
  ingredient_id TEXT NOT NULL,
  amount  TEXT NOT NULL,
  PRIMARY KEY (recipe_id, ingredient_id),
  FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
  FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_steps (
  id TEXT PRIMARY KEY NOT NULL,
  recipe_id TEXT NOT NULL,
  step TEXT NOT NULL,
  FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS user_liked_recipe;
DROP TABLE IF EXISTS ingredients_for_recipe;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS recipe_steps;
