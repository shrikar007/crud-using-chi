package contollers

import (
	"crud-using-chi/common"
	"crud-using-chi/models"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type (
	Users struct {
		Conf   *viper.Viper
		Logger *logrus.Logger
		DB     *sqlx.DB
	}
)

func NewUser(conf *viper.Viper, logger *logrus.Logger, db *sqlx.DB) (u *Users) {
	u = &Users{}
	u.Conf = conf
	u.Logger = logger
	u.DB = db
	return
}

func (u *Users) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := models.Login{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&login); err != nil {
			common.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		user := models.User{}
		err := u.DB.Get(&user, "SELECT * FROM users WHERE email=? AND password=?", login.Email, login.Password)
		if err != nil {
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		loginResponse := models.LoginResponse{Token: common.GenerateJWT(user,u.Conf.GetString("jwt.jwt_signing_key"))}

		common.RespondJSON(w, http.StatusOK, loginResponse)
	}
}

func (u *Users) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			common.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		_, err := u.DB.Exec(`INSERT INTO users (first_name,middle_name,password,last_name, email,gender,mobile) VALUES (?,?,?,?,?,?,?)`, user.FirstName, user.MiddleName, user.Password, user.LastName, user.Email, user.Gender, user.Mobile)
		if err != nil {
			common.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		common.RespondJSON(w, http.StatusCreated, "Signed Up successful")
	}
}

func (u *Users) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		Id := chi.URLParam(r, "id")
		err := u.DB.Get(&user, "SELECT * FROM users WHERE id=?", Id)
		if err != nil {
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		common.RespondJSON(w, http.StatusOK, user)
	}
}
