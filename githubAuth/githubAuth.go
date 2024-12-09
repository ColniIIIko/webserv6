package githubauth

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
	"golang.org/x/oauth2/github"
)

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

var (
	clientID     = os.Getenv("GITHUB_CLIENT_ID")
	clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURL  = os.Getenv("GITHUB_REDIRECT_URL")
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	oauthStateString = "random"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(oauthStateString, oauth2.SetAuthURLParam("prompt", "login"))
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
	userInfoResponse, err := client.Get("https://api.github.com/user")
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
			Surname: "",
			Method:  "GitHub",
		}

		if sessionUser.Name == "" {
			sessionUser.Name = user.Login
		}

		sessionID := uuid.New()
		session.SetSession(w, sessionID.String(), sessionUser)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Could not get user info", http.StatusBadRequest)
	}
}
