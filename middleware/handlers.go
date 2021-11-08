package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/models"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func Register(w http.ResponseWriter, r *http.Request) {

	// create an empty user of type models.User
	var user models.User

	// Checking headers
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		err = decoder.Decode(&user, r.PostForm)
		if err != nil {
			log.Fatal(err)
		}
	} else if contentType == "application/json" {
		// decode the json request to user
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Log.Error("Unable to decode the request body. " + err.Error())
		}
	}

	sqlQuery := `INSERT INTO users(email, first_name, last_name, user_picture) VALUES ($1, $2, $3, $4) RETURNING user_id`

	var id string

	err := DB.QueryRow(sqlQuery, user.Email, user.FirstName, user.LastName, user.Picture).Scan(&id)

	if err != nil {
		w.Write([]byte(err.Error()))
		logger.Log.Error("Unable to execute the query " + err.Error())
	}

	/* json response object */
	var res models.Response

	res.ID = id
	res.Email = user.Email
	res.FirstName = user.FirstName
	res.LastName = user.LastName
	res.UserPicture = user.Picture

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), 500)
	}

	logger.Log.Info("UUID returned: " + id)
}
