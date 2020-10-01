package contollers

import (
	"crud-using-chi/internal/models"
	"crud-using-chi/pkg/common"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
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
			u.Logger.Error(err)
			common.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := modelUser.GetUser(login)
		if err != nil {
			u.Logger.Error(err)
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		loginSession := models.Session{Token: common.GenerateJWT(user, u.Conf.GetString("jwt.jwt_signing_key"))}
		err = modelAuth.StoreAuth(loginSession)
		if err != nil {
			u.Logger.Error(err)
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		common.RespondJSON(w, http.StatusOK, map[string]string{"Token": loginSession.Token})
	}
}
func (u *Users) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			modelAuth = models.NewAuthModel(u.DB)
		)
		token := r.Header["Token"]
		session, err := modelAuth.GetSessionByToken(token[0])
		if err != nil {
			u.Logger.Error(err)
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		err = modelAuth.UpdateAuthById(session.Id)
		if err != nil {
			u.Logger.Error(err)
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		common.RespondJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
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
			u.Logger.Error(err)
			common.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
		err := modelUser.StoreUser(user)
		if err != nil {
			u.Logger.Error(err)
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
			u.Logger.Error(err)
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
			u.Logger.Error(err)
			common.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			u.Logger.Error(err)
			common.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
		err = modelUser.UpdateProfile(user, Id)
		if err != nil {
			u.Logger.Error(err)
			common.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		common.RespondJSON(w, http.StatusCreated, map[string]string{"message": "updated"})
	}
}

func (u *Users) UserReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			users     = []models.User{}
			modelUser = models.NewUserModel(u.DB)
		)
		users, err := modelUser.GetProfiles()
		if users == nil || err != nil {
			u.Logger.Error("no data found")
			common.RespondError(w, http.StatusNotFound, errors.New("no data found").Error())
			return
		}
		header := []string{"ID", "FirstName", "MiddleName", "LastName", "Email", "Gender", "Mobile"}
		data := [][]string{header}
		for _, row := range users {
			var user = []string{}
			user = append(user, strconv.FormatUint(uint64(row.Id), 10))
			user = append(user, row.FirstName)
			user = append(user, row.MiddleName)
			user = append(user, row.LastName)
			user = append(user, row.Email)
			user = append(user, row.Gender)
			user = append(user, row.Mobile)
			data = append(data, user)
		}
		err = common.GenerateCSV("users", data)
		if err != nil {
			u.Logger.Error(err)
			common.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		message := []string{"CSV Generated Successfully"}
		common.RespondJSON(w, http.StatusOK, message)
	}
}
