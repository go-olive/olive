package config

import (
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	l "github.com/go-olive/olive/src/log"
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
	OutTmpl      string
	Parser       string
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
		l.Logger.Fatal(err)
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

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		l.Logger.Fatal("ffmpeg needs to be installed first")
	}

	for _, v := range APP.Shows {
		if v.Platform == "youtube" {
			if _, err := exec.LookPath("streamlink"); err != nil {
				l.Logger.Fatal("streamlink needs to be installed first")
			}
		}

		if v.Parser != "" {
			if _, err := exec.LookPath(v.Parser); err != nil {
				l.Logger.Fatalf("%s needs to be installed first", v.Parser)
			}
		}
	}

	if APP.UploadConfig.Enable {
		if _, err := exec.LookPath(APP.UploadConfig.ExecPath); err != nil {
			l.Logger.Fatal("biliup needs to be installed first")
		}
		path, err := os.Getwd()
		if err != nil {
			l.Logger.Fatal(err)
		}
		cookiesFilePath := filepath.Join(path, "cookies.json")
		if _, err := os.Stat(cookiesFilePath); errors.Is(err, os.ErrNotExist) {
			l.Logger.Fatal("biliup: please put cookies.json file at the current path")
		}
	}
}
