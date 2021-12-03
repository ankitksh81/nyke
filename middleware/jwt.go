package middleware

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

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
