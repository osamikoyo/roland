package ui

import (
	"bytes"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#baacc7")).Align(lipgloss.Top)
	QueryStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

type Window struct {
	listening bool
	Query     string
}

func (w *Window) Init() tea.Cmd {
	return nil
}

func (w *Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return w, tea.Quit
		case "m":
			w.listening = !w.listening
			return w, nil
		}
	}

	return w, nil
}

func (w *Window) View() string {
	var buff bytes.Buffer

	title := "Roland\n"

	title = TitleStyle.Render()

	buff.WriteString(title)

	query := QueryStyle.Render(w.Query)

	return buff.String()
}
