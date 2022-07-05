package uploader

import (
	"os/exec"
	"sync"

	l "github.com/go-olive/olive/src/log"
	"github.com/sirupsen/logrus"
)

type Uploader interface {
	proc()
	stop()
	done() <-chan struct{}
}

type TaskGroup struct {
	Filepath string
	PostCmds []*exec.Cmd
}

type uploader struct {
	taskGroup *TaskGroup
	cmd       *exec.Cmd
	closeOnce sync.Once
	stopChan  chan struct{}
	doneChan  chan struct{}
}

func NewUploader(taskGroup *TaskGroup) Uploader {
	return &uploader{
		taskGroup: taskGroup,
		stopChan:  make(chan struct{}),
		doneChan:  make(chan struct{}),
	}
}

func (u *uploader) proc() {
	defer close(u.doneChan)

	for _, postCmd := range u.taskGroup.PostCmds {
		select {
		case <-u.stopChan:
			return
		default:
			handler := DefaultTaskMux.MustGetHandler(postCmd.Path)
			err := handler.Process(
				&Task{
					Filepath: u.taskGroup.Filepath,
					StopChan: u.stopChan,
					Cmd:      postCmd,
				},
			)
			if err != nil {
				l.Logger.WithFields(logrus.Fields{
					"postCmd":  postCmd.String(),
					"filepath": u.taskGroup.Filepath,
				}).Error(err)
				return
			}
		}
	}
}

func (u *uploader) stop() {
	u.closeOnce.Do(func() {
		close(u.stopChan)
	})
}

func (u *uploader) done() <-chan struct{} {
	return u.doneChan
}
