package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ConnectDb(lg *logrus.Logger,conf *viper.Viper) (DB *sqlx.DB, err error) {
	DB, err = sqlx.Connect(conf.GetString("database.db_driver"),
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
			conf.GetString("database.db_user"), conf.GetString("database.db_pass"),
			conf.GetString("database.db_host"), conf.GetString("database.db_port"), conf.GetString("database.db_name"),
		))
	if err != nil {
      lg.Error(err)
      return
	}
	return
}

