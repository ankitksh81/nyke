package main

import (
	"log"
	"net/http"

	auth "github.com/ankitksh81/nyke/auth"
	configs "github.com/ankitksh81/nyke/config"
	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/middleware"
	"github.com/ankitksh81/nyke/routes"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/spf13/viper"
)

func main() {

	configs.InitializeViper()          // Initialize Viper across the application
	logger.InitializeZapCustomLogger() // Initialize Logger across the application
	auth.InitializeOAuthGoogle()       // Initialize Oauth2 Services

	middleware.CreateConnection() // Create db connection

	router := mux.NewRouter()

	// Initialize CORS for the application
	c := cors.New(cors.Options{
		AllowedHeaders:   []string{"X-Requested-With"},
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	routes.SetupRoutes(router)

	logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), handler))
}
