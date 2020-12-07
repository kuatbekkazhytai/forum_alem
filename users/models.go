package users

import (
	"fmt"
	"net/http"

	config "../config"
)

var dbUsers = map[string]User{} // user ID, user
var dbSessions = map[string]string{}

type User struct {
	ID           int
	UserName     string
	Password     []byte
	First        string
	Last         string
	SessionToken string
}

func GetUsers() ([]User, error) {

	rows, err := config.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var us []User // ps for posts
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.UserName, &user.Password, &user.First, &user.Last, &user.SessionToken)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		us = append(us, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return us, nil
}

func GetUser(w http.ResponseWriter, r *http.Request) (bool, User) {
	c, err := r.Cookie("session")
	var user User
	if err != nil {
		return false, user
	}
	sessionToken := c.Value
	err = config.DB.QueryRow("SELECT id, username, firstname, lastname FROM users WHERE session=?",
		sessionToken).Scan(&user.ID, &user.UserName, &user.First, &user.Last)
	// fmt.Println(user)
	// return false, user
	// checkInternalServerError(err, w)
	if err != nil {
		panic(err)
		return false, user
	}
	fmt.Println(user)
	return true, user
}
