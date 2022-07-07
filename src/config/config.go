package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	l "github.com/go-olive/olive/src/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	AppVersion string
	APP        *appConfig
	defaultAPP = appConfig{}
)

type appConfig struct {
	LogLevel        logrus.Level
	SnapRestSeconds uint32

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
	SaveDir      string
	PostCmds     []*exec.Cmd
}

// fix parser
func (s *Show) checkAndFix() {
	if s.Parser != "" {
		return
	}
	switch s.Platform {
	case "youtube",
		"twitch":
		s.Parser = "streamlink"
	default:
		s.Parser = "flv"
	}
}

type UploadConfig struct {
	Enable   bool
	ExecPath string
	Filepath string
}

type PlatformConfig struct {
	DouyinCookie   string
	KuaishouCookie string
}

func init() {
	var (
		appCfgFilePath string
		version        bool
	)
	flag.BoolVar(&version, "v", version, "print olive version")
	flag.StringVar(&appCfgFilePath, "c", appCfgFilePath, "set config.toml filepath")
	flag.Parse()

	if version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	var Usage = func() {
		fmt.Printf("Powered by go-olive/olive %s\n", AppVersion)
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}
	if appCfgFilePath == "" {
		Usage()
		os.Exit(0)
	}

	viper.SetConfigFile(appCfgFilePath)
	if err := viper.ReadInConfig(); err != nil {
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
	l.Logger.SetLevel(APP.LogLevel)

	if APP.UploadConfig == nil {
		APP.UploadConfig = &UploadConfig{}
	}
	if APP.PlatformConfig == nil {
		APP.PlatformConfig = &PlatformConfig{}
	}

	for _, s := range APP.Shows {
		s.checkAndFix()
		if s.Parser == "flv" {
			continue
		}

		if _, err := exec.LookPath(s.Parser); err != nil {
			l.Logger.Fatalf("%s needs to be installed first", s.Parser)
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
