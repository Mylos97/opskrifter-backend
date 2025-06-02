-- +goose Up
CREATE TABLE
  IF NOT EXISTS user_liked_recipe (
    recipe_id TEXT,
    user_id TEXT,
    PRIMARY KEY (recipe_id, user_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE
    FOREIGN KEY (user_id) REFERENCES id (id) ON DELETE CASCADE
  );

-- +goose Down
DROP TABLE IF EXISTS user_liked_recipe;