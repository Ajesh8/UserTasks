package models

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"" json:"name"`
	Email    string `json:"email"`
	Username string `json:username`
	Password string `json:password`
	Tasks    []Task `json:"tasks"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var JwtKey = []byte("my_secret_key")

func (b *User) CreateUser() (*User, error) {
	if !validateUser(*b) {
		return nil, errors.New("user' mandatory fields not found")
	}
	user := &User{}
	db.Where("username=?", b.Username).Find(&user)
	if user.ID != 0 {
		return nil, errors.New("username already exists")
	}
	db.Create(&b)
	return b, nil
}

func validateUser(user User) bool {
	if user.Name == "" || user.Email == "" || user.Password == "" || user.Username == "" {
		return false
	}
	return true
}

func CheckUserCredentials(credentials *Credentials) bool {
	user := FindUserByUserName(credentials.Username)
	return user.Password == credentials.Password
}

func FindUserByUserName(username string) *User {
	user := &User{}
	db.Where("username = ?", username).First(&user)
	return user
}
