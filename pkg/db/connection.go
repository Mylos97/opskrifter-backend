package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS recipe_cuisines (
			id TEXT PRIMARY KEY,
			name TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT,
			name TEXT,
			createdAt TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS recipes (
			id TEXT PRIMARY KEY,
			name TEXT,
			minutes INTEGER,
			description TEXT,
			likes INTEGER,
			comments INTEGER,
			image TEXT,
			recipe_cuisine TEXT,
			user TEXT,
			FOREIGN KEY (recipe_cuisine) REFERENCES recipe_cuisines(id) ON DELETE SET NULL
			FOREIGN KEY (user) REFERENCES users(id) ON DELETE SET NULL
		);`,
		`CREATE TABLE IF NOT EXISTS ingredients (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT NOT NULL UNIQUE,
    	unit TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS ingredient_amounts (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	recipe_id TEXT NOT NULL,
    	ingredient_id INTEGER NOT NULL,
    	amount TEXT NOT NULL,
    	FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    	FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE,
    	UNIQUE(recipe_id, ingredient_id)
		);`,
		`CREATE TABLE IF NOT EXISTS recipe_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id TEXT,
			name TEXT,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS cookbooks (
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			likes INTEGER,
			user TEXT,
			FOREIGN KEY (user) REFERENCES users(id) ON DELETE SET NULL
		);`,
		`CREATE TABLE IF NOT EXISTS cookbook_recipes (
			cookbook_id TEXT,
			recipe_id TEXT,
			PRIMARY KEY (cookbook_id, recipe_id),
			FOREIGN KEY (cookbook_id) REFERENCES cookbooks(id) ON DELETE CASCADE,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			log.Fatalf("Failed to execute query: %v\nQuery: %s\n", err, query)
		}
	}
}