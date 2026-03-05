package ui

import (
	"fmt"
	"io"
	"roland/logger"
	"strings"

	tea "charm.land/bubbletea/v2"
	"go.uber.org/zap"
)

type Window struct {
	ui       *Tui
	programm *tea.Program
	logger   *logger.Logger
}

func NewWindow(logger *logger.Logger) *Window {
	tui := newTui()

	return &Window{
		ui:       tui,
		programm: tea.NewProgram(tui),
		logger:   logger,
	}
}

func (w *Window) SetQuery(query string) {
	w.ui.query = query
}

func (w *Window) IsListening() bool {
	return w.ui.listening
}

func (w *Window) NewSession(key string, stderr, stdout io.ReadWriter) error {
	w.ui.Tabs = append(w.ui.Tabs, key)

	var builder strings.Builder

	out, err := io.ReadAll(stdout)
	if err != nil {
		w.logger.Error("failed read from command stdout",
			zap.Error(err))

		return fmt.Errorf("failed read from command stdout %w", err)
	}

	builder.Write(out)

	cerr, err := io.ReadAll(stderr)
	if err != nil {
		w.logger.Error("failed read from command stderr",
			zap.Error(err))

		return fmt.Errorf("failed read from command stderr %w", err)
	}

	builder.Write(cerr)

	w.ui.TabContent = append(w.ui.TabContent, builder.String())

	w.programm.Send(nil)

	return nil
}

func (w *Window) CloseSession(key string) {
	tabs := make([]string, len(w.ui.Tabs)-1)
	contents := make([]string, len(w.ui.TabContent)-1)

	index := 0

	 j := 0

	for i, tab := range w.ui.Tabs{
		if tab == key {
			index = i
		} else {
			tabs[j] = tab
			j++
		}
	}

	j = 0

	for i, content := range w.ui.TabContent {
		if i != index {
			contents[j] = content
			j++ 
		}
	}
}
