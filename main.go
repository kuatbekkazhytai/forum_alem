package main

import (
	"net/http"

	posts "../Forum_alem/posts"
	users "../Forum_alem/users"
)

func main() {
	http.HandleFunc("/", index)

	http.HandleFunc("/signup", users.Signup)
	http.HandleFunc("/logout", users.Logou)
	http.HandleFunc("/users", users.AllUsers)
	http.HandleFunc("/login", users.Login)
	http.HandleFunc("/posts", posts.Index)
	http.HandleFunc("/posts/create", posts.Create)
	http.HandleFunc("/posts/show", posts.Show)
	http.HandleFunc("/comments/", posts.CreateCommentsProcess)
	http.HandleFunc("/like/", posts.CreateLikesProcess)
	http.HandleFunc("/posts/create/process", posts.CreateProcess)
	http.HandleFunc("/posts/update", posts.Update)
	http.HandleFunc("/posts/update/process", posts.UpdateProcess)
	http.HandleFunc("/posts/delete/process", posts.DeleteProcess)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8081", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
