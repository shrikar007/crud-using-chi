package logger

import (
	"github.com/sirupsen/logrus"
)

func InitLogger() (lg *logrus.Logger){
	lg = logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "01-01-2019 15:10:10"
	formatter.FullTimestamp = true
	lg.SetFormatter(formatter)
	return
}
