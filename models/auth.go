package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

type MyCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

var Claims *MyCustomClaims

type AuthResource struct {
	DB *sqlx.DB
}

func NewAuthModel(db *sqlx.DB) (auth *AuthResource) {
	auth = &AuthResource{}
	auth.DB = db
	return
}
func (auth *AuthResource) StoreAuth(login LoginResponse) (err error) {
	_, err = auth.DB.Exec(`INSERT INTO sessions (token,valid) VALUES (?,?)`, login.Token, true)
	if err != nil {
		return
	}
	return
}

func (u *UserResource) GetAuthByToken(token string) (user User, err error) {
	err = u.DB.Get(&user, "SELECT * FROM sessions WHERE token=?", token)
	if err != nil {
		return
	}
	return
}
