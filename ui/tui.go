package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Top).
			MarginTop(1).
			MarginLeft(10).
			Foreground(lipgloss.Color("5"))

	MutedStyle = lipgloss.NewStyle().
			Bold(true).
			MarginTop(1)

	QueryStyle = lipgloss.NewStyle().
	Align(lipgloss.Left).
	MarginTop(1)
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
	var builder strings.Builder

	title := "Roland\n"
	var listening string

	if w.listening {
		listening = MutedStyle.Foreground(lipgloss.Color("7")).Render("Unmuted\n")
	} else {
		listening = MutedStyle.Foreground(lipgloss.Color("1")).Render("muted\n")
	}

	builder.WriteString(TitleStyle.Render(title))

	builder.WriteString(listening)

	builder.WriteString(QueryStyle.Render(w.Query))

	return builder.String()
}
