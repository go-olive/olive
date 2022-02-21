package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luxcgo/lifesaver/engine"
	_ "github.com/luxcgo/lifesaver/internal"
	"github.com/luxcgo/lifesaver/monitor"
	"github.com/luxcgo/lifesaver/recorder"
)

var done = make(chan struct{})

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	s, err := engine.NewShow("huya", "92852")
	if err != nil {
		println(err)
		return
	}
	s.AddMonitor()
	go listenSignal()
	<-done
}

func listenSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	for sig := range ch {
		log.Printf("收到结束信号(%s)，准备结束进程\n", sig.String())
		go func() {
			select {
			case <-time.After(time.Duration(time.Second * 5)):
				log.Printf("超时，强制退出~\n")
				done <- struct{}{}
				return
			}
		}()
		recorder.RecorderManager.Stop()
		monitor.MonitorManager.Stop()
		done <- struct{}{}
	}
}
