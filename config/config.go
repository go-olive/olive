package config

import (
	"github.com/fsnotify/fsnotify"
	l "github.com/luxcgo/lifesaver/log"
	"github.com/spf13/viper"
)

var (
	APP        *appConfig
	defaultAPP = appConfig{}
)

type appConfig struct {
	Shows []*Show
}

type Show struct {
	Platform     string
	RoomID       string
	StreamerName string
}

func init() {
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		l.Logger.WithField("err", err.Error()).
			Fatal("load config file failed")
	}
	viper.Unmarshal(&APP)
	verify()

	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.Unmarshal(&APP)
		verify()
	})
	viper.WatchConfig()
}

func verify() {
	if APP == nil {
		l.Logger.Info("use default APP config")
		APP = &defaultAPP
	}
}
