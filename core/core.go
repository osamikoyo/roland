package core

import (
	"context"
	"fmt"
	"roland/config"
	"roland/logger"
	"roland/parser"
	"roland/router"
	"roland/ui"

	"go.uber.org/zap"
)

type Core struct {
	inbound chan string
	parser  *parser.Parser
	router  *router.WorkerRouter
	window  *ui.Window
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
	}, nil
}

func (c *Core) Start(ctx context.Context) {
	
}
