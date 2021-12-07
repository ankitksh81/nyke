package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/middleware"
	"github.com/ankitksh81/nyke/models"
	"github.com/gorilla/mux"
)

// create user endpoint
// api: /api/register
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middleware.Error500(w)
		logger.Log.Error("Unable to decode the request body. " + err.Error())
		return
	}

	// validate email
	if !middleware.IsEmailValid(user.Email) {
		err := &models.Error{
			Message: "Please enter a valid email.",
		}
		middleware.JSONError(w, err, 400)
		return
	}

	/*
		// validate password
		if !isPasswordValid(user.Password) {
			err := &models.Error{
				Message: "Password doesn't follow the rules.",
			}
			JSONError(w, err, 400)
			return
		}
	*/

	// check if user already exist
	found, err := middleware.CheckUserExist(user.Email, w)
	if err != nil {
		middleware.Error500(w)
		return
	} else if found {
		middleware.ErrorUserExist(w)
		return
	}

	hashedPassword, err := middleware.HashPassword(user.Password)
	if err != nil {
		middleware.Error500(w)
		logger.Log.Error("Unable to hash the password. " + err.Error())
		return
	}

	// Insert data in database
	sqlQuery := `INSERT INTO users(email, first_name, last_name, password, user_picture) 
			VALUES ($1, $2, $3, $4, $5) RETURNING user_id`

	var user_id string
	err = middleware.DB.QueryRow(sqlQuery, user.Email, user.FirstName, user.LastName,
		hashedPassword, user.Picture).Scan(&user_id)
	if err != nil {
		middleware.Error500(w)
		logger.Log.Error("Unable to execute the query " + err.Error())
		return
	}

	middleware.GenerateResponse(w, &user, user_id)
}

// get user with user_id endpoint
// api: /api/user/{id}
func GetUser(w http.ResponseWriter, r *http.Request) {
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)
	user_id := params["id"]

	// call getUser function with user_id to retrieve a single user
	user, err := getUser(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.Error404(w)
		} else {
			middleware.Error500(w)
		}
		logger.Log.Error("Could not get user from the database. " + err.Error())
		return
	}

	json.NewEncoder(w).Encode(user)
}

func getUser(id string) (models.ResponseUser, error) {
	var user models.ResponseUser

	// create the select sql query
	sqlQuery := `SELECT user_id, email, first_name, last_name, user_picture 
		FROM users WHERE user_id = $1`

	row := middleware.DB.QueryRow(sqlQuery, id)
	// unmarshal the row object to user
	err := row.Scan(&user.ID, &user.Email, &user.FirstName,
		&user.LastName, &user.UserPicture)
	if err != nil {
		return user, err
	}
	return user, nil
}

// login user using oauth
