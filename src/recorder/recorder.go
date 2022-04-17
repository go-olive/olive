package recorder

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/go-olive/olive/src/config"
	"github.com/go-olive/olive/src/engine"
	"github.com/go-olive/olive/src/enum"
	l "github.com/go-olive/olive/src/log"
	"github.com/go-olive/olive/src/parser"
	"github.com/go-olive/olive/src/uploader"
	"github.com/sirupsen/logrus"
)

type Recorder interface {
	Start() error
	Stop()
	StartTime() time.Time
	Done() <-chan struct{}
}

type recorder struct {
	status    enum.StatusID
	show      engine.Show
	stop      chan struct{}
	startTime time.Time
	parser    parser.Parser
	done      chan struct{}
}

func NewRecorder(show engine.Show) (Recorder, error) {
	parser, err := show.NewParser()
	if err != nil {
		return nil, err
	}
	return &recorder{
		status:    enum.Status.Starting,
		show:      show,
		stop:      make(chan struct{}),
		startTime: time.Now(),
		parser:    parser,
		done:      make(chan struct{}),
	}, nil
}

func (r *recorder) Start() error {
	if !atomic.CompareAndSwapUint32(&r.status, enum.Status.Starting, enum.Status.Pending) {
		return nil
	}
	defer atomic.CompareAndSwapUint32(&r.status, enum.Status.Pending, enum.Status.Running)
	go r.run()

	l.Logger.WithFields(logrus.Fields{
		"pf": r.show.GetPlatform(),
		"id": r.show.GetRoomID(),
	}).Info("recorder start")

	return nil
}

func (r *recorder) Stop() {
	if !atomic.CompareAndSwapUint32(&r.status, enum.Status.Running, enum.Status.Stopping) {
		return
	}
	close(r.stop)
	r.parser.Stop()
}

func (r *recorder) StartTime() time.Time {
	return r.startTime
}

func (r *recorder) record() {
	var out string
	defer func() {
		fi, err := os.Stat(out)
		if err != nil {
			return
		}
		const tenMB = 1e7
		if fi.Size() < tenMB {
			os.Remove(out)
			return
		}

		if config.APP.UploadConfig.Enable {
			r.SubmitUploadTask(out)
		}
	}()

	r.show.Refresh()
	streamUrl, ok := r.show.StreamUrl()
	roomName, _ := r.show.RoomName()
	if !ok {
		l.Logger.WithFields(logrus.Fields{
			"pf": r.show.GetPlatform(),
			"id": r.show.GetRoomID(),
		}).Debug("fail to get StreamURL")
		time.Sleep(5 * time.Second)
		return
	}

	const format = "2006-01-02 15-04-05"
	out = fmt.Sprintf("[%s][%s][%s].flv", r.show.GetStreamerName(), roomName, time.Now().Format(format))

	l.Logger.WithFields(logrus.Fields{
		"pf": r.show.GetPlatform(),
		"id": r.show.GetRoomID(),
	}).Info("record start")

	err := r.parser.Parse(streamUrl, out)

	l.Logger.WithFields(logrus.Fields{
		"pf": r.show.GetPlatform(),
		"id": r.show.GetRoomID(),
	}).Infof("record stop: %+v", err)

}

func (r *recorder) run() {
	for {
		select {
		case <-r.stop:
			close(r.done)
			l.Logger.WithFields(logrus.Fields{
				"pf": r.show.GetPlatform(),
				"id": r.show.GetRoomID(),
			}).Info("recorder stop")
			return
		default:
			r.record()
		}
	}
}

func (r *recorder) Done() <-chan struct{} {
	return r.done
}

func (r *recorder) SubmitUploadTask(filepath string) {
	uploader.UploaderWorkerPool.AddTask(&uploader.UploadTask{
		Filepath: filepath,
		Tryout:   2,
	})
}
