package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Ajesh8/UserTasks/pkg/models"
	"github.com/dgrijalva/jwt-go"
)

func ParseBody(r *http.Request, x interface{}) error {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func CheckUserAuthorisation(w http.ResponseWriter, r *http.Request) (bool, *models.User) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Print(err.Error())
			return false, nil
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return false, nil
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &models.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return models.JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Print(err.Error())
			return false, nil
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(err.Error())
		return false, nil
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Print(err.Error())
		return false, nil
	}
	return true, models.FindUserByUserName(claims.Username)
}
