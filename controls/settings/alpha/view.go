package alpha

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawAlphaOptions() string {
	title := style.DimmedTitle.Copy().PaddingLeft(1).Render("Enable:")

	yesButton := style.NormalButtonNode
	if m.IsActive && m.focus == AlphaYes {
		yesButton = style.FocusButtonNode
	} else if m.useAlpha == true {
		yesButton = style.ActiveButtonNode
	}
	yesNode := yesButton.Render("Yes")
	yesNode = lipgloss.NewStyle().PaddingLeft(1).Render(yesNode)

	noButton := style.NormalButtonNode
	if m.IsActive && m.focus == AlphaNo {
		noButton = style.FocusButtonNode
	} else if m.useAlpha == false {
		noButton = style.ActiveButtonNode
	}
	noNode := noButton.Render("No")
	noNode = lipgloss.NewStyle().PaddingLeft(1).Render(noNode)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, yesNode, noNode)
}
