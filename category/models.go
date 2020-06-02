package category

import (
	config "../config"
	"net/http"
)

type SessionData struct {
	IndexUser  User
	LoggedIn   bool
	Categories []Category
	Posts      []Post
	Category   string
}

type Category struct {
	ID   int64
	Name string
}

func getCategories(w http.ResponseWriter) []Category {
	rows, err := config.DB.Query("SELECT * FROM categories")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	var categories []Category
	var category Category
	for rows.Next() {
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
		categories = append(categories, category)
	}
	return categories
}
