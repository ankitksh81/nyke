package auth

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ankitksh81/nyke/Auth/helper/pages"
	"github.com/ankitksh81/nyke/logger"
	"golang.org/x/oauth2"
)

// HandleLogin function

func HandleLogin(w http.ResponseWriter, r *http.Request, oauthConfig *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConfig.Endpoint.AuthURL)
	if err != nil {
		logger.Log.Error("Parse: " + err.Error())
	}
	logger.Log.Info(URL.String())
	parameters := url.Values{}
	parameters.Add("client_id", oauthConfig.ClientID)
	parameters.Add("scope", strings.Join(oauthConfig.Scopes, " "))
	parameters.Add("redirect_uri", oauthConfig.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)

	URL.RawQuery = parameters.Encode()
	url := URL.String()
	logger.Log.Info(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleMain function renders the index page when the index route is called

func HandleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pages.IndexPage))
}
