package config

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	var err error
	// DB, err = sql.Open("postgres", "postgres://bond:password@localhost/forum?sslmode=disable")
	DB, err = sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You are connected to DB")
	statement, _ := DB.Prepare(`CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY, 
		title TEXT, 
		description TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,			
		author_id INTEGER NOT NULL, 
		category_id INTEGER NOT NULL, 
		FOREIGN KEY(author_id) REFERENCES users(id), 
		FOREIGN KEY(category_id) REFERENCES categories(id))`)
	statement.Exec()

	// statement,_ := DB.Prepare("DROP TABLE IF EXISTS posts")
	// statement.Exec()
	statement, _ = DB.Prepare(`
	CREATE TABLE IF NOT EXISTS users( 
		id INTEGER PRIMARY KEY,
		 username TEXT, 
		 password TEXT,
		 firstname TEXT, 
		 lastname TEXT,
		 session TEXT
		 )
	`)
	statement.Exec()

	statement, _ = DB.Prepare(`
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY, 
		name TEXT, 
		UNIQUE(name)
	)
	`)
	statement.Exec()

	var categories = []string{"sport", "it", "lifestyle"}
	for _, category := range categories {
		_, err = DB.Exec(`INSERT OR IGNORE INTO categories(name) VALUES(?)`, category)

		if err != nil {
			panic(err)
		}
	}
	statement, _ = DB.Prepare(`
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY,
		text TEXT, 
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
		author_id INTEGER NOT NULL, 
		post_id INTEGER NOT NULL, 
		FOREIGN KEY(author_id) REFERENCES users(id), 
		FOREIGN KEY(post_id) REFERENCES posts(id)
	)
`)
	statement.Exec()

}
