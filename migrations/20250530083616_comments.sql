-- +goose Up
CREATE TABLE
  IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipe_id TEXT,
    user TEXT,
    comment TEXT,
    FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE
    FOREIGN KEY (user) REFERENCES user (id) ON DELETE CASCADE
  );

-- +goose Down
DROP TABLE IF EXISTS comments;