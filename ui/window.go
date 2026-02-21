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
	tui := newTui()
	
	return &Window{
		ui: tui,
		programm: tea.NewProgram(tui),
		logger: logger,
	}
}

func (w *Window) SetQuery(query string) {
	w.ui.Query = query
} 

func (w *Window) IsListening() bool {
	return w.ui.listening
}