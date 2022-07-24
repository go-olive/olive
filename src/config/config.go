package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	l "github.com/go-olive/olive/src/log"
	"github.com/go-olive/tv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	AppVersion string
	APP        = &appConfig{}
	defaultAPP = &appConfig{
		LogLevel:          logrus.DebugLevel,
		SnapRestSeconds:   15,
		SplitRestSeconds:  60,
		CommanderPoolSize: 1,
		UploadConfig:      &UploadConfig{},
		PlatformConfig: &PlatformConfig{
			DouyinCookie:   "__ac_nonce=062c84d05004a461cf7f2; __ac_signature=_02B4Z6wo00f01NAKk1QAAIDBqMR4UNttsUjQKpfAAFbQjTrG-JmICfTUMMVzKe3crg5Fk4y4e4DGURjAV4VW2B6WwXdqq3UC1c0waQMKIjhZn5Ve1LxiGmyDuVlSBN7aRhuGfEIwwfxXcYhA4e;",
			KuaishouCookie: "did=web_e7367cd09866e986065dd2e0f93630fe",
		},
	}

	appCfgFilePath string
	version        bool
	url            string
	cookie         string
	usage          = func() {
		fmt.Printf("Powered by go-olive/olive %s\n", AppVersion)
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}
)

type appConfig struct {
	LogLevel          logrus.Level
	SnapRestSeconds   uint
	SplitRestSeconds  uint
	CommanderPoolSize uint

	*UploadConfig
	*PlatformConfig
	Shows []*Show
}

func (this *appConfig) checkAndFix() {
	if this.LogLevel == 0 {
		this.LogLevel = defaultAPP.LogLevel
	}
	if this.SnapRestSeconds == 0 {
		this.SnapRestSeconds = defaultAPP.SnapRestSeconds
	}
	if this.SplitRestSeconds == 0 {
		this.SplitRestSeconds = defaultAPP.SplitRestSeconds
	}
	if this.CommanderPoolSize == 0 {
		this.CommanderPoolSize = defaultAPP.CommanderPoolSize
	}
	if this.UploadConfig == nil {
		this.UploadConfig = &UploadConfig{}
	}
	if this.PlatformConfig == nil {
		this.PlatformConfig = &PlatformConfig{}
	}
	if this.DouyinCookie == "" {
		this.DouyinCookie = defaultAPP.DouyinCookie
	}
	if this.KuaishouCookie == "" {
		this.KuaishouCookie = defaultAPP.KuaishouCookie
	}
	if cookie != "" {
		this.DouyinCookie = cookie
		this.KuaishouCookie = cookie
	}
}

type Show struct {
	Platform     string
	RoomID       string
	StreamerName string
	OutTmpl      string
	Parser       string
	SaveDir      string
	PostCmds     []*exec.Cmd
	SplitRule    *SplitRule
}

type SplitRule struct {
	FileSize int64

	Duration       string
	parsedDuration time.Duration
}

func (this *SplitRule) IsValid() bool {
	if this == nil {
		return false
	}
	if this.parsedDuration <= 0 && this.FileSize <= 0 {
		return false
	}
	return true
}

func (this *SplitRule) Satisfy(startTime time.Time, out string) bool {
	if !this.IsValid() {
		return false
	}

	if this.parsedDuration > 0 {
		if time.Since(startTime) >= this.parsedDuration {
			return true
		}
	}

	if this.FileSize > 0 {
		if fi, err := os.Stat(out); err == nil {
			if fi.Size() >= this.FileSize {
				return true
			}
		}
	}

	return false
}

// fix parser
func (s *Show) checkAndFix() {
	if s.Parser != "" {
		return
	}
	switch s.Platform {
	case "youtube",
		"twitch",
		"streamlink":
		s.Parser = "streamlink"
	default:
		s.Parser = "flv"
	}

	if s.SplitRule != nil {
		if s.SplitRule.Duration != "" {
			var err error
			s.SplitRule.parsedDuration, err = time.ParseDuration(s.SplitRule.Duration)
			if err != nil {
				l.Logger.Error(err)
			}
		}
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
	flag.BoolVar(&version, "v", version, "print olive version")
	flag.StringVar(&appCfgFilePath, "f", appCfgFilePath, "set config.toml filepath")

	flag.StringVar(&url, "u", url, "room url")
	flag.StringVar(&cookie, "c", "", "http cookie")

	flag.Parse()

	if version {
		fmt.Println(AppVersion)
		os.Exit(0)
	} else if url != "" {
		t, err := tv.NewWithUrl(url, tv.SetCookie(cookie))
		if err != nil {
			l.Logger.Fatal(err)
		}
		site, _ := tv.Sniff(t.SiteID)
		APP.Shows = []*Show{
			{
				StreamerName: site.Name(),
				Platform:     t.SiteID,
				RoomID:       t.RoomID,
			},
		}
		APP.verify()
	} else {
		if appCfgFilePath == "" {
			usage()
			os.Exit(0)
		}

		viper.SetConfigFile(appCfgFilePath)
		if err := viper.ReadInConfig(); err != nil {
			l.Logger.WithField("err", err.Error()).
				Fatal("load config file failed")
		}
		if err := viper.Unmarshal(&APP); err != nil {
			l.Logger.WithField("err", err.Error()).
				Fatal("load config file failed")
		}
		APP.verify()

		viper.OnConfigChange(func(e fsnotify.Event) {
			viper.Unmarshal(&APP)
			APP.verify()
		})
		viper.WatchConfig()
	}
}

func (this *appConfig) verify() {
	this.checkAndFix()
	l.Logger.SetLevel(this.LogLevel)

	for _, s := range this.Shows {
		s.checkAndFix()
		if s.Parser == "flv" {
			continue
		}

		if _, err := exec.LookPath(s.Parser); err != nil {
			l.Logger.Fatalf("%s needs to be installed first", s.Parser)
		}
	}

	if this.UploadConfig.Enable {
		if _, err := exec.LookPath(this.UploadConfig.ExecPath); err != nil {
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
