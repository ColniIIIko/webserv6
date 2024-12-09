package yandexauth

import (
	"context"
	"encoding/json"
	"lab/db"
	"lab/session"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "d97bcc05d418408d83a097e707a3f7a5",
		ClientSecret: "3b490899475e4461ab83256d554d993a",
		RedirectURL:  "http://localhost:8080/yandex/callback",
		Scopes:       []string{"login:info"},
		Endpoint:     yandex.Endpoint,
	}
	oauthStateString = "random"
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"first_name"`
	Surname string `json:"last_name"`
	Email   string `json:"default_email,omitempty"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != oauthStateString {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Could not get token", http.StatusBadRequest)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	userInfoResponse, err := client.Get("https://login.yandex.ru/info?format=json&oauth_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Could not get user info", http.StatusBadRequest)
		return
	}

	defer userInfoResponse.Body.Close()

	if userInfoResponse.StatusCode == http.StatusOK {
		var user User
		if err := json.NewDecoder(userInfoResponse.Body).Decode(&user); err != nil {
			http.Error(w, "Could not parse user info", http.StatusBadRequest)
			return
		}

		sessionUser := db.User{
			Name:    user.Name,
			Surname: user.Surname,
			Method:  "Yandex",
		}

		sessionID := uuid.New()
		session.SetSession(w, sessionID.String(), sessionUser)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Could not get user info", http.StatusBadRequest)
	}
}
