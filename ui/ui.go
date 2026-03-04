package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	query       lipgloss.Style
	doc         lipgloss.Style
	highlight   lipgloss.Style
	inactiveTab lipgloss.Style
	activeTab   lipgloss.Style
	window      lipgloss.Style
}

func newStyles(bgIsDark bool) *styles {
	lightDark := lipgloss.LightDark(bgIsDark)

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	highlightColor := lightDark(lipgloss.Color("#874BFD"), lipgloss.Color("#7D56F4"))

	s := new(styles)
	s.doc = lipgloss.NewStyle().
		Padding(1, 2, 1, 2)
	s.inactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	s.activeTab = s.inactiveTab.
		Border(activeTabBorder, true)
	s.window = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	return s
}

type Tui struct {
	Tabs       []string
	TabContent []string
	styles     *styles
	activeTab  int
	listening bool

	query string
}

func newTui() *Tui {
	return &Tui{
		styles: newStyles(true),
	}
}

func (t *Tui) Init() tea.Cmd {
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg == nil {
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return t, tea.Quit
		case "right", "l", "n", "tab":
			t.activeTab = min(t.activeTab+1, len(t.Tabs)-1)
			return t, nil
		case "left", "h", "p", "shift+tab":
			t.activeTab = max(t.activeTab-1, 0)
			return t, nil
		}
	}

	return t, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (t *Tui) View() tea.View {
	if t.styles == nil {
		return tea.NewView("")
	}

	doc := strings.Builder{}
	s := t.styles

	var (
		renderedTabs []string
		row          string
	)
	if len(t.Tabs) > 1 {
		for i, tab := range t.Tabs {
			var style lipgloss.Style
			isFirst, isLast, isActive := i == 0, i == len(t.Tabs)-1, i == t.activeTab
			if isActive {
				style = s.activeTab
			} else {
				style = s.inactiveTab
			}
			border, _, _, _, _ := style.GetBorder()
			if isFirst && isActive {
				border.BottomLeft = "│"
			} else if isFirst && !isActive {
				border.BottomLeft = "├"
			} else if isLast && isActive {
				border.BottomRight = "│"
			} else if isLast && !isActive {
				border.BottomRight = "┤"
			}
			style = style.Border(border)
			renderedTabs = append(renderedTabs, style.Render(tab))
		}
	}

	row = lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	doc.WriteString(t.styles.query.Render(t.query))
	doc.WriteString("\n")
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(s.window.Width((lipgloss.Width(row))).Render(t.TabContent[t.activeTab]))
	return tea.NewView(s.doc.Render(doc.String()))
}
