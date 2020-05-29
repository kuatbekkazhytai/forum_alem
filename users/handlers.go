package users

import (
	config "../config"
	"database/sql"
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(w, req, "/posts", http.StatusSeeOther)
		return
	}
	var u User
	// process form submission
	if req.Method == http.MethodPost {
		// var u user
		// get form values
		u.UserName = req.FormValue("username")
		p := req.FormValue("password")
		u.First = req.FormValue("firstname")
		u.Last = req.FormValue("lastname")

		err1 := config.DB.QueryRow("SELECT username, password FROM users WHERE username=?",
			u.UserName).Scan(&u.UserName, &u.Password)
		err2 := config.DB.QueryRow("SELECT username, password FROM users WHERE password=?",
			u.Password).Scan(&u.UserName, &u.Password)

		if err1 == sql.ErrNoRows && err2 == sql.ErrNoRows {

			// store user in dbUsers
			bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
			u.Password = bs
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			_, err = config.DB.Exec(`INSERT INTO users(username, password, firstname, lastname) VALUES(?, ?, ?, ?)`,
				u.UserName, u.Password, u.First, u.Last)
			if err != nil {
				fmt.Println(err)
			}
			createSession(w, u.UserName)
			http.Redirect(w, req, "/posts", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, req, "/signup", http.StatusSeeOther)
			return
		}

	}
	config.TPL.ExecuteTemplate(w, "signup.html", u)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}
	// process form submission
	if r.Method == http.MethodPost {

		un := r.FormValue("username")
		password := r.FormValue("password")

		var u User
		err := config.DB.QueryRow("SELECT username, password FROM users WHERE username=?",
			un).Scan(&u.UserName, &u.Password)

		if err != nil {
			// errLogin = true
			http.Redirect(w, r, "/login", 303)
		}
		// validate password
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			// errLogin = true
			http.Redirect(w, r, "/login", 303)
		}
		// // create session
		createSession(w, un)

		http.Redirect(w, r, "/posts", 303)
	}

	config.TPL.ExecuteTemplate(w, "login.html", nil)
}

// func Logou(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("im inside logou")
// 	if alreadyLoggedIn(r) {
// 		ok, user := GetUser(w, r)
// 		fmt.Println(user)
// 		fmt.Println(ok)
// 		deleteSession(w, user.ID)
// 		http.Redirect(w, r, "/users", 301)
// 	} else {
// 		fmt.Println("im here")
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}
// 	http.Redirect(w, r, "/users", 301)
// }
func Logout(w http.ResponseWriter, r *http.Request) {
	// if !alreadyLoggedIn(r) {
	// 	fmt.Println("im here")
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }
	// ok, user := GetUser(w, r)
	// fmt.Println(ok)
	// fmt.Println(user)
	// deleteSession(w, user.ID)

	// http.Redirect(w, r, "/users", 301)
	// fmt.Println("users")
	http.Redirect(w, r, "/helloworld", http.StatusSeeOther)
}
func Logou(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/helloworld", http.StatusSeeOther)
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	us, err := GetUsers()
	// fmt.Println("im here")

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	config.TPL.ExecuteTemplate(w, "users.html", us)
}
