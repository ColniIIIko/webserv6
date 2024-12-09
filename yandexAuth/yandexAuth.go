package yandexauth

import (
	"context"
	"encoding/json"
	"fmt"
	"lab/db"
	"lab/session"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

var (
	clientID     = os.Getenv("YANDEX_CLIENT_ID")
	clientSecret = os.Getenv("YANDEX_CLIENT_SECRET")
	redirectURL  = os.Getenv("YANDEX_REDIRECT_URL")
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
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
		fmt.Println(err)
		http.Error(w, "Could not get token", http.StatusBadRequest)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	userInfoResponse, err := client.Get("https://login.yandex.ru/info?format=json&oauth_token=" + token.AccessToken)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not get user info", http.StatusBadRequest)
		return
	}

	defer userInfoResponse.Body.Close()

	if userInfoResponse.StatusCode == http.StatusOK {
		var user User
		if err := json.NewDecoder(userInfoResponse.Body).Decode(&user); err != nil {
			fmt.Println(err)
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
