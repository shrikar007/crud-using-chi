package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

type MyCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type Session struct {
	Id    string `json:"id" db:"id"`
	Token string `json:"token" db:"token"`
	Valid bool   `json:"valid" db:"valid"`
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
func (auth *AuthResource) StoreAuth(login Session) (err error) {
	_, err = auth.DB.Exec(`INSERT INTO sessions (token,valid) VALUES (?,?)`, login.Token, true)
	if err != nil {
		return
	}
	return
}

func (auth *AuthResource) GetSessionByToken(token string) (session Session, err error) {
	err = auth.DB.Get(&session, "SELECT * FROM sessions WHERE token=?", token)
	if err != nil {
		return
	}
	return
}
func (auth *AuthResource) UpdateAuthById(Id string) (err error) {
	_, err = auth.DB.Exec("UPDATE sessions SET valid=? WHERE id = ?", false, Id)
	if err != nil {
		return
	}
	return
}
