package posts

import (
	"fmt"
	config "../config"
	users "../users"
	"database/sql"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	// fmt.Println("i got here")
	var err error
	var indexpagedata IndexPageData
	indexpagedata.LoggedIn = users.AlreadyLoggedIn(r)
	// fmt.Println(indexpagedata.LoggedIn)
	indexpagedata.Posts, err = AllPosts()
	// fmt.Println(indexpagedata.Posts)
	indexpagedata.Posts = AddDataToPost(w, indexpagedata.Posts)
	
	 // to add category name and author name
	// _, indexpagedata.IndexUser = users.GetUser(w, r)
	// fmt.Println(indexpagedata.IndexUser)
	// fmt.Println("new posts")
	fmt.Println(indexpagedata.Posts)
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
	post, err := OnePost(r)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	config.TPL.ExecuteTemplate(w, "show.html", post)
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
