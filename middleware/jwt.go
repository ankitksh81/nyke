package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// function to generate jwt token
func GenerateJWT(user_id string) (string, error) {
	secretKey := viper.GetString("jwtSecret")
	var mySigningKey = []byte(secretKey)
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = user_id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Printf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

// function to verify token
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		var mySigningKey = viper.GetString("jwtSecret")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error in parsing")
				}
				return []byte(mySigningKey), nil
			})
			if err != nil {
				errMsg := &models.Error{
					Message: "Your token has been expired!",
				}
				JSONError(w, errMsg, 401)
				logger.Log.Error("Token Error: " + err.Error())
				return
			}

			if token.Valid {
				handler.ServeHTTP(w, r)
				return
			} else {
				err := &models.Error{
					Message: "You are not authorized!",
				}
				JSONError(w, err, 401)
			}
		} else {
			err := &models.Error{
				Message: "Invalid token!",
			}
			JSONError(w, err, 401)
		}

	}
}
