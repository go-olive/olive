package uploader

import (
	"path/filepath"

	"github.com/luxcgo/lifesaver/config"
	l "github.com/luxcgo/lifesaver/log"
)

var UploaderWorkerPool = NewWorkerPool(1)

func init() {
	if !config.APP.UploadConfig.Enable {
		return
	}

	files, err := filepath.Glob(filepath.Join("*.flv"))
	if err != nil {
		l.Logger.Fatal(err)
	}
	tasks := make([]*UploadTask, len(files))
	for i, filepath := range files {
		tasks[i] = &UploadTask{
			Filepath: filepath,
			Tryout:   2,
		}
	}
	UploaderWorkerPool.AddTask(tasks...)
}

type WorkerPool struct {
	concurrency uint
	workers     []*worker
	uploadTasks chan *UploadTask
}

func NewWorkerPool(concurrency uint) *WorkerPool {
	wp := &WorkerPool{
		concurrency: concurrency,
		uploadTasks: make(chan *UploadTask, 1024),
	}
	for i := uint(0); i < wp.concurrency; i++ {
		w := newWorker(i)
		wp.workers = append(wp.workers, w)
	}
	return wp
}

func (wp *WorkerPool) AddTask(tasks ...*UploadTask) {
	for _, t := range tasks {
		wp.uploadTasks <- t
	}
}

func (wp *WorkerPool) Run() {
	for _, worker := range wp.workers {
		go worker.start(wp.uploadTasks)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.uploadTasks)
	for _, worker := range wp.workers {
		worker.stop()
		<-worker.done()
	}
}
