package router

import (
	"bytes"
	"errors"
	"fmt"
	"roland/entity/request"
	"roland/logger"
	"strings"

	"go.uber.org/zap"
)

var (
	ErrNotFoundParameter = errors.New("not found parameter by key")
)

type WorkerRouter struct {
	logger *logger.Logger
	cmds   map[string]map[string]string
}

func NewWorkerRouter(logger *logger.Logger) *WorkerRouter {
	return &WorkerRouter{
		logger: logger,
	}
}

func (wr *WorkerRouter) RouteRequest(req *request.Request) error {
	actions := wr.cmds[req.Category]
	cmd := actions[req.Action]

	newcmd, err := setupCMD(cmd, req.Parameters)
	if err != nil {
		wr.logger.Error("failed setup cmd",
			zap.String("cmd", cmd),
			zap.Error(err))

		return fmt.Errorf("failed setup cmd: %w", err)
	}

}

func setupCMD(cmd string, parameters map[string]string) (string, error) {
	tokens := strings.Split(cmd, " ")

	newTokens := make([]string, len(tokens))

	for i, token := range tokens {
		if strings.HasPrefix(token, "$") {
			trimmedKey := strings.TrimPrefix(token, "$")

			value, ok := parameters[trimmedKey]
			if !ok {
				return "", ErrNotFoundParameter
			}

			newTokens[i] = value
		} else {
			newTokens[i] = token
		}
	}

	var buffer bytes.Buffer

	for _, token := range newTokens {
		buffer.WriteString(token)

		buffer.WriteRune(' ')
	}

	respcmd := buffer.String()
	respcmd = respcmd[:len(respcmd)]

	return respcmd, nil
}
