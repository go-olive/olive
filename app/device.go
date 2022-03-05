package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luxcgo/lifesaver/engine"
	l "github.com/luxcgo/lifesaver/log"
	"github.com/luxcgo/lifesaver/monitor"
	"github.com/luxcgo/lifesaver/recorder"
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	s, err := engine.NewShow("huya", "yingying8808", "沐莹莹")
	if err != nil {
		println(err)
		return
	}
	s.AddMonitor()

	s2, err := engine.NewShow("youtube", "UCwV9VXgUFpCKbf8SUsE6OSw", "domado")
	if err != nil {
		println(err)
		return
	}
	s2.AddMonitor()

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
		log.Printf("超时，强制退出~\n")
		d.done <- struct{}{}
	}()
	recorder.RecorderManager.Stop()
	monitor.MonitorManager.Stop()
	close(d.done)
}
