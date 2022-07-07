package uploader

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-olive/olive/src/config"
	l "github.com/go-olive/olive/src/log"
	"github.com/sirupsen/logrus"
)

const (
	olivetrash   = "olivetrash"
	olivearchive = "olivearchive"
	olivebiliup  = "olivebiliup"
	oliveshell   = "oliveshell"
)

var DefaultHandlerFunc = TaskHandlerFunc(OliveDefault)

func init() {
	DefaultTaskMux.RegisterHandler(olivetrash, TaskHandlerFunc(OliveTrash))
	DefaultTaskMux.RegisterHandler(olivearchive, TaskHandlerFunc(OliveArchive))
	DefaultTaskMux.RegisterHandler(olivebiliup, TaskHandlerFunc(OliveBiliup))
	DefaultTaskMux.RegisterHandler(oliveshell, DefaultHandlerFunc)
}

func OliveTrash(t *Task) error {
	return os.Remove(t.Filepath)
}

func OliveArchive(t *Task) error {
	if err := os.MkdirAll("archive", os.ModePerm); err != nil {
		return err
	}
	return t.move("archive")
}

func (t *Task) move(dest string) error {
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dest, os.ModePerm)
		return err
	}

	base := filepath.Base(t.Filepath)
	dest = filepath.Join(dest, base)
	err := os.Rename(t.Filepath, dest)
	return err
}

func OliveBiliup(t *Task) error {
	l.Logger.WithFields(logrus.Fields{
		"filepath": t.Filepath,
	}).Info("upload start")

	doneChan := make(chan struct{})
	defer close(doneChan)

	var cmd *exec.Cmd
	if _, err := os.Stat(config.APP.UploadConfig.Filepath); errors.Is(err, os.ErrNotExist) {
		cmd = exec.Command(
			config.APP.UploadConfig.ExecPath,
			"upload",
			// "--tag=olive",
			"--limit=1",
			"--tid=21",
			t.Filepath,
		)
	} else {
		cmd = exec.Command(
			config.APP.UploadConfig.ExecPath,
			"upload",
			// "--tag=olive",
			"-c",
			config.APP.UploadConfig.Filepath,
			t.Filepath,
		)
	}

	go func() {
		select {
		case <-t.StopChan:
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			return
		case <-doneChan:
			return
		}
	}()

	resp, err := cmd.CombinedOutput()

	if err != nil {
		l.Logger.Debug("upload fail: ", err)
		OliveArchive(t)
		return err
	}

	if strings.Contains(string(resp), "投稿成功") {
		l.Logger.WithFields(logrus.Fields{
			"filepath": t.Filepath,
		}).Info("upload succeed")
		return nil
	}

	return nil
}

func OliveDefault(t *Task) error {
	doneChan := make(chan struct{})
	defer close(doneChan)

	cmd := exec.Command(t.Cmd.Args[0], t.Cmd.Args[1:]...)

	envFilepath := "FILE_PATH=" + t.Filepath
	cmd.Env = append([]string{envFilepath}, t.Cmd.Env...)
	cmd.Dir = t.Cmd.Dir

	go func() {
		select {
		case <-t.StopChan:
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			return
		case <-doneChan:
			return
		}
	}()

	resp, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	l.Logger.Infof("oliveshell success: %s", resp)
	return nil
}
