package process

import (
	"strings"
	"bufio"
	"github.com/charmbracelet/lipgloss"
)

const (
	AlphaPlaceholder string = "ALPHA"
	MagicTransparentPixel string = "[39;2;0;0;0;49;2;0;0;0m [0m"
)

func (m Renderer) outputStrings(rows ...string) (string) {
	content := ""
	if m.Settings.Alpha.ShouldOutputAlpha() {
		// replace ALPHA placeholder with a blank square (space)
		contentAlpha := strings.ReplaceAll(lipgloss.JoinVertical(lipgloss.Left, rows...), AlphaPlaceholder, MagicTransparentPixel)
		// iterate through the return of JoinVertical, separating by lines, trimming whitespace, and then recombining
		reader := strings.NewReader(contentAlpha)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			content += strings.TrimSpace(scanner.Text()) + "\n"
		}
	} else {
		content += lipgloss.JoinVertical(lipgloss.Left, rows...)
	}
	return content
}