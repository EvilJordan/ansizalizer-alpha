package alpha

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	Input State = iota
	Browser
	AlphaYes
	AlphaNo
	UseAlpha
)

type Model struct {
	focus State
	Browser browser.Model
	doUseAlpha bool
	useAlpha bool
	ShouldUnfocus bool
	IsActive bool
	width int
}

func New(w int) Model {

	return Model{
		focus: AlphaYes,
		Browser: browser.New(nil, w-2),
		doUseAlpha: true,
		useAlpha: true,
		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focus {
	case AlphaYes:
		m.useAlpha = true
	case AlphaNo:
		m.useAlpha = false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	content := make([]string, 0, 5)

	if m.doUseAlpha {
		content = append(content, m.drawAlphaOptions())
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m Model) ShouldOutputAlpha() bool {
	return m.useAlpha
}
