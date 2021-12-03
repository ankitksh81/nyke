package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/models"
)

/* Register user using registration form */
func RegisterFromForm(w http.ResponseWriter, r *http.Request) {
	// create an empty user of type models.User
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.Log.Error("Unable to decode the request body. " + err.Error())
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		logger.Log.Error("Unable to hash the password. " + err.Error())
	}

	/* Inserting data in database */
	sqlQuery := `INSERT INTO users(email, first_name, last_name, password) VALUES ($1, $2, $3, $4) RETURNING user_id`

	var user_id string

	err = DB.QueryRow(sqlQuery, user.Email, user.FirstName, user.LastName, hashedPassword).Scan(&user_id)
	if err != nil {

		errMsg := err.Error()

		json.NewEncoder(w).Encode(errMsg)
		logger.Log.Error("Unable to execute the query " + err.Error())
	}

	/* json response object */
	var res models.Response

	res.User.ID = user_id
	res.User.Email = user.Email
	res.User.FirstName = user.FirstName
	res.User.LastName = user.LastName
	res.User.UserPicture = user.Picture
	res.JwtToken, err = GenerateJWT(user_id)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), 500)
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Log.Error(err.Error())
		// w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), 500)
	}

	logger.Log.Info("UUID returned: " + user_id)
}
