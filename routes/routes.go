package routes

import (
	auth "github.com/ankitksh81/nyke/auth"
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {

	// Google authenticatioin routes
	r.HandleFunc("/", auth.HandleMain)
	r.HandleFunc("/login-gl", auth.HandleGoogleLogin)
	r.HandleFunc("/callback", auth.CallBackFromGoogle)
}
