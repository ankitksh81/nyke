package routes

import (
	auth "github.com/ankitksh81/nyke/auth"
	"github.com/ankitksh81/nyke/middleware"
	"github.com/ankitksh81/nyke/services/api"
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {

	// Google authentication routes
	// r.HandleFunc("/", middleware.Homepage)
	r.HandleFunc("/oauth", auth.HandleMain)
	r.HandleFunc("/login-gl", auth.HandleGoogleLogin)
	r.HandleFunc("/callback", auth.CallBackFromGoogle)

	// users routes
	r.HandleFunc("/api/register", api.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{id}", api.GetUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/login", middleware.Login).Methods("POST")
	r.HandleFunc("/protected", middleware.IsAuthorized(middleware.ProtectedEndpoint)).Methods("GET")

	// products routes
	r.HandleFunc("/products", api.GetProducts).Methods("GET")
	r.HandleFunc("/product/{id}", api.GetProductByID).Methods("GET")
}
