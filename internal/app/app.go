package app

import (
	"crud-using-chi/config"
	"crud-using-chi/internal/database"
	"crud-using-chi/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Initialize() (r chi.Router, conf *viper.Viper) {
	var (
		lg  *logrus.Logger
		DB  *sqlx.DB
		err error
	)
	lg = logger.InitLogger()
	lg.Info("Logger initialized")

	if conf, err = config.InitConfig(); err != nil {
		lg.Fatal(err)
	}
	lg.Info("Config initialized")
	lg.Info("Connecting to database")
	DB, err = database.ConnectDb(lg, conf)
	if err != nil {
		lg.Error(err)
		return
	}
	lg.Info("Successfully connected to database")
	r = RegisterRoutes(lg, DB, conf)
	return
}
