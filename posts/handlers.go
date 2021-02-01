package posts

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	config "../config"
	users "../users"
)

//Index function
func Index(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	var err error
	var indexpagedata IndexPageData

	var isLoggedIn = users.AlreadyLoggedIn(r)

	var user users.User
	if isLoggedIn {
		var noErr, data = users.GetUser(w, r)
		if noErr {
			user = data
		}
	}

	indexpagedata.LoggedIn = isLoggedIn
	indexpagedata.Posts, err = AllPosts(isLoggedIn, user)

	indexpagedata.Posts = AddDataToPost(w, indexpagedata.Posts)
	if indexpagedata.LoggedIn {
		_, indexpagedata.IndexUser = users.GetUser(w, r)
	}

	// for _, post := range indexpagedata.Posts {
	// 	category := GetCategoryName(w, post.Category)
	// 	indexpagedata.Categories = append(indexpagedata.Categories, category)
	// }

	categories := getCategories(w)
	for _, j := range categories {
		indexpagedata.Categories = append(indexpagedata.Categories, j)
	}

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "posts.html", indexpagedata)
}

func Show(w http.ResponseWriter, r *http.Request) {
	var postpagedata PostPageData
	var err error
	postpagedata.LoggedIn = users.AlreadyLoggedIn(r)
	_, postpagedata.UserData = users.GetUser(w, r)
	// fmt.Println(postpagedata.UserData.UserName)
	postpagedata.ThisPost, err = OnePost(r)
	postpagedata.ThisPost = AddToPost(w, postpagedata.ThisPost)
	fmt.Println(postpagedata.UserData.UserName, postpagedata.ThisPost.AuthorName)
	postpagedata.Comments = getComments(w, postpagedata.ThisPost.Id)
	postpagedata.Comments = AddDataToComments(w, postpagedata.Comments)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	postID := postpagedata.ThisPost.Id
	likesCount := getLikes(w, postID)
	dislikesCount := getDislikes(w, postID)

	userLiked := false
	userDisliked := false
	postpagedata.Likes = likesCount
	postpagedata.Dislikes = dislikesCount
	postpagedata.UserLiked = userLiked
	postpagedata.UserDisliked = userDisliked

	_, user := users.GetUser(w, r)
	if users.AlreadyLoggedIn(r) {
		userLiked = getUserLike(w, postID, int(user.ID))
		userDisliked = getUserDislike(w, postID, int(user.ID))
	}
	config.TPL.ExecuteTemplate(w, "show.html", postpagedata)
}

func Create(w http.ResponseWriter, r *http.Request) {

	if !(users.AlreadyLoggedIn(r)) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	categories := getCategories(w)
	_, user := users.GetUser(w, r)
	var templateData SessionData
	templateData.Categories = categories
	templateData.IndexUser = user
	templateData.LoggedIn = users.AlreadyLoggedIn(r)
	config.TPL.ExecuteTemplate(w, "create.html", templateData)
}
func CreateProcess(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, user := users.GetUser(w, r)
	_, err := PutPost(r, user)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

// Update ...
func Update(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	post, err := OnePost(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "update.html", post)
}

func UpdateProcess(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	err := UpdatePost(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)

}

func DeleteProcess(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	err := DeletePost(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
func CreateCommentsProcess(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")
	parameters := strings.Split(r.URL.Path, "/")
	postString := ""

	if len(parameters) == 3 && parameters[2] != "" {
		postString = parameters[2]
	}

	id, err := strconv.Atoi(postString)
	if err != nil {
		panic(err)
	}

	_, user := users.GetUser(w, r)
	alreadyloggedin := users.AlreadyLoggedIn(r)
	if comment != "" {
		if alreadyloggedin {
			_, err := config.DB.Exec(`INSERT INTO comments(text, author_id, post_id) VALUES(?, ?, ?)`,
				comment, user.ID, id)
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			// /posts/show?id=10
			http.Redirect(w, r, "/posts/show?id="+postString, 301)
		} else {
			http.Redirect(w, r, "/login", 301)
		}
	}
}

func CreateLikesProcess(w http.ResponseWriter, r *http.Request) {
	like := r.FormValue("likeordislike")

	params := strings.Split(r.URL.Path, "/")
	// fmt.Println(params)
	postIdString := ""
	if len(params) == 3 && params[2] != "" {
		postIdString = params[2]
	} else {
		http.NotFound(w, r)
		return
	}

	postID, err := strconv.Atoi(postIdString)

	_, user := users.GetUser(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if like != "" {
		if users.AlreadyLoggedIn(r) {
			if like == "like" {
				addLike(w, postID, int(user.ID))
			} else if like == "dislike" {
				addDislike(w, postID, int(user.ID))
			}
			http.Redirect(w, r, "/post/"+postIdString, 301)

		} else {
			http.Redirect(w, r, "/login", 301)
		}
	}
}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {

	parameters := strings.Split(r.URL.Path, "/")
	param := ""
	fmt.Printf("Parameters: %v", parameters)
	if len(parameters) == 3 && parameters[2] != "" {
		param = parameters[2]
	} else {
		http.NotFound(w, r)
		return
	}

	categoryID, err := strconv.Atoi(param)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	posts, err := getCategoryPosts(w, categoryID)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	category, err := getCategoryName(w, categoryID)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Println(category)
	categories := getCategories(w)

	alreadyloggedin := users.AlreadyLoggedIn(r)

	var templateData IndexPageData
	templateData.Categories = categories
	templateData.Posts = posts
	templateData.Posts = AddDataToPost(w, templateData.Posts)
	// templateData.IndexUser = user
	templateData.LoggedIn = alreadyloggedin
	// templateData.Categories = category
	fmt.Println(templateData)
	config.TPL.ExecuteTemplate(w, "postsOfCategory.html", templateData)

}

func getCategoryPosts(w http.ResponseWriter, categoryID int) ([]Post, error) {

	var posts []Post
	var post Post
	var err error

	rows, err := config.DB.Query("SELECT * FROM posts WHERE category_id=?", categoryID)

	defer rows.Close()

	for rows.Next() {
		// var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Description, &post.Timestamp, &post.Author, &post.Category)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return posts, err
		}

		// err = config.DB.QueryRow("SELECT * FROM posts WHERE id=?", postID).Scan(&post.Id, &post.Title, &post.Description, &post.Timestamp, &post.Author)
		posts = append(posts, post)
	}
	fmt.Println(posts)

	return posts, err
}

func getCategoryName(w http.ResponseWriter, categoryID int) (string, error) {
	categoryName := ""
	var err error
	err = config.DB.QueryRow("SELECT name FROM categories WHERE id=?",
		categoryID).Scan(&categoryName)

	return categoryName, err
}

// func formatPosts(w http.ResponseWriter, posts []Post) []Post {
// 	var err error

// 	for i, post := range posts {
// 		err = config.DB.QueryRow("SELECT username FROM users WHERE id=?",
// 			post.Author).Scan(&posts[i].AuthorName)
// 		if err != nil {
// 			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 			return posts
// 		}
// 		tempTimeArray := strings.Split(post.Timestamp, "T")
// 		posts[i].Timestamp = tempTimeArray[0]

// 		tempContentArray := strings.Split(post.Description, " ")
// 		tempContentString := ""
// 		if len(tempContentArray) > 20 {
// 			tempContentArray = tempContentArray[:20]
// 		}
// 		for i, str := range tempContentArray {
// 			if i != 0 {
// 				tempContentString += " "
// 			}
// 			tempContentString += str
// 		}
// 		posts[i].Description = tempContentString

// 		var category Category
// 		var categories []Category

// 		categoriesOfPost, err := config.DB.Query("SELECT category_id FROM postcategories WHERE post_id=?", post.Id)
// 		for categoriesOfPost.Next() {
// 			err = categoriesOfPost.Scan(&category.ID)
// 			if err != nil {
// 				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 				return posts
// 			}

// 			err = config.DB.QueryRow("SELECT name FROM categories WHERE id=?",
// 				category.ID).Scan(&category.Name)
// 			if err != nil {
// 				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 				return posts
// 			}

// 			categories = append(categories, category)

// 		}

// 	}

// 	return posts
// }
