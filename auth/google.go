package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/middleware"
	"github.com/ankitksh81/nyke/models"
	"github.com/ankitksh81/nyke/services/api"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfigGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	oauthStateStringGl string
)

// InitializeOAuthGoogle Function
func InitializeOAuthGoogle() {
	oauthConfigGl.ClientID = viper.GetString("google.clientID")
	oauthConfigGl.ClientSecret = viper.GetString("google.clientSecret")
	oauthStateStringGl = viper.GetString("oauthStateString")
}

// HandleGoogleLogin Function
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	HandleLogin(w, r, oauthConfigGl, oauthStateStringGl)
}

type Response struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// CallBackFromGoogle Function
const oauth_url = "https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token="

func CallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
	// logger.Log.Info("Callback-gl..")
	state := r.FormValue("state")
	// logger.Log.Info(state)
	if state != oauthStateStringGl {
		logger.Log.Info("Invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	// logger.Log.Info(code)

	if code == "" {
		logger.Log.Warn("Code not found..")
		// w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			// w.Write([]byte("User has denied Permission.."))
			errMsg := "User has denied Permission!"
			middleware.JSONError(w, errMsg, 401)
		} else {
			errMsg := "Could not get the access token."
			middleware.JSONError(w, errMsg, 400)
		}
	} else {
		token, err := oauthConfigGl.Exchange(oauth2.NoContext, code)
		if err != nil {
			logger.Log.Error("oauthConfigGl.Exchange() failed with " + err.Error() + "\n")
			errMsg := "Could not get the access token."
			middleware.JSONError(w, errMsg, 400)
			return
		}
		// logger.Log.Info("TOKEN>> AccessToken>> " + token.AccessToken)
		// logger.Log.Info("TOKEN>> Expiration Time>> " + token.Expiry.String())
		// logger.Log.Info("TOKEN>> RefreshToken>> " + token.RefreshToken)

		resp, err := http.Get(oauth_url + url.QueryEscape(token.AccessToken))
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

		// logger.Log.Info("parseResponseBody: " + string(response) + "\n")

		// unmarshal response
		var authRes models.AuthResponse
		json.Unmarshal(response, &authRes)

		// check if user already exists
		userFound, err := middleware.CheckUserExist(authRes.Email)
		if err != nil {
			middleware.Error500(w)
		}

		// if user not found, register the user
		if !userFound {
			api.CreateAuthUser(w, &authRes)
		} else {
			sqlQuery := `SELECT user_id FROM users WHERE email = $1`
			var user_id string
			err := middleware.DB.QueryRow(sqlQuery, authRes.Email).Scan(&user_id)
			if err != nil {
				middleware.Error500(w)
				return
			}
			jwtToken, err := middleware.GenerateJWT(user_id)
			if err != nil {
				middleware.Error500(w)
				return
			}

			token := &middleware.JwtToken{
				Token: jwtToken,
			}
			json.NewEncoder(w).Encode(token)
		}
		return
	}
}
