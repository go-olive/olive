package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luxcgo/lifesaver/engine"
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

	s, err := engine.NewShow("huya", "yingying8808")
	if err != nil {
		println(err)
		return
	}
	s.AddMonitor()
	go d.listenSignal()
	<-d.done
}

func (d *device) listenSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	for sig := range ch {
		log.Printf("收到结束信号(%s)，准备结束进程\n", sig.String())
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
