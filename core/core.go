package core

import (
	"bytes"
	"context"
	"fmt"
	"roland/config"
	"roland/logger"
	"roland/parser"
	"roland/router"
	"roland/ui"

	"go.uber.org/zap"
)

const MaxSessionName = 32

type Core struct {
	parser   *parser.Parser
	router   *router.WorkerRouter
	window   *ui.Window
	logger   *logger.Logger
	sessions []string
}

func NewCore(cfg *config.Config, logger *logger.Logger, inbound chan string) (*Core, error) {
	logger.Info("setup parser")

	parser, err := parser.NewParser(cfg, logger)
	if err != nil {
		logger.Error("failed setup parser",
			zap.Error(err))

		return nil, fmt.Errorf("failed setup parser: %w", err)
	}

	logger.Info("setup router")

	router := router.NewWorkerRouter(logger)

	logger.Info("setup window")

	window := ui.NewWindow(logger)

	return &Core{
		window: window,
		router: router,
		parser: parser,
		logger: logger,
	}, nil
}

func (c *Core) Start(ctx context.Context, inbound chan string) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("core started")
		case phrase := <-inbound:
			c.logger.Info("new phrase",
				zap.String("phrase", phrase))

			if err := c.routePhrase(phrase); err != nil {
				c.logger.Error("failed route phrase",
					zap.String("phrase", phrase),
					zap.Error(err))
			}
		}
	}
}

func (c *Core) routePhrase(phrase string) error {
	c.logger.Info("route phrase",
		zap.String("phrase", phrase))

	c.logger.Info("parse phrase to request",
		zap.String("phrase", phrase))

	req, err := c.parser.Parse(phrase)
	if err != nil {
		c.logger.Error("failed parse phrase",
			zap.String("phrase", phrase),
			zap.Error(err))

		return fmt.Errorf("failed parse phrase %s: %w", phrase, err)
	}

	session := ""

	if len(phrase) > MaxSessionName {
		session = phrase[MaxSessionName:]
	} else {
		session = phrase
	}

	c.logger.Info("route request",
		zap.Any("request", req))

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)

	if err = c.router.RouteRequest(session, req, &stderr, &stdout); err != nil {
		c.logger.Error("failed route request",
			zap.String("session", session),
			zap.Any("req", req),
			zap.Error(err))

		return fmt.Errorf("failed route request: %w", err)
	}

	c.sessions = append(c.sessions, session)

	c.window.NewSession(session, &stderr, &stdout)

	return nil
}

func (c *Core) stopSession(session string) {
	if session == "" {
		return
	}

	c.router.StopSession(session)

	sessions := make([]string, len(c.sessions)-1)
	i := 0

	for _, s := range c.sessions {
		if s != session {
			sessions[i] = c.sessions[i]
			i++
		}
	}

	c.sessions = sessions
}
