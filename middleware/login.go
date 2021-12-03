package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ankitksh81/nyke/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.UserLogin
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		err := &models.Error{
			Message: "Email cannot be empty",
		}
		http.Error(w, err.Message, 400)
		return
	}

	if user.Password == "" {
		err := &models.Error{
			Message: "Password cannot be empty",
		}
		http.Error(w, err.Message, 400)
		return
	}

	sqlQuery := `SELECT password, user_id FROM users WHERE email = $1`

	row := DB.QueryRow(sqlQuery, user.Email)
	var password_hash, user_id string
	err := row.Scan(&password_hash, &user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			err := &models.Error{
				Message: "Email is incorrect",
			}
			http.Error(w, err.Message, 404)
		} else {
			err := &models.Error{
				Message: "Internal server error",
			}
			http.Error(w, err.Message, 500)
		}
		return
	}

	ok := CheckPasswordHash(user.Password, password_hash)
	if !ok {
		err := &models.Error{
			Message: "Password is incorrect",
		}
		http.Error(w, err.Message, 401)
		return
	}

	jwtToken, err := GenerateJWT(user_id)
	if err != nil {
		err := &models.Error{
			Message: "Internal server error",
		}
		http.Error(w, err.Message, 500)
		return
	}

	json.NewEncoder(w).Encode(jwtToken)
}
