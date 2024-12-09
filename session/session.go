package session

import (
	"fmt"
	"lab/db"

	"net/http"
	"time"
)

func SetSession(w http.ResponseWriter, sessionID string, user db.User) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	user.Token = sessionID
	db.CreateUserSession(user)
}

func GetSession(r *http.Request) (*db.User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}
	user, err := db.GetUserSession(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}
	return &user, nil
}

func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
	})
}
