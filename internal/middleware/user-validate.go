package middleware

import (
	"crud-using-chi/internal/models"
	"crud-using-chi/pkg/common"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type (
	MiddlewareUser struct {
		Conf   *viper.Viper
		Logger *logrus.Logger
		DB     *sqlx.DB
	}
)

func NewMiddlerwareUser(conf *viper.Viper, logger *logrus.Logger, db *sqlx.DB) (mu *MiddlewareUser) {
	mu = &MiddlewareUser{}
	mu.Conf = conf
	mu.Logger = logger
	mu.DB = db
	return
}

func (mu *MiddlewareUser) IsAuthorized(endpoint http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			modelAuth = models.NewAuthModel(mu.DB)
			session   models.Session
		)
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Something went wrong")
				}
				return []byte(mu.Conf.GetString("jwt.jwt_signing_key")), nil
			})
			if err != nil {
				mu.Logger.Errorf("Error occured : %v", err)
			}
			session, err = modelAuth.GetSessionByToken(token.Raw)
			if token.Valid && session.Valid {
				token, _ := jwt.ParseWithClaims(token.Raw, &models.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(mu.Conf.GetString("jwt.jwt_signing_key")), nil
				})

				models.Claims = token.Claims.(*models.MyCustomClaims)
				endpoint.ServeHTTP(w, r)

			} else {
				common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
		} else {
			common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
	})
}
