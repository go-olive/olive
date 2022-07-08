package uploader

import (
	"os/exec"
	"path/filepath"

	"github.com/go-olive/olive/src/config"
	l "github.com/go-olive/olive/src/log"
)

var UploaderWorkerPool = NewWorkerPool(config.APP.CommanderPoolSize)

func init() {
	if !config.APP.UploadConfig.Enable {
		return
	}

	files, err := filepath.Glob(filepath.Join("*.flv"))
	if err != nil {
		l.Logger.Fatal(err)
	}
	tasks := make([]*TaskGroup, len(files))
	for i, filepath := range files {
		tasks[i] = &TaskGroup{
			Filepath: filepath,
			PostCmds: []*exec.Cmd{
				{Path: olivebiliup},
				{Path: olivetrash},
			},
		}
	}
	UploaderWorkerPool.AddTask(tasks...)
}

type WorkerPool struct {
	concurrency uint
	workers     []*worker
	uploadTasks chan *TaskGroup
	stopChan    chan struct{}
}

func NewWorkerPool(concurrency uint) *WorkerPool {
	wp := &WorkerPool{
		concurrency: concurrency,
		uploadTasks: make(chan *TaskGroup, 1024),
		stopChan:    make(chan struct{}),
	}
	for i := uint(0); i < wp.concurrency; i++ {
		w := newWorker(i)
		wp.workers = append(wp.workers, w)
	}
	return wp
}

func (wp *WorkerPool) AddTask(tasks ...*TaskGroup) {
	for _, t := range tasks {
		select {
		case <-wp.stopChan:
			return
		default:
			wp.uploadTasks <- t
		}
	}
}

func (wp *WorkerPool) Run() {
	for _, worker := range wp.workers {
		go worker.start(wp.uploadTasks)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.stopChan)
	close(wp.uploadTasks)
	for _, worker := range wp.workers {
		worker.stop()
		<-worker.done()
	}
}
