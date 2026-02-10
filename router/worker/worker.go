package worker

import (
	"errors"
	"roland/logger"
	"sync"

	"go.uber.org/zap"
)

var ErrNotFound = errors.New("not found session")

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

func (w *Worker) StopCmd(sessionName string) error {
	w.mu.Lock()

	job, ok := w.jobs[sessionName]
	if !ok {
		w.logger.Error("not found session",
			zap.String("session_name", sessionName))

		w.mu.Unlock()

		return ErrNotFound
	}

	job.cancel()

	w.mu.Unlock()

	return nil
}
