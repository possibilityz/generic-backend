package services

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"example.com/internal/logger"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConfGh = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{"public_profile"},
		Endpoint:     github.Endpoint,
	}
	oauthStateStringGh = ""
)

/*
InitializeOAuthGithub Function
*/
func InitializeOAuthGithub() {
	oauthConfGh.ClientID = viper.GetString("github.clientID")
	oauthConfGh.ClientSecret = viper.GetString("github.clientSecret")
	oauthConfGh.RedirectURL = viper.GetString("github.redirectUri")
	oauthStateStringGh = viper.GetString("oauthStateString")
}

/*
HandleGithubLogin Function
*/
func HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	HandleLogin(w, r, oauthConfGh, oauthStateStringGh)
}

/*
CallBackFromGithub Function
*/
func CallBackFromGithub(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Callback-github..")

	state := r.FormValue("state")
	logger.Log.Info(state)
	if state != oauthStateStringGh {
		logger.Log.Info("invalid oauth state, expected " + oauthStateStringGh + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	logger.Log.Info(code)

	if code == "" {
		logger.Log.Warn("Code not found..")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		// User has denied access..
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		// Getting authentication tokens
		token, err := oauthConfGh.Exchange(oauth2.NoContext, code)
		if err != nil {
			logger.Log.Error("oauthConfGh.Exchange() failed with " + err.Error() + "\n")
			return
		}
		logger.Log.Info("TOKEN>> AccessToken>> " + token.AccessToken)
		logger.Log.Info("TOKEN>> Expiration Time>> " + token.Expiry.String())
		logger.Log.Info("TOKEN>> RefreshToken>> " + token.RefreshToken)

		// Getting user information with the given access token
		logger.Log.Info("https://api.github.com/user?access_token=" + url.QueryEscape(token.AccessToken))
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
		req.Header.Set("Authorization", "token "+url.QueryEscape(token.AccessToken))
		resp, err := client.Do(req)

		if err != nil {
			logger.Log.Error("Get: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Log.Error("ReadAll: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		logger.Log.Info("parseResponseBody: " + string(response) + "\n")

		w.Write([]byte("Hello, I'm protected\n"))
		w.Write([]byte(string(response)))
		return
	}
}
