package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-olive/olive/src/config"
	"github.com/go-olive/olive/src/engine"
	l "github.com/go-olive/olive/src/log"
	"github.com/go-olive/olive/src/monitor"
	"github.com/go-olive/olive/src/recorder"
	"github.com/go-olive/olive/src/uploader"
	"github.com/sirupsen/logrus"
)

type IDevice interface {
	Run()
	Stop()
}

type device struct {
	done chan struct{}
}

func NewDevice() IDevice {
	return &device{
		done: make(chan struct{}),
	}
}

func (d *device) Run() {
	l.Logger.Infof("Powered by go-olive/olive %s", config.AppVersion)

	for _, v := range config.APP.Shows {
		s, err := engine.NewShow(v.Platform, v.RoomID,
			engine.WithStreamerName(v.StreamerName),
			engine.WithOutTmpl(v.OutTmpl),
			engine.WithParser(v.Parser),
			engine.WithSaveDir(v.SaveDir),
			engine.WithPostCmds(v.PostCmds),
			engine.WithSplitRule(v.SplitRule),
		)
		if err != nil {
			l.Logger.WithFields(logrus.Fields{
				"pf": v.Platform,
				"id": v.RoomID,
			}).Error(err)
			continue
		}
		s.AddMonitor()
	}
	uploader.UploaderWorkerPool.Run()
	go d.listenSignal()
	<-d.done
}

func (d *device) listenSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	for sig := range ch {
		l.Logger.WithField("signal", sig.String()).
			Info("handle request")
		d.Stop()
		return
	}
}

func (d *device) Stop() {
	go func() {
		<-time.After(time.Duration(time.Second * 5))
		l.Logger.Info("timeout, force quit")
		d.done <- struct{}{}
	}()
	recorder.RecorderManager.Stop()
	monitor.MonitorManager.Stop()
	uploader.UploaderWorkerPool.Stop()
	close(d.done)
}
