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
		m.ShouldUnfocus = true
	}
	return m, nil
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	if focus == AlphaYes {
		m.useAlpha = true
		m.focus = AlphaYes
	}

	if focus == AlphaNo {
		m.useAlpha = false
		m.focus = AlphaNo
	}

	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case AlphaYes:
		m.useAlpha = true
	case AlphaNo:
		m.useAlpha = false
	}
	return m, event.StartRenderToViewCmd
}
