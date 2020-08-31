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

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	var err error
	var indexpagedata IndexPageData
	indexpagedata.LoggedIn = users.AlreadyLoggedIn(r)
	indexpagedata.Posts, err = AllPosts()
	indexpagedata.Posts = AddDataToPost(w, indexpagedata.Posts)

	_, indexpagedata.IndexUser = users.GetUser(w, r)
	for _, post := range indexpagedata.Posts {
		category := GetCategoryName(w, post.Category)
		indexpagedata.Categories = append(indexpagedata.Categories, category)
	}
	fmt.Println()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	config.TPL.ExecuteTemplate(w, "posts.html", indexpagedata)
}
func Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	var postpagedata PostPageData
	var err error
	postpagedata.LoggedIn = users.AlreadyLoggedIn(r)
	_, postpagedata.UserData = users.GetUser(w, r)
	postpagedata.ThisPost, err = OnePost(r)
	postpagedata.ThisPost = AddToPost(w, postpagedata.ThisPost)
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

	fmt.Println(like)
}
