package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"lab/db"
	githubauth "lab/githubAuth"
	googleauth "lab/googleAuth"
	"lab/session"
	yandexauth "lab/yandexAuth"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetSession(r)

	if err != nil {
		fmt.Fprint(w, `<html><body><h1>You are not logged in!</h1></body></html>`)
		return
	}

	fmt.Fprintf(w, `<html><body><h1>Welcome, %s %s!</h1><h2>You are logged in with %s</h2></body></html>`, user.Name, user.Surname, user.Method)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return
	}
	db.DeleteUserSession(cookie.Value)
	session.ClearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/google/login", googleauth.HandleLogin)
	http.HandleFunc("/google/callback", googleauth.HandleCallback)
	http.HandleFunc("/github/login", githubauth.HandleLogin)
	http.HandleFunc("/github/callback", githubauth.HandleCallback)
	http.HandleFunc("/yandex/login", yandexauth.HandleLogin)
	http.HandleFunc("/yandex/callback", yandexauth.HandleCallback)
	http.HandleFunc("/logout", handleLogout)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}