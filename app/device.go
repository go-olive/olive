package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luxcgo/lifesaver/config"
	"github.com/luxcgo/lifesaver/engine"
	l "github.com/luxcgo/lifesaver/log"
	"github.com/luxcgo/lifesaver/monitor"
	"github.com/luxcgo/lifesaver/recorder"
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
	for _, v := range config.APP.Shows {
		s, err := engine.NewShow(v.Platform, v.RoomID, v.StreamerName)
		if err != nil {
			l.Logger.WithFields(logrus.Fields{
				"pf": v.Platform,
				"id": v.RoomID,
			}).Error(err)
			continue
		}
		s.AddMonitor()
	}

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
	close(d.done)
}
