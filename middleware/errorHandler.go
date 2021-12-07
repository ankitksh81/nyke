package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/ankitksh81/nyke/models"
)

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func Error500(w http.ResponseWriter) {
	err := &models.Error{
		Message: "500 - Something bad happend!",
	}
	JSONError(w, err, 500)
}

func ErrorWrongPassword(w http.ResponseWriter) {
	err := &models.Error{
		Message: "Incorrect password!",
	}
	JSONError(w, err, 401)
}

func ErrorEmptyInput(w http.ResponseWriter, field string) {
	err := &models.Error{
		Message: field + " cannot be empty!",
	}
	JSONError(w, err, 400)
}

func ErrorUserNotFound(w http.ResponseWriter) {
	err := &models.Error{
		Message: "User not found!",
	}
	JSONError(w, err, 404)
}

func ErrorUserExist(w http.ResponseWriter) {
	err := &models.Error{
		Message: "User already exist!",
	}
	JSONError(w, err, 409)
}
