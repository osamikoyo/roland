package router

import (
	"errors"
	"fmt"
	"roland/entity/request"
	"roland/logger"
	"roland/router/worker"
	"strings"

	"go.uber.org/zap"
)

var (
	ErrNotFoundParameter = errors.New("not found parameter by key")
)

type WorkerRouter struct {
	logger   *logger.Logger
	worker   *worker.Worker
	cmds     map[string]map[string]string
}

func NewWorkerRouter(logger *logger.Logger) *WorkerRouter {
	return &WorkerRouter{
		logger: logger,
	}
}

func (wr *WorkerRouter) RouteRequest(session string, req *request.Request) error {
	actions := wr.cmds[req.Category]
	cmd := actions[req.Action]

	chunks, err := setupCMD(cmd, req.Parameters)
	if err != nil {
		wr.logger.Error("failed setup cmd",
			zap.String("cmd", cmd),
			zap.Error(err))

		return fmt.Errorf("failed setup cmd: %w", err)
	}

	go wr.worker.StartCmd(session, chunks)

	return nil
}

func setupCMD(cmd string, parameters map[string]string) ([]string, error) {
	tokens := strings.Split(cmd, " ")

	newTokens := make([]string, len(tokens))

	for i, token := range tokens {
		if strings.HasPrefix(token, "$") {
			trimmedKey := strings.TrimPrefix(token, "$")

			value, ok := parameters[trimmedKey]
			if !ok {
				return nil, ErrNotFoundParameter
			}

			newTokens[i] = value
		} else {
			newTokens[i] = token
		}
	}

	return newTokens, nil
}

func (wr *WorkerRouter) StopSession(session string) {
	wr.worker.StopCmd(session)
}