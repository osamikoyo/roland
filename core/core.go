package core

import (
	"roland/config"
	"roland/logger"
	"roland/parser"
	"roland/router"
	"roland/ui"
)

type Core struct{
	inbound chan string
	parser *parser.Parser
	router *router.WorkerRouter
	window *ui.Window
}

func NewCore(cfg *config.Config, logger *logger.Logger) (*Core, error) {
	
}