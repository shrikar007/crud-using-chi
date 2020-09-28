package app

import (
	"crud-using-chi/config"
	"crud-using-chi/contollers"
	"crud-using-chi/database"
	"crud-using-chi/logger"
	user_middleware "crud-using-chi/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Initialize() (r chi.Router, conf *viper.Viper) {
	var (
		lg    *logrus.Logger
		DB    *sqlx.DB
		user  *contollers.Users
		err   error
		muser *user_middleware.MiddlewareUser
	)
	lg = logger.InitLogger()
	lg.Info("Logger initialized")

	if conf, err = config.InitConfig(); err != nil {
		lg.Fatal(err)
	}
	lg.Info("Config initialized")
	lg.Info("Connecting to database")
	DB = database.ConnectDb(lg, conf)
	lg.Info("Successfully connected to database")

	user = contollers.NewUser(conf, lg, DB)
	muser = user_middleware.NewMiddlerwareUser(conf, lg, DB)

	lg.Info("Routes initialized")
	r = chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/", func(r chi.Router) {
		r.Post("/signup", user.SignUp())
		r.Post("/login", user.Login())
	})
	r.Group(func(r chi.Router) {
		r.Use(muser.IsAuthorized)
		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", user.GetProfile())
			//r.Put("/", custMiddleware.UserAuthorized(controller.UpdateUser()))
		})
	})
	return
}
