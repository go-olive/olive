package uploader

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-olive/olive/src/config"
	l "github.com/go-olive/olive/src/log"
	"github.com/sirupsen/logrus"
)

type Uploader interface {
	proc()
	stop()
	done() <-chan struct{}
}

type UploadTask struct {
	Filepath string
	Tryout   int64
	PostCmds []*exec.Cmd
}

type uploader struct {
	task      *UploadTask
	cmd       *exec.Cmd
	closeOnce sync.Once
	stopChan  chan struct{}
	doneChan  chan struct{}
}

func NewUploader(task *UploadTask) Uploader {
	return &uploader{
		task:     task,
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

func (u *uploader) proc() {
	resp, err := u.upload()
	if err != nil {
		l.Logger.Debug("upload fail: ", err)
		// return
	}
	if strings.Contains(string(resp), "投稿成功") {
		l.Logger.WithFields(logrus.Fields{
			"filepath": u.task.Filepath,
		}).Info("upload succeed")
		u.moveToTrash()
		return
	}

	if u.task.Tryout > 0 {
		UploaderWorkerPool.AddTask(u.task)
		return
	}

	l.Logger.WithFields(logrus.Fields{
		"filepath": u.task.Filepath,
	}).Info("upload fail and no more tryout")
	u.moveToArchive()
}

func (u *uploader) upload() (resp []byte, err error) {
	defer close(u.doneChan)

	u.task.Tryout--

	dir, _ := os.Getwd()
	dir = filepath.Join(dir, u.task.Filepath)

	l.Logger.WithFields(logrus.Fields{
		"filepath": dir,
	}).Info("upload start")

	if _, err := os.Stat(config.APP.UploadConfig.Filepath); errors.Is(err, os.ErrNotExist) {
		u.cmd = exec.Command(
			config.APP.UploadConfig.ExecPath,
			"upload",
			// "--tag=olive",
			"--limit=1",
			"--tid=21",
			u.task.Filepath,
		)
	} else {
		u.cmd = exec.Command(
			config.APP.UploadConfig.ExecPath,
			"upload",
			// "--tag=olive",
			"-c",
			config.APP.UploadConfig.Filepath,
			u.task.Filepath,
		)
	}

	go func() {
		select {
		case <-u.stopChan:
			if u.cmd.Process != nil {
				u.cmd.Process.Kill()
			}
			return
		case <-u.done():
			return
		}
	}()
	return u.cmd.CombinedOutput()
}

func (u *uploader) stop() {
	u.closeOnce.Do(func() {
		close(u.stopChan)
	})
}

func (u *uploader) done() <-chan struct{} {
	return u.doneChan
}

func (u *uploader) moveToArchive() {
	os.MkdirAll("archive", os.ModePerm)
	u.move("archive")
}

func (u *uploader) moveToTrash() {
	// os.MkdirAll("trash", os.ModePerm)
	// u.move("trash")
	os.Remove(u.task.Filepath)
}

func (u *uploader) move(dest string) {
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dest, os.ModePerm)
		if err != nil {
			l.Logger.Debug(err)
			return
		}
	}

	base := filepath.Base(u.task.Filepath)
	dest = filepath.Join(dest, base)
	err := os.Rename(u.task.Filepath, dest)
	if err != nil {
		l.Logger.Debug(err)
	}
}
