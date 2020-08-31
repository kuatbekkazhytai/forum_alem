package config

import (
	"database/sql"
	"fmt"
	"net/http"

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
	createLikes, _ := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS likes (
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			UNIQUE(user_id, post_id),
			FOREIGN KEY(user_id) REFERENCES users(id), 
			FOREIGN KEY(post_id) REFERENCES posts(id)
		)
	`)
	createLikes.Exec()

	createDislikes, _ := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS dislikes (
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			UNIQUE(user_id, post_id),
			FOREIGN KEY(user_id) REFERENCES users(id), 
			FOREIGN KEY(post_id) REFERENCES posts(id)
		)
	`)
	createDislikes.Exec()

	createCommentLikes, _ := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS commentlikes (
			user_id INTEGER NOT NULL,
			comment_id INTEGER NOT NULL,
			UNIQUE(user_id, comment_id),
			FOREIGN KEY(user_id) REFERENCES users(id), 
			FOREIGN KEY(comment_id) REFERENCES comments(id)
		)
	`)
	createCommentLikes.Exec()

	createCommentDislikes, _ := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS commentdislikes (
			user_id INTEGER NOT NULL,
			comment_id INTEGER NOT NULL,
			UNIQUE(user_id, comment_id),
			FOREIGN KEY(user_id) REFERENCES users(id), 
			FOREIGN KEY(comment_id) REFERENCES comments(id)
		)
	`)
	createCommentDislikes.Exec()

}

func addLike(w http.ResponseWriter, postID int, userID int) {
	// Add like
	likeRows, err := DB.Query("SELECT * FROM likes WHERE post_id=? AND user_id=?", postID, userID)

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	count := 0
	for likeRows.Next() {
		count++
	}
	if count >= 1 {
		// Remove like if exists
		_, err = DB.Exec(`
			DELETE from likes WHERE user_id=? AND post_id=?
		`, userID, postID)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
		//
	} else {
		// Ad like if not exists
		_, err = DB.Exec(`
			INSERT OR IGNORE INTO likes (user_id, post_id) VALUES (?, ?)
		`, userID, postID)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
	}

	// Remove dislike
	_, err = DB.Exec(`
		DELETE from dislikes WHERE user_id=? AND post_id=?
	`, userID, postID)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

func addDislike(w http.ResponseWriter, postID int, userID int) {

	// Add dislike
	dislikeRows, err := DB.Query("SELECT * FROM dislikes WHERE post_id=? AND user_id=?", postID, userID)

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	count := 0
	for dislikeRows.Next() {
		count++
	}
	if count >= 1 {
		// Remove dislike if exists
		_, err = DB.Exec(`
			DELETE from dislikes WHERE user_id=? AND post_id=?
		`, userID, postID)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
	} else {
		// Add dislike if not exists
		_, err = DB.Exec(`
			INSERT OR IGNORE INTO dislikes (user_id, post_id) VALUES (?, ?)
		`, userID, postID)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
	}

	// Remove like
	_, err = DB.Exec(`
		DELETE from likes WHERE user_id=? AND post_id=?
	`, userID, postID)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}
