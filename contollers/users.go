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
		var (
			login     = models.Login{}
			modelUser = models.NewUserModel(u.DB)
			modelAuth = models.NewAuthModel(u.DB)
		)
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&login); err != nil {
			common.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := modelUser.GetUser(login)
		if err != nil {
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		loginResponse := models.LoginResponse{Token: common.GenerateJWT(user, u.Conf.GetString("jwt.jwt_signing_key"))}
		err = modelAuth.StoreAuth(loginResponse)
		if err != nil {
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		common.RespondJSON(w, http.StatusOK, loginResponse)
	}
}

func (u *Users) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user      = models.User{}
			modelUser = models.NewUserModel(u.DB)
		)
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			common.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
		err := modelUser.StoreUser(user)
		if err != nil {
			common.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		common.RespondJSON(w, http.StatusCreated, map[string]string{"message": "Signed Up successful"})
	}
}

func (u *Users) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user      models.User
			modelUser = models.NewUserModel(u.DB)
			err       error
		)
		Id := chi.URLParam(r, "id")
		user, err = modelUser.GetProfileById(Id)
		if err != nil {
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		common.RespondJSON(w, http.StatusOK, user)
	}
}

func (u *Users) UpdateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user      models.User
			modelUser = models.NewUserModel(u.DB)
			err       error
		)
		Id := chi.URLParam(r, "id")
		user, err = modelUser.GetProfileById(Id)
		if err != nil {
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			common.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
		err = modelUser.UpdateProfile(user, Id)
		if err != nil {
			common.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		common.RespondJSON(w, http.StatusCreated, map[string]string{"message": "updated"})
	}
}
