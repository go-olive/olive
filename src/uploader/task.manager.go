package uploader

var SharedTaskMux = &defaultTaskMux

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

func (this *TaskMux) GetHandler(pattern string) {

}

type TaskHandler interface {
	Process(*UploadTask) error
}

type TaskHandlerFunc func(t *UploadTask) error

func (f TaskHandlerFunc) Process(t *UploadTask) error {
	return f(t)
}
