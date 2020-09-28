package models

import (
	"github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}
var Claims *MyCustomClaims
