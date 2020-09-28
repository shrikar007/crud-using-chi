package common

import (
	"crud-using-chi/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateJWT(user models.User,mySigningKey string) string {
	claims := models.MyCustomClaims{
		Name: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return tokenString
}
