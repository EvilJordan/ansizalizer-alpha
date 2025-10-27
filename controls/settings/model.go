package settings

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/settings/advanced"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/colors"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
	"github.com/Zebbeni/ansizalizer/controls/settings/alpha"
)

type Model struct {
	active State
	focus  State

	Colors     colors.Model
	Characters characters.Model
	Size       size.Model
	Advanced   advanced.Model
	Alpha      alpha.Model

	ShouldUnfocus bool
	ShouldClose   bool

	width int
}

func New(w int) Model {
	return Model{
		active: None,
		focus:  Colors,

		Colors:     colors.New(w),
		Characters: characters.New(w - 2),
		Size:       size.New(),
		Advanced:   advanced.New(w - 2),
		Alpha:      alpha.New(w - 2),

		ShouldUnfocus: false,
		ShouldClose:   false,

		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case Colors:
		return m.handleColorsUpdate(msg)
	case Characters:
		return m.handleCharactersUpdate(msg)
	case Size:
		return m.handleSizeUpdate(msg)
	case Advanced:
		return m.handleAdvancedUpdate(msg)
	case Alpha:
		return m.handleAlphaUpdate(msg)
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m.handleKeyMsg(keyMsg)
}

func (m Model) View() string {
	colorCtrls := m.Colors.View()
	charCtrls := m.Characters.View()
	sizeCtrls := m.Size.View()
	sampCtrls := m.Advanced.View()
	alfCtrls := m.Alpha.View()

	col := m.renderWithBorder(colorCtrls, Colors)
	char := m.renderWithBorder(charCtrls, Characters)
	siz := m.renderWithBorder(sizeCtrls, Size)
	sam := m.renderWithBorder(sampCtrls, Advanced)
	alf := m.renderWithBorder(alfCtrls, Alpha)

	return lipgloss.JoinVertical(lipgloss.Top, col, char, siz, sam, alf)
}
