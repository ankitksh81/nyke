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
	// Initialize Viper across the application
	configs.InitializeViper()

	// Initialize Logger across the application
	logger.InitializeZapCustomLogger()

	// Initialize Oauth2 Services
	auth.InitializeOAuthGoogle()

	// Routes for the application
	router := mux.NewRouter()

	router.HandleFunc("/", auth.HandleMain)
	router.HandleFunc("/login-gl", auth.HandleGoogleLogin)
	router.HandleFunc("/callback", auth.CallBackFromGoogle)

	logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), router))
}
