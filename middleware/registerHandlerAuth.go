package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/models"
)

/* Register user using OAuth2.0 */
func RegisterFromAuth(w http.ResponseWriter, authRes *models.AuthResponse) {
	var user models.User

	/* Getting details from OAuth2.0 response */
	name := strings.Split(authRes.Name, " ")

	user.Email = authRes.Email
	user.FirstName = name[0]
	user.LastName = name[1]
	user.Picture = authRes.UserPicture

	/* Inserting data in database */
	sqlQuery := `INSERT INTO users(email, first_name, last_name, user_picture) VALUES ($1, $2, $3, $4) RETURNING user_id`

	var user_id string

	err := DB.QueryRow(sqlQuery, user.Email, user.FirstName, user.LastName, user.Picture).Scan(&user_id)
	if err != nil {
		w.Write([]byte(err.Error()))
		logger.Log.Error("Unable to execute the query " + err.Error())
	}

	/* json response object */
	var res models.Response

	res.User.ID = user_id
	res.User.Email = user.Email
	res.User.FirstName = user.FirstName
	res.User.LastName = user.LastName
	res.User.UserPicture = user.Picture

	token, err := GenerateJWT(user_id)
	if err != nil {
		log.Printf("Error generating jwt token, %s", err.Error())
	}

	fmt.Println("JWT token: " + token)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), 500)
	}

	// logger.Log.Info("UUID returned: " + user_id)
}
