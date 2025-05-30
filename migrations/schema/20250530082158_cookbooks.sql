-- +goose Up
CREATE TABLE
  IF NOT EXISTS cookbooks (
    id TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    likes INTEGER,
    user TEXT,
    FOREIGN KEY (user) REFERENCES users (id) ON DELETE SET NULL
  );

CREATE TABLE
  IF NOT EXISTS cookbook_recipes (
    cookbook_id TEXT,
    recipe_id TEXT,
    PRIMARY KEY (cookbook_id, recipe_id),
    FOREIGN KEY (cookbook_id) REFERENCES cookbooks (id) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE
  );

-- +goose Down
DROP TABLE IF EXISTS cookbooks;
DROP TABLE IF EXISTS cookbook_recipes;