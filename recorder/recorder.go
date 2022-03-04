package recorder

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/luxcgo/lifesaver/engine"
	"github.com/luxcgo/lifesaver/enum"
	l "github.com/luxcgo/lifesaver/log"
	"github.com/luxcgo/lifesaver/parser"
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
	u, err := r.show.StreamURL()
	if err != nil {
		l.Logger.Error(err)
		time.Sleep(5 * time.Second)
		return
	}
	t := time.Now().Format("[2006-01-02 15-04-05].flv")

	l.Logger.WithFields(logrus.Fields{
		"pf": r.show.GetPlatform(),
		"id": r.show.GetRoomID(),
	}).Info("record start")

	err = r.parser.Parse(u, t)

	if err != nil {
		log.Println(err.Error())
		l.Logger.WithFields(logrus.Fields{
			"pf": r.show.GetPlatform(),
			"id": r.show.GetRoomID(),
		}).Infof("record stop fail: %+v", err)
	} else {
		l.Logger.WithFields(logrus.Fields{
			"pf": r.show.GetPlatform(),
			"id": r.show.GetRoomID(),
		}).Info("record stop success")
	}
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
