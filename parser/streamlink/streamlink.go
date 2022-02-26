package streamlink

import (
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/luxcgo/lifesaver/parser"
)

func init() {
	parser.SharedManager.Register(
		new(streamlink),
	)
}

type streamlink struct {
	cmd      *exec.Cmd
	cmdStdIn io.WriteCloser

	closeOnce sync.Once
	stop      chan struct{}
}

func (s *streamlink) New() parser.Parser {
	return &streamlink{
		stop: make(chan struct{}),
	}
}

func (s *streamlink) Stop() {
	s.closeOnce.Do(func() {
		close(s.stop)
	})
}

func (s *streamlink) Type() string {
	return "streamlink"
}

// streamlink -o "a.mp4"  https://www.huya.com/631275 best -f
func (s *streamlink) Parse(streamURL string, out string) (err error) {
	log.Println(streamURL)
	log.Println("work")
	s.cmd = exec.Command(
		"streamlink",
		"-o", out,
		streamURL,
		"best",
		"-f",
	)
	if s.cmdStdIn, err = s.cmd.StdinPipe(); err != nil {
		return err
	}
	if err = s.cmd.Start(); err != nil {
		s.cmd.Process.Kill()
		return err
	}
	go func() {
		<-s.stop
		s.cmd.Process.Kill()
	}()
	return s.cmd.Wait()
}
