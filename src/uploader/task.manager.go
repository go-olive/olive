package uploader

import "os/exec"

var DefaultTaskMux = &defaultTaskMux

var defaultTaskMux TaskMux

type TaskMux struct {
	m map[string]muxEntry
}

type muxEntry struct {
	h       TaskHandler
	pattern string
}

func (mux *TaskMux) RegisterHandler(pattern string, handler TaskHandler) {
	if pattern == "" {
		panic("task: invalid pattern")
	}
	if handler == nil {
		panic("task: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("task: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e
}

func (mux *TaskMux) MustGetHandler(pattern string) TaskHandler {
	if handler, ok := mux.m[pattern]; ok {
		return handler.h
	}
	return DefaultHandlerFunc
}

type Task struct {
	Filepath string
	StopChan chan struct{}
	Cmd      *exec.Cmd
}

type TaskHandler interface {
	Process(t *Task) error
}

type TaskHandlerFunc func(t *Task) error

func (f TaskHandlerFunc) Process(t *Task) error {
	return f(t)
}
