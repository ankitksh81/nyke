package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"unicode"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/models"
)

// generate json response with jwt token
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

// check if user exists
func CheckUserExist(email string) (bool, error) {
	sqlQuery := `SELECT user_id FROM users WHERE email = $1`

	row := DB.QueryRow(sqlQuery, email)

	var user_id string
	err := row.Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Check if email is valid
func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

// Check if password is valid
func IsPasswordValid(p string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(p) >= 7 {
		hasMinLen = true
	}

	for _, char := range p {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
