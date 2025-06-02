-- +goose Up
CREATE TABLE
  IF NOT EXISTS cookbooks (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    likes INTEGER DEFAULT 0,
    user_id TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
  );

CREATE TABLE
  IF NOT EXISTS cookbook_recipes (
    cookbook_id TEXT NOT NULL,
    recipe_id TEXT NOT NULL,
    PRIMARY KEY (cookbook_id, recipe_id),
    FOREIGN KEY (cookbook_id) REFERENCES cookbooks (id) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE
  );

-- +goose Down
DROP TABLE IF EXISTS cookbooks;
DROP TABLE IF EXISTS cookbook_recipes;