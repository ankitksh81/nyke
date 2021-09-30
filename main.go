package main

import (
	"log"
	"net/http"

	auth "github.com/ankitksh81/nyke/Auth"
	configs "github.com/ankitksh81/nyke/config"
	"github.com/ankitksh81/nyke/logger"
	"github.com/gorilla/mux"

	"github.com/spf13/viper"
)

func main() {

	configs.InitializeViper()          // Initialize Viper across the application
	logger.InitializeZapCustomLogger() // Initialize Logger across the application
	auth.InitializeOAuthGoogle()       // Initialize Oauth2 Services

	router := mux.NewRouter()

	router.HandleFunc("/", auth.HandleMain)
	router.HandleFunc("/login-gl", auth.HandleGoogleLogin)
	router.HandleFunc("/callback", auth.CallBackFromGoogle)

	logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), router))
}
