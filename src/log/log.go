package l

import (
	"io"
	"log"
	"os"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	f, err := os.OpenFile("a.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal(err)
	}

	Logger = &logrus.Logger{
		Out: io.MultiWriter(os.Stderr),
		Formatter: &logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.TraceLevel,
		// ReportCaller: true,
	}

	Logger.AddHook(lfshook.NewHook(
		f,
		// &logrus.JSONFormatter{
		// 	TimestampFormat: "2006-01-02 15:04:05",
		// },
		&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	))

}
