package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Ajesh8/UserTasks/pkg/models"
	"github.com/Ajesh8/UserTasks/pkg/utils"
	"github.com/dgrijalva/jwt-go"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	CreateUser := &models.User{}
	err := utils.ParseBody(r, CreateUser)
	if err != nil {
		log.Printf("Could not parse body:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := CreateUser.CreateUser()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(&Message{
			Msg: err.Error(),
		})
		w.Write(res)
		return
	}
	res, err := json.Marshal(user)
	if err != nil {
		log.Printf("Marshalling error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	creds := &models.Credentials{}
	err := utils.ParseBody(r, creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorized := models.CheckUserCredentials(creds)
	if !authorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(models.JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
