package alpha

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Up Direction = iota
	Left
	Right
)

var (
	navMap = map[Direction]map[State]State{
		Right: {AlphaYes: AlphaNo},
		Left:  {AlphaNo: AlphaYes},
	}
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldUnfocus = true
	m.IsActive = false
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Up):
		m.IsActive = false
		m.ShouldUnfocus = true
	}
	return m, nil
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	if m.focus == AlphaYes {
		m.useAlpha = true
	}

	if m.focus == AlphaNo {
		m.useAlpha = false
	}

	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus
	switch m.focus {
	case AlphaYes:
		m.useAlpha = true
	case AlphaNo:
		m.useAlpha = false
	}
	return m, event.StartRenderToViewCmd // on enter, event.StartRenderToViewCmd is not firing when character type is switched, and then retuning to alpha
}
