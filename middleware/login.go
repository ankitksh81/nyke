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
		ErrorEmptyInput(w, "Email")
		return
	}

	if user.Password == "" {
		ErrorEmptyInput(w, "Password")
		return
	}

	sqlQuery := `SELECT password, user_id FROM users 
			WHERE email = $1`

	row := DB.QueryRow(sqlQuery, user.Email)
	var password_hash, user_id string
	err := row.Scan(&password_hash, &user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorUserNotFound(w)
		} else {
			Error500(w)
		}
		return
	}

	err = CheckPasswordHash(user.Password, password_hash)
	if err != nil {
		ErrorWrongPassword(w)
		return
	}

	jwtToken, err := GenerateJWT(user_id)
	if err != nil {
		Error500(w)
		return
	}

	json.NewEncoder(w).Encode(jwtToken)
}
