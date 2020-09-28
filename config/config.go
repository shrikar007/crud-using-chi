package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"path"
	"path/filepath"
)

func InitConfig() (conf *viper.Viper, err error) {
	conf = viper.New()
	conf.SetConfigType("toml")
	conf.SetConfigName("config")

	absPath, err := filepath.Abs("")

	if err != nil {
		log.Fatal(err)
		return
	}

	directory := path.Join(absPath, "config")

	conf.AddConfigPath(directory)

	go conf.WatchConfig()

	conf.OnConfigChange(func(_ fsnotify.Event) {
		log.Printf("config changed..reloading config\n")
	})

	err = conf.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	return
}
