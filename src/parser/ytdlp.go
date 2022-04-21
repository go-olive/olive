package parser

import (
	"io"
	"os/exec"
	"path/filepath"
	"sync"

	l "github.com/go-olive/olive/src/log"
	"github.com/sirupsen/logrus"
)

func init() {
	SharedManager.Register(
		new(ytdlp),
	)
}

type ytdlp struct {
	cmd      *exec.Cmd
	cmdStdIn io.WriteCloser

	closeOnce sync.Once
	stop      chan struct{}
}

func (p *ytdlp) New() Parser {
	return &ytdlp{
		stop: make(chan struct{}),
	}
}

func (p *ytdlp) Stop() {
	p.closeOnce.Do(func() {
		close(p.stop)
	})
}

func (p *ytdlp) Type() string {
	return "yt-dlp"
}

// yt-dlp -f "bv[height=1080]+ba/b" https://www.youtube.com/watch?v=f6PdkucL1hk
func (p *ytdlp) Parse(streamURL string, out string) (err error) {
	ext := filepath.Ext(out)
	out = out[0:len(out)-len(ext)] + ".mp4"

	l.Logger.WithFields(logrus.Fields{
		// "streamURL": streamURL,
		"out": out,
	}).Debug("yt-dlp working")

	p.cmd = exec.Command(
		"yt-dlp",
		"-o", out,
		"-f", "bv[height=1080]+ba/b",
		streamURL,
	)
	// s.cmd.Stderr = os.Stderr
	if p.cmdStdIn, err = p.cmd.StdinPipe(); err != nil {
		return err
	}
	if err = p.cmd.Start(); err != nil {
		p.cmd.Process.Kill()
		return err
	}
	go func() {
		<-p.stop
		p.cmdStdIn.Write([]byte("\x03"))
	}()
	return p.cmd.Wait()
}
