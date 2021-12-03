package routes

import (
	auth "github.com/ankitksh81/nyke/auth"
	"github.com/ankitksh81/nyke/middleware"
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {

	// Google authentication routes
	r.HandleFunc("/", auth.HandleMain)
	r.HandleFunc("/login-gl", auth.HandleGoogleLogin)
	r.HandleFunc("/callback", auth.CallBackFromGoogle)

	// Database related routes
	r.HandleFunc("/api/register", middleware.RegisterFromForm).Methods("POST")
	r.HandleFunc("/products", middleware.SendProducts).Methods("GET")
	r.HandleFunc("/api/login", middleware.Login).Methods("POST")
}
