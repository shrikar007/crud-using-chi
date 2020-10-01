package app

import (
	"crud-using-chi/internal/contollers"
	user_middleware "crud-using-chi/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterRoutes(lg *logrus.Logger, DB *sqlx.DB, conf *viper.Viper) (r chi.Router){
	var (
		muser *user_middleware.MiddlewareUser
		user  *contollers.Users
	)
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
		r.Post("/logout", user.Logout())
		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", user.GetProfile())
			r.Put("/{id}", user.UpdateProfile())
			r.Get("/report", user.UserReport())
		})
	})
	return
}
