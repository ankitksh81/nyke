package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/middleware"
	"github.com/ankitksh81/nyke/models"
	"github.com/gorilla/mux"
)

// create user endpoint
// @ /api/register
func CreateUser(w http.ResponseWriter, r *http.Request) {
	middleware.SetContentJSON(w)
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
		if !middleware.IsPasswordValid(user.Password) {
			err := &models.Error{
				Message: "Password doesn't follow the rules.",
			}
			middleware.JSONError(w, err, 400)
			return
		}
	*/

	// check if user already exist
	found, err := middleware.CheckUserExist(user.Email)
	if err != nil {
		middleware.Error500(w)
		return
	} else if found {
		passwordSet, err := isPasswordSet(user.Email)
		if err != nil {
			middleware.Error500(w)
			return
		} else if !passwordSet {
			// do something
			middleware.SetContentJSON(w)
			json.NewEncoder(w).Encode("Set a new password")
			return
		} else {
			middleware.ErrorUserExist(w)
			return
		}
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

func isPasswordSet(e string) (bool, error) {
	sqlQuery := `SELECT password FROM users WHERE email=$1`
	var password sql.NullString
	err := middleware.DB.QueryRow(sqlQuery, e).Scan(&password)
	if err != nil {
		logger.Log.Error("Query Error: " + err.Error())
		return false, err
	} else if !password.Valid {
		return false, nil
	}
	return true, nil
}

// get user by user_id endpoint
// @ /api/user/{id}
func GetUser(w http.ResponseWriter, r *http.Request) {
	middleware.SetContentJSON(w)
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

// create user using oauth
func CreateAuthUser(w http.ResponseWriter, authRes *models.AuthResponse) {
	middleware.SetContentJSON(w)
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
	err := middleware.DB.QueryRow(sqlQuery, user.Email, user.FirstName, user.LastName, user.Picture).Scan(&user_id)
	if err != nil {
		middleware.Error500(w)
		logger.Log.Error("Unable to execute the query " + err.Error())
		return
	}
	middleware.GenerateResponse(w, user, user_id)
}

// login user using oauth
// login user using password
