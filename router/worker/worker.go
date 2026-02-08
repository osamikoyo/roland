package worker

import (
	"roland/logger"
	"sync"

	"go.uber.org/zap"
)

type Worker struct {
	logger *logger.Logger

	mu   sync.Mutex
	jobs map[string]*Job
}

func NewWorker(logger *logger.Logger) *Worker {
	return &Worker{
		logger: logger,
		jobs:   make(map[string]*Job),
	}
}

func (w *Worker) StartCmd(sessionName string, chunks []string) {
	job := SetupJob(chunks)

	w.mu.Lock()

	w.jobs[sessionName] = job

	w.mu.Unlock()

	if err := job.Run(); err != nil {
		w.logger.Error("failed start cmd",
			zap.Strings("chunks", chunks),
			zap.Error(err))

		return
	}
}

func (w *Worker) StopCmd(sessionName string) {
	w.mu.Lock()

	w.jobs[sessionName].cancel()

	w.mu.Unlock()
}
