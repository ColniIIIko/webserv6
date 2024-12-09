package googleauth

import (
	"context"
	"encoding/json"
	"lab/db"
	"lab/session"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "336825104577-gk6614uv2ek46b5jib6j4767smog5lr9.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	oauthStateString = "random"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
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
	userInfoResponse, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
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
			Name:    user.GivenName,
			Surname: user.FamilyName,
			Method:  "Google",
		}

		sessionID := uuid.New()
		session.SetSession(w, sessionID.String(), sessionUser)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Could not get user info", http.StatusBadRequest)
	}
}
