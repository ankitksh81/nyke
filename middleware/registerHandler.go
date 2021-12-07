package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/models"
)

// Register user using registration form
func RegisterFromForm(w http.ResponseWriter, r *http.Request) {
	// user model
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Error500(w)
		logger.Log.Error("Unable to decode the request body. " + err.Error())
		return
	}

	// check if user already exist
	found, err := CheckUserExist(user.Email, w)
	if err != nil {
		Error500(w)
		return
	} else if found {
		ErrorUserExist(w)
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		Error500(w)
		logger.Log.Error("Unable to hash the password. " + err.Error())
		return
	}

	// Insert data in database
	sqlQuery := `INSERT INTO users(email, first_name, last_name, password) 
			VALUES ($1, $2, $3, $4) RETURNING user_id`

	var user_id string
	err = DB.QueryRow(sqlQuery, user.Email, user.FirstName, user.LastName, hashedPassword).Scan(&user_id)
	if err != nil {
		Error500(w)
		logger.Log.Error("Unable to execute the query " + err.Error())
		return
	}

	GenerateResponse(w, &user, user_id)
}

// Register user using OAuth2
func RegisterFromAuth(w http.ResponseWriter, authRes *models.AuthResponse) {
	// Getting details from OAuth2.0 response
	name := strings.Split(authRes.Name, " ")
	user := &models.User{
		Email:     authRes.Email,
		FirstName: name[0],
		LastName:  name[1],
		Picture:   authRes.UserPicture,
	}

	// Insert data in database
	sqlQuery := `INSERT INTO users(email, first_name, last_name, user_picture) 
			VALUES ($1, $2, $3, $4) RETURNING user_id`

	var user_id string
	err := DB.QueryRow(sqlQuery, user.Email, user.FirstName, user.LastName, user.Picture).Scan(&user_id)
	if err != nil {
		Error500(w)
		logger.Log.Error("Unable to execute the query " + err.Error())
		return
	}

	GenerateResponse(w, user, user_id)
}

func GenerateResponse(w http.ResponseWriter, user *models.User, user_id string) {

	token_string, err := GenerateJWT(user_id)
	if err != nil {
		Error500(w)
		logger.Log.Error("Error generating jwt token" + err.Error())
		return
	}

	resUser := &models.ResponseUser{
		ID:          user_id,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		UserPicture: user.Picture,
	}
	res := &models.Response{
		User:     *resUser,
		JwtToken: token_string,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		Error500(w)
		logger.Log.Error(err.Error())
	}
}

func CheckUserExist(email string, w http.ResponseWriter) (bool, error) {
	sqlQuery := `SELECT user_id FROM users WHERE email = $1`

	row := DB.QueryRow(sqlQuery, email)

	var user_id string
	err := row.Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			Error500(w)
		}
		return false, err
	}

	return true, nil
}
