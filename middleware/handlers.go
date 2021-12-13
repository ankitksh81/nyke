package middleware

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Function to hash password
func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPass), err
}

// Function to compare password hash
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

// Function to set content JSON
func SetContentJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
