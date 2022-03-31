package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	l "github.com/luxcgo/lifesaver/log"
	"github.com/spf13/viper"
)

var (
	APP        *appConfig
	defaultAPP = appConfig{}
)

type appConfig struct {
	*UploadConfig
	*PlatformConfig
	Shows []*Show
}

type Show struct {
	Platform     string
	RoomID       string
	StreamerName string
}

type UploadConfig struct {
	Enable   bool
	ExecPath string
	Filepath string
}

type PlatformConfig struct {
	DouyinCookie string
}

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	appCfgFilePath := filepath.Join(path, "config.toml")

	flag.StringVar(&appCfgFilePath, "c", appCfgFilePath, "config.toml配置文件存放路径")
	flag.Parse()

	viper.SetConfigFile(appCfgFilePath)
	err = viper.ReadInConfig()
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
	if APP.UploadConfig == nil {
		APP.UploadConfig = &UploadConfig{}
	}
	if APP.PlatformConfig == nil {
		APP.PlatformConfig = &PlatformConfig{}
	}
}
