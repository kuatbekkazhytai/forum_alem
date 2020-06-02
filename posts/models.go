package posts

import (
	config "../config"
	user "../users"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Post struct {
	Id           int
	Title        string
	Description  string
	Timestamp    string
	Author       int
	AuthorName   string
	Category     int
	CategoryName string
}

type SessionData struct {
	IndexUser  user.User
	LoggedIn   bool
	Categories []Category
	Posts      []Post
	Category   string
}

type IndexPageData struct {
	IndexUser  user.User
	LoggedIn   bool
	Categories []string
	Posts      []Post
}
type PostPageData struct {
	LoggedIn     bool
	UserData     user.User
	ThisPost     Post
	UserLiked    bool
	UserDisliked bool
	Likes        int
	Dislikes     int
	Comments     []Comment
}
type Category struct {
	ID   int64
	Name string
}
type Comment struct {
	ID           int64  `json:"id"`
	Text         string `json:"name"`
	Timestamp    string `json:"timestamp"`
	Author       int64  `json:"autor"`
	AuthorName   string `json:"author_name"`
	Post         int64  `json:"post"`
	Likes        string `json:"likes"`
	UserLiked    bool   `json:"userliked"`
	UserDisliked bool   `json:"userdisliked"`
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
func GetCategoryName(w http.ResponseWriter, categoryid int) string {
	categoryName := ""
	err := config.DB.QueryRow("SELECT name FROM categories WHERE id=?",
		categoryid).Scan(&categoryName)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	return categoryName
}

func getComments(w http.ResponseWriter, postID int) []Comment {

	commentRows, err := config.DB.Query("SELECT * FROM comments WHERE post_id=?", postID)

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	var comments []Comment
	var comment Comment
	for commentRows.Next() {
		err = commentRows.Scan(&comment.ID, &comment.Text, &comment.Timestamp, &comment.Author, &comment.Post)
		comments = append(comments, comment)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
	}
	return comments
}
func AddDataToComments(w http.ResponseWriter, comments []Comment) []Comment {

	for i, comment := range comments {
		err := config.DB.QueryRow("SELECT username FROM users WHERE id=?",
			comment.Author).Scan(&comments[i].AuthorName)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}

		tempTimeArray := strings.Split(comment.Timestamp, "T")
		comments[i].Timestamp = tempTimeArray[0]

	}

	return comments
}

func AllPosts() ([]Post, error) {

	rows, err := config.DB.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ps []Post // ps for posts
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Description, &post.Timestamp, &post.Author, &post.Category)
		if err != nil {
			return nil, err
		}
		ps = append(ps, post)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ps, nil
}
func AddDataToPost(w http.ResponseWriter, posts []Post) []Post {

	for i, post := range posts {
		err := config.DB.QueryRow("SELECT username FROM users WHERE id=?",
			post.Author).Scan(&posts[i].AuthorName)
		if err != nil {
			return nil
		}
		tempTimeArray := strings.Split(post.Timestamp, "T")
		posts[i].Timestamp = tempTimeArray[0]

		err = config.DB.QueryRow("SELECT name FROM categories WHERE id=?",
			post.Category).Scan(&posts[i].CategoryName)
		if err != nil {
			return nil
		}
	}

	return posts
}

func OnePost(r *http.Request) (Post, error) {
	var post Post
	id := r.FormValue("id")
	if id == "" {
		return post, errors.New("400. Bad Request.")
	}
	row := config.DB.QueryRow("SELECT * FROM posts where id = ?", id)
	err := row.Scan(&post.Id, &post.Title, &post.Description, &post.Timestamp, &post.Author, &post.Category)
	if err != nil {
		return post, err
	}
	return post, nil
}

func PutPost(r *http.Request, user user.User) (Post, error) {
	var post Post
	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	categoryid := r.FormValue("category")
	fmt.Println(post.Title)
	fmt.Println(post.Description)
	fmt.Println(user.ID)
	fmt.Println(categoryid)

	if post.Title == "" || post.Description == "" {
		return post, errors.New("400. Bad request. All fields must be complete.")
	}
	statement, err := config.DB.Prepare("INSERT INTO posts (title, description,author_id, category_id) VALUES (?, ?, ?, ?)") //, post.Title, post.Description)
	statement.Exec(post.Title, post.Description, user.ID, categoryid)
	if err != nil {
		fmt.Println(err)
		return post, errors.New("500. Internal Server Error." + err.Error())
	}

	return post, nil
}
func UpdatePost(r *http.Request) error {
	var post Post
	id := r.FormValue("id")
	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	if post.Title == "" || post.Description == "" || id == "" {
		return errors.New("400. Bad request. All fields must be complete.")
	}

	post.Id, _ = strconv.Atoi(id)
	// _, err := config.DB.Exec("UPDATe posts SET title = $2, description = $3 WHERE id = $1", post.Id, post.Title, post.Description )
	statement, err := config.DB.Prepare("UPDATE posts SET title = ?, description = ? WHERE id = ?") //, post.Id, post.Title, post.Description )
	statement.Exec(post.Title, post.Description, post.Id)
	if err != nil {
		return err

	}
	return nil
}

func DeletePost(r *http.Request) error {
	id := r.FormValue("id")
	if id == "" {
		return errors.New("400. Bad Request.")
	}
	// _, err := config.DB.Exec("DELETE FROM posts where id = $1", id)
	statement, err := config.DB.Prepare("DELETE FROM posts where id = ?") //, id)
	statement.Exec(id)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
