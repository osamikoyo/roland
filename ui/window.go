package ui

import (
	"roland/logger"

	tea "github.com/charmbracelet/bubbletea"
)

type Window struct {
	ui       *Tui
	programm *tea.Program
	logger   *logger.Logger
}

func NewWindow(logger *logger.Logger) *Window {
}
