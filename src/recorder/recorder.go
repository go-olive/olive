package recorder

import (
	"bytes"
	"fmt"
	"os"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/go-dora/filenamify"
	"github.com/go-olive/olive/src/config"
	"github.com/go-olive/olive/src/engine"
	"github.com/go-olive/olive/src/enum"
	l "github.com/go-olive/olive/src/log"
	"github.com/go-olive/olive/src/parser"
	"github.com/go-olive/olive/src/uploader"
	"github.com/sirupsen/logrus"
)

var (
	nameFuncMap = func() template.FuncMap {
		m := sprig.TxtFuncMap()
		return m
	}()

	defaultOutTmpl = template.Must(template.New("filename").Funcs(nameFuncMap).
			Parse(`[{{ .StreamerName }}][{{ .RoomName }}][{{ now | date "2006-01-02 15-04-05"}}].flv`))
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

	r.show.Snap()
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

	info := &struct {
		StreamerName string
		RoomName     string
	}{
		StreamerName: r.show.GetStreamerName(),
		RoomName:     roomName,
	}

	tmpl := defaultOutTmpl
	if r.show.GetOutTmpl() != "" {
		_tmpl, err := template.New("user_defined_filename").Funcs(nameFuncMap).Parse(r.show.GetOutTmpl())
		if err == nil {
			tmpl = _tmpl
		} else {
			l.Logger.Error(err)
		}
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, info); err != nil {
		l.Logger.Error(err)
		const format = "2006-01-02 15-04-05"
		out = fmt.Sprintf("[%s][%s][%s].flv", r.show.GetStreamerName(), roomName, time.Now().Format(format))
	} else {
		out = buf.String()
	}

	out = filenamify.FilenamifyMustCompile(out)

	l.Logger.WithFields(logrus.Fields{
		"pf": r.show.GetPlatform(),
		"id": r.show.GetRoomID(),
		"rn": roomName,
	}).Info("record start")

	err := r.parser.Parse(streamUrl, out)

	l.Logger.WithFields(logrus.Fields{
		"pf": r.show.GetPlatform(),
		"id": r.show.GetRoomID(),
	}).Infof("record stop: %+v", err)

	if err != nil {
		time.Sleep(5 * time.Second)
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

func (r *recorder) SubmitUploadTask(filepath string) {
	uploader.UploaderWorkerPool.AddTask(&uploader.UploadTask{
		Filepath: filepath,
		Tryout:   2,
	})
}
