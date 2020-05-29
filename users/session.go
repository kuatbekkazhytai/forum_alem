package users

import (
	config "../config"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

func AlreadyLoggedIn(req *http.Request) bool {

	c, err := req.Cookie("session")
	var user User
	if err != nil {
		//User is not logged in
		return false
	}
	sessionToken := c.Value
	err = config.DB.QueryRow("SELECT id, username, firstname, lastname FROM users WHERE session=?",
		sessionToken).Scan(&user.ID, &user.UserName, &user.First, &user.Last)

	// checkInternalServerError(err, w)
	if err != nil {
		return false
	}
	return true
}

func createSession(w http.ResponseWriter, username string) {
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name:    "session",
		Value:   sID.String(),
		Expires: time.Now().Add(2 * 60 * 60 * time.Second),
	}
	http.SetCookie(w, c)
	addSession, err := config.DB.Prepare(`
	UPDATE users SET session=?
	WHERE username=?
	`)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	_, err = addSession.Exec(sID, username)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

func deleteSession(w http.ResponseWriter, userID int) {

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	})
	fmt.Println("cookie deleted")
	deleteSession, err := config.DB.Prepare(`
		UPDATE users SET session=?
		WHERE id=?
	`)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	_, err = deleteSession.Exec("logout", userID)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}
